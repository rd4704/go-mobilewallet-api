package transfer

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Module transfer module
type Module struct {
	Router *mux.Router
	DB     *sql.DB
}

// GetUser godoc
// @Summary Get a transfer information by id
// @Description Get a transfer information by id
// @Tags transfer
// @Accept  json
// @Produce  json
// @Success 200 {array} Transfer
// @Router /api/transfer/{id} [get]
// @Security ApiKeyAuth
// @param X-Session-Token header string true "X-Session-Token"
// @param id path int true "id"
func (a *Module) getTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transfer ID")
		return
	}

	t := Transfer{ID: id}
	if err := t.GetTransfer(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Transfer not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

// GetUser godoc
// @Summary Get all transfers history of a user
// @Description Get all transfers history of a user
// @Tags transfer
// @Accept  json
// @Produce  json
// @Success 200 {array} Transfer
// @Router /api/user-transfers/{userId} [get]
// @Security ApiKeyAuth
// @param X-Session-Token header string true "X-Session-Token"
// @param userId path int true "userId"
func (a *Module) getUserTransfers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["userId"])
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	transfers, err := getUserTransfers(a.DB, userId, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, transfers)
}

// GetUser godoc
// @Summary Makes a transfer to a specified walletId
// @Description Makes a transfer to a specified walletId
// @Tags transfer
// @Accept  json
// @Produce  json
// @Success 200 {array} Transfer
// @Router /api/make-transfer [post]
// @Security ApiKeyAuth
// @param X-Session-Token header string true "X-Session-Token"
// @param body body Transfer true "body"
func (a *Module) makeTransfer(w http.ResponseWriter, r *http.Request) {
	var t Transfer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := t.MakeTransfer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, t)
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
	a.Router.HandleFunc("/transfer/{id:[0-9]+}", a.getTransfer).Methods("GET")
	a.Router.HandleFunc("/user-transfers/{userId:[0-9]+}", a.getUserTransfers).Methods("GET")
	a.Router.HandleFunc("/make-transfer", a.makeTransfer).Methods("POST")
}

// New init transfer module
func New(db *sql.DB, router *mux.Router) *Module {
	m := &Module{DB: db, Router: router}
	m.initializeRoutes()
	return m
}
