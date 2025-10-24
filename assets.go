package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}
func (cfg apiConfig) getAssetsPath(key, mediaType string) string{
// func (cfg apiConfig) getAssetsPath(videoId uuid.UUID, mediaType string) string{
	var ext string
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		ext = ".bin"
	} else {
		ext = fmt.Sprintf("%s.%s", key, parts[1])
	}
	return ext
}
func (cfg apiConfig) getAssetsDiskPath(assetPath string)string{
	return filepath.Join(cfg.assetsRoot, assetPath)
}
func (cfg apiConfig) getAssetsURL(assetPath string) string{
	return  fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}
func (cfg apiConfig) getObjectURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
}
// func (cfg apiConfig) getPresignedObjectURL(key string) string {
// 	return fmt.Sprintf("%s,%s", cfg.s3Bucket, key)
// }
// func (cfg apiConfig) generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
// 	presignClient := s3.NewPresignClient(s3Client)
// 	resp, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(key),
// 	}, s3.WithPresignExpires(expireTime))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
// 	}
// 	return resp.URL, nil
// }
func getAssetPath(mediaType string) string {
	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes")
	}
	id := base64.RawURLEncoding.EncodeToString(base)

	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", id, ext)
}
func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}