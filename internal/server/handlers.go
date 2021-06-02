// Package server implements http handlers.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/events"
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

// IsAuthorized is a middleware that checks user authorization.
func (s *Server) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			res.RespondErr(w, http.StatusUnauthorized, errors.New("malformed token"))
			return
		}

		jwtToken := authHeader[1]
		claimUserID, err := s.tokenMgr.ParseJWT(jwtToken)
		if err != nil {
			res.RespondErr(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.ContextWithUserID(r.Context(), claimUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RegisterRoutes registers application rotes.
func (s *Server) RegisterRoutes(r *mux.Router) {
	api := r.NewRoute().Subrouter()
	api.Use(s.IsAuthorized)

	// r.HandleFunc("/user/signup", s.SignUp).Methods("POST")
	r.HandleFunc("/user/login", s.LogIn).Methods("POST")
	// api.HandleFunc("/docs", s.ShowDoc).Methods("GET")
	api.HandleFunc("/", s.TestRequest).Methods("GET")
	api.HandleFunc("/conversion", s.ConversionRequest).Methods("POST")
	api.HandleFunc("/request_history", s.RequestHistory).Methods("GET")
	api.HandleFunc("/download_audio/{id}", s.Download).Methods("GET")
}

// ShowDoc shows service documentation.
func (s *Server) ShowDoc(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Showing documentation",
		StatusCode: 200,
	}, nil
}

func (s *Server) TestRequest(w http.ResponseWriter, r *http.Request) {
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["sub"] == nil {
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	userID := requestContext.Authorizer["sub"].(string)
	fmt.Fprintf(w, "Hello %s from Go", userID)
}

func createResponse(statusCode int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode:        statusCode,
		Headers:           nil,
		MultiValueHeaders: nil,
		Body:              body,
		IsBase64Encoded:   false,
	}, nil
}

var S *Server

// SignUp implements user's signing up.
func SignUp(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := json.Unmarshal([]byte(r.Body), &req); err != nil {
		return createResponse(500, fmt.Sprintf("can't parse body, %s", err.Error()))
	}

	if err := ValidateUserCredentials(req.Username, req.Password); err != nil {
		return createResponse(400, fmt.Sprintf("Invalid user credentials: %v", err))
	}

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		return createResponse(500, fmt.Sprintf("Hashed passwords don't match: %v", err))
	}

	fmt.Println("hash ", hash)

	fmt.Println("here01")
	userID, err := S.repo.InsertUser(req.Username, hash)
	fmt.Println("here02")

	if err == repository.ErrUserAlreadyExists {
		return createResponse(409, fmt.Sprintf("can't insert user: %v", err))
	}
	if err != nil {
		return createResponse(500, fmt.Sprintf("can't insert user: %v", err))
	}

	return createResponse(200, fmt.Sprintf("Hi, %s", userID))
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

	fileURL, err := s.storage.GetDownloadURL(audioInfo.Location, audioInfo.Format)
	if err != nil {
		res.RespondErr(w, http.StatusInternalServerError, err)
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
