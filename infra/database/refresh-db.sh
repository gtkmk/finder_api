docker exec -i finder_db mariadb -u root -p$1 -D dbapp < db-fixtures/pre-fixtures.sql

docker exec -i finder_api_container migrate -path ./infra/database/migrations/mariaDB -database "mysql://root:$1@tcp(finder_db)/dbfinder" down

docker exec -i finder_api_container migrate -path ./infra/database/migrations/mariaDB -database "mysql://root:$1@tcp(finder_db)/dbfinder" up

docker exec -i finder_db mariadb -u root -p$1 -D dbfinder < db-fixtures/fixtures.sql
