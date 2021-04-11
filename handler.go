package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// hold application with db and router
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// connect to db
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.handleRequests()
}

// start application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

/*

user
GET /bihack/rest/users/[firebase_id_or_whatever]/ -- returns overall score and in progress deliveries even if user is new

items
GET /bihack/rest/item/[bar-code]/ -- add new item (bar-code and type) or return its known type

history
POST /bihack/rest/history/ -- add new record to history

completeness
POST /bihack/rest/complete/ -- send start date and turn complete to all earlier

residentials
POST /bihack/rest/new_residential/ -- add new residence to db

gps
POST /bihack/rest/coords/ -- check if residential is close return 1 and upsate if yes and 0 otherwise
*/

func (a *App) handleRequests() {
	a.Router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			(w).Header().Set("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		})
	})
	a.Router.HandleFunc("/bihack/rest/users/{firebase_id_or_whatever/user}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/bihack/rest/item/{bar-code/item}", a.getItem).Methods("GET")
	a.Router.HandleFunc("/bihack/rest/history/", a.addHistoryRecord).Methods("POST")
	//a.Router.HandleFunc("/bihack/rest/complete/", a.setComplete).Methods("POST")
	a.Router.HandleFunc("/bihack/rest/new_residential/", a.addResidence).Methods("POST")
	a.Router.HandleFunc("/bihack/rest/coords/{COORD_STRING}", a.getGPSCoords).Methods("GET")
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["user"]

	var u user
	if err := u.dbGetUserScore(a.DB, uid); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	icode := vars["item"]
	i := item{code: icode}

	// check if empty
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := i.dbGetItemByCode(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, i)
}

func (a *App) addHistoryRecord(w http.ResponseWriter, r *http.Request) {
	var rec record

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rec); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := rec.dbCreateRecord(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, 0)
}

func (a *App) setComplete(w http.ResponseWriter, r *http.Request) {

}

func (a *App) addResidence(w http.ResponseWriter, r *http.Request) {
	var res Residence

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&res); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := res.dbAddNewResidence(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, 0)
}

func (a *App) getGPSCoords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coords := vars["coordinate"]

	if err := dbCompareCoords(a.DB, coords); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
