package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katiasuya/audio-conversion-service/internal/domain"
)

func userValidation(w http.ResponseWriter, r *http.Request) (domain.User, error) {

	var u domain.User
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return domain.User{}, err
	}

	if err := u.ValidateCredentials(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return domain.User{}, err
	}

	return u, nil
}

func authValidation(w http.ResponseWriter, u domain.User) error {
	if err := u.IsAuthorized(); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}
	return nil
}

func showDoc(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	if err := authValidation(w, u); err != nil {
		return
	}
	fmt.Fprint(w, "Showing documentation")
}

func signUp(w http.ResponseWriter, r *http.Request) {
	u, err := userValidation(w, r)
	if err != nil {
		return
	}

	if err := u.Exists(); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	userID := "1fa85f64-5717-4562-b3fc-2c963f66afa5"
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(userID))
}

func logIn(w http.ResponseWriter, r *http.Request) {
	u, err := userValidation(w, r)
	if err != nil {
		return
	}

	if err := u.IncorrectPwd(); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	generatedToken := "eyJhbGciOiJIUzI1NiIs..."
	fmt.Fprint(w, generatedToken)
}

func convert(w http.ResponseWriter, r *http.Request) {

	var u domain.User
	if err := authValidation(w, u); err != nil {
		return
	}

	var a domain.Audio
	err := json.NewDecoder(r.Body).Decode(&a)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestID := "2fa85f64-5717-4562-b3fc-2c963f66afa5"
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(requestID))
}

func showHistory(w http.ResponseWriter, r *http.Request) {

	var u domain.User
	if err := authValidation(w, u); err != nil {
		return
	}

	rr := []domain.Request{}
	fmt.Fprint(w, fmt.Sprint(rr))
}

func download(w http.ResponseWriter, r *http.Request) {

	var u domain.User
	if err := authValidation(w, u); err != nil {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := domain.AudioExists(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	audio := domain.Audio{}
	fmt.Fprint(w, fmt.Sprint(audio))
}
