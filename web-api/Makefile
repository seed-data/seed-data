
clean:
	docker rmi -f web-api || true

build: clean
	docker build -t web-api --no-cache .

stop:
	docker stop web-api

start:
	docker run -p 127.0.0.1:5000:80 --name web-api -t web-api:latest
