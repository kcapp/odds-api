package models

type BetMatch struct {
	ID           int `json:"id,omitempty"`
	UserId       int `json:"user_id,omitempty"`
	TournamentId int `json:"tournament_id,omitempty"`
	MatchId      int `json:"match_id,omitempty"`
	Bet1         int `json:"bet_1,omitempty"`
	BetX         int `json:"bet_x,omitempty"`
	Bet2         int `json:"bet_2,omitempty"`
}
