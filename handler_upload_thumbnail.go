package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
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
	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)
	const maxMemory = 10>>20
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
	defer multiPartFile.Close()
	mediaType:=multipPartHeader.Header.Get("Content-Type")
	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", nil)
		return
	}
	// var data []byte
	data, err:= io.ReadAll(multiPartFile)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse data", err)
		return
	}
	
	metaData,err:=cfg.db.GetVideo(videoID)
	if metaData.UserID!=userID{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized User ", err)
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't fetch video", err)
		return
	}
	var thumbnail thumbnail
	thumbnail.data=data
	thumbnail.mediaType=mediaType
	videoThumbnails[videoID]= thumbnail
	thumbnail_url:= fmt.Sprintf("http://localhost:%s/api/thumbnails/%s", cfg.port, videoID)
	metaData.ThumbnailURL=&thumbnail_url
	err=cfg.db.UpdateVideo(metaData)
	if err != nil {
		delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}
	// respondWithJSON(w,r.Response.StatusCode,metaData)
	respondWithJSON(w, http.StatusOK, metaData)
}
