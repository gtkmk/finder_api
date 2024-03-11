docker exec -i finder_api_container migrate -path ./infra/database/migrations/mariaDB -database "mysql://root:$1@tcp(finder_db)/dbfinder" up
