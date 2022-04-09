package data

import (
	"errors"
	"github.com/kcapp/odds-api/models"
)

func GetUserTournamentGamesBets(userId, tournamentId int) ([]*models.BetMatch, error) {
	rows, err := models.DB.Query(`
		SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2
		FROM bets_games bm
		WHERE bm.user_id = ? and bm.tournament_id = ?`, userId, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]*models.BetMatch, 0)
	for rows.Next() {
		b := new(models.BetMatch)
		err := rows.Scan(&b.ID, &b.UserId, &b.MatchId, &b.TournamentId, &b.Bet1, &b.BetX, &b.Bet2)
		if err != nil {
			return nil, err
		}

		bets = append(bets, b)
	}

	if err != nil {
		return nil, err
	}

	return bets, nil
}

func AddBet(bet models.BetMatch) (int64, error) {
	var s string
	var bm *models.BetMatch
	var err error

	if bet.ID == 0 {
		s = `INSERT INTO bets_games (id, user_id, match_id, tournament_id, player1, player2, bet1, betx, bet2, odds1, oddsx, odds2) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		// We are placing new bet, but we might be passing 0s only
		if bet.Bet1+bet.Bet2+bet.BetX == 0 {
			return 0, errors.New("can't place an empty bet")
		}
	} else {
		s = `REPLACE INTO bets_games (id, user_id, match_id, tournament_id, player1, player2, bet1, betx, bet2, odds1, oddsx, odds2)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		// This is an update of existing row,
		// the coins value are (coins now + row's value) - new bet values
		bm, err = GetUserBetById(bet.ID)
		if err != nil {
			return 0, err
		}
	}

	err = ValidateInput(bet)
	if err != nil {
		return 0, err
	}

	// We can have an existing bet with 0s passed, then we should remove the row from the db
	if bet.Bet1+bet.Bet2+bet.BetX == 0 {
		s = `DELETE FROM bets_games
				WHERE match_id = ? and user_id = ? AND outcome IS NULL
				AND match_id NOT IN (SELECT gm.match_id from games_metadata gm WHERE gm.bets_off = 1)`
		args := make([]interface{}, 0)
		args = append(args, bet.MatchId, bet.UserId)

		lid, err := RunTransaction(s, args...)

		return lid, err
	}

	tx, err := models.DB.Begin()
	if err != nil {
		return 0, errors.New("error creating transaction")
	}

	res, err := tx.Exec(s, bet.ID, bet.UserId, bet.MatchId, bet.TournamentId, bet.Player1, bet.Player2, bet.Bet1, bet.BetX, bet.Bet2, bet.Odds1, bet.OddsX, bet.Odds2)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	var newBetCoins int

	// if bet.ID == 0, this was a new bet, so we can safely substract the bet amount from coins
	// coins = coins - (bet1 + betx + bet2)
	if bet.ID == 0 {
		newBetCoins = bet.Bet1 + bet.BetX + bet.Bet2
		err := UpdateUserCoins(bet.UserId, bet.TournamentId, newBetCoins)
		if err != nil {
			return 0, err
		}
	} else {
		newBetCoins := -(bm.Bet1 + bm.BetX + bm.Bet2) + (bet.Bet1 + bet.BetX + bet.Bet2)
		err = UpdateUserCoins(bet.UserId, bet.TournamentId, newBetCoins)
		if err != nil {
			return 0, err
		}
	}

	return lid, err
}

func UpdateUserCoins(userId int, tournamentId int, bets int) error {
	s := `UPDATE user_coins uc SET uc.coins = uc.coins - ?
		  WHERE uc.user_id = ? AND uc.tournament_id = ?`

	tx, err := models.DB.Begin()
	if err != nil {
		return errors.New("error creating transaction")
	}

	res, err := tx.Exec(s, bets, userId, tournamentId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func CheckBetOff(matchId int) int {
	var bo int
	_ = models.DB.QueryRow(`
				select COALESCE(mg.bets_off, 0)
				from games_metadata mg
				where mg.match_id = ?`, matchId).Scan(&bo)

	return bo
}

func GetUserActiveBets(bm models.BetMatch) (*models.UserActiveBets, error) {
	var err error

	uab := new(models.UserActiveBets)
	if bm.ID != 0 {
		err = models.DB.QueryRow(`
						select uc.user_id, uc.tournament_id,
					   coalesce(sum(bm.bet1 + bm.betx + bm.bet2),0) as betsTotal,
					   uc.coins as availableCoins, coalesce(sum(bmc.bet1 + bmc.betx + bmc.bet2),0) currentSavedBet
				from user_coins uc
				left join bets_games bm on uc.user_id = bm.user_id and uc.tournament_id = bm.tournament_id
				left join bets_games bmc on uc.user_id = bmc.user_id and uc.tournament_id = bmc.tournament_id and bmc.id = ?
				and bm.id != ? and bm.outcome is null
				where uc.user_id = ? and uc.tournament_id = ?`, bm.ID, bm.ID, bm.UserId, bm.TournamentId).
			Scan(&uab.UserId, &uab.TournamentId, &uab.BetsTotal, &uab.AvailableCoins, &uab.CurrentSavedBet)
	} else {
		err = models.DB.QueryRow(`
				select uc.user_id, uc.tournament_id,
					   coalesce(sum(bm.bet1 + bm.betx + bm.bet2),0) as bets,
					   uc.coins
				from user_coins uc
				left join bets_games bm on uc.user_id = bm.user_id and uc.tournament_id = bm.tournament_id and bm.outcome is null
				where uc.user_id = ? and uc.tournament_id = ? 
		`, bm.UserId, bm.TournamentId).
			Scan(&uab.UserId, &uab.TournamentId, &uab.BetsTotal, &uab.AvailableCoins)
	}

	if err != nil {
		return nil, err
	}

	return uab, nil
}

func GetUserBetById(betId int) (*models.BetMatch, error) {
	var err error
	bm := new(models.BetMatch)
	{
		err = models.DB.QueryRow(`
		SELECT bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2 
		FROM bets_games bm
		WHERE bm.id = ?`, betId).
			Scan(&bm.ID, &bm.UserId, &bm.MatchId, &bm.TournamentId, &bm.Bet1, &bm.BetX, &bm.Bet2)
	}

	if err != nil {
		return nil, err
	}

	return bm, nil
}

func ValidateInput(bm models.BetMatch) error {
	betsOff := CheckBetOff(bm.MatchId)
	if betsOff == 1 {
		return errors.New("bet are off for this match")
	}

	uab, err := GetUserActiveBets(bm)
	if err != nil {
		return errors.New("can't fetch coin data")
	}

	var newBets int

	// how much can I bet?
	availableBet := (uab.BetsTotal + uab.AvailableCoins) - (uab.BetsTotal - uab.CurrentSavedBet)
	newBets = bm.Bet1 + bm.BetX + bm.Bet2

	// new bet amount can be greater than available amount
	if newBets > availableBet {
		return errors.New("not enough coins")
	}

	return nil
}
