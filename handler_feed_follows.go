package main

import (
	"encoding/json"
	"fmt"
	"github.com/Caaki/rssproject/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (api *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsnig JSON: %s", err))
		return
	}
	feedFollow, err := api.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Conldn't create feed follow: %s", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToUserFeedFollow(feedFollow))
}

func (api *apiConfig) handlerGetFeedFollowsByUserId(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := api.DB.GetFeedFollowsByUserId(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed folows for you: %v", err))
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiConfig *apiConfig) handlerGetFeedFollowsByUserApiKey(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiConfig.DB.GetFeedFollowsByUserApiKey(r.Context(), user.ApiKey)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed folows for you: %v", err))
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiConfig *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowID := chi.URLParam(r, "feedFollowID")
	feedFollowUUID, err := uuid.Parse(feedFollowID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsnig UUID: %s", err))
		return
	}
	err = apiConfig.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowUUID,
		ApiKey: user.ApiKey,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed folow: %s", err))
	}

	respondWithJSON(w, 200, fmt.Sprintf("Unsubscribed from feed"))

}
