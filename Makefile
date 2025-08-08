watch_server:
	docker container logs server
build:
	docker-compose up -d --build
down:
	docker-compose down -v
exec_db:
	docker container exec -it postgres psql -U postgres -d dvdrental

