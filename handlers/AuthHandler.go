package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kcapp/odds-api/data"
	"github.com/kcapp/odds-api/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func ChangePass(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var authdetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		JSONError(w, "Username does not exist", http.StatusInternalServerError)
		return
	}

	var authuser *models.User
	authuser, err = data.GetUserByLogin(authdetails.Login)
	if authuser == nil {
		JSONError(w, "Username does not exist", http.StatusInternalServerError)
		return
	}
	if authuser.Login == "" {
		JSONError(w, "Username does not exist", http.StatusInternalServerError)
		return
	}

	_, err = data.ChangePassword(authdetails)
	if err != nil {
		JSONError(w, "Password change error", http.StatusInternalServerError)
		return
	}

	JSONSuccess(w, "Password changed", http.StatusOK)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var authdetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		JSONError(w, "User not found", http.StatusInternalServerError)
		return
	}

	var authuser *models.User
	authuser, err = data.GetUserByLogin(authdetails.Login)
	if authuser == nil {
		JSONError(w, "User not found", http.StatusInternalServerError)
		return
	}
	if authuser.Login == "" {
		JSONError(w, "User not found", http.StatusInternalServerError)
		return
	}

	ds, err := base64.StdEncoding.DecodeString(authdetails.Password)
	check := CheckPasswordHash(string(ds), authuser.Password)

	if !check {
		JSONError(w, "Username or password incorrect", http.StatusInternalServerError)
		return
	}

	validToken, err := GenerateJWT(authuser.Login)
	if err != nil {
		JSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	var response models.Response
	response.UserId = authuser.ID
	response.Login = authuser.Login
	response.TokenString = validToken
	response.RequiresChange = authuser.RequiresChange
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
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
