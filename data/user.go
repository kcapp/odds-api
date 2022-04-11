package data

import (
	"encoding/base64"
	"github.com/kcapp/odds-api/models"
	"golang.org/x/crypto/bcrypt"
)

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

func GetUserTournamentBalance(userId int, tournamentId int) (*models.UserTournamentBalance, error) {
	u := new(models.UserTournamentBalance)
	err := models.DB.QueryRow(`
		SELECT uc.id, uc.user_id, u.first_name, u.last_name, uc.tournament_id, uc.coins, uc.tournament_coins
		FROM user_coins uc
		JOIN users u on uc.user_id = u.id
		WHERE uc.user_id = ? AND uc.tournament_id = ?`, userId, tournamentId).
		Scan(&u.ID, &u.UserId, &u.FirstName, &u.LastName, &u.TournamentId, &u.Coins, &u.TournamentCoins)
	if err != nil {
		return nil, err
	}

	return u, nil
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
