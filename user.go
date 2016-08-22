package main

type User struct {
	Username string `db:"username" json:"username"`
	Pass     string `db:"password" json:"password"`
}
