#!/bin/bash

#create ssl directory if it doesn't exist
mkdir -p /etc/nginx/ssl

# Generate self-signed certificate and key in PEM format
openssl req -x509 -newkey rsa:2048 -keyout /etc/nginx/ssl/key.pem -out /etc/nginx/ssl/cert.pem -days 365 -nodes \
    -subj "/C=IL/ST=Tel Aviv/L=Tel Aviv/O=Organization/CN=localhost"

echo "Self-signed certificate generated:"
echo "  Certificate: cert.pem"
echo "  Private Key: key.pem"