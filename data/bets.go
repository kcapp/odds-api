package data

import (
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
		//b := new(models.BetMatch)
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
	if bet.ID == nil {
		s = `INSERT INTO bets_matches (id, user_id, match_id, tournament_id, bet1, betx, bet2) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	} else {
		s = `REPLACE INTO bets_matches (id, user_id, match_id, tournament_id, bet1, betx, bet2)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	}

	tx, err := models.DB.Begin()
	if err != nil {
		return 0, err
	}

	res, err := tx.Exec(s, bet.ID, bet.UserId, bet.MatchId, bet.TournamentId, bet.Bet1, bet.BetX, bet.Bet2)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	betId, err := res.LastInsertId()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	// log.Printf("Created / updated bet (%d)", betId)
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return betId, err
}
