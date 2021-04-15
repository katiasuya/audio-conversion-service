// Package server implements http handlers.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/server/context"
	"github.com/katiasuya/audio-conversion-service/internal/server/response"
	res "github.com/katiasuya/audio-conversion-service/internal/server/response"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"github.com/katiasuya/audio-conversion-service/pkg/hash"
	log "github.com/sirupsen/logrus"
)

var (
	errInvalidUsernameOrPassword = errors.New("invalid username or password")
	errCantGetUserIDFomContext   = errors.New("can't retrieve user id from context")
)

// Server represents application server.
type Server struct {
	repo      *repository.Repository
	storage   *storage.Storage
	converter *converter.Converter
	tokenMgr  *auth.TokenManager
	logger    *log.Entry
}

// New creates new application server.
func New(repo *repository.Repository, storage *storage.Storage, converter *converter.Converter, tokenMgr *auth.TokenManager, logger *log.Entry) *Server {
	return &Server{
		repo:      repo,
		storage:   storage,
		converter: converter,
		tokenMgr:  tokenMgr,
		logger:    logger,
	}
}

// IsAuthorized is a middleware that checks user authorization.
func (s *Server) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			s.logAndRespondErr(w, "", errors.New("malformed token"), http.StatusUnauthorized)
			return
		}

		jwtToken := authHeader[1]
		claimUserID, err := s.tokenMgr.ParseJWT(jwtToken)
		if err != nil {
			s.logAndRespondErr(w, "can't parse JWT: ", err, http.StatusUnauthorized)
			return
		}

		s.logger.Debugln("JWT was parsed and validated successfully")

		ctx := context.ContextWithUserID(r.Context(), claimUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	api := r.NewRoute().Subrouter()
	api.Use(s.IsAuthorized)

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
		s.logAndRespondErr(w, "can't decode the request body: ", err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := ValidateUserCredentials(req.Username, req.Password); err != nil {
		s.logAndRespondErr(w, "invalid user credentials: ", err, http.StatusBadRequest)
		return
	}
	s.logger.Debugln("user's credentials validated successfully")

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		s.logAndRespondErr(w, "can't hash password: ", err, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("password hashed successfully")

	userID, err := s.repo.InsertUser(req.Username, hash)
	if err == repository.ErrUserAlreadyExists {
		s.logAndRespondErr(w, "can't insert user: ", err, http.StatusConflict)
		return
	}
	if err != nil {
		s.logAndRespondErr(w, "can't insert user: ", err, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("user inserted successfully")

	resp := response{
		ID: userID,
	}

	s.logger.WithField("userID", userID).Infoln("user signed up successfully")
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
		s.logAndRespondErr(w, "can't decode the request body: ", err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	userID, hashedPwd, err := s.repo.GetIDAndPasswordByUsername(req.Username)
	if err == repository.ErrNoSuchUser {
		s.logAndRespondErr(w, "can't get user id and password: ", err, http.StatusUnauthorized)
		return
	}
	if err != nil {
		s.logAndRespondErr(w, "can't get user id and password: ", err, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("user's id and password got successfully")

	if !hash.CheckPasswordHash(req.Password, hashedPwd) {
		s.logAndRespondErr(w, "", errInvalidUsernameOrPassword, http.StatusUnauthorized)
		return
	}
	s.logger.Debugln("hashed passwords match")

	jwtToken, err := s.tokenMgr.NewJWT(userID)
	if err != nil {
		s.logAndRespondErr(w, "can't create JWT: ", err, http.StatusUnauthorized)
		return
	}
	s.logger.Debugln("JWT was created successfully")

	resp := response{
		Token: jwtToken,
	}

	s.logger.WithField("userID", userID).Infoln("user authenticated successfully")
	res.Respond(w, http.StatusCreated, resp)
}

// ConversionRequest creates a request for audio conversion.
func (s *Server) ConversionRequest(w http.ResponseWriter, r *http.Request) {
	sourceFile, header, err := r.FormFile("file")
	if err != nil {
		s.logAndRespondErr(w, "can't get file from the form: ", err, http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	sourceContentType := header.Header.Values("Content-type")
	sourceFormat := strings.ToLower(r.FormValue("sourceFormat"))
	targetFormat := strings.ToLower(r.FormValue("targetFormat"))
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)

	if err := ValidateRequest(filename, sourceFormat, targetFormat, sourceContentType[0]); err != nil {
		s.logAndRespondErr(w, "invalid request: ", err, http.StatusBadRequest)
		return
	}
	s.logger.Debugln("request validated successfully")

	fileID, err := s.storage.UploadFile(sourceFile, sourceFormat)
	if err != nil {
		s.logAndRespondErr(w, "can't upload the file: ", err, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("file uploaded successfully")

	userID, ok := context.UserIDFromContext(r.Context())
	if !ok {
		s.logAndRespondErr(w, "", errCantGetUserIDFomContext, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("user id retrieved from the context successfully")

	requestID, err := s.repo.MakeRequest(filename, sourceFormat, targetFormat, fileID, userID)
	if err != nil {
		s.logAndRespondErr(w, "can't make conversion request: ", err, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("conversion request made successfully")

	go s.converter.Convert(fileID, filename, sourceFormat, targetFormat, requestID)

	type response struct {
		ID string `json:"id"`
	}
	convertResp := response{
		ID: requestID,
	}

	s.logger.WithField("requestID", requestID).Infoln("requestID created successfully")
	res.Respond(w, http.StatusAccepted, convertResp)
}

// RequestHistory shows request history of a user.
func (s *Server) RequestHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := context.UserIDFromContext(r.Context())
	if !ok {
		s.logAndRespondErr(w, "", errCantGetUserIDFomContext, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("user id retrieved from the context successfully")

	resp, err := s.repo.GetRequestHistory(userID)
	if err != nil {
		s.logAndRespondErr(w, "can't get request history: ", err, http.StatusInternalServerError)
		return
	}

	s.logger.WithField("userID", userID).Infoln("request history got successfully")
	res.Respond(w, http.StatusOK, resp)
}

// Download implements audio downloading.
func (s *Server) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	audioInfo, err := s.repo.GetAudioByID(id)
	if err == repository.ErrNoSuchAudio {
		s.logAndRespondErr(w, "can't get audio: ", err, http.StatusNotFound)
		return
	}
	if err != nil {
		s.logAndRespondErr(w, "can't get audio: ", err, http.StatusInternalServerError)
		return
	}
	s.logger.Debugln("audio found successfully by id")

	fileURL, err := s.storage.GetDownloadURL(audioInfo.Location, audioInfo.Format)
	if err != nil {
		s.logAndRespondErr(w, "can't download the file: ", err, http.StatusInternalServerError)
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

func (s *Server) logAndRespondErr(w http.ResponseWriter, wrapper string, err error, code int) {
	errMsg := fmt.Errorf(wrapper+"%w", err)
	s.logger.Errorln(errMsg)
	response.RespondErr(w, code, errMsg)
}
