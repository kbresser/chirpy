package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something has gone wrong", err)
		return
	}

	mainChirps := []Chirp{}

	for i := range chirps {
		chirp := Chirp{
			ID:        chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body:      chirps[i].Body,
			UserID:    chirps[i].UserID,
		}
		mainChirps = append(mainChirps, chirp)
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
