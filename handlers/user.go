package handlers

import (
	"context"
	"errors"
	"go-microservices-example/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUsersHandler(ctx context.Context,
	collection *mongo.Collection) *UsersHandler {
	return &UsersHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation GET /users/{username} users getUserInformation
// Get a user information by username
// ---
// produces:
// - application/json
// parameters:
//   - name: username
//     in: path
//     description: user name
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: User not found
//     '200':
//         description: Internal error
func (handler *UsersHandler) GetUserInformationHandler(c *gin.Context) {
	username := c.Param("username")
	cursor := handler.collection.FindOne(handler.ctx, bson.M{"username": username})
	if cursor == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't retrieve value"})

		return
	}
	var user models.User
	err := cursor.Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":      user.Username,
		"registered_at": user.RegisteredAt,
	})
}
