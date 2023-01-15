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
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database
var ctx context.Context
var err error
var collection *mongo.Collection
var client *mongo.Client

const PORT = ":3000"
var recipes []Recipe;
type Recipe struct {
	ID				primitive.ObjectID		`json:"id" bson:"_id"`
	Name 			string 					`json:"name" bson:"name"`
	Tags 			[]string 				`json:"tags" bson:"tags"`
	Ingredients 	[]string 				`json:"ingredients" bson:"ingredients"`
	Instructions 	[]string 				`json:"instructions" bson:"instructions"`
	PublishedAt 	time.Time 				`json:"publishedAt" bson:"publishedAt"`
};

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
}


// MAIN FUNCTION
func main() {
	router := gin.Default()
	recipes := router.Group("/api/v1/recipes")

	recipes.POST("", newRecipeHandler)
	recipes.GET("", getAllRecipesHandler)
	recipes.PUT("/:id", updateRecipeHandler)
	// recipes.DELETE("/:id", deleteRecipeHandler)
	recipes.GET("/search", searchRecipeHandler)


	router.Run(PORT)
	fmt.Printf("Server is running on port %s", PORT)
}

// swagger:route POST /recipes recipes newRecipe
//
// Create a new recipe
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
// Responses:
//   '200': Recipe
//   '400':
//   description: Bad request
	func newRecipeHandler(c *gin.Context) {
		var recipe Recipe

		err := c.ShouldBindJSON(&recipe);
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		recipe.ID = primitive.NewObjectID()
		recipe.PublishedAt = time.Now()
		_, err = collection.InsertOne(ctx, recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return 
		}

		c.JSON(http.StatusOK, recipe)
	}

// swagger:route GET /recipes recipes getAllRecipes
//
// Get all recipes
//
// This will show all available recipes.
//
// Produces:
// - application/json
//
// Responses:
//   '200': 
//	recipes: []Recipe
	func getAllRecipesHandler(c *gin.Context) {
		// collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return 
		}
		defer cur.Close(ctx)
		recipes := make([]Recipe, 0)

		for cur.Next(ctx){
			var recipe Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		c.JSON(http.StatusOK, recipes)
	}

// swagger:route PUT /recipes/{id} recipes updateRecipe
//
// Update a recipe
// ---
// Parameters:
// + name: id
//   in: path
//   description: ID of the recipe to update
// 	 required: true
// 	 type: string
//
// Produces:
// - application/json
// Responses:
//   '200': Recipe
//   '400': 
//   description: Bad request
//   '404':
//   description: Recipe not found
	func updateRecipeHandler(c *gin.Context) {
		id := c.Param("id")
		var recipe Recipe
		err := c.ShouldBindJSON(&recipe)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		objectId, _ := primitive.ObjectIDFromHex(id)
		cur, err := collection.UpdateOne(ctx, bson.M{
			"_id": objectId,
			}, bson.D{{"$set", bson.D{
					{"name", recipe.Name},
					{"instructions", recipe.Instructions},
					{"ingredients", recipe.Ingredients},
					{"tags", recipe.Tags},
				}}})

			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return;
			}
		fmt.Println("", cur)
		c.JSON(http.StatusOK, recipe)
	}
// delete a recipe by id
// swagger:route DELETE /recipes/{id} recipes deleteRecipe
//
// Delete a recipe
//
// Parameters:
// + name: id
//   in: path
//   required: true
//   type: string
//   description: ID of the recipe
//
// Produces:
// - application/json
// Responses:
//   '200': Recipe
//   '404':
//   description: Recipe not found

	func deleteRecipeHandler(c *gin.Context) {
		id := c.Param("id")
		objectId, _:= primitive.ObjectIDFromHex(id)
		
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
	}

// swagger:route GET /recipes/search recipes searchRecipe
//
// Get a recipe by tag
//
// Parameters:
// + name: tag
//   in: query
//   required: true
//   type: string
//   description: Tag of the recipe
//
// Produces:
// - application/json
// Responses:
//   '200': []Recipe
//   '404':
//   description: Recipe not found
	func searchRecipeHandler(c *gin.Context) {
		tag := c.Query("tag")
		listOfRecipes := make([]Recipe, 0)

		for _, recipe := range recipes {
			found := false
			for _, t := range recipe.Tags {
				if strings.EqualFold(t, tag) {
					found = true
					break
				}
			}
			if found {
				listOfRecipes = append(listOfRecipes, recipe)
			}
		}
		if len(listOfRecipes) < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
		}
		c.JSON(http.StatusOK, listOfRecipes)
	}
	