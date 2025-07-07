package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kbresser/chirpy/internal/auth"
	"github.com/kbresser/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expires := 3600

	token, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(expires)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get token", err)
		return
	}

	refreshTokenDB := database.RegRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}

	err = cfg.db.RegRefreshToken(r.Context(), refreshTokenDB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create refresh token", err)
		return
	}

	respondWithJSON(w, 200, returnVals{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
