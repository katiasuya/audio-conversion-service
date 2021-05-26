// Package server implements http handlers.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/auth"
	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/mycontext"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
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
		logger := mycontext.LoggerFromContext(r.Context())

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logAndRespondErr(r.Context(), logger, w, "", errors.New("malformed token"), http.StatusUnauthorized)
			return
		}

		jwtToken := authHeader[1]
		claimUserID, err := s.tokenMgr.ParseJWT(jwtToken)
		if err != nil {
			logAndRespondErr(r.Context(), logger, w, "can't parse JWT: ", err, http.StatusUnauthorized)
			return
		}

		ctx := mycontext.ContextWithUserID(r.Context(), claimUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) AssignRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rqID, err := uuid.NewRandom()
		if err != nil {
			errMsg := fmt.Errorf("can't generate request id"+": %w", err)
			log.Errorln(errMsg)
			res.RespondErr(r.Context(), w, http.StatusInternalServerError, errMsg)
			return
		}
		rqCtx := mycontext.ContextWithRequestID(r.Context(), rqID.String())

		next.ServeHTTP(w, r.WithContext(rqCtx))
	})
}

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	r.Use(s.AssignRequestID)
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
	res.Respond(r.Context(), w, http.StatusOK, "Showing documentation")
}

// SignUp implements user's signing up.
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	logger := mycontext.LoggerFromContext(r.Context())

	type request struct {
		Username string
		Password string
	}
	type response struct {
		ID string `json:"id"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't decode request body: ", err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := ValidateUserCredentials(req.Username, req.Password); err != nil {
		logAndRespondErr(r.Context(), logger, w, "invalid user credentials: ", err, http.StatusBadRequest)
		return
	}
	logger.Debugln("user's credentials validated")

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't hash password: ", err, http.StatusInternalServerError)
		return
	}

	userID, err := s.repo.InsertUser(req.Username, hash)
	if err == repository.ErrUserAlreadyExists {
		logAndRespondErr(r.Context(), logger, w, "can't insert user: ", err, http.StatusConflict)
		return
	}
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't insert user: ", err, http.StatusInternalServerError)
		return
	}

	resp := response{
		ID: userID,
	}

	logger.WithField("userID", userID).Infoln("user signed up successfully")

	res.Respond(r.Context(), w, http.StatusCreated, resp)
}

// LogIn implements user's logging in.
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	logger := mycontext.LoggerFromContext(r.Context())

	type request struct {
		Username string
		Password string
	}
	type response struct {
		Token string `json:"token"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't decode request body: ", err, http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	userID, hashedPwd, err := s.repo.GetIDAndPasswordByUsername(req.Username)
	if err == repository.ErrNoSuchUser {
		logAndRespondErr(r.Context(), logger, w, "can't get user id and password: ", err, http.StatusUnauthorized)
		return
	}
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't get user id and password: ", err, http.StatusInternalServerError)
		return
	}

	if !hash.CheckPasswordHash(req.Password, hashedPwd) {
		logAndRespondErr(r.Context(), logger, w, "", errInvalidUsernameOrPassword, http.StatusUnauthorized)
		return
	}
	logger.Debugln("log in: hashed passwords match")

	jwtToken, err := s.tokenMgr.NewJWT(userID)
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't create JWT: ", err, http.StatusUnauthorized)
		return
	}

	resp := response{
		Token: jwtToken,
	}

	logger.WithField("userID", userID).Infoln("user authenticated and authorized successfully")

	res.Respond(r.Context(), w, http.StatusCreated, resp)
}

// ConversionRequest creates a request for audio conversion.
func (s *Server) ConversionRequest(w http.ResponseWriter, r *http.Request) {
	logger := mycontext.LoggerFromContext(r.Context())

	sourceFile, header, err := r.FormFile("file")
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't get file from the form: ", err, http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()
	logger.Debugln("source file recieved from form successfully")

	sourceContentType := header.Header.Values("Content-type")
	sourceFormat := strings.ToLower(r.FormValue("sourceFormat"))
	targetFormat := strings.ToLower(r.FormValue("targetFormat"))
	filename := strings.TrimSuffix(header.Filename, "."+sourceFormat)

	if err = ValidateRequest(filename, sourceFormat, targetFormat, sourceContentType[0]); err != nil {
		logAndRespondErr(r.Context(), logger, w, "invalid request: ", err, http.StatusBadRequest)
		return
	}
	logger.Debugln("conversion request validated")

	fileID, err := s.storage.UploadFile(sourceFile, sourceFormat)
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't upload file: ", err, http.StatusInternalServerError)
		return
	}
	logger.Debugln("source file uploaded successfully")

	userID, ok := mycontext.UserIDFromContext(r.Context())
	if !ok {
		logAndRespondErr(r.Context(), logger, w, "", errCantGetUserIDFomContext, http.StatusInternalServerError)
		return
	}

	requestID, err := s.repo.MakeRequest(filename, sourceFormat, targetFormat, fileID, userID)
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't make conversion request: ", err, http.StatusInternalServerError)
		return
	}
	logger.WithField("convReqID", requestID).Infoln("conversion request created successfully")

	go s.converter.Convert(r.Context(), fileID, filename, sourceFormat, targetFormat, requestID)

	type response struct {
		ID string `json:"id"`
	}
	convertResp := response{
		ID: requestID,
	}

	res.Respond(r.Context(), w, http.StatusAccepted, convertResp)
}

// RequestHistory shows request history of a user.
func (s *Server) RequestHistory(w http.ResponseWriter, r *http.Request) {
	logger := mycontext.LoggerFromContext(r.Context())

	userID, ok := mycontext.UserIDFromContext(r.Context())
	if !ok {
		logAndRespondErr(r.Context(), logger, w, "", errCantGetUserIDFomContext, http.StatusInternalServerError)
		return
	}

	resp, err := s.repo.GetRequestHistory(userID)
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't get request history: ", err, http.StatusInternalServerError)
		return
	}

	logger.Infoln("request history got successfully")

	res.Respond(r.Context(), w, http.StatusOK, resp)
}

// Download implements audio downloading.
func (s *Server) Download(w http.ResponseWriter, r *http.Request) {
	logger := mycontext.LoggerFromContext(r.Context())

	vars := mux.Vars(r)
	id := vars["id"]

	audioInfo, err := s.repo.GetAudioByID(id)
	if err == repository.ErrNoSuchAudio {
		logAndRespondErr(r.Context(), logger, w, "can't get audio: ", err, http.StatusNotFound)
		return
	}
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't get audio: ", err, http.StatusInternalServerError)
		return
	}

	fileURL, err := s.storage.GetDownloadURL(audioInfo.Location, audioInfo.Format)
	if err != nil {
		logAndRespondErr(r.Context(), logger, w, "can't get download URL: ", err, http.StatusInternalServerError)
		return
	}

	type response struct {
		FileURL string `json:"fileURL"`
	}
	downloadResp := response{
		FileURL: fileURL,
	}

	logger.WithField("audioID", id).Infoln("download URL created successfully")
	res.Respond(r.Context(), w, http.StatusOK, downloadResp)
}

func logAndRespondErr(ctx context.Context, logger *log.Entry, w http.ResponseWriter, wrapper string, err error, code int) {
	errMsg := fmt.Errorf(wrapper+"%w", err)
	logger.Errorln(errMsg)
	res.RespondErr(ctx, w, code, errMsg)
}
