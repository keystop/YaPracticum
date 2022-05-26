package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

//CheckDBConnection trying connect to db.
func (s *ServerRepo) Close() {
	s.cancel()
	s.db.Close()
}

func (s *ServerRepo) flushDBuf(ctx context.Context) error {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelfunc()

	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	// q := `UPDATE urls SET for_delete = true WHERE correlation_id=$1 and user_id = $2`

	q := `UPDATE urls SET for_delete = true WHERE shorten_url=$1 and user_id = $2`

	pc, err := t.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer pc.Close()

	for _, v := range s.dBuf {
		if _, err := pc.ExecContext(ctx,
			v.url,
			v.id,
		); err != nil {
			return err
		}
	}

	t.Commit()

	s.dBuf = s.dBuf[:0]

	return nil

}

func (s *ServerRepo) addDBuf(ctx context.Context, v delBufRow) {
	s.dBuf = append(s.dBuf, v)
	if cap(s.dBuf) == len(s.dBuf) {
		s.flushDBuf(ctx)
	}
}

func (s *ServerRepo) saveUrlsToDB(ctx context.Context, us []urlInfo, baseURL, userID string) error {
	db := s.db
	ctx, cancelfunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelfunc()

	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	q := `INSERT INTO urls (
		correlation_id,
		shorten_url,
		url,
		base_url,
		user_id
	  ) VALUES ($1,$2,$3,$4,(SELECT COALESCE(id, 0) FROM users where user_enc_id=$5))`
	if len(us) > 1 {
		q += ` ON CONFLICT (url) DO NOTHING`
	}

	pc, err := t.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer pc.Close()

	for _, u := range us {
		if _, err := pc.ExecContext(ctx,
			u.CorID,
			u.Shorten,
			u.Original,
			baseURL,
			userID,
		); err != nil {
			return err
		}
	}

	t.Commit()

	return nil

}

func (s *ServerRepo) delUrls(ctx context.Context) {
	db := s.db
	ctx, cancelFunc := context.WithTimeout(ctx, 20*time.Second)
	defer cancelFunc()

	t, err := db.Begin()
	if err != nil {
		fmt.Println("Ошибка удаления URL: ", err.Error())
	}
	defer t.Rollback()
	q := `DELETE FROM URLS WHERE for_delete=true`

	if _, err := t.ExecContext(ctx, q); err != nil {
		fmt.Println("Ошибка удаления URL: ", err.Error())
	}

	t.Commit()

}

//CheckDBConnection trying connect to db.
func (s *ServerRepo) CheckDBConnection(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerRepo) createTables(ctx context.Context) error {
	db := s.db
	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
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
		date_add TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		for_delete BOOLEAN DEFAULT false
	)`
	if _, err := db.ExecContext(ctx, q); err != nil {
		return err
	}
	return nil
}

func NewServerRepo(ctx context.Context, c string) (*ServerRepo, error) {
	db, err := sql.Open("postgres", c)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	tP := 1 * time.Second
	timer := time.NewTimer(tP)

	sr := &ServerRepo{
		connStr: c,
		db:      db,
		cancel:  cancel,
		dBuf:    make([]delBufRow, 0, 1000),
		delCh:   make(chan delBufRow, 100),
		timer:   timer,
		dur:     tP,
	}
	if err := sr.CheckDBConnection(ctx); err != nil {
		return nil, err
	}

	if err := sr.createTables(ctx); err != nil {
		return nil, err
	}

	go sr.addURLToDel(ctx)
	return sr, nil
}
