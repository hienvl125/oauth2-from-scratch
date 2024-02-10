package models

type User struct {
	ID             string `db:"id"`
	Email          string `db:"email"`
	HashedPassword string `db:"hashed_password"`
}
