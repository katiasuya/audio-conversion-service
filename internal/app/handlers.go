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

func handlerShowDoc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Showing documentation")
}

func handlerSignUp(w http.ResponseWriter, r *http.Request) {
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

func handlerLogIn(w http.ResponseWriter, r *http.Request) {
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

func handlerConvert(w http.ResponseWriter, r *http.Request) {
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

func handlerShowHistory(w http.ResponseWriter, r *http.Request) {
	rr := []domain.Request{}
	fmt.Fprint(w, fmt.Sprint(rr))
}

func handlerDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := domain.AudioExists(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	audio := domain.Audio{}
	fmt.Fprint(w, fmt.Sprint(audio))
}
