package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/checkpw", routes.checkpw)

	// Start the server
	if err := r.Run(); err != nil {
		log.Fatal("Server Run Failed:", err)
	} else {
		fmt.Println("Server is running on port 8080")
	}
}
