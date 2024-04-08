package api

import (
	"chatbot/pkg/chat"
	"chatbot/pkg/db"
	"chatbot/pkg/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Handler holds dependencies for the HTTP server.
type Handler struct {
	DB     *db.Database
	ChatAI *chat.OpenAI
}

func NewHandler(db *db.Database, chatAI *chat.OpenAI) *Handler {
	return &Handler{DB: db, ChatAI: chatAI}
}

// UploadImage handles image uploads.
func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // Max 10 MB file
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	identifier := r.FormValue("identifier")
	if identifier == "" {
		http.Error(w, "Identifier is required", http.StatusBadRequest)
		return
	}

	// Save the file
	filePath := filepath.Join("uploads", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	// Save metadata in the database
	err = h.DB.SaveImage(&model.ImageInfo{
		Identifier: identifier,
		URL:        filePath,
	})
	if err != nil {
		http.Error(w, "Error saving image metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully: %s", filePath)
}

// RetrieveImage handles image metadata retrieval.
func (h *Handler) RetrieveImage(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")
	if identifier == "" {
		http.Error(w, "Identifier is required", http.StatusBadRequest)
		return
	}

	image, err := h.DB.GetImage(identifier)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}
