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

func GetUserLeaderboards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, err := strconv.Atoi(params["tournamentId"])
	if err != nil {
		log.Println("Invalid tournament id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	pid, err := strconv.Atoi(params["userId"])
	if err != nil {
		log.Println("Invalid user id")
		http.Error(writer, "Invalid user id", http.StatusBadRequest)
		return
	}

	SetHeaders(writer)
	lb, err := data.GetUserGameLeaderboards(tid, pid)
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.LeaderBoard))
		return
	} else if err != nil {
		log.Println("Unable to get leaderboards", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(lb)
}
