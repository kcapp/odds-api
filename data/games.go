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

	// mark game metadata as bets off
	s = `INSERT INTO odds.games_metadata (match_id, bets_off) VALUES (?, 1);`

	args = make([]interface{}, 0)
	args = append(args, gf.MatchId)
	_, err = RunTransaction(s, args...)

	if err != nil {
		return 0, errors.New("error setting game finished")
	}

	return lid, err
}

func GetGameMetadata(id int) (*models.GameMetadata, error) {
	md := new(models.GameMetadata)
	err := models.DB.QueryRow(`
			SELECT
			gm.match_id, gm.bets_off
			FROM games_metadata gm
			WHERE match_id = ?`, id).Scan(&md.MatchId, &md.BetsOff)

	if err != nil {
		return nil, err
	}

	return md, nil
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
