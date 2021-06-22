// Package server implements http handlers.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/appcontext"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	res "github.com/katiasuya/audio-conversion-service/internal/server/response"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"github.com/katiasuya/audio-conversion-service/pkg/hash"
)

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

// IsAuthorized is a middleware that checks user authorization.
func (s *Server) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logAndRespondErr(r.Context(), w, "", errors.New("malformed token"), http.StatusUnauthorized)
			return
		}

		jwtToken := authHeader[1]
		claimUserID, err := s.tokenMgr.ParseJWT(jwtToken)
		if err != nil {
			logAndRespondErr(r.Context(), w, "can't parse JWT", err, http.StatusUnauthorized)
			return
		}

		ctx := appcontext.AddUserID(r.Context(), claimUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AddLogger creates logger and adds it to the context.
func (s *Server) AddLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.AddToContext(r.Context(), logger.New())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	r.Use(s.AddLogger)
	api := r.NewRoute().Subrouter()
	api.Use(s.IsAuthorized)

	r.HandleFunc("/signup", s.SignUp).Methods("POST")
	r.HandleFunc("/login", s.LogIn).Methods("POST")
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
		logAndRespondErr(r.Context(), w, "can't decode request body", err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := ValidateUserCredentials(req.Username, req.Password); err != nil {
		res.RespondErr(w, http.StatusBadRequest, fmt.Errorf("invalid user credentials: %w", err))
		return
	}

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't hash password", err, http.StatusInternalServerError)
		return
	}

	userID, err := s.repo.InsertUser(req.Username, hash)
	if err == repository.ErrUserAlreadyExists {
		res.RespondErr(w, http.StatusConflict, fmt.Errorf("can't insert user: %w", err))
		return
	}
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't insert user", err, http.StatusInternalServerError)
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
		logAndRespondErr(r.Context(), w, "can't decode request body", err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	userID, hashedPwd, err := s.repo.GetIDAndPasswordByUsername(req.Username)
	if err == repository.ErrNoSuchUser {
		res.RespondErr(w, http.StatusUnauthorized, fmt.Errorf("can't get user id and password: %w", err))
		return
	}
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't get user id and password", err, http.StatusInternalServerError)
		return
	}

	if !hash.CheckPasswordHash(req.Password, hashedPwd) {
		res.RespondErr(w, http.StatusUnauthorized, errors.New("invalid username or password"))
		return
	}

	jwtToken, err := s.tokenMgr.NewJWT(userID)
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't create JWT", err, http.StatusUnauthorized)
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
		res.RespondErr(w, http.StatusBadRequest, fmt.Errorf("can't get file from the form: %w", err))
		return
	}
	defer sourceFile.Close()

	sourceContentType := header.Header.Values("Content-type")
	sourceFormat := strings.ToLower(r.FormValue("sourceFormat"))
	targetFormat := strings.ToLower(r.FormValue("targetFormat"))
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)

	if err = ValidateRequest(filename, sourceFormat, targetFormat, sourceContentType[0]); err != nil {
		res.RespondErr(w, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		return
	}

	fileID, err := s.storage.UploadFile(sourceFile, sourceFormat)
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't upload file", err, http.StatusInternalServerError)
		return
	}

	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		logAndRespondErr(r.Context(), w, "", errors.New("can't get user id from context"), http.StatusInternalServerError)
		return
	}

	requestID, err := s.repo.MakeRequest(filename, sourceFormat, targetFormat, fileID, userID)
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't make conversion request", err, http.StatusInternalServerError)
		return
	}

	go s.converter.Convert(r.Context(), fileID, filename, sourceFormat, targetFormat, requestID)

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
	userID, ok := appcontext.GetUserID(r.Context())
	if !ok {
		logAndRespondErr(r.Context(), w, "", errors.New("can't get user id from context"), http.StatusInternalServerError)
		return
	}

	resp, err := s.repo.GetRequestHistory(userID)
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't get request history", err, http.StatusInternalServerError)
		return
	}

	res.Respond(w, http.StatusOK, resp)
}

// Download implements audio downloading.
func (s *Server) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	audioID := vars["id"]

	audioInfo, err := s.repo.GetAudioByID(audioID)
	if err == repository.ErrNoSuchAudio {
		res.RespondErr(w, http.StatusNotFound, fmt.Errorf("can't get audio: %w", err))
		return
	}
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't get audio", err, http.StatusInternalServerError)
		return
	}

	fileURL, err := s.storage.GetDownloadURL(audioInfo.Location, audioInfo.Format)
	if err != nil {
		logAndRespondErr(r.Context(), w, "can't get download URL", err, http.StatusInternalServerError)
		return
	}

	type response struct {
		FileURL string `json:"fileURL"`
	}
	downloadResp := response{
		FileURL: fileURL,
	}

	res.Respond(w, http.StatusOK, downloadResp)
}

func logAndRespondErr(ctx context.Context, w http.ResponseWriter, wrapper string, err error, code int) {
	errMsg := fmt.Errorf(wrapper+": %w", err)
	logger.Error(ctx, errMsg)
	res.RespondErr(w, code, errMsg)
}
