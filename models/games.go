package models

type GameMetadata struct {
	MatchId int `json:"match_id"`
	BetsOff int `json:"bets_off"`
}

type GameFinish struct {
	MatchId  int `json:"match_id"`
	WinnerId int `json:"winner_id"`
}
