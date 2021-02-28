package app

import (
	"fmt"
	"net/http"
)

func handlerShowDoc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Showing documentation")
}
