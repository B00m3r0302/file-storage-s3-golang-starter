package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	// Const in Bytes
	const maxMemory = 10 << 20

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	metadata, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get video", err)
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

	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}

	// Get image data from multipart form
	fileData, fileHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get image data from multipart form", err)
		return
	}

	defer fileData.Close()
	mediaType := fileHeader.Header.Get("Content-Type")

	if metadata.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "You can't upload a thumbnail for this video", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	fileType, params, err := mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse media type", err)
		return
	}

	if fileType != "image/jpeg" && fileType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}

	fmt.Println("file type", fileType)
	fmt.Println("params", params)
	parts := strings.Split(mediaType, "/")
	// parts[0] = "image", parts[1] = "png"
	ext := "." + parts[1]
	key := make([]byte, 32)
	rand.Read(key)
	id := base64.RawURLEncoding.EncodeToString(key)
	filePath := filepath.Join(cfg.assetsRoot, id+ext)
	fileName := fmt.Sprintf("%s%s", id, ext)

	file, err := os.Create(filePath)
	defer file.Close()

	_, err = io.Copy(file, fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save thumbnail", err)
		return
	}

	thumbnailUrl := fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, fileName)
	metadata.ThumbnailURL = &thumbnailUrl
	err = cfg.db.UpdateVideo(metadata)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, metadata)
}
