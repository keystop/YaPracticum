package models

import (
	"errors"
)

var ErrConflictInsert error = errors.New("conflict")
