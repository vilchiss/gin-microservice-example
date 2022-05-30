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
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// Login with username and password
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

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "User signed in",
	})
}

// swagger:operation POST /signup auth signup
// Register a new user
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid data
//     '500':
//         description: Internal server error
func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	user.ID = primitive.NewObjectID()
	user.RegisteredAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not crate user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered",
	})
}

// swagger:operation POST /refresh auth refresh
// Get new token in exchange for and old one
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

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out... Bye!",
	})
}
