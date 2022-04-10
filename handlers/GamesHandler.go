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

func StartGame(writer http.ResponseWriter, reader *http.Request) {
	params := mux.Vars(reader)
	gameId, err := strconv.Atoi(params["gameId"])

	if err != nil {
		log.Println("Unable to get match id json", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	gm := models.GameMetadata{
		MatchId: gameId,
		BetsOff: 1,
	}

	lid, err := data.StartGame(gm)
	if err != nil {
		log.Println("Unable to start game", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(writer)
	err = json.NewEncoder(writer).Encode(lid)
	if err != nil {
		return
	}
}

func FinishGame(writer http.ResponseWriter, reader *http.Request) {
	var gf models.GameFinish
	err := json.NewDecoder(reader.Body).Decode(&gf)

	if err != nil {
		log.Println("Unable to parse game finish", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	lid, err := data.FinishGame(gf)

	if err != nil {
		log.Println("Unable to finish game", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	SetHeaders(writer)
	err = json.NewEncoder(writer).Encode(lid)
	if err != nil {
		return
	}
}

func GetGamesMetadata(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetGamesMetadata()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.GameMetadata))
		return
	} else if err != nil {
		log.Println("Unable to get metadata", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(md)
}
