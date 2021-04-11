package main

import "database/sql"

type user struct {
	score       int
	in_progress map[string]int
}

type item struct {
	code  string
	class string
}

type record struct {
	user      string
	residence string
	class     string
	amount    int
}

func (u *user) dbGetUserScore(db *sql.DB, uid string) error {
	var err error
	return err
}

func (i *item) dbGetItemByCode(db *sql.DB) error {

}

func (r *record) dbCreateRecord(db *sql.DB) error {

}

func dbCompareCoords(db *sql.DB, coords string) error {
	// if close enough update status in db
	// else do nothings
}
