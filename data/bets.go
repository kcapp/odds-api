package data

import (
	"errors"
	"fmt"
	"github.com/kcapp/odds-api/models"
)

func GetUserTournamentGamesBets(userId, tournamentId int) ([]*models.BetMatch, error) {
	rows, err := models.DB.Query(`
		SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2
		FROM bets_matches bm
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

func AddBet(bet models.BetMatch) error {
	var s string
	if bet.ID == nil {
		s = `INSERT INTO bets_matches (id, user_id, match_id, tournament_id, bet1, betx, bet2, odds1, oddsx, odds2) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	} else {
		s = `REPLACE INTO bets_matches (id, user_id, match_id, tournament_id, bet1, betx, bet2, odds1, oddsx, odds2)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	err := ValidateInput(bet)
	if err != nil {
		return err
	}

	tx, err := models.DB.Begin()
	if err != nil {
		return errors.New("error creating transaction")
	}

	res, err := tx.Exec(s, bet.ID, bet.UserId, bet.MatchId, bet.TournamentId, bet.Bet1, bet.BetX, bet.Bet2, bet.Odds1, bet.OddsX, bet.Odds2)
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

	// log.Printf("Created / updated bet (%d)", betId)
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

func GetUserActiveBets(bm models.BetMatch) (*models.UserActiveBets, error) {
	var err error
	uab := new(models.UserActiveBets)
	if bm.ID != nil {
		err = models.DB.QueryRow(`
		SELECT bm.user_id, bm.tournament_id, sum(bm.bet1) + sum(bm.bet2) as bets, uc.coins from
		bets_matches bm
		JOIN user_coins uc on bm.user_id = uc.user_id and bm.tournament_id = uc.tournament_id
		WHERE bm.user_id = ? AND bm.tournament_id = ? AND bm.id != ?
		and bm.outcome is null`, bm.UserId, bm.TournamentId, bm.ID).
			Scan(&uab.UserId, &uab.TournamentId, &uab.Bets, &uab.Coins)
	} else {
		err = models.DB.QueryRow(`
		SELECT bm.user_id, bm.tournament_id, sum(bm.bet1) + sum(bm.bet2) as bets, uc.coins from
		bets_matches bm
		JOIN user_coins uc on bm.user_id = uc.user_id and bm.tournament_id = uc.tournament_id
		WHERE bm.user_id = ? AND bm.tournament_id = ?
		and bm.outcome is null`, bm.UserId, bm.TournamentId).
			Scan(&uab.UserId, &uab.TournamentId, &uab.Bets, &uab.Coins)
	}

	if err != nil {
		return nil, err
	}

	return uab, nil
}

func ValidateInput(bm models.BetMatch) error {
	uab, err := GetUserActiveBets(bm)

	totalBets := uab.Bets + bm.Bet1 + bm.BetX + bm.Bet2
	if err != nil {
		return errors.New("can't fetch coin data")
	}

	fmt.Println(totalBets, uab.Coins)

	if totalBets > uab.Coins {
		return errors.New("not enough coins")
	}

	return nil
}
