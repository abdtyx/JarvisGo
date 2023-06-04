package main

import (
	"log"

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

	r.POST("/", h.Handle)

	r.Run(":5701")
}
