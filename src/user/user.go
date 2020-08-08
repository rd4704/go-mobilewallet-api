package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Module user module
type Module struct {
	Router *mux.Router
	DB     *sql.DB
}

// GetUsers godoc
// @Summary Get the list of all users and their wallets
// @Description Get the list of all users and their wallets
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Router /api/users [get]
func (a *Module) getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting users")

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	users, err := GetUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary Get the a user's details and wallets information
// @Description Get the a user's details and wallets information
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Router /api/user/{id} [get]
// @Security ApiKeyAuth
// @param X-Session-Token header string true "X-Session-Token"
// @param id path int true "id"
func (a *Module) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := User{ID: id}
	if err := u.GetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
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

func (a *Module) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
}

// New init user module
func New(db *sql.DB, router *mux.Router) *Module {
	m := &Module{DB: db, Router: router}
	m.initializeRoutes()
	return m
}
