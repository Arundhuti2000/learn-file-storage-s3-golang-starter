package main

import (
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