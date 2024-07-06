package models

import (
	"time"
)

type User struct {
	Id               int       `db:"id"`
	Email            string    `db:"email"`
	Password         string    `db:"password"`
	RegistrationTime time.Time `db:"registration_time"`
}
