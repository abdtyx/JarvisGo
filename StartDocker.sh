#!/bin/bash
echo '**Warning**: This is an one-time docker container'
sudo docker stop Jarvis
sudo docker rm Jarvis
sudo docker run -d -t --name Jarvis -p 5700:5700 -p 5701:5701 jarvisgo
