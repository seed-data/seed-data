
clean:
	docker rmi -f web-frontend || true

build: clean
	docker build -t web-frontend --no-cache .

stop:
	docker stop web-frontend

start:
	docker run -p 127.0.0.1:5001:80 --name web-frontend -t web-frontend:latest
