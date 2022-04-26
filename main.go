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
	"go-microservices-example/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx            context.Context
	err            error
	client         *mongo.Client
	collection     *mongo.Collection
	recipesHandler *handlers.RecipesHandler
)

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.CreateRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetRecipeByIDHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.Run()
}
