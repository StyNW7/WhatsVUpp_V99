# Create directories
mkdir -p auth data

# Create htpasswd file with your credentials
# Username: [INITIAL], Password: G4c0r!
docker run --rm --entrypoint htpasswd httpd:2 -Bbn 'NW' 'G4c0r!' > auth/htpasswd

echo "Registry setup complete!"
echo "Start with: docker-compose up -d"
echo "Registry will be available at: localhost:5000"