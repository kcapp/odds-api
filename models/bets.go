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
	Outcome      *int    `json:"outcome"`
	BetsOff      int     `json:"bets_off"`
}

type UserTournamentBalance struct {
	UserId                int     `json:"user_id"`
	FirstName             string  `json:"first_name"`
	LastName              string  `json:"last_name"`
	TournamentId          int     `json:"tournament_id"`
	BetsPlaced            int     `json:"bets_placed"`
	CoinsBetsOpen         float32 `json:"coins_bets_open"`
	CoinsBetsClosed       float32 `json:"coins_bets_closed"`
	CoinsWon              float32 `json:"coins_won"`
	TournamentCoinsOpen   float32 `json:"tournament_coins_open"`
	TournamentCoinsClosed float32 `json:"tournament_coins_closed"`
	StartCoins            float32 `json:"start_coins"`
	CoinsAvailable        float32 `json:"coins_available"`
}

type SortBalanceByCoinsWon []*UserTournamentBalance
type SortBalanceByCoinsAvailable []*UserTournamentBalance

func (a SortBalanceByCoinsWon) Len() int           { return len(a) }
func (a SortBalanceByCoinsWon) Less(i, j int) bool { return a[i].CoinsWon > a[j].CoinsWon }
func (a SortBalanceByCoinsWon) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a SortBalanceByCoinsAvailable) Len() int { return len(a) }
func (a SortBalanceByCoinsAvailable) Less(i, j int) bool {
	return a[i].CoinsAvailable > a[j].CoinsAvailable
}
func (a SortBalanceByCoinsAvailable) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type UserActiveBets struct {
	UserId          int     `json:"user_id"`
	TournamentId    int     `json:"tournament_id"`
	BetsTotal       float32 `json:"bets"`
	AvailableCoins  float32 `json:"coins"`
	CurrentSavedBet int     `json:"current_saved_bet,omitempty"`
}

type CoinBalance struct {
	Coins float32 `json:"coins"`
}
