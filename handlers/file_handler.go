package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/naresh2002/8901135580_Backend/db"
	"github.com/naresh2002/8901135580_Backend/utils"
)

type FileHandler struct {
	lg *log.Logger
	db *db.Database
	mu sync.RWMutex
}

func NewFileHandler(lg *log.Logger, db *db.Database) *FileHandler {
	return &FileHandler{lg: lg, db: db}
}

// UploadFile handles file uploads
func (f *FileHandler) UploadFile(rw http.ResponseWriter, req *http.Request) {
	f.lg.Println("Upload file endpoint hit")

	// Extract user ID from JWT token
	userID, err := utils.ExtractUserID(req)
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form
	err = req.ParseMultipartForm(10 << 20) // 10 MB max upload size
	if err != nil {
		http.Error(rw, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(rw, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// File details
	filename := header.Filename
	fileExt := filepath.Ext(filename)[1:]         // Extract file extension without dot
	uploadDate := time.Now().Format("2006-01-02") // Format: YYYY-MM-DD

	// Create user-specific folder
	uploadDir := fmt.Sprintf("media/%d/", userID)
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		http.Error(rw, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// Save file to storage
	filePath := filepath.Join(uploadDir, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(rw, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(rw, "Failed to write file", http.StatusInternalServerError)
		return
	}

	// Get file size in KB
	fileInfo, err := outFile.Stat()
	if err != nil {
		http.Error(rw, "Failed to get file info", http.StatusInternalServerError)
		return
	}
	fileSize := fileInfo.Size() / 1024 // Convert bytes to KB

	// Insert into files table
	var fileID int
	err = f.db.Conn.QueryRow(context.Background(),
		"INSERT INTO files (filename, file_path, uploaded_by, file_size) VALUES ($1, $2, $3, $4) RETURNING id",
		filename, filePath, userID, fileSize,
	).Scan(&fileID)
	if err != nil {
		http.Error(rw, "Failed to insert file record", http.StatusInternalServerError)
		return
	}

	// Insert into files_metadata table
	_, err = f.db.Conn.Exec(context.Background(),
		"INSERT INTO files_metadata (file_id, uploaded_by, filename, upload_date, file_type) VALUES ($1, $2, $3, $4, $5)",
		fileID, userID, filename, uploadDate, fileExt,
	)
	if err != nil {
		http.Error(rw, "Failed to insert metadata", http.StatusInternalServerError)
		return
	}

	// Generate temporary URL (random string)
	tempURL := generateTempURL()
	expireAt := time.Now().Add(72 * time.Hour) // 3 days expiry

	// Insert into files_public_url table
	_, err = f.db.Conn.Exec(context.Background(),
		"INSERT INTO files_public_url (file_id, uploaded_by, temporary_url, expire_at) VALUES ($1, $2, $3, $4)",
		fileID, userID, tempURL, expireAt,
	)
	if err != nil {
		http.Error(rw, "Failed to insert temporary URL", http.StatusInternalServerError)
		return
	}

	// Response
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"message":       "File uploaded successfully",
		"file_id":       fmt.Sprintf("%d", fileID),
		"file_path":     filePath,
		"temporary_url": fmt.Sprintf("http://localhost:8000/file?url=%s", tempURL),
		"expires_at":    expireAt.Format(time.RFC3339),
	})
}

// ServeFile serves a file using its temporary URL
func (f *FileHandler) ServeFile(rw http.ResponseWriter, req *http.Request) {
	tempURL := req.URL.Query().Get("url")
	if tempURL == "" {
		http.Error(rw, "Missing temporary URL", http.StatusBadRequest)
		return
	}

	var filePath string
	var expireAt time.Time

	// Fetch file details from database
	err := f.db.Conn.QueryRow(context.Background(),
		"SELECT files.file_path, files_public_url.expire_at FROM files_public_url "+
			"JOIN files ON files.id = files_public_url.file_id WHERE files_public_url.temporary_url = $1",
		tempURL,
	).Scan(&filePath, &expireAt)

	if err != nil {
		http.Error(rw, "File not found", http.StatusNotFound)
		return
	}

	// Check expiration time
	if time.Now().After(expireAt) {
		http.Error(rw, "Temporary URL expired", http.StatusGone)
		return
	}

	http.ServeFile(rw, req, filePath)
}

// generateTempURL creates a short random string for temporary access
func generateTempURL() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
