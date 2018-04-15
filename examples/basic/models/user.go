package models

import (
	"errors"
	"strings"
)

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Rating int64  `json:"rating"`
}

var (
	ErrNameTooShort = errors.New("name should be >= 3 symbols")
	ErrWrongEmail   = errors.New("email looks bad")
)

func NewUser(name, email string) (*User, error) {
	if len(name) <= 3 {
		return nil, ErrNameTooShort
	}

	if !strings.Contains(email, "@") {
		return nil, ErrWrongEmail
	}

	return &User{
		Name:  name,
		Email: email,
	}, nil
}
