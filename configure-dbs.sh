wal_level="logical"
max_number_of_replicas=4
max_wal_senders=8

# Define all PostgreSQL instances
databases=("pgmaster" "pgslave_1" "pgslave_2")

for db in "${databases[@]}"; do
    # Apply configuration changes
    ...
done

