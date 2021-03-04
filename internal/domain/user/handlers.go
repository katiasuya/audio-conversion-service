package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/web"
)

type userHandler struct {
	db DataBase
}

// SignUp implements user's signing up.
func (uh *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if err := u.validate(); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := uh.db.InsertUser(&u); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	type signupResponse struct {
		ID string `json:"id"`
	}
	resp := signupResponse{
		ID: "1fa85f64-5717-4562-b3fc-2c963f66afa5",
	}
	web.Respond(w, http.StatusCreated, resp)
}

// LogIn implements user's logging in.
func (uh *userHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		web.RespondErr(w, http.StatusUnauthorized, err)
		return
	}
	defer r.Body.Close()

	if !uh.db.ComparePasswords(&u) {
		web.RespondErr(w, http.StatusUnauthorized, errors.New("wrong password"))
		return
	}

	type loginResponse struct {
		Token string `json:"token"`
	}
	resp := loginResponse{
		Token: "eyJhbGciOiJIUzI1NiIs...",
	}
	web.Respond(w, http.StatusCreated, resp)

}

// RegisterUserRoutes registers /user endpoints.
func RegisterUserRoutes(db *sql.DB, r *mux.Router) {
	h := &userHandler{db: DataBase{dbase: *db}}
	r.HandleFunc("/user/signup", h.SignUp).Methods("POST")
	r.HandleFunc("/user/login", h.LogIn).Methods("POST")
}
