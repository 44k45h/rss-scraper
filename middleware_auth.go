package main

import (
	"fmt"
	"net/http"

	"github.com/44k45h/rss-scraper/internal/auth"
	"github.com/44k45h/rss-scraper/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlerwareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	}
}
