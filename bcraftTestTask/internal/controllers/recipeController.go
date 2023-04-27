package controllers

import (
	"bcraftTestTask/internal/models"
	"bcraftTestTask/internal/repository"
	u "bcraftTestTask/internal/utils"
	"encoding/json"
	"net/http"
)

type RecipeController struct {
	rp *repository.RecipeRepository
}

func NewRecipeController(recipeRepository_ *repository.RecipeRepository) *RecipeController {
	return &RecipeController{
		rp: recipeRepository_,
	}
}

// GetRecipe function of getting recipe by id
func (c *RecipeController) GetRecipes(w http.ResponseWriter, r *http.Request) {
	ingredients := r.FormValue("ingredients")
	sortBy := r.FormValue("sortBy")
	time := r.FormValue("time")

	recipes, err := c.rp.GetRecipes(ingredients, sortBy, time)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"GetRecipes"), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "success",
		"GetRecipes")
	resp["data"] = recipes

	u.Respond(w, resp, http.StatusOK)
}

// PutRecipe add/change recipe function
func (c *RecipeController) PutRecipe(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var recipe models.Recipe
	err := decoder.Decode(&recipe)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"PutRecipe"), http.StatusBadRequest)
		return
	}

	status, err := c.rp.PutRecipe(r.Context(), recipe)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"PutRecipe"), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "success",
		"PutRecipe")

	u.Respond(w, resp, status)
}

func (c *RecipeController) GetRecipeById(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIdFromUrl(r.URL.Path, "recipe")
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"GetRecipeById"), http.StatusBadRequest)
		return
	}

	recipes, err := c.rp.GetRecipeById(id)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"GetRecipeById"), http.StatusBadRequest)
		return
	}

	if recipes.Id == 0 {
		u.Respond(w, u.Message(false, "No such recipe",
			"GetRecipeById"), http.StatusNotFound)
		return
	}

	resp := u.Message(true, "success",
		"GetRecipeById")
	resp["data"] = recipes

	u.Respond(w, resp, http.StatusOK)
}

func (c *RecipeController) DeleteRecipeById(w http.ResponseWriter, r *http.Request) {
	id, err := u.GetIdFromUrl(r.URL.Path, "recipe")
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"DeleteRecipeById"), http.StatusBadRequest)
		return
	}

	recipes, err := c.rp.DeleteRecipeById(id)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"DeleteRecipeById"), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "success",
		"DeleteRecipeById")
	resp["data"] = recipes

	u.Respond(w, resp, http.StatusOK)
}
