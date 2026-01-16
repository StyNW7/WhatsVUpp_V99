# Amen lulus TPA Network 25-1

# 🚀 Multi-Service Application with Go, Rust, and Kubernetes

This is a full-stack microservices project composed of backend services written in Go and Rust, a frontend built with Nginx, and orchestrated using Docker Compose and Kubernetes. The project supports local development with Docker and CI/CD integration for deployment.

---

## 📁 Project Structure

```
.
├── .github/workflows/       # GitHub Actions for CI/CD
├── backend-go/              # Backend service written in Go
├── backend-rust/            # Backend service written in Rust (Debian Bookworm-based)
├── docker/                  # Docker scripts and private registry setup
├── frontend/                # Frontend served via Nginx
├── init-db/                 # Database initialization scripts
├── k8s/                     # Kubernetes manifests
├── notes/                   # Notes and helper scripts
├── .gitignore               # Git ignored files
├── README.md                # Project documentation (this file)
├── ci-cd.yml                # CI/CD pipeline configuration
├── docker-compose.yml       # Docker Compose for local development
```

---

## 🛠️ Technologies Used

* **Go** for backend service
* **Rust** for an alternative backend service
* **Nginx** for serving frontend
* **Docker & Docker Compose** for containerization and local dev
* **Kubernetes** for orchestration
* **Prometheus** for metrics and monitoring
* **GitHub Actions** for CI/CD
* **Private Docker Registry** setup

---

## 🚧 Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/your-repo/project-name.git
cd project-name
```

### 2. Build and run using Docker Compose

Make sure Docker is installed and running:

```bash
docker-compose up --build
```

This will:

* Build Go and Rust backends
* Set up the frontend via Nginx
* Initialize the database
* Connect all services using Docker Compose network

### 3. Database Initialization

Ensure the database is seeded correctly using scripts in the `init-db` folder. This runs automatically via Docker Compose.

---

## 📦 Backend Services

### Go Backend (`backend-go`)

* Connects to the initialized database
* Updated to remove unused ports
* Kubernetes-ready with metrics exposed for Prometheus

### Rust Backend (`backend-rust`)

* Uses `debian:bookworm-slim` base image
* Dockerfile optimized and tested
* Integrated into the Docker Compose and Kubernetes setups

---

## 🌐 Frontend

### Nginx-based frontend (`frontend`)

* Serves static content
* Configured via `nginx.conf`
* Connected with backend services through reverse proxy if needed

---

## ☸️ Kubernetes Setup

Inside the `k8s/` folder:

* Deployments and services for both backend services
* Prometheus configuration for monitoring
* Secrets and config maps for environment management

To apply to your cluster:

```bash
kubectl apply -f k8s/
```

---

## 🔐 Private Docker Registry

Set up scripts located in:

* `.github/workflows/`
* `docker/`
* `notes/`

Run the setup script:

```bash
sh docker/setup-private-registry.sh
```

Push your images:

```bash
docker tag backend-go localhost:5000/backend-go
docker push localhost:5000/backend-go
```

---

## ✅ CI/CD

CI/CD pipeline is configured via:

* `.github/workflows/`
* `ci-cd.yml`

Pipeline performs:

* Build and test of backend/frontend
* Docker build and push
* Optional: deploy to staging/production via Kubernetes

---

## 🙋‍♂️ Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## 📬 Contact

Developed by: NW, NV, AT, NP, VW, KF For questions or collaboration, please contact via email or open an issue.
