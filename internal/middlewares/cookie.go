package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/keystop/YaPracticum.git/internal/models"
)

var Repo models.Repository

func SetCookieUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("UserID")
		cv := ""
		if err == nil {
			cv = c.Value
		}
		if ok := Repo.FindUser(r.Context(), cv); !ok {
			cv, err = Repo.CreateUser(r.Context())
			if err != nil {
				fmt.Println("Can't create cookie", err)
				next.ServeHTTP(w, r)
				return
			}
		}

		c = &http.Cookie{
			Name:  "UserID",
			Value: cv,
		}
		http.SetCookie(w, c)

		ctx := context.WithValue(r.Context(), models.UserKey, cv)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewCookie(repo models.Repository) {
	Repo = repo
}
