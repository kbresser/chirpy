package main

import (
	"net/http"
	"time"

	"github.com/kbresser/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	result, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	JWTtoken, err := auth.MakeJWT(result.UserID, cfg.Secret, time.Duration(60)*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: JWTtoken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
