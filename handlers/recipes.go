package handlers


import (
	"fmt"
	"net/http"
	"time"
	"context"
	"recipes/models"

	"github.com/gin-gonic/gin"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
// "golang.org/x/net/context"


type RecipesHandler struct {
	collection *mongo.Collection
	ctx context.Context
}
func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
	collection: collection,
	ctx: ctx,
	}
}


/// End ponts


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
	func(handler *RecipesHandler) GetAllRecipesHandler(c *gin.Context) {
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return 
		}
		defer cur.Close(handler.ctx)
		recipes := make([]models.Recipe, 0)

		for cur.Next(handler.ctx){
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		c.JSON(http.StatusOK, recipes)
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
	func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
		var recipe models.Recipe

		err := c.ShouldBindJSON(&recipe);
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		recipe.ID = primitive.NewObjectID()
		recipe.PublishedAt = time.Now()
		_, err = handler.collection.InsertOne(handler.ctx, recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return 
		}

		c.JSON(http.StatusOK, recipe)
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
	func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
		id := c.Param("id")
		var recipe models.Recipe
		err := c.ShouldBindJSON(&recipe)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		objectId, _ := primitive.ObjectIDFromHex(id)
		cur, err := handler.collection.UpdateOne(handler.ctx, bson.M{
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

	func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {

		id := c.Param("id")
		objectId, err:= primitive.ObjectIDFromHex(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		res, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectId})

		if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error deleting recipe",
		})
		return

	}
		if res.DeletedCount == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Recipe not found",
			})
			return
		}
		c.JSON(http.StatusOK, res)
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
	func (handler *RecipesHandler)SearchRecipeHandler(c *gin.Context) {
		tag := c.Query("tag")
		listOfRecipes := make([]models.Recipe, 0)
		res, err := handler.collection.Find(handler.ctx, bson.D{{"tags", bson.D{{"$in", bson.A{tag}}}}})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return;
		}
		defer res.Close(handler.ctx)

		for res.Next(handler.ctx){
			var recipe models.Recipe
			res.Decode(&recipe)
			listOfRecipes = append(listOfRecipes, recipe)
		}
		c.JSON(http.StatusOK, listOfRecipes)
	}