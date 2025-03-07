package models

import "time"

// Users table
type Users struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Files table
type Files struct {
	ID         int        `json:"id"`
	Filename   string     `json:"filename"`
	FilePath   string     `json:"file_path"`
	UploadedBy int        `json:"uploaded_by"`
	IsPublic   bool       `json:"is_public"`
	FileSize   int        `json:"file_size"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// FilesMetadata table
type FilesMetadata struct {
	ID         int        `json:"id"`
	FileID     int        `json:"file_id"`
	UploadedBy int        `json:"uploaded_by"`
	Filename   string     `json:"filename"`
	UploadDate time.Time  `json:"upload_date"`
	FileType   string     `json:"file_type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// FilesPublicURL table
type FilesPublicURL struct {
	ID         int       `json:"id"`
	FileID     int       `json:"file_id"`
	UploadedBy int       `json:"uploaded_by"`
	TempURL    string    `json:"temporary_url"`
	CreatedAt  time.Time `json:"created_at"`
	ExpireAt   time.Time `json:"expire_at"`
}
