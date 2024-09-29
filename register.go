package main

import (
	// "fmt"

	// "github.com/golang-jwt/jwt"
	// "gorm.io/gorm"
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"errors"
	"strings"
	
)

type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}


func HashPassword(PasswordHash string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(PasswordHash), 14)
	return string(bytes), err
}

func HashCompare(compare string, passwordhash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordhash), []byte(compare))
	return err
}

func register(r *gin.Engine, db *gorm.DB) {
	r.POST("/register",  func(c *gin.Context) {
		var user Users

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if rec := db.First(&user).Where("username = ?", user.Username); rec != nil {
			c.JSON(http.StatusConflict, gin.H{"message": "Username already taken"})
			return
		}

		var err error
		user.PasswordHash, err = HashPassword(user.PasswordHash)

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		user.Role = "User"
		result := db.Create(&user)

		if result.RowsAffected == 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "Didn't work"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User " + user.Username + " succesfuly created"})
	})
}

func createToken(user Users) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["userID"] = user.ID
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["role"] = "user"
	err := godotenv.Load()
	if err != nil {
		panic("No .env file found")
	}
	key := os.Getenv("SECRETKEY")
	if key == "" {
		return "", fmt.Errorf("SECRETKEY environment variable is not set")
	}
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func jwtMiddleWare(c *gin.Context) error {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "No or invalid authorization given"})
		return errors.New("no or invalid authorization given")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// Parse the token
	key := os.Getenv("SECRETKEY")
	if key == "" {
		return fmt.Errorf("SECRETKEY environment variable is not set")
	}
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(key), nil
    })

	if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
        return err
    }

	if !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
        return errors.New("invalid token")
    }

	c.Set("claims", claims)

    // Continue to the next handler
    return nil

}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve the user's role from the context
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "No role found"})
            c.Abort()
            return
        }

        // Check if the user's role matches one of the allowed roles
        for _, role := range roles {
            if userRole == role {
                c.Next() // Role is valid, proceed to the next handler
                return
            }
        }

        // If we reach this point, the role is not allowed
        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
        c.Abort()
    }
}


func jwtAuthMiddleware(c *gin.Context) {
    if err := jwtMiddleWare(c); err != nil {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }
    c.Next() // If everything is fine, proceed to the next handler
}

func login(r *gin.Engine, db *gorm.DB) {
	r.POST("/login", func(c *gin.Context) {
		var user Users

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mdp := user.PasswordHash
		
		if rec := db.Where("username = ?", user.Username).First(&user); rec == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist"})
			return
		}

		if HashCompare(mdp, user.PasswordHash) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong password"})
			return
		}

		token, err := createToken(user)
		
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Access Token wasn't generated"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}

