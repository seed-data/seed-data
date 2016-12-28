
clean:
	docker rmi -f seed-data || true

build: clean
	docker-compose build -t seed-data --no-cache .

stop:
	docker-compose down

start:
	docker-compose up
