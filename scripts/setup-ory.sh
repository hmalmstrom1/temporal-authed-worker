#!/bin/bash

set -e

echo "Waiting for Hydra to be ready..."
until curl -s http://localhost:4445/health/ready > /dev/null; do
  sleep 2
done

echo "Hydra is ready."

echo "Creating OAuth2 client for Temporal Worker..."
CLIENT_ID="temporal-worker"
CLIENT_SECRET="secret-temporal-worker"

# Delete if exists
docker compose exec hydra hydra clients delete $CLIENT_ID --endpoint http://127.0.0.1:4445 || true

# Create client
docker compose exec hydra hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id $CLIENT_ID \
    --secret $CLIENT_SECRET \
    --grant-types client_credentials \
    --response-types token,code,id_token \
    --scope openid,offline,worker \
    --audience temporal

echo "Client created."
echo "Client ID: $CLIENT_ID"
echo "Client Secret: $CLIENT_SECRET"

echo "Performing Client Credentials Flow to get a token..."
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:4444/oauth2/token \
    -u "$CLIENT_ID:$CLIENT_SECRET" \
    -d "grant_type=client_credentials" \
    -d "scope=openid,worker")

echo "Token Response:"
echo $TOKEN_RESPONSE | jq .

ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r .access_token)

if [ "$ACCESS_TOKEN" != "null" ]; then
    echo "Successfully obtained access token!"
else
    echo "Failed to obtain access token."
    exit 1
fi
