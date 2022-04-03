package models

type BetMatch struct {
	ID           *int `json:"id"`
	UserId       int  `json:"user_id"`
	TournamentId int  `json:"tournament_id"`
	MatchId      int  `json:"match_id"`
	Bet1         int  `json:"bet_1"`
	BetX         int  `json:"bet_x,omitempty"`
	Bet2         int  `json:"bet_2"`
}

type UserTournamentBalance struct {
	ID              int `json:"id"`
	UserId          int `json:"user_id"`
	TournamentId    int `json:"tournament_id"`
	Coins           int `json:"coins"`
	TournamentCoins int `json:"tournament_coins"`
}
