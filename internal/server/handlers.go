// Package server implements http handlers.
package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/pkg/hash"
)

// Server represents application server.
type Server struct {
	repo *repository.Repository
}

// New creates new application server.
func New(repo *repository.Repository) *Server {
	return &Server{repo: repo}
}

// SignUp implements user's signing up.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	type signupRequest struct {
		Username string
		Password string
	}
	type signupResponse struct {
		ID string `json:"id"`
	}

	var sr signupRequest
	if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
		RespondErr(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if err := ValidateUserCredentials(sr.Username, sr.Password); err != nil {
		RespondErr(w, http.StatusBadRequest, err)
		return
	}

	var err error
	sr.Password, err = hash.HashPassword(sr.Password)
	if err != nil {
		RespondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := s.repo.InsertUser(sr.Username, sr.Password); err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}
	sr.Password = ""

	resp := signupResponse{
		ID: "1fa85f64-5717-4562-b3fc-2c963f66afa5",
	}
	Respond(w, http.StatusCreated, resp)
}

// LogIn implements user's logging in.
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Username string
		Password string
	}
	type loginResponse struct {
		Token string `json:"token"`
	}

	var lr loginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		RespondErr(w, http.StatusUnauthorized, err)
		return
	}
	defer r.Body.Close()

	hashedPwd, err := s.repo.GetUserPassword(lr.Username)
	if err != nil {
		RespondErr(w, http.StatusUnauthorized, err)
		return
	}

	if !hash.CheckPasswordHash(lr.Password, hashedPwd) {
		RespondErr(w, http.StatusUnauthorized, errors.New("wrong password"))
		return
	}

	resp := loginResponse{
		Token: "eyJhbGciOiJIUzI1NiIs...",
	}
	Respond(w, http.StatusCreated, resp)

}

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/signup", s.SignUp).Methods("POST")
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
}
