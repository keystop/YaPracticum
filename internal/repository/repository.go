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

type delBufRow struct {
	url string
	id  int
}

type delBuf []delBufRow

//UrlsData repository of urls. Realize Repository interface.
type ServerRepo struct {
	connStr string
	db      *sql.DB
	cancel  context.CancelFunc
	dBuf    delBuf
	delCh   chan delBufRow
	timer   *time.Timer
	dur     time.Duration
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

func (s *ServerRepo) SaveURL(ctx context.Context, url, baseURL, userID string) (string, error) {

	r := shorter.MakeShortner(url)
	u := urlInfo{
		Shorten:  r,
		Original: url,
		CorID:    uuid.New().String(),
	}
	us := []urlInfo{u}
	if err := s.saveUrlsToDB(ctx, us, baseURL, userID); err != nil {
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

func (s *ServerRepo) GetURL(ctx context.Context, id string) (string, error) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()
	q := `SELECT url, for_delete FROM urls WHERE shorten_url=$1`
	var url string
	var forD bool
	row := db.QueryRowContext(ctx, q, id)

	if err := row.Scan(&url, &forD); err != nil {
		return "", err
	}
	if forD {
		return "", models.ErrURLSetToDel
	}
	return url, nil
}

func (s *ServerRepo) SetURLsToDel(ctx context.Context, d []string, userID string) error {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelfunc()

	q := `SELECT id FROM users WHERE user_enc_id=$1`
	var id int
	row := db.QueryRowContext(ctx, q, userID)

	if err := row.Scan(&id); err != nil {
		return err
	}
	go func() {
		for _, v := range d {
			s.delCh <- delBufRow{v, id}
		}
	}()

	return nil
}

func (s *ServerRepo) addURLToDel(ctx context.Context) {
	timerCounter := 0
	for {
		select {
		case <-s.timer.C:
			timerCounter += 1
			if timerCounter == 4 {
				s.delUrls(ctx)
				timerCounter = 0
			}
			s.flushDBuf(ctx)
			s.timer.Reset(s.dur)
		case v := <-s.delCh:
			s.addDBuf(ctx, v)
		}
	}
}

// func (s *ServerRepo) AddTo

func (s *ServerRepo) SaveURLs(ctx context.Context, u map[string]string, baseURL string, userID string) (map[string]string, error) {
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
	if err := s.saveUrlsToDB(ctx, us, baseURL, userID); err != nil {

		return u, err
	}
	return u, nil
}

func (s *ServerRepo) GetUserURLs(ctx context.Context, userEncID string) ([]models.URLs, error) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
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

func (s *ServerRepo) FindUser(ctx context.Context, userEncID string) (finded bool) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
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

func (s *ServerRepo) CreateUser(ctx context.Context) (string, error) {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 5*time.Second)
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
