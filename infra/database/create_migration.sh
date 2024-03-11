docker exec -i finder_api_container migrate create -ext sql -dir ./infra/database/migrations/mariaDB $1

user=$(whoami)

sudo chown -R $user:$user ./migrations
