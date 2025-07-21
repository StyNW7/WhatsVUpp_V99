# Set Up SSL Certificates (Required for Production)

mkdir -p certs
openssl req -newkey rsa:4096 -nodes -sha256 -keyout certs/domain.key \
  -x509 -days 365 -out certs/domain.crt \
  -subj "/CN=myregistry.example.com" \
  -addext "subjectAltName=DNS:myregistry.example.com,DNS:localhost,IP:127.0.0.1"

# Configure Docker Daemon to Trust Your Registry