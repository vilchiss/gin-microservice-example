// Recipes API
//
// This is a sample recipes API
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Luis Oropeza <luis.oropeza@gmail.com> https://www.gitlab.com/luisfi
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"go-microservices-example/handlers"
	"log"
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/square/go-jose.v2"
)

var (
	ctx            context.Context
	err            error
	client         *mongo.Client
	collection     *mongo.Collection
	collectionUser *mongo.Collection
	redisClient    *redis.Client
	recipesHandler *handlers.RecipesHandler
	authHandler    *handlers.AuthHandler
)

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collectionUser = client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, collectionUser)

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth0Domain = "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
		client := auth0.NewJWKClient(
			auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"},
			nil,
		)
		configuration := auth0.NewConfiguration(client, []string{os.Getenv("AUTH0_API_IDENTIFIER")}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			c.Abort()

			return
		}

		c.Next()
	}
}

func main() {
	store, err := redisStore.NewStore(10, "tcp", os.Getenv("REDIS_URI"), "", []byte(os.Getenv("REDIS_SECRET")))
	if err != nil {
		panic(fmt.Sprintf("failed to create redis instance: %s", err.Error()))
	}
	router := gin.Default()
	router.Use(sessions.Sessions("recipes_api", store))
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signup", authHandler.SignUpHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	authorized.POST("/recipes", recipesHandler.CreateRecipeHandler)
	authorized.GET("/recipes/:id", recipesHandler.GetRecipeByIDHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.Run()
}
