
clean:
	docker rmi -f worker || true

build: clean
	docker build -t worker --no-cache .

stop:
	docker stop worker

start:
	docker run -p 127.0.0.1:5002:80 --name worker -t worker:latest
