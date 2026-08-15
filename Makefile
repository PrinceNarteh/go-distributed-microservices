up:
	docker-compose up --build

down:
	docker-compose down

logs:
	docker-compose logs -f

build:
	docker-compose build

restart:
	docker-compse down && docker-compose up --build
