package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kcapp/odds-api/data"
	"log"
	"net/http"
)

func GetUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	login := params["login"]
	if login == "" {
		log.Println("Invalid login")
		http.Error(writer, "Invalid login", http.StatusBadRequest)
		return
	}
	user, err := data.GetUser(login)
	if err != nil {
		log.Println("Unable to get user", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(user)
}
