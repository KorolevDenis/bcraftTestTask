package controllers

import (
	"bcraftTestTask/internal/models"
	"bcraftTestTask/internal/repository"
	u "bcraftTestTask/internal/utils"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	ar *repository.AuthRepository
}

func NewAuthController(authRepository_ *repository.AuthRepository) *AuthController {
	return &AuthController{
		ar: authRepository_,
	}
}

func (c *AuthController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := models.Account{}

	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request",
			"Jira Analyzer REST API Create Account"), http.StatusBadRequest)
		return
	}

	err = c.ar.CreateAccount(account)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"GetRecipeById"), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "success",
		"CreateAccount")

	u.Respond(w, resp, http.StatusOK)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	account := models.Account{}
	resAccount := &models.Account{}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request",
			"Jira Analyzer REST API Create Account"), http.StatusBadRequest)
		return
	}

	resAccount, err = c.ar.Login(account.Email, account.Password)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error(),
			"GetRecipeById"), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "success",
		"Login")
	resp["data"] = resAccount

	u.Respond(w, resp, http.StatusOK)
}
