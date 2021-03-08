// Package server implements http handlers.
package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/signup", s.SignUp).Methods("POST")
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
	r.HandleFunc("/conversion", s.ConversionRequest).Methods("POST")
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

	userID, err := s.repo.InsertUser(sr.Username, sr.Password)
	if err != nil {
		RespondErr(w, http.StatusInternalServerError, err)
		return
	}
	sr.Password = ""

	resp := signupResponse{
		ID: userID,
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

// ConversionRequest creates a request for audio conversion.
func (s *Server) ConversionRequest(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		RespondErr(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	name := header.Filename
	sourceFormat := strings.ToUpper(r.FormValue("sourceFormat"))
	targetFormat := strings.ToUpper(r.FormValue("targetFormat"))

	if err := ValidateRequest(name, sourceFormat, targetFormat); err != nil {
		RespondErr(w, http.StatusBadRequest, err)
		return
	}

	requestID, err := s.repo.MakeRequest(name, sourceFormat, targetFormat, "some location", "992dee5c-b4e3-49f8-9d4c-8903fa2284c9")
	if err != nil {
		RespondErr(w, http.StatusBadRequest, err)
		return
	}

	type conversionResponse struct {
		ID string `json:"id"`
	}
	convertResp := conversionResponse{
		ID: requestID,
	}

	Respond(w, http.StatusCreated, convertResp)
}
