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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const PORT = ":3000"
var recipes []Recipe;
type Recipe struct {
	ID				string		`json:"id"`
	Name 			string 		`json:"name"`
	Tags 			[]string 	`json:"tags"`
	Ingredients 	[]string 	`json:"ingredients"`
	Instructions 	[]string 	`json:"instructions"`
	PublishedAt 	time.Time 	`json:"publishedAt"`
};
	func init(){
	recipes = make([]Recipe,0 )
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes);
	}
// MAIN FUNCTION
func main() {
	router := gin.Default()
	recipes := router.Group("/api/v1/recipes")

	recipes.POST("", newRecipeHandler)
	recipes.GET("", getAllRecipesHandler)
	recipes.PUT("/:id", updateRecipeHandler)
	recipes.DELETE("/:id", deleteRecipeHandler)
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
	func newRecipeHandler(ctx *gin.Context) {
		var recipe Recipe

		err := ctx.ShouldBindJSON(&recipe);
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		recipe.ID = xid.New().String()
		recipe.PublishedAt = time.Now()
		recipes = append(recipes, recipe)

		ctx.JSON(http.StatusOK, recipe)
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
	func getAllRecipesHandler(ctx *gin.Context) {
		if len(recipes) < 1 {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Recipe not found",})
				return
			}
		ctx.JSON(http.StatusOK, recipes)
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
	func updateRecipeHandler(ctx *gin.Context) {
		id := ctx.Param("id")
		var recipe Recipe
		err := ctx.ShouldBindJSON(&recipe)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		for index, reci := range recipes {
			if reci.ID == id {
				recipes[index] = recipe;
				ctx.JSON(http.StatusOK, recipe)
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
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

	func deleteRecipeHandler(ctx *gin.Context) {
		id := ctx.Param("id")
		for index, recipe := range recipes {
			if recipe.ID == id {
				recipes = append(recipes[:index], recipes[index+1:]...)
				ctx.JSON(http.StatusOK, recipe)
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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
	func searchRecipeHandler(ctx *gin.Context) {
		tag := ctx.Query("tag")
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
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
		}
		ctx.JSON(http.StatusOK, listOfRecipes)
	}
	