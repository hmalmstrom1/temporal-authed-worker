# Temporal Secure Worker

## Overview
Warning: this is an expiremental implementation and is not recommended for production use.
It's just a POC to see if I can get an OAuth 2.0 flow working with Temporal Workers paired with TLS.

This repository contains a secure worker for Temporal that uses OAuth 2.0 for authentication.


## Running locally


### Docker Compose Components

- Ory Hydra
- Ory Kratos
- Temporal
- Worker

### Docker compose usage:

To start the services:

```bash
docker compose up --build
```

This will create a self signed certificate for the worker as well from the scripts/gen-certs.sh script.

To stop the services:

```bash
docker compose down
```

Ory Kratos and Hydra are available at:

- Kratos: http://localhost:4433
- Hydra: http://localhost:4444

Temporal is available at:

- Temporal: http://localhost:7233

The worker will be using the OAuth 2.0 flow to authenticate with Temporal.

Temporal will utilize the OAuth 2.0 flow to authenticate with the worker through an Authorizer and ClaimMapper.

Creating new identities in Ory Kratos and Hydra is done through the Ory Kratos and Hydra admin APIs (see scripts/setup-ory.sh where we add the worker)

## WorkerEnvironment Variables

- TEMPORAL_ADDRESS: The address of the Temporal server.
- TEMPORAL_NAMESPACE: The namespace of the Temporal server.
- OAUTH_CLIENT_ID: The client ID of the OAuth 2.0 client.
- OAUTH_CLIENT_SECRET: The client secret of the OAuth 2.0 client.
- OAUTH_TOKEN_URL: The token URL of the OAuth 2.0 server.
- TLS_CERT_PATH: The path to the TLS certificate of the worker.
- TLS_KEY_PATH: The path to the TLS key of the worker.
- TLS_SERVER_ROOT_CA: The path to the TLS root CA of the server.

Example running the worker with this configuration:

```bash
export TEMPORAL_ADDRESS="localhost:7233"
export TEMPORAL_NAMESPACE="default"
export OAUTH_CLIENT_ID="temporal-worker"
export OAUTH_CLIENT_SECRET="secret-temporal-worker"
export OAUTH_TOKEN_URL="http://localhost:4444/oauth2/token"
export TLS_CERT_PATH="./certs/client.cert"
export TLS_KEY_PATH="./certs/client.key"
export TLS_SERVER_ROOT_CA="./certs/ca.cert"

cd worker-csharp

dotnet run
```


