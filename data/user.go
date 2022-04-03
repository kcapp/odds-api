package data

import (
	"github.com/kcapp/odds-api/models"
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
		SELECT uc.id, uc.user_id, uc.tournament_id, uc.coins, uc.tournament_coins
		FROM user_coins uc
		WHERE uc.user_id = ? AND uc.tournament_id = ?`, userId, tournamentId).
		Scan(&u.ID, &u.UserId, &u.TournamentId, &u.Coins, &u.TournamentCoins)
	if err != nil {
		return nil, err
	}

	return u, nil
}
