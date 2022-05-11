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

func GetTournamentRanking(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	tid, err := strconv.Atoi(params["tournamentId"])
	if err != nil {
		log.Println("Unable to get tournament id", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(writer)
	ranking, err := data.GetTournamentRanking(tid)
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.BetMatch))
		return
	} else if err != nil {
		log.Println("Unable to get bets", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(ranking)
}

func GetTournamentGameRanking(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	tid, err := strconv.Atoi(params["tournamentId"])
	if err != nil {
		log.Println("Unable to get tournament id", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(writer)
	ranking, err := data.GetTournamentGameRanking(tid)
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.BetMatch))
		return
	} else if err != nil {
		log.Println("Unable to get bets", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(ranking)
}

func StartTournament(writer http.ResponseWriter, reader *http.Request) {
	params := mux.Vars(reader)
	tournamentId, err := strconv.Atoi(params["tournamentId"])

	if err != nil {
		log.Println("Unable to get tournament id json", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	gm := models.TournamentMetadata{
		TournamentId: tournamentId,
		BetsOff:      1,
	}

	lid, err := data.StartTournament(gm)
	if err != nil {
		log.Println("Unable to start tournament", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(writer)
	err = json.NewEncoder(writer).Encode(lid)
	if err != nil {
		return
	}
}

func GetTournamentOutcomes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	tid, err := strconv.Atoi(params["tournamentId"])
	if err != nil {
		log.Println("Unable to get tournament id", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(writer)
	ranking, err := data.GetTournamentOutcomes(tid)
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.BetMatch))
		return
	} else if err != nil {
		log.Println("Unable to get outcomes", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(ranking)
}
