package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/44k45h/rss-scraper/internal/auth"
	"github.com/44k45h/rss-scraper/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}

	responseWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	ApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("Authentication error: %s", err))
		return
	}

	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), ApiKey)
	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("Couldn't find user: %s", err))
		return
	}

	responseWithJSON(w, 200, databaseUserToUser(user))
}
