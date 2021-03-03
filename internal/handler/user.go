package handler

import (
	"encoding/json"
	"net/http"

	"github.com/katiasuya/audio-conversion-service/internal/web"
)

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ur *userRequest) validateCredentials() error {
	// return  errors.New("invalid login/password")
	return nil
}

func (ur *userRequest) exists() error {
	// return  errors.New("the user already exists")
	return nil
}

// SignUp implements user's signing up.
func SignUp(w http.ResponseWriter, r *http.Request) {
	var ur userRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	defer r.Body.Close()
	if err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := ur.validateCredentials(); err != nil {
		web.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := ur.exists(); err != nil {
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
func LogIn(w http.ResponseWriter, r *http.Request) {
	var ur userRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	defer r.Body.Close()
	if err != nil {
		web.RespondErr(w, http.StatusUnauthorized, err)
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
