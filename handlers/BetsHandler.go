package handlers

import (
	"database/sql"
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

	SetHeaders(writer)
	bets, err := data.GetUserTournamentGamesBets(uid, tid)
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.BetMatch))
		return
	} else if err != nil {
		log.Println("Unable to get bets", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(bets)
}

// AddVisit will add the visit to the database
func AddBet(writer http.ResponseWriter, reader *http.Request) {
	var bet models.BetMatch
	err := json.NewDecoder(reader.Body).Decode(&bet)
	if err != nil {
		log.Println("Unable to deserialize bet json", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	betId, err := data.AddBet(bet)
	if err != nil {
		log.Println("Unable to add bet", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(betId)
}
