package handlers

import (
	"io"
	"net/http"

	"github.com/keystop/YaPracticum.git/internal/repository"
)

// HandlerUrlPost saves url from request body to repository.
func HandlerURLPost(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		textBody, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(400)
			return
		}
		retURL := "http://" + r.Host + "/" + repo.SaveURL(textBody)
		w.WriteHeader(201)
		io.WriteString(w, retURL)

	}
}

// HandlerUrlGet returns url from repository to resp.Head - "Location".
func HandlerURLGet(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := ctx.Value(repository.Key("id")).(string)
		val, err := repo.GetURL(id)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}
		w.Header().Add("Location", val)
		w.WriteHeader(307)
	}
}
