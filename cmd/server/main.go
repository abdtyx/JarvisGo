package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdtyx/JarvisGo/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize handler
	h, err := handler.InitHandler()
	if err != nil {
		log.Fatalln("main: Failed to initialize handler")
	}

	// Start listening
	r := gin.Default()

	// bot router register
	r.POST("/", h.Handle)

	// Process run
	go r.Run(":5701")

	// Gracefully handle signal
	sigCh := make(chan os.Signal, 1)

	// Registered signals
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Waiting for signal
	<-sigCh

	// Shutdown
	h.Shutdown()
}
