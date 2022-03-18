package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kcapp/odds-api/data"
	"log"
	"net/http"
	"strconv"
)

func GetUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := data.GetUser(id)
	if err != nil {
		log.Println("Unable to get user", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(user)
}
