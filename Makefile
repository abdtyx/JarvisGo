all: build start
	
build:
	sudo docker build -t jarvisgo .

start:
	echo '**Warning**: This is an one-time docker container'
	sudo docker run -d -t --name Jarvis -p 5700:5700 -p 5701:5701 -p 3306:3306 -p 6379:6379 jarvisgo

restart: build stop clean start

stop:
	sudo docker stop Jarvis

clean:
	sudo docker rm Jarvis

prune:
	sudo docker image prune

run:
	go run cmd/server/main.go