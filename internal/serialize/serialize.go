package serialize

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/keystop/YaPracticum.git/internal/models"
)

type reader struct {
	file    *os.File
	decoder *gob.Decoder
}

type writer struct {
	file    *os.File
	encoder *gob.Encoder
}

func (w *writer) Close() {
	w.file.Close()
}

func (r *reader) Close() {
	r.file.Close()
}

func newWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	file.Sync()
	if err != nil {
		return nil, errors.New("не удалось найти файл " + fileName)
	}
	return &writer{
		file:    file,
		encoder: gob.NewEncoder(file),
	}, nil

}

func newReader(fileName string) (*reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, errors.New("не удалось найти файл " + fileName)
	}
	return &reader{
		file:    file,
		decoder: gob.NewDecoder(file),
	}, nil

}

// SaveRepoToFile save local db to file.
func SaveRepoToFile(rep models.Repository, fileName string) {
	w, err := newWriter(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.encoder.Encode(rep)
	w.Close()
}

// ReadRepoFromFile read from local file to local base.
func ReadRepoFromFile(rep models.Repository, fileName string) {
	r, err := newReader(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	r.decoder.Decode(rep)
	r.Close()
}
