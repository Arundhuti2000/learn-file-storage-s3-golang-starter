package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("uploading video", videoID, "by user", userID)
	const maxByte = 1<<30
	r.Body = http.MaxBytesReader(w, r.Body, maxByte)
	//body, err := io.ReadAll(r.Body)
	video,err:= cfg.db.GetVideo(videoID)
	if video.UserID!=userID{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized User ", err)
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't fetch video", err)
		return
	}
	err= r.ParseMultipartForm(maxByte)
	if err!=nil{
		respondWithError(w, http.StatusBadRequest, "Couldn't Prse Multippart file", err)
		return 
	}
	multiPartFile, multipPartHeader,err:= r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn'tcreate a multipart form", err)
		return
	}
	defer multiPartFile.Close()
	mediaType,_, err:= mime.ParseMediaType(multipPartHeader.Header.Get("Content-Type"))
	// mediaType,_, err:= mime.ParseMediaType("video/mp4")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse media type", err)
		return
	}
	if mediaType != "video/mp4" {
		fmt.Sprintf("invalid MIME type: expected 'video/mp4', got %s", mediaType)
		return 
	}
	// assetPath:= cfg.getAssetsPath(videoID.String(),mediaType)
	// assetDiskPath:=cfg.getAssetsDiskPath(assetPath)
	// fmt.Println("Saving thumbnail to:", assetDiskPath)
	key := make([]byte, 32)
	rand.Read(key)
	encodedKey := base64.RawURLEncoding.EncodeToString(key)
	dst, err:=os.CreateTemp("", "tubely-upload.mp4")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't Create Temp file", err)
		return
	}
	defer os.Remove("tubely-upload.mp4")
	defer dst.Close()
	io.Copy(dst, multiPartFile)
	dst.Seek(0, io.SeekStart)
	input:= &s3.PutObjectInput{
		Bucket: aws.String(cfg.s3Bucket),
		Key: aws.String(encodedKey),
		Body: dst,
		ContentType: aws.String(mediaType),
	}
	_,err= cfg.s3Client.PutObject(context.Background(),input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upload to S3", err)
		return
	}
	videoURL:= fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",cfg.s3Bucket,cfg.s3Region,encodedKey)
	video.VideoURL=&videoURL
	err=cfg.db.UpdateVideo(video)
	if err != nil {
		// delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}
	respondWithJSON(w, http.StatusOK, video)
}
