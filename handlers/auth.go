package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"go-microservices-example/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AuthorizationHeader = "Authorization"
	EnvJWTSecret        = "JWT_SECRET"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// swagger:model jwtOutput
type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation POST /signin auth signin
// Authenticate user
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//         schema:
//             "$ref": "#/definitions/jwtOutput"
//     '401':
//         description: Invalid username or password
//     '404':
//         description: Invalid data
//     '500':
//         description: Internal server error
func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	hash := sha256.Sum256([]byte(user.Password))

	cursor := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": fmt.Sprintf("%x", hash),
	})
	if cursor.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})

		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv(EnvJWTSecret)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

// swagger:operation POST /refresh auth refresh
// Refresh token
// ---
// parameters:
// - name: Authorization
//   in: header
//   description: token
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '401':
//         description: Invalid token
//     '404':
//         description: Token is not expired yet
//     '500':
//         description: Internal server error
func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader(AuthorizationHeader)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenValue, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(EnvJWTSecret)), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})

		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is not expired yet",
		})

		return
	}

	expiratonTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expiratonTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(os.Getenv(EnvJWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	JWTOutput := JWTOutput{
		Token:   newTokenString,
		Expires: expiratonTime,
	}
	c.JSON(http.StatusOK, JWTOutput)
}
