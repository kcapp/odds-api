package models

type TournamentMetadata struct {
	TournamentId int `json:"tournament_id"`
	BetsOff      int `json:"bets_off"`
}

type TopGameResult struct {
	UserId  int     `json:"user_id"`
	MatchId int     `json:"match_id"`
	Amount  float64 `json:"amount"`
}

type TournamentStatistics struct {
	TournamentId        int             `json:"tournament_id"`
	TournamentBetsCount int             `json:"tournament_bets_count"`
	TournamentUserCount int             `json:"tournament_user_count"`
	BankWinnings        float64         `json:"bank_winnings"`
	UserWinnings        float64         `json:"user_winnings"`
	BiggestWins         []TopGameResult `json:"biggest_wins"`
	BiggestLosses       []TopGameResult `json:"biggest_losses"`
}
