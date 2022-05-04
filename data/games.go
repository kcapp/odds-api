package data

import (
	"errors"
	"github.com/kcapp/odds-api/models"
)

func StartGame(gm models.GameMetadata) (int64, error) {
	var s string
	var err error

	s = `REPLACE INTO games_metadata (match_id, bets_off) VALUES (?, ?)`
	args := make([]interface{}, 0)
	args = append(args, gm.MatchId, 1)

	lid, err := RunTransaction(s, args...)

	return lid, err
}

func FinishGame(gf models.GameFinish) (int64, error) {
	var s string
	var err error

	// set outcome of the game
	s = `UPDATE bets_games bg SET bg.outcome = ? WHERE bg.match_id = ?`

	args := make([]interface{}, 0)
	args = append(args, gf.WinnerId, gf.MatchId)
	lid, err := RunTransaction(s, args...)

	if err != nil {
		return 0, errors.New("error setting outcome")
	}

	if err != nil {
		return 0, errors.New("error updating coins")
	}

	return lid, err
}

func GetGamesMetadata() ([]*models.GameMetadata, error) {
	rows, err := models.DB.Query(`
			SELECT
			gm.match_id, gm.bets_off
			FROM games_metadata gm`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]*models.GameMetadata, 0)
	for rows.Next() {
		g := new(models.GameMetadata)
		err := rows.Scan(&g.MatchId, &g.BetsOff)
		if err != nil {
			return nil, err
		}

		bets = append(bets, g)
	}

	if err != nil {
		return nil, err
	}

	return bets, nil
}
