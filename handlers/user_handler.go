package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/naresh2002/8901135580_Backend/db"
	"github.com/naresh2002/8901135580_Backend/models"
	"github.com/naresh2002/8901135580_Backend/utils"
)

type UserHandler struct {
	lg *log.Logger
	db *db.Database
	mu sync.RWMutex
}

func NewUserHandler(lg *log.Logger, db *db.Database) *UserHandler {
	return &UserHandler{lg: lg, db: db}
}

// Signup handles user registration
func (u *UserHandler) Signup(rw http.ResponseWriter, req *http.Request) {
	u.lg.Println("Signup endpoint hit")

	var user models.Users
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(rw, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Store user in the database
	_, err = u.db.Conn.Exec(context.Background(),
		"INSERT INTO users (email, password) VALUES ($1, $2)",
		user.Email, hashedPassword,
	)
	if err != nil {
		http.Error(rw, "Failed to create user", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user authentication
func (u *UserHandler) Login(rw http.ResponseWriter, req *http.Request) {
	u.lg.Println("Login endpoint hit")

	var user models.Users
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetching user from DB
	var dbUser models.Users
	err = u.db.Conn.QueryRow(context.Background(),
		"SELECT id, email, password FROM users WHERE email=$1", user.Email).
		Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password)

	if err == pgx.ErrNoRows {
		http.Error(rw, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(rw, "Database error", http.StatusInternalServerError)
		return
	}

	// Verifing password
	if !utils.ComparePasswords(dbUser.Password, user.Password) {
		http.Error(rw, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generating JWT token
	token, err := utils.GenerateToken(dbUser.ID)
	if err != nil {
		http.Error(rw, "Error generating token", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"token": token})
}
