package models

import (
	"errors"
)

var ErrConflictInsert error = errors.New("conflict")
var ErrURLSetToDel error = errors.New("url set to delete")
