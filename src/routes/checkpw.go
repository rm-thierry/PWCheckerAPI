package routes

import "github.com/gin-gonic/gin"

func checkpw(gin *gin.Context) {
	// Get the password from the request
	password := gin.Query("password")

	// Check if the password is correct
	if password == "your_password" {
		gin.JSON(200, gin.H{
			"message": "Password is correct",
		})
	} else {
		gin.JSON(401, gin.H{
			"message": "Password is incorrect",
		})
	}
}
