package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

//CheckDBConnection trying connect to db.
func (s *ServerRepo) Close() {
	s.db.Close()
}

func (s *ServerRepo) saveUrlsToDB(urls []urlInfo, baseURL, userID string) error {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancelfunc()

	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	q := `INSERT INTO urls (
		correlation_id,
		shorten_url,
		url,
		base_url,
		user_id
	  ) VALUES ($1,$2,$3,$4,(SELECT COALESCE(id, 0) FROM users where user_enc_id=$5))`
	if len(urls) > 1 {
		q += ` ON CONFLICT (url) DO NOTHING`
	}

	pc, err := transaction.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer pc.Close()

	for _, url := range urls {
		if _, err := pc.ExecContext(ctx,
			url.CorID,
			url.Shorten,
			url.Original,
			baseURL,
			userID,
		); err != nil {
			return err
		}
	}

	transaction.Commit()

	return nil

}

//CheckDBConnection trying connect to db.
func (s *ServerRepo) CheckDBConnection() error {
	err := s.db.PingContext(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerRepo) createTables() error {
	db := s.db
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	q := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_uuid VARCHAR(36),
		user_enc_id VARCHAR(36),
		date_add timestamp
	)`
	if _, err := db.ExecContext(ctx, q); err != nil {
		return err
	}

	q = `CREATE TABLE IF NOT EXISTS urls (
		id SERIAL NOT NULL,
		correlation_id VARCHAR(36),  
		shorten_url VARCHAR(32),
		url VARCHAR(255) UNIQUE,
		base_url VARCHAR(255),
		user_id INTEGER REFERENCES users (id),
		date_add TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	if _, err := db.ExecContext(ctx, q); err != nil {
		return err
	}
	return nil
}

func NewServerRepo(connectionString string) (*ServerRepo, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	sr := &ServerRepo{
		connStr: connectionString,
		db:      db,
		ctx:     context.Background(),
	}
	if err := sr.CheckDBConnection(); err != nil {
		return nil, err
	}

	if err := sr.createTables(); err != nil {
		return nil, err
	}
	return sr, nil
}
