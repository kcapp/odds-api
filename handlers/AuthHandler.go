package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kcapp/odds-api/data"
	"github.com/kcapp/odds-api/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var authdetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		http.Error(w, "Error in request body", http.StatusBadRequest)
		return
	}

	var authuser *models.User
	authuser, err = data.GetUser(authdetails.Login)
	if authuser == nil {
		http.Error(w, "Username does not exist", http.StatusInternalServerError)
		return
	}
	if authuser.Login == "" {
		http.Error(w, "Username or password is incorrect", http.StatusNotFound)
		return
	}
	check := CheckPasswordHash(authdetails.Password, authuser.Password)

	if !check {
		http.Error(w, "Username or password is incorrect", http.StatusNotFound)
		return
	}

	validToken, err := GenerateJWT(authuser.Login)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	var token models.Token
	token.Login = authuser.Login
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		return
	}
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(login string) (string, error) {
	var mySigningKey = []byte("secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["login"] = login
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
