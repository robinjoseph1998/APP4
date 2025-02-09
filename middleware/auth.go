package middleware

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// )

// func GenToken(userID string, email string, c *gin.Context) (string, error) {
// 	SecretKey := os.Getenv("JWT_SECRET_KEY")
// 	claims := YourClaims{
// 		UserID: userID,
// 		Email:  email,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
// 			Issuer:    "user",
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(SecretKey))
// 	if err != nil {
// 		return "", err
// 	}

// 	// Set both cookie and header
// 	maxAge := 60000
// 	c.SetCookie(
// 		"Authorization",       // name
// 		"Bearer "+tokenString, // value (include Bearer prefix)
// 		maxAge,                // max age
// 		"/",                   // path
// 		"",                    // domain
// 		false,                 // secure (set to true in production with HTTPS)
// 		true,                  // httpOnly
// 	)

// 	// Also set header for API requests
// 	c.Header("Authorization", "Bearer "+tokenString)

// 	// Store token in session if needed
// 	if session := c.MustGet("session"); session != nil {
// 		session.(*sessions.Session).Values["token"] = tokenString
// 		session.(*sessions.Session).Save(c.Request, c.Writer)
// 	}

// 	return tokenString, nil
// }

// // Middleware to check token
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Try to get token from different sources
// 		var tokenString string

// 		// Check Authorization header
// 		if auth := c.GetHeader("Authorization"); auth != "" {
// 			if len(auth) > 7 && auth[:7] == "Bearer " {
// 				tokenString = auth[7:]
// 			}
// 		}

// 		// If not in header, check cookie
// 		if tokenString == "" {
// 			cookie, err := c.Cookie("Authorization")
// 			if err == nil && len(cookie) > 7 && cookie[:7] == "Bearer " {
// 				tokenString = cookie[7:]
// 			}
// 		}

// 		// Validate token
// 		if tokenString == "" {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
// 			return
// 		}

// 		token, err := jwt.ParseWithClaims(tokenString, &YourClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
// 			return
// 		}

// 		// Set claims to context
// 		if claims, ok := token.Claims.(*YourClaims); ok {
// 			c.Set("userID", claims.UserID)
// 			c.Set("email", claims.Email)
// 		}

// 		c.Next()
// 	}
// }

// func DeleteCookie(c *gin.Context) error {
// 	c.SetCookie("Authorization", "", 0, "", "", false, true)
// 	fmt.Println("Cookie Deleted")
// 	return nil
// }
