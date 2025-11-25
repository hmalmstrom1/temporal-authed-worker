#!/bin/bash

set -eu -o pipefail

# Wait for Cassandra to be ready (handled by Docker Compose healthcheck)
echo "Cassandra is ready."

# Drop keyspaces if they exist
# echo Drop keyspaces if they exist \(skipped because cqlsh is not available\)
# cqlsh cassandra 9042 -e "DROP KEYSPACE IF EXISTS temporal; DROP KEYSPACE IF EXISTS temporal_visibility;"

# Create keyspaces
echo "Creating keyspace 'temporal'..."
temporal-cassandra-tool --endpoint cassandra create -k temporal --replication-factor 1

echo "Creating keyspace 'temporal_visibility'..."
temporal-cassandra-tool --endpoint cassandra create -k temporal_visibility --replication-factor 1

# Setup schema
echo "Setting up schema for 'temporal'..."
temporal-cassandra-tool --endpoint cassandra -k temporal setup-schema -v 0.0
temporal-cassandra-tool --endpoint cassandra -k temporal update-schema -d /etc/temporal/schema/cassandra/temporal/versioned

echo "Setting up schema for 'temporal_visibility'..."
# temporal_visibility schema is not present in the repo for Cassandra in v1.29.1
# We will use ElasticSearch for visibility instead.
# But we still create the keyspace just in case.
temporal-cassandra-tool --endpoint cassandra -k temporal_visibility setup-schema -v 0.0

# Setup ElasticSearch indices
echo "Setting up ElasticSearch indices..."
ES_SERVER="http://elasticsearch:9200"
ES_USER=""
ES_PWD=""
ES_VIS_INDEX="temporal_visibility_v1_dev"

# Wait for ES to be ready
until curl --silent --fail "${ES_SERVER}" > /dev/null; do
    echo "Waiting for Elasticsearch to start up..."
    sleep 5
done

echo "Elasticsearch started."

SETTINGS_FILE="/etc/temporal/schema/elasticsearch/visibility/cluster_settings_v7.json"
TEMPLATE_URL="${ES_SERVER}/_template/temporal_visibility_v1_template"
SCHEMA_FILE="/etc/temporal/schema/elasticsearch/visibility/index_template_v7.json"
INDEX_URL="${ES_SERVER}/${ES_VIS_INDEX}"

# Apply settings (might fail if already applied, ignore error)
curl -X PUT "${ES_SERVER}/_cluster/settings" -H "Content-Type: application/json" --data-binary "@${SETTINGS_FILE}" || true

# Apply template
curl --fail -X PUT "${TEMPLATE_URL}" -H 'Content-Type: application/json' --data-binary "@${SCHEMA_FILE}"

# Create index
curl -X PUT "${INDEX_URL}" || true

echo "Schema setup complete."
