// Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// A wrapper for our data models
type Models struct {
	Fitness FitnessModel
	Tokens TokenModel
	Users UserModel
}

// NewModels() allows us to create a new Models
func NewModels(db *sql.DB) Models {
	return Models{
		Fitness: FitnessModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users: UserModel{DB: db},
	}
}
