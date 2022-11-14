//FILENAME: Interal/data/users.go

package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64 	`json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password password `json:"-"`
	Activated bool `json:"activated"`
	Version int `json:"-"`
}

//Create a customer password type
type password struct {
	plaintext *string
	hash []byte 
}

//The set() Method stores the hash of the plaintexxt password
func (p *password) Set(plaintextPassword string) error {
	hash,err := bcrypt.GenerateFromPassword([]btye(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

//Check