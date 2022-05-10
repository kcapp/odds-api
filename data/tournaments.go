package data

import "github.com/kcapp/odds-api/models"

func StartTournament(gm models.TournamentMetadata) (int64, error) {
	var s string
	var err error

	s = `REPLACE INTO tournaments_metadata (tournament_id, bets_off) VALUES (?, ?)`
	args := make([]interface{}, 0)
	args = append(args, gm.TournamentId, 1)

	lid, err := RunTransaction(s, args...)

	return lid, err
}
