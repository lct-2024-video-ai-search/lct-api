mkdir data
docker-compose -d
docker exec -i clickhouse-server clickhouse-client < init.sql
