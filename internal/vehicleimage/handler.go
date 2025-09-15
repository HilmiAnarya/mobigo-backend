package vehicleimage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	// Note: We are attaching the image upload route to the /vehicles path for logical grouping.
	r := router.PathPrefix("/api/vehicles/{vehicleID}/images").Subrouter()
	r.Use(authMiddleware)
	r.HandleFunc("", h.uploadImageHandler).Methods("POST")

	// Routes for a specific image (by its own ID)
	imageRouter := router.PathPrefix("/api/images/{id}").Subrouter()
	imageRouter.Use(authMiddleware)
	imageRouter.HandleFunc("", h.deleteImageHandler).Methods("DELETE")
	imageRouter.HandleFunc("/primary", h.setPrimaryImageHandler).Methods("PUT")
}

func (h *Handler) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get vehicleID from URL
	vars := mux.Vars(r)
	vehicleID, err := strconv.ParseInt(vars["vehicleID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	// 2. Parse the multipart form data (the image)
	// 10 << 20 specifies a maximum upload size of 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// 3. Retrieve the file from the form data
	file, handler, err := r.FormFile("image") // "image" is the name of the form field
	if err != nil {
		http.Error(w, "Invalid image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 4. Create the uploads directory if it doesn't exist
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// 5. Create a new file on the server
	// We generate a unique filename to avoid collisions
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), handler.Filename)
	filePath := filepath.Join(uploadDir, uniqueFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 6. Copy the uploaded file's content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	// 7. Save the image path to the database
	// We will store the path that the frontend can use to access the image.
	imageURL := "/uploads/" + uniqueFileName
	_, err = h.service.CreateVehicleImage(r.Context(), vehicleID, imageURL, false) // Default isPrimary to false
	if err != nil {
		// Attempt to remove the saved file if DB entry fails
		os.Remove(filePath)
		http.Error(w, "Could not save image metadata to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteVehicleImage(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) setPrimaryImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	err = h.service.SetPrimaryImage(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to set primary image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
