package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
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

	userID, err := cfg.getAuthenticatedUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something has gone wrong", err)
		return
	}

	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "asc"
	}

	switch sortParam {
	case "asc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	mainChirps := []Chirp{}
	s := r.URL.Query().Get("author_id")

	if s == "" {
		for i := range chirps {
			returnChirp := Chirp{
				ID:        chirps[i].ID,
				CreatedAt: chirps[i].CreatedAt,
				UpdatedAt: chirps[i].UpdatedAt,
				Body:      chirps[i].Body,
				UserID:    chirps[i].UserID,
			}
			mainChirps = append(mainChirps, returnChirp)
		}
	} else {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something has gone wrong", err)
			return
		}
		for i, chirp := range chirps {
			if chirp.UserID == authorID {
				returnChirp := Chirp{
					ID:        chirps[i].ID,
					CreatedAt: chirps[i].CreatedAt,
					UpdatedAt: chirps[i].UpdatedAt,
					Body:      chirps[i].Body,
					UserID:    chirps[i].UserID,
				}

				mainChirps = append(mainChirps, returnChirp)
			}
		}
	}

	respondWithJSON(w, http.StatusOK, mainChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Chirp
	}

	chirpID := r.PathValue("chirpID")
	log.Println("Extracted chirpID:", chirpID)
	parsedUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	log.Println("Extracted chirpID:", chirpID)
	parsedUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	userID, err := cfg.getAuthenticatedUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	if chirp.UserID == userID {
		err = cfg.db.DeleteChirp(r.Context(), chirp.ID)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "Could not delete chirp", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		respondWithError(w, http.StatusForbidden, "This is not your chirp", err)
	}

}
