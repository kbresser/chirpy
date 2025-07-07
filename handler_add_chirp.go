package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kbresser/chirpy/internal/auth"
	"github.com/kbresser/chirpy/internal/database"
)

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	paramsDB := database.AddChirpParams{
		Body:   cleaned,
		UserID: userID,
	}

	chirp, err := cfg.db.AddChirp(r.Context(), paramsDB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt insert chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
