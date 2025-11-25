#!/bin/bash
set -e

# Wait for Hydra to be ready
echo "Waiting for Hydra Admin API..."
until curl -s http://localhost:4445/health/ready > /dev/null; do
  sleep 2
done

echo "Creating OAuth2 Client 'temporal-worker'..."

curl -X POST http://localhost:4445/clients \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "temporal-worker",
    "client_secret": "secret-temporal-worker",
    "grant_types": ["client_credentials"],
    "response_types": ["token"],
    "scope": "openid offline_access temporal:writer",
    "audience": ["temporal"],
    "token_endpoint_auth_method": "client_secret_post"
  }'

echo ""
echo "Client created successfully."
