#!/bin/bash

docker-compose up -d

sleep 30

master_db="pgmaster"
replicas=("pgslave_1" "pgslave_2")

# Configure the master database
echo "Configuring $master_db as master..."
docker exec -i "$master_db" psql -U postgres -d social_network -c "create role replicator with login replication password 'pass';"
docker exec -i "$master_db" psql -U postgres -d social_network -c "GRANT CONNECT ON DATABASE postgres TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;"
docker exec -i "$master_db" psql -U postgres -d social_network -c "DROP PUBLICATION IF EXISTS my_publication;"
docker exec -i "$master_db" psql -U postgres -d social_network -c "CREATE PUBLICATION my_publication FOR ALL TABLES;"

# Configure each replica to follow the master
for replica in "${replicas[@]}"; do
    echo "Configuring $replica to replicate from $master_db..."

    # Construct the connection string
    conninfo="host=$master_db port=5431 user=postgres password=postgres dbname=social_network"

    # Configure the subscription on the replica
    docker exec -i "$replica" psql -U postgres -d social_network -c "DROP SUBSCRIPTION IF EXISTS ${replica}_subscription;"
    docker exec -i "$replica" psql -U postgres -d social_network -c "CREATE SUBSCRIPTION ${replica}_subscription CONNECTION 'dbname=social_network host=$master_db user=replicator password=pass' PUBLICATION my_publication;"

    echo "$replica configured to replicate from $master_db."
done

echo "Master and replicas configured successfully for logical replication."