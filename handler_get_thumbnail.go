package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)
const maxMemory =10<<20
func (cfg *apiConfig) handlerThumbnailGet(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid video ID", err)
		return
	}
	err= r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Parsing Failed", err)
		return
	}
	multiPartFile, multipPartHeader,err:= r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn'tcreate a multipart form", err)
		return
	}
	mediaType:=multipPartHeader.Header["Content-Type"]

	// var data []byte
	data, err:= io.ReadAll(multiPartFile)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse data", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}
	cfg.db.GetVideo(userID)
	var thumbnail thumbnail
	thumbnail.data=data
	thumbnail.mediaType=mediaType[0]
	videoThumbnails[videoID]= thumbnail
	tn, ok := videoThumbnails[videoID]
	if !ok {
		respondWithError(w, http.StatusNotFound, "Thumbnail not found", nil)
		return
	}

	w.Header().Set("Content-Type", tn.mediaType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(tn.data)))

	_, err = w.Write(tn.data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing response", err)
		return
	}
}
