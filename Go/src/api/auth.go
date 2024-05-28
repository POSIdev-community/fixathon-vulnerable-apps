package api

import (
	"context"
	"fmt"
	"net/http"
	"phdays-app/src/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const JwtCookie = "jwt"

var secret = []byte("secret")
var logger = logrus.WithContext(context.Background())

const loggedInUser = "userId"

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// login is a handler that parses a form and checks for specific data.
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	redirect := c.PostForm("redirect_to")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match
	user, err := models.GetUser(username, password)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	tokenString, err := createToken(fmt.Sprint(user.UserId))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating token")
		return
	}

	logger.Info(fmt.Sprintf("Token created: %s\n", tokenString))
	c.SetCookie(JwtCookie, tokenString, 3600, "/", "localhost", false, false)
	c.Redirect(http.StatusFound, redirect)
}

// Function to verify JWT tokens
func AuthRequired(c *gin.Context) {
	// Retrieve the token from the cookie
	tokenString, err := c.Cookie(JwtCookie)
	if err != nil {
		logger.Error("Token missing in cookie")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Verify the token
	token, err := verifyToken(tokenString)
	if err != nil {
		logger.Error("Token verification failed")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Print information about the verified token
	logger.Info("Token verified successfully.")
	userId, _ := token.Claims.GetSubject()
	c.Set(loggedInUser, userId)
	// Continue with the next middleware or route handler
	c.Next()
}

// Logout is the handler called for the user to log out.
func Logout(c *gin.Context) {
	c.Set(loggedInUser, "")
	c.SetCookie(JwtCookie, "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"status": "You are logged out"})
}

// Function to create JWT tokens with claims
func createToken(userId string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,                           // Subject (user identifier)
		"iss": "phdays-app",                     // Issuer
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secret)
	if err != nil {
		return "", err
	}

	logger.Infof("Token claims added: %v", claims)
	return tokenString, nil
}

// Function to verify JWT tokens
func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}

func isAuthorized(c *gin.Context) bool {
	tokenString, err := c.Cookie(JwtCookie)
	if err != nil {
		return false
	}

	// Verify the token
	token, _ := verifyToken(tokenString)
	return token != nil
}
