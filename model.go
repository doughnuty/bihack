package main

import (
	"database/sql"

	geo "github.com/kellydunn/golang-geo"
)

type user struct {
	Score       int            `json:"score"`
	In_progress map[string]int `json:"in_progress"`
}

type user_data struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Firebase string `json:"fid"`
}

type item struct {
	Code  string `json:"code"`
	Class string `json:"type"`
}

type record struct {
	User      string `json:"user"`
	Residence string `json:"residence"`
	Class     string `json:"type"`
	Amount    int    `json:"amount"`
	GPS_id    int    `json:"gps_id"`
}

// type residence struct {
// 	name        string
// 	coordinates string
// }

type GPS struct {
	ID   string  `json:"id"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (u *user_data) dbInitUser(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (fid, name, surname) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", u.Firebase, u.Name, u.Surname)
	if err != nil {
		return err
	}

	return nil
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
		u.Score += score
		u.In_progress[status] += score
	}

	return nil
}

func (i *item) dbGetItemByCode(db *sql.DB) error {
	err := db.QueryRow("SELECT types.name FROM items LEFT JOIN types ON items.type_id=types.id WHERE scan=$1", i.Code).Scan(i.Class)
	if err != nil {
		return err
	}

	return nil
}

func (i *item) dbAddItem(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO items (scan, type_id) SELECT $1, t.id FROM types t WHERE t.name=$2", i.Code, i.Class)
	if err != nil {
		return err
	}

	return nil
}

func (r *record) dbCreateRecord(db *sql.DB) error {
	// insert into history all data
	_, err := db.Exec(`INSERT INTO history (uid, rid, type_id, amount, date_start, status_id, gps_id) 
						SELECT u.id, r.id, t.id, $1, now(), 1, 0 FROM users u, residentials r, types t 
						WHERE u.fid=$2, r.name=$3, t.name=$4`,
		r.Amount, r.User, r.Residence, r.Class)
	if err != nil {
		return err
	}

	return nil
}

func dbCompareCoords(db *sql.DB, coords GPS) error {
	// if close enough update status in db
	// else do nothings
	rows, err := db.Query(`SELECT lat, long, id FROM residentials`)
	if err != nil {
		return err
	}
	defer rows.Close()

	gps_point := geo.NewPoint(coords.Lat, coords.Long)
	for rows.Next() {
		var lat, long float64
		id := 0
		rows.Scan(&lat, &long, &id)
		row_point := geo.NewPoint(lat, long)

		dist := gps_point.GreatCircleDistance(row_point)
		if dist <= 10 && id != 1 {
			_, err := db.Exec("UPDATE history SET status_id=2, gps_id=$1 WHERE rid=$2 AND status_id=1", coords.ID, id)
			if err != nil {
				return err
			}
		} else if dist <= 10 && id == 1 {
			_, err := db.Exec("UPDATE history SET status_id=3 WHERE gps_id=$1 AND status_id=2", coords.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dbGetUsers(db *sql.DB) ([]string, error) {
	userids := make([]string, 0, 10)

	rows, err := db.Query("SELECT fid from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		id := ""
		rows.Scan(&id)
		userids = append(userids, id)
	}

	return userids, nil
}

func (u *user_data) dbGetUserData(db *sql.DB, fid string) error {
	err := db.QueryRow("SELECT name, surname FROM users WHERE fid=$1", fid).Scan(&u.Name, &u.Surname)

	return err
}
