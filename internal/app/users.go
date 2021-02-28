package app

import (
	"encoding/json"
	"net/http"
)

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	ID string `json:"id"`
}

func (ur *userRequest) validateCredentials() error {
	// return  errors.New("invalid login/password")
	return nil
}

func (ur *userRequest) exists() error {
	// return  errors.New("the user already exists")
	return nil
}

func handlerSignUp(w http.ResponseWriter, r *http.Request) {
	var ur userRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	defer r.Body.Close()
	if err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := ur.validateCredentials(); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	if err := ur.exists(); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	userResp := userResponse{
		ID: "1fa85f64-5717-4562-b3fc-2c963f66afa5",
	}
	respond(w, http.StatusCreated, userResp)
}

func handlerLogIn(w http.ResponseWriter, r *http.Request) {
	var ur userRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	defer r.Body.Close()
	if err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	generatedToken := "eyJhbGciOiJIUzI1NiIs..."
	respond(w, http.StatusOK, generatedToken)
}
