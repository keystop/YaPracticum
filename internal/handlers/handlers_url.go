package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/keystop/YaPracticum.git/internal/models"
)

var Repo models.Repository
var BaseURL string
var Opt models.Options

func HandlerUserPostURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(models.UserKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ud, err := Repo.GetUserURLs(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(ud) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	res, err := json.Marshal(ud)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandlerAPIURLsPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(models.UserKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	text, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type uJ struct {
		CorID     string `json:"correlation_id"`
		OriginURL string `json:"original_url"`
	}

	var uJs []uJ

	err = json.Unmarshal(text, &uJs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uts := make(map[string]string)
	for _, u := range uJs {
		uts[u.CorID] = u.OriginURL
	}

	uts, err = Repo.SaveURLs(uts, BaseURL, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type uJR struct {
		CorID    string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}

	var uJsR []uJR

	for key, value := range uts {
		u := uJR{
			CorID:    key,
			ShortURL: value,
		}
		uJsR = append(uJsR, u)
	}

	res, err := json.Marshal(&uJsR)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func HandlerCheckDBConnect(w http.ResponseWriter, r *http.Request) {
	if err := Repo.CheckDBConnection(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandlerURLPost saves url from request body to repository.
func HandlerURLPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(models.UserKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	textBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	retURL, err := Repo.SaveURL(string(textBody), BaseURL, userID)
	if err != nil {
		if err == models.ErrConflictInsert {
			w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(retURL))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(retURL))
}

//HandlerAPIURLPost saves url from body request.
func HandlerAPIURLPost(w http.ResponseWriter, r *http.Request) {
	aSuccessCode := http.StatusCreated
	ctx := r.Context()

	userID, ok := ctx.Value(models.UserKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tURLJson := &struct {
		URLLong string `json:"url"`
	}{}
	textBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(textBody, tURLJson)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	su, err := Repo.SaveURL(tURLJson.URLLong, BaseURL, userID)
	if err != nil {
		switch err {
		case models.ErrConflictInsert:
			aSuccessCode = http.StatusConflict
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
	tResJSON := &struct {
		URLShorten string `json:"result"`
	}{}

	tResJSON.URLShorten = su

	res, err := json.Marshal(tResJSON)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(aSuccessCode)
	w.Write(res)
}

// HandlerURLGet returns url from repository to resp.Head - "Location".
func HandlerURLGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := ctx.Value(models.URLID).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	val, err := Repo.GetURL(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Add("Location", val)
	w.Header().Add("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func NewHandlers(repo models.Repository, opt models.Options) {
	Repo = repo
	BaseURL = opt.RespBaseURL() + "/"
	Opt = opt
}
