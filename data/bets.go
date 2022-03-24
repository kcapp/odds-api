package data

import "github.com/kcapp/odds-api/models"

func GetUserTournamentGamesBets(userId, tournamentId int) (*models.BetMatch, error) {
	b := new(models.BetMatch)
	err := models.DB.QueryRow(`
		SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2
		FROM bets_matches bm
		WHERE bm.user_id = ? and bm.tournament_id = ?`, userId, tournamentId).
		Scan(&b.ID, &b.UserId, &b.MatchId, &b.TournamentId, &b.Bet1, &b.BetX, &b.Bet2)
	if err != nil {
		return nil, err
	}

	return b, nil
}
