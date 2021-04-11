package main

import "database/sql"

type user struct {
	score       int
	in_progress map[string]int
}

func (u *user) dbGetUserScore(db *sql.DB, uid int) error {
	var err error
	return err
}
