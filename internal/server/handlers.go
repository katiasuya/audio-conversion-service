// Package server implements http handlers.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/server/context"
	res "github.com/katiasuya/audio-conversion-service/internal/server/response"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"github.com/katiasuya/audio-conversion-service/pkg/hash"
)

var errInvalidUsernameOrPassword = errors.New("invalid username or password")

// Server represents application server.
type Server struct {
	repo      *repository.Repository
	storage   *storage.Storage
	converter *converter.Converter
	tokenMgr  *auth.TokenManager
}

// New creates new application server.
func New(repo *repository.Repository, storage *storage.Storage, converter *converter.Converter, tokenMgr *auth.TokenManager) *Server {
	return &Server{
		repo:      repo,
		storage:   storage,
		converter: converter,
		tokenMgr:  tokenMgr,
	}
}

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	api := r.NewRoute().Subrouter()
	api.Use(s.tokenMgr.IsAuthorized)

	r.HandleFunc("/user/signup", s.SignUp).Methods("POST")
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
	api.HandleFunc("/docs", s.ShowDoc).Methods("GET")
	api.HandleFunc("/conversion", s.ConversionRequest).Methods("POST")
	api.HandleFunc("/request_history", s.RequestHistory).Methods("GET")
	api.HandleFunc("/download_audio/{id}", s.Download).Methods("GET")
}

// ShowDoc shows service documentation.
func (s *Server) ShowDoc(w http.ResponseWriter, r *http.Request) {
	res.Respond(w, http.StatusOK, "Showing documentation")
}

// SignUp implements user's signing up.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string
		Password string
	}
	type response struct {
		ID string `json:"id"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.RespondErr(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if err := ValidateUserCredentials(req.Username, req.Password); err != nil {
		res.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	userID, err := s.repo.InsertUser(req.Username, hash)
	if err == repository.ErrUserAlreadyExists {
		res.RespondErr(w, http.StatusConflict, err)
		return
	}
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	resp := response{
		ID: userID,
	}

	res.Respond(w, http.StatusCreated, resp)
}

// LogIn implements user's logging in.
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string
		Password string
	}
	type response struct {
		Token string `json:"token"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.RespondErr(w, http.StatusUnauthorized, err)
		return
	}
	defer r.Body.Close()

	userID, hashedPwd, err := s.repo.GetIDAndPasswordByUsername(req.Username)
	if err == repository.ErrNoSuchUser {
		res.RespondErr(w, http.StatusUnauthorized, err)
		return
	}
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	if !hash.CheckPasswordHash(req.Password, hashedPwd) {
		res.RespondErr(w, http.StatusUnauthorized, errInvalidUsernameOrPassword)
		return
	}

	jwtToken, err := s.tokenMgr.NewJWT(userID)
	if err != nil {
		res.RespondErr(w, http.StatusUnauthorized, err)
		return
	}

	resp := response{
		Token: jwtToken,
	}

	res.Respond(w, http.StatusCreated, resp)
}

// ConversionRequest creates a request for audio conversion.
func (s *Server) ConversionRequest(w http.ResponseWriter, r *http.Request) {
	sourceFile, header, err := r.FormFile("file")
	if err != nil {
		res.RespondErr(w, http.StatusBadRequest, err)
		return
	}
	defer sourceFile.Close()

	sourceContentType := header.Header.Values("Content-type")
	sourceFormat := strings.ToLower(r.FormValue("sourceFormat"))
	targetFormat := strings.ToLower(r.FormValue("targetFormat"))
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)

	if err := ValidateRequest(filename, sourceFormat, targetFormat, sourceContentType[0]); err != nil {
		res.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	fileID, err := s.storage.UploadFile(sourceFile, sourceFormat)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	userID, ok := context.UserIDFromContext(r.Context())
	if !ok {
		res.RespondErr(w, http.StatusInternalServerError, fmt.Errorf("can't retrieve user id from context"))
		return
	}

	requestID, err := s.repo.MakeRequest(filename, sourceFormat, targetFormat, fileID, userID)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	go s.converter.Convert(fileID, filename, sourceFormat, targetFormat, requestID)

	type response struct {
		ID string `json:"id"`
	}
	convertResp := response{
		ID: requestID,
	}

	res.Respond(w, http.StatusAccepted, convertResp)
}

// RequestHistory shows request history of a user.
func (s *Server) RequestHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := context.UserIDFromContext(r.Context())
	if !ok {
		res.RespondErr(w, http.StatusInternalServerError, fmt.Errorf("can't retrieve user id from context"))
		return
	}

	resp, err := s.repo.GetRequestHistory(userID)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	res.Respond(w, http.StatusOK, resp)
}

// Download implements audio downloading.
func (s *Server) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	audioInfo, err := s.repo.GetAudioByID(id)
	if err == repository.ErrNoSuchAudio {
		res.RespondErr(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	file, err := s.storage.DownloadFile(audioInfo.Location, audioInfo.Format)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", formats[audioInfo.Format])
	w.Header().Set("Content-Disposition", "attachment; filename="+audioInfo.Name+"."+audioInfo.Format)

	_, err = io.Copy(w, file)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
		return
	}
}
