package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/env"
	"server/internal/jwt"
	"server/schema"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserConnectionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUserSession struct {
	Session string `json:"session"`
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

		user.Name = strings.TrimSpace(user.Name)
		user.Email = strings.TrimSpace(user.Email)
		user.Password = strings.TrimSpace(user.Password)

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

func UserConnectionHandler(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var anon UserConnectionRequest

		if err := json.NewDecoder(r.Body).Decode(&anon); err != nil {
			log.Printf("[v1-auth] decode error: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		anon.Email = strings.TrimSpace(anon.Email)
		anon.Password = strings.TrimSpace(anon.Password)

		// --- Simple Validation ---
		if anon.Email == "" || anon.Password == "" {
			http.Error(w, "Email and password are required", http.StatusUnprocessableEntity)
			return
		}

		existingUsr, err := db.GetUserByMail(ctx, anon.Email)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Email or Password Incorrect", http.StatusUnauthorized)
				return
			}

			log.Printf("[v1-auth] db error: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		comp_err := bcrypt.CompareHashAndPassword([]byte(existingUsr.PasswordHash), []byte(anon.Password))

		if comp_err != nil {
			http.Error(w, "Email or Password Incorrect", http.StatusForbidden)
			return
		}

		access_tk, gen_err := jwt.GenAccessToken(env.Load().JWT_ACCESS_SECRET, existingUsr.ID)

		if gen_err != nil {
			log.Printf("[v1-auth] access token gen error: %v", gen_err)
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(NewUserSession{
			Session: access_tk,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("[v1-auth] failed json encoding on access token : %v", gen_err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func HashPassword(password string) (string, error) {
	// GenerateFromPassword returns a []byte, usually stored as a string
	// '10' is the cost factor (standard for most apps)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
