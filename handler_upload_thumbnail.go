package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

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
	mediaType, _, err := mime.ParseMediaType(multipPartHeader.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", err)
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}
	// var data []byte
	// data, err:= io.ReadAll(multiPartFile)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Couldn't parse data", err)
	// 	return
	// }
	key := make([]byte, 32)
	rand.Read(key)
	encodedKey := base64.RawURLEncoding.EncodeToString(key)
	// assetPath:= cfg.getAssetsPath(videoID,mediaType)
	assetPath:= cfg.getAssetsPath(encodedKey,mediaType)
	assetDiskPath:=cfg.getAssetsDiskPath(assetPath)
	fmt.Println("Saving thumbnail to:", assetDiskPath)
	dst, err:=os.Create(assetDiskPath)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError,"Unable to create file on server",err)
		return
	}
	defer dst.Close()
	if _,err= io.Copy(dst,multiPartFile);err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving file", err)
		return
	}
	
	video,err:=cfg.db.GetVideo(videoID)
	if video.UserID!=userID{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized User ", err)
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't fetch video", err)
		return
	}
	// encodedimageData:=base64.StdEncoding.EncodeToString(data)
	
	
	// var thumbnail thumbnail
	// thumbnail.data=data
	// thumbnail.mediaType=mediaType
	// videoThumbnails[videoID]= thumbnail
	// thumbnail_url:= fmt.Sprintf("http://localhost:%s/api/thumbnails/%s", cfg.port, videoID)
	// thumbnail_url:= fmt.Sprintf("data:%s;base64,%v", mediaType,encodedimageData)
	thumbnail_url:= cfg.getAssetsURL(assetPath)
	video.ThumbnailURL=&thumbnail_url
	err=cfg.db.UpdateVideo(video)
	if err != nil {
		// delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}
	// respondWithJSON(w,r.Response.StatusCode,video)
	respondWithJSON(w, http.StatusOK, video)
}
