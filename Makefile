all: build start
	
build:
	sudo docker build -t jarvisgo .

start:
	echo '**Warning**: This is an one-time docker container'
	sudo docker run -d -t --name Jarvis --network host jarvisgo

restart: build stop clean start

stop:
	sudo docker stop Jarvis

clean:
	sudo docker rm Jarvis

prune:
	sudo docker image prune

run:
	go run cmd/server/main.go