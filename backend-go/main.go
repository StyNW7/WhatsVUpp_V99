package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	ID        int       `json:"id" db:"id"`
	Sender    string    `json:"sender" db:"sender"`
	Content   string    `json:"content" db:"content"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Expose metrics
	promhttp.Handler().ServeHTTP(w, r)
}

func instrumentHandler(inner http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start time to measure request duration
		start := time.Now()

		// Create a response writer to capture status code
		rw := &statusCodeResponseWriter{ResponseWriter: w}
		inner(rw, r)

		// Track the request duration
		duration := time.Since(start).Seconds()

		// Increment the counters
		httpRequests.WithLabelValues(r.Method, fmt.Sprint(rw.statusCode)).Inc()
		httpDuration.WithLabelValues(r.Method, fmt.Sprint(rw.statusCode)).Observe(duration)
	}
}

// Create a custom response writer to capture the status code
type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusCodeResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func main() {
	var err error
	connStr := "postgres://postgres:password@db:6543/chatdb?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error creating table users: ", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Error creating table messages: ", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/register", instrumentHandler(registerHandler)).Methods("POST")
	router.HandleFunc("/api/login", instrumentHandler(loginHandler)).Methods("POST")
	router.HandleFunc("/api/messages", instrumentHandler(getMessages)).Methods("GET")
	router.HandleFunc("/api/messages", instrumentHandler(postMessage)).Methods("POST")
	router.HandleFunc("/metrics", metricsHandler).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}
func getMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, sender, content, timestamp FROM messages ORDER BY timestamp ASC")
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.Timestamp); err != nil {
			http.Error(w, "Error scanning message", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO messages (sender, content) VALUES ($1, $2)", msg.Sender, msg.Content)
	if err != nil {
		http.Error(w, "Failed to store message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
	if err != nil || exists {
		http.Error(w, "User already exists or database error.", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password.", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error creating user.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found.", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error.", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid password.", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Error generating token.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
