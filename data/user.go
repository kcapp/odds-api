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
	utb, _ := GetUserTournamentRanking(tournamentId, userId)

	uabt, _ := GetUserTournamentTournamentCoinsOpen(userId, tournamentId, skipGameId)
	ucat, _ := GetUserTournamentTournamentCoinsClosed(userId, tournamentId)
	ucwt, _ := GetUserTournamentTournamentCoinsWon(userId, tournamentId)

	utb.TournamentCoinsOpen = uabt.Coins
	utb.TournamentCoinsClosed = ucat.Coins
	utb.TournamentCoinsWon = ucwt.Coins

	return utb, nil
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
