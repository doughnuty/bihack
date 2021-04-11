package main

import (
	"database/sql"

	geo "github.com/kellydunn/golang-geo"
)

type user struct {
	Score       int            `json:"score"`
	In_progress map[string]int `json:"in_progress"`
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

	// 	err := db.QueryRow("SELECT id FROM types WHERE name=$1", i.Class).Scan(&id)
	// 	if err != nil {
	// 		return err
	// 	}

	// } else {
	// 	var id int
	// 	_, err = db.Exec("INSERT INTO items (scan, type_id) $1, $2)", i.Code, id)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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
