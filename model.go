package main

import "database/sql"

type user struct {
	score       int	`json: "score"`
	in_progress map[string]int `json: "in_progress"`
}

type item struct {
	code  string `json: "code"`
	class string `json: "type"`
}

type record struct {
	user      string `json: "user"`
	residence string `json: "residence"`
	class     string `json: "type"`
	amount    int `json: "amount"`
}

type residence struct {
	name        string
	coordinates string
}

func (u *user) dbGetUserScore(db *sql.DB, uid string) error {
	rows, err := 
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
