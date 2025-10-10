package main

import (
	"fmt"
	"io"
	"net/http"

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
	contentType:=multipPartHeader.Header["Content-Type"]

	var data []byte
	data,err = io.ReadAll(multiPartFile)
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
