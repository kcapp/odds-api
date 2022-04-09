package models

type BetMatch struct {
	ID           int     `json:"id"`
	UserId       int     `json:"user_id"`
	TournamentId int     `json:"tournament_id"`
	MatchId      int     `json:"match_id"`
	Player1      int     `json:"player_1"`
	Player2      int     `json:"player_2"`
	Bet1         int     `json:"bet_1"`
	BetX         int     `json:"bet_x"`
	Bet2         int     `json:"bet_2"`
	Odds1        float64 `json:"odds_1"`
	OddsX        float64 `json:"odds_x"`
	Odds2        float64 `json:"odds_2"`
	Outcome      int     `json:"outcome,omitempty"`
}

type UserTournamentBalance struct {
	ID              int     `json:"id"`
	UserId          int     `json:"user_id"`
	TournamentId    int     `json:"tournament_id"`
	Coins           float64 `json:"coins"`
	TournamentCoins float64 `json:"tournament_coins"`
}

type UserActiveBets struct {
	UserId          int `json:"user_id"`
	TournamentId    int `json:"tournament_id"`
	BetsTotal       int `json:"bets"`
	AvailableCoins  int `json:"coins"`
	CurrentSavedBet int `json:"current_saved_bet,omitempty"`
}
