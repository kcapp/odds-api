package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kcapp/odds-api/data"
	"github.com/kcapp/odds-api/models"
	"log"
	"net/http"
	"strconv"
)

func GetUserTournamentsGamesBets(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	uid, err := strconv.Atoi(params["userId"])
	if err != nil {
		log.Println("Invalid user")
		http.Error(writer, "Invalid user", http.StatusBadRequest)
		return
	}

	tid, err := strconv.Atoi(params["tournamentId"])
	if err != nil {
		log.Println("Unable to get tournament id", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	bets, err := data.GetUserTournamentGamesBets(uid, tid)
	if err != nil {
		match := models.BetMatch{}
		json.NewEncoder(writer).Encode(match)
		return
	}

	json.NewEncoder(writer).Encode(bets)
}
