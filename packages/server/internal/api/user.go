package api

import (
	"encoding/json"
	"log"
	"net/http"
	"server/internal/db"
	"server/schema"

	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var user UserRegistrationRequest

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("[v1-auth] decode error: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// --- Simple Validation ---
		if user.Email == "" || user.Password == "" {
			http.Error(w, "Email and password are required", http.StatusUnprocessableEntity)
			return
		}

		if existingUsr, _ := db.GetUserByMail(ctx, user.Email); existingUsr.Email != "" {
			http.Error(w, "Email already exist", http.StatusForbidden)
			return
		}

		usrPassword, err := HashPassword(user.Password)
		if err != nil {
			log.Printf("[v1-auth] hash error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = db.CreateUser(ctx, schema.CreateUserParams{
			Name:         user.Name,
			Email:        user.Email,
			PasswordHash: usrPassword,
		})

		if err != nil {
			// Check if this is a "Duplicate Entry" error (MySQL Error 1062)
			// If using standard database/sql:
			log.Printf("[v1-auth] database error: %v", err)
			http.Error(w, "User registration failed", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated) // 201 is better for creation than 200
	}
}

func HashPassword(password string) (string, error) {
	// GenerateFromPassword returns a []byte, usually stored as a string
	// '10' is the cost factor (standard for most apps)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
