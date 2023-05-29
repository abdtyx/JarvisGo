package main

import (
	"github.com/abdtyx/JarvisGo/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/", handler.Handler)

	r.Run(":8000")
}
