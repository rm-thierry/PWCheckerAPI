package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rm-thierry/PWCheckerAPI/routes"
)

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/checkpw", routes.CheckPW)

	// Start the server on port 2025
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed:", err)
	} else {
		fmt.Println("Server is running on port 2025")
	}
}
