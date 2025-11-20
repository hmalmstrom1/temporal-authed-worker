#!/bin/bash
set -e

# Create certs directory if it doesn't exist
mkdir -p ../certs
cd ../certs

echo "Generating Root CA..."
# Generate Root CA key
openssl genrsa -out ca.key 4096
# Generate Root CA certificate
openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=State/L=City/O=Organization/CN=Temporal Root CA" -days 365 -out ca.cert

echo "Generating Server Certificate..."
# Generate Server key
openssl genrsa -out server.key 4096
# Generate Server CSR
openssl req -new -key server.key -out server.csr -config <(cat /etc/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:localhost,DNS:temporal,IP:127.0.0.1")) -subj "/C=US/ST=State/L=City/O=Organization/CN=temporal" -reqexts SAN
# Sign Server Certificate with Root CA
openssl x509 -req -in server.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out server.cert -days 365 -sha256 -extfile <(cat /etc/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:localhost,DNS:temporal,IP:127.0.0.1")) -extensions SAN

echo "Generating Client Certificate..."
# Generate Client key
openssl genrsa -out client.key 4096
# Generate Client CSR
openssl req -new -key client.key -out client.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=client"
# Sign Client Certificate with Root CA
openssl x509 -req -in client.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out client.cert -days 365 -sha256

echo "Certificates generated in certs/"
