// Recipe API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
// Schemes: http
//
// Host: localhost:3000
//
// BasePath: /api/v1/recipes
//
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"recipes/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database
var ctx context.Context
var err error
var collection *mongo.Collection
var client *mongo.Client
// Handlers
var recipesHandler *handlers.RecipesHandler


const PORT = ":3000"


func init(){
	err = godotenv.Load()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
		return 
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}


// MAIN FUNCTION
func main() {
	router := gin.Default()
	recipes := router.Group("/api/v1/recipes")

	recipes.POST("", recipesHandler.NewRecipeHandler)
	recipes.GET("", recipesHandler.GetAllRecipesHandler)
	recipes.PUT("/:id", recipesHandler.UpdateRecipeHandler)
	recipes.DELETE("/:id", recipesHandler.DeleteRecipeHandler)
	recipes.GET("/search", recipesHandler.SearchRecipeHandler)


	router.Run(PORT)
	fmt.Printf("Server is running on port %s", PORT)
}
	