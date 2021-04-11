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
	rows, err := db.Query(`SELECT amount, statuses.name FROM history 
							LEFT JOIN users ON history.uid=users.id 
							LEFT JOIN statuses ON status_id=statuses.id 
							WHERE users.fid=$1`, uid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var score int
		var status string
		if err := rows.Scan(&score, &status); err != nil {
			return err
		}
		u.score += score
		u.in_progress[status] += score
	}

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
