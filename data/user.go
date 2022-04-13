package data

import (
	"encoding/base64"
	"github.com/kcapp/odds-api/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUserById(id int) (*models.User, error) {
	u := new(models.User)
	err := models.DB.QueryRow(`
		SELECT
			u.id, u.first_name, u.last_name, u.login, u.password, u.requires_change
		FROM users u
		WHERE u.id = ?`, id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Login, &u.Password, &u.RequiresChange)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserByLogin(login string) (*models.User, error) {
	u := new(models.User)
	err := models.DB.QueryRow(`
		SELECT
			u.id, u.first_name, u.last_name, u.login, u.password, u.requires_change
		FROM users u
		WHERE u.login = ?`, login).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Login, &u.Password, &u.RequiresChange)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserTournamentBalance(userId int, tournamentId int, skipGameId int) (*models.UserTournamentBalance, error) {
	u, _ := GetUserById(userId)
	uab, _ := GetUserTournamentCoinsOpen(userId, tournamentId, skipGameId)
	uca, _ := GetUserTournamentCoinsClosed(userId, tournamentId)
	ucw, _ := GetUserTournamentCoinsWon(userId, tournamentId)
	tb, _ := GetUserTournamentGamesBets(userId, tournamentId)

	var startCoins float32 = 1000

	utb := models.UserTournamentBalance{
		UserId:                u.ID,
		FirstName:             u.FirstName,
		LastName:              u.LastName,
		TournamentId:          tournamentId,
		BetsPlaced:            len(tb),
		CoinsBetsOpen:         uab.Coins,
		CoinsBetsClosed:       uca.Coins,
		CoinsWon:              ucw.Coins,
		TournamentCoinsOpen:   1000,
		TournamentCoinsClosed: 1000,
		StartCoins:            startCoins,
		CoinsAvailable:        startCoins - uab.Coins - uca.Coins + ucw.Coins,
	}

	return &utb, nil
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ChangePassword(ad models.Authentication) (int64, error) {
	var s string
	var err error

	ds, err := base64.StdEncoding.DecodeString(ad.Password)
	np, err := GenerateHashPassword(string(ds))

	// update password and change field
	s = `UPDATE users SET password = ?, requires_change = 0 WHERE login = ?`
	args := make([]interface{}, 0)
	args = append(args, np, ad.Login)
	lid, err := RunTransaction(s, args...)

	return lid, err
}
