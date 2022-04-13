package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kcapp/odds-api/data"
	"log"
	"net/http"
	"strconv"
)

func GetUserByLogin(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	login := params["login"]
	if login == "" {
		log.Println("Invalid login")
		http.Error(writer, "Invalid login", http.StatusBadRequest)
		return
	}
	user, err := data.GetUserByLogin(login)
	if err != nil {
		log.Println("Unable to get user", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(user)
}

func GetUserTournamentBalance(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId, err1 := strconv.Atoi(params["userId"])
	tournamentId, err2 := strconv.Atoi(params["tournamentId"])
	if err1 != nil || err2 != nil {
		log.Println("Invalid data")
		http.Error(writer, "Invalid data", http.StatusBadRequest)
		return
	}
	user, err := data.GetUserTournamentBalance(userId, tournamentId, 0)
	if err != nil {
		log.Println("Unable to get user balance data", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	SetHeaders(writer)
	json.NewEncoder(writer).Encode(user)
}
