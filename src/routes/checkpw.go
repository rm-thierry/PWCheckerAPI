package routes

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	passwordSet   map[string]struct{}
	passwordMutex sync.RWMutex
	initialized   bool
	initMutex     sync.Mutex
)

func initializePasswordSet() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	passwordSet = make(map[string]struct{})

	path := filepath.Join(".", "data", "pw.txt")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passwordSet[scanner.Text()] = struct{}{}
	}

	initialized = true
	return scanner.Err()
}

func CheckPW(c *gin.Context) {
	var requestBody struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if requestBody.Password == "" {
		c.JSON(400, gin.H{"error": "Password cannot be empty"})
		return
	}

	if !initialized {
		if err := initializePasswordSet(); err != nil {
			fmt.Println("Error initializing password set:", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
	}

	// Initialize strength score
	strengthScore := 0

	// Check if the password is less than 8 characters
	if len(requestBody.Password) < 8 {
		c.JSON(200, gin.H{"result": "Password is weak", "strength": strengthScore})
		return
	}

	// Check if the password contains only digits
	onlyDigits := true
	for _, char := range requestBody.Password {
		if char < '0' || char > '9' {
			onlyDigits = false
			break
		}
	}
	if onlyDigits {
		// Only digits passwords get minimal score of 10
		strengthScore = 10
		c.JSON(200, gin.H{"result": "Password is weak", "strength": strengthScore})
		return
	}

	// Calculate strength score
	// Check length (up to 30 points)
	lengthScore := len(requestBody.Password)
	if lengthScore > 30 {
		lengthScore = 30
	}
	strengthScore += lengthScore

	// Check for uppercase letters (10 points)
	for _, char := range requestBody.Password {
		if char >= 'A' && char <= 'Z' {
			strengthScore += 10
			break
		}
	}

	// Check for lowercase letters (10 points)
	for _, char := range requestBody.Password {
		if char >= 'a' && char <= 'z' {
			strengthScore += 10
			break
		}
	}

	// Check for special characters (20 points)
	for _, char := range requestBody.Password {
		if (char >= '!' && char <= '/') || (char >= ':' && char <= '@') || (char >= '[' && char <= '`') || (char >= '{' && char <= '~') {
			strengthScore += 20
			break
		}
	}

	// Check for numbers (10 points)
	for _, char := range requestBody.Password {
		if char >= '0' && char <= '9' {
			strengthScore += 10
			break
		}
	}

	// Return early with strength percentage if score is calculated
	if strengthScore < 30 {
		c.JSON(200, gin.H{"result": "Password is very weak", "strength": strengthScore})
		return
	} else if strengthScore < 50 {
		c.JSON(200, gin.H{"result": "Password is weak", "strength": strengthScore})
		return
	} else if strengthScore < 70 {
		c.JSON(200, gin.H{"result": "Password is moderate", "strength": strengthScore})
		return
	} else {
		// Check if the password is in the set
		passwordMutex.RLock()
		_, found := passwordSet[requestBody.Password]
		passwordMutex.RUnlock()

		if found {
			c.JSON(200, gin.H{"result": "Password is weak", "strength": strengthScore})
		} else {
			c.JSON(200, gin.H{"result": "Password is strong", "strength": strengthScore})
		}
	}
}
