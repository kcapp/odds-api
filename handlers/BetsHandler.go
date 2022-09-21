package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kcapp/odds-api/data"
	"github.com/kcapp/odds-api/models"
)

func GetUserGamesBets(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	uid, err := strconv.Atoi(params["userId"])
	if err != nil {
		log.Println("Invalid user")
		http.Error(writer, "Invalid user", http.StatusBadRequest)
		return
	}

	SetHeaders(writer)
	bets, err := data.GetUserGamesBets(uid)
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

func GetUserTournamentGamesBets(writer http.ResponseWriter, request *http.Request) {
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

func GetUserTournamentTournamentBets(writer http.ResponseWriter, request *http.Request) {
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
	bets, err := data.GetUserTournamentTournamentsBets(uid, tid)
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.BetTournament))
		return
	} else if err != nil {
		log.Println("Unable to get bets", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(bets)
}

func GetGameBets(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	gid, err := strconv.Atoi(params["gameId"])
	if err != nil {
		log.Println("Invalid game id")
		http.Error(writer, "Invalid game id", http.StatusBadRequest)
		return
	}

	SetHeaders(writer)
	bets, err := data.GetGameBets(gid)
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

func validateToken(userId int, token string) bool {
	user, err := data.GetUserById(userId)
	if err != nil {
		log.Fatal("error:", err)
		return false
	}

	split := strings.Split(token, ".")

	d, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		log.Fatal("error:", err)
		return false
	}

	var td models.TokenData
	err = json.Unmarshal(d, &td)

	if td.Login == user.Login {
		return true
	}

	return false
}

// AddBet AddVisit will add the visit to the database
func AddBet(writer http.ResponseWriter, reader *http.Request) {
	var bet models.BetMatch
	err := json.NewDecoder(reader.Body).Decode(&bet)
	if err != nil {
		log.Println("Unable to deserialize bet json", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	md, err := data.GetGameMetadata(bet.MatchId)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Unable to get game metadata", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if md != nil && md.BetsOff == 1 {
		log.Printf("Bets are off for this game")
		http.Error(writer, "Bets are off for this game", http.StatusBadRequest)
		return
	}

	//v := validateToken(bet.UserId, bet.Token)
	//if !v {
	//	log.Println("Invalid token.", err)
	//	http.Error(writer, err.Error(), http.StatusForbidden)
	//	return
	//}

	// todo - check if the bet was not changed
	// get current odds from db
	// get new odds from kcapp api

	if bet.Bet1 < 0 || bet.Bet2 < 0 || bet.BetX < 0 {
		// Someone is doing something funny, don't allow it
		log.Printf("Unable to add bet with negative value for user %d", bet.UserId)
		http.Error(writer, "Don't do that", http.StatusBadRequest)
		return
	}

	lid, err := data.AddBet(bet)
	if err != nil {
		log.Println("Unable to add bet", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Added bet %d for player %d", lid, bet.UserId)
	SetHeaders(writer)
	json.NewEncoder(writer).Encode(lid)
}

// AddTournamentBet will add the visit to the database
func AddTournamentBet(writer http.ResponseWriter, reader *http.Request) {
	var bet models.BetOutcome
	err := json.NewDecoder(reader.Body).Decode(&bet)
	if err != nil {
		log.Println("Unable to deserialize outcome bet json", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if bet.Bet1 < 0 || bet.Bet2 < 0 || bet.BetX < 0 {
		// Someone is doing something funny, don't allow it
		log.Printf("Unable to add bet with negative value for user %d", bet.UserId)
		http.Error(writer, "Don't do that", http.StatusBadRequest)
		return
	}

	lid, err := data.AddTournamentBet(bet)
	if err != nil {
		log.Println("Unable to add tournament bet", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Added tournament bet %d for player %d", lid, bet.UserId)

	SetHeaders(writer)
	json.NewEncoder(writer).Encode(lid)
}
