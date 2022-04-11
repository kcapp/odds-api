package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kcapp/odds-api/handlers"
	"github.com/kcapp/odds-api/models"
	"log"
	"net/http"
	"os"
)

func main() {
	var configFileParam string

	if len(os.Args) > 1 {
		configFileParam = os.Args[1]
	}

	config, err := models.GetConfig(configFileParam)
	if err != nil {
		panic(err)
	}
	models.InitDB(config.GetMysqlConnectionString())

	router := mux.NewRouter()
	// User
	router.HandleFunc("/user/login", handlers.SignIn).Methods("POST", "OPTIONS")
	//router.HandleFunc("/user/changepass", handlers.ChangePass).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/{userId}/tournament/{tournamentId}/balance", handlers.GetUserTournamentBalance).Methods("GET")
	// Bets
	router.HandleFunc("/user/{userId}/tournament/{tournamentId}/bets", handlers.GetUserTournamentsGamesBets).Methods("GET")
	router.HandleFunc("/bets/{gameId}", handlers.AddBet).Methods("POST", "OPTIONS")

	router.HandleFunc("/games/{gameId}/start", handlers.StartGame).Methods("POST", "OPTIONS")
	router.HandleFunc("/games/{gameId}/finish", handlers.FinishGame).Methods("POST", "OPTIONS")
	router.HandleFunc("/games/meta", handlers.GetGamesMetadata).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.APIConfig.Port), router))
}
