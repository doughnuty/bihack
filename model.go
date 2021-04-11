package main

import "database/sql"

type user struct {
	score       int            `json: "score"`
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
	amount    int    `json: "amount"`
}

type residence struct {
	name        string
	coordinates string
}

func (u *user) dbGetUserScore(db *sql.DB, uid string) error {
	// select amount from history left join users where user.name=$1 and status=1
	// select amount from history left join users where user.name=$1 and status=2
	// select amount from history left join users where user.name=$1 and status=3
	// add to map in progress each with key type name and amount as value
	// sum all the values and add them to score

	return err
}

func (i *item) dbGetItemByCode(db *sql.DB) error {
	// insert into types if conflict select from its type id
	// or check if type empty: y - select, n - insert
}

func (r *record) dbCreateRecord(db *sql.DB) error {
	// insert into history all data
}

func dbCompareCoords(db *sql.DB, coords string) error {
	// if close enough update status in db
	// else do nothings
}
