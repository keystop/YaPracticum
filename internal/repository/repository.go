package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	encription "github.com/keystop/YaPracticum.git/internal/Encription"
	"github.com/keystop/YaPracticum.git/internal/models"
	"github.com/keystop/YaPracticum.git/internal/shorter"
	"github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
)

var serializeURLRepo func(models.Repository)

//UrlsData repository of urls. Realize Repository interface.
type ServerRepo struct {
	connStr string
	db      *sql.DB
	ctx     context.Context
}

type UsersRepo struct {
	Data      map[string]int
	CurrentID int
}

type urlInfo struct {
	Shorten  string
	Original string
	CorID    string
}

func (s *ServerRepo) SaveURL(url, baseURL, userID string) (string, error) {

	r := shorter.MakeShortner(url)
	u := urlInfo{
		Shorten:  r,
		Original: url,
		CorID:    uuid.New().String(),
	}
	us := []urlInfo{u}
	if err := s.saveUrlsToDB(us, baseURL, userID); err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code == pgerrcode.UniqueViolation {
				return baseURL + r, models.ErrConflictInsert
			}
		}
		return "", err
	}
	return baseURL + r, nil
}

func (s *ServerRepo) GetURL(id string) (string, error) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancelfunc()
	q := `SELECT url FROM urls WHERE shorten_url=$1`
	var url string
	row := db.QueryRowContext(ctx, q, id)

	if err := row.Scan(&url); err != nil {
		return "", err
	}
	return url, nil
}

func (s *ServerRepo) SaveURLs(u map[string]string, baseURL string, userID string) (map[string]string, error) {
	var us []urlInfo
	for k, v := range u {
		r := shorter.MakeShortner(v)
		ui := urlInfo{
			Shorten:  r,
			Original: v,
			CorID:    k,
		}
		u[k] = baseURL + r
		us = append(us, ui)
	}
	if err := s.saveUrlsToDB(us, baseURL, userID); err != nil {

		return u, err
	}
	return u, nil
}

func (s *ServerRepo) GetUserURLs(userEncID string) ([]models.URLs, error) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancelfunc()
	m := make([]models.URLs, 0)
	q := `SELECT url, base_url || shorten_url from urls as u
		  INNER JOIN users as us ON u.user_id=us.id
		  where us.user_enc_id=$1`
	rows, err := db.QueryContext(ctx, q, userEncID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var g models.URLs
		if err := rows.Scan(&g.OriginalURL, &g.ShortURL); err != nil {
			return m, err
		}
		m = append(m, g)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *ServerRepo) FindUser(userEncID string) (finded bool) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancelfunc()
	q := `SELECT id FROM users WHERE user_enc_id=$1`
	var id int
	row := db.QueryRowContext(ctx, q, userEncID)

	if err := row.Scan(&id); err != nil {
		return false
	}
	if id == 0 {
		return false
	}
	return true
}

func (s *ServerRepo) CreateUser() (string, error) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancelfunc()

	ur := uuid.New()
	urEnc, err := encription.EncriptStr(ur.String())
	if err != nil {
		return "", err
	}
	q := `INSERT INTO users (user_uuid, user_enc_id) VALUES ($1, $2)`

	if _, err := db.ExecContext(ctx, q, ur, urEnc); err != nil {
		return "", err
	}

	return urEnc, nil
}
