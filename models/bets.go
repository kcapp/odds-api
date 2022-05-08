package models

type Market struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TournamentOutcome struct {
	ID             int     `json:"id"`
	TournamentId   int     `json:"tournament_id"`
	MarketId       int     `json:"market_id"`
	MarketName     string  `json:"market_name"`
	MarketTypeId   int     `json:"market_type_id"`
	MarketTypeName string  `json:"market_type_name"`
	OutcomeValue   float64 `json:"outcome_value"`
	Odds1          float64 `json:"odds_1"`
	Odds2          float64 `json:"odds_2"`
	OddsX          float64 `json:"odds_x"`
	PlayerName     *string `json:"player_name"`
}

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

type BetTournament struct {
	ID           int      `json:"id"`
	UserId       int      `json:"user_id"`
	TournamentId int      `json:"tournament_id"`
	OutcomeId    int      `json:"outcome_id"`
	Bet1         int      `json:"bet_1"`
	BetX         int      `json:"bet_x"`
	Bet2         int      `json:"bet_2"`
	Outcome      *float64 `json:"outcome"` // final outcome
	OutcomeValue float64  `json:"outcome_value"`
	Odds1        float64  `json:"odds_1"`
	OddsX        float64  `json:"odds_x"`
	Odds2        float64  `json:"odds_2"`
	MarketId     int      `json:"market_id"`
	MarketTypeId int      `json:"market_type_id"`
}

type BetOutcome struct {
	ID           int `json:"id"`
	UserId       int `json:"user_id"`
	TournamentId int `json:"tournament_id"`
	OutcomeId    int `json:"outcome_id"`
	Bet1         int `json:"bet_1"`
	BetX         int `json:"bet_x"`
	Bet2         int `json:"bet_2"`
}

type UserTournamentBalance struct {
	UserId                int     `json:"user_id"`
	FirstName             string  `json:"first_name"`
	LastName              string  `json:"last_name"`
	TournamentId          int     `json:"tournament_id"`
	BetsPlaced            int     `json:"bets_placed"`
	BetsClosed            int     `json:"bets_closed"`
	CoinsBetsOpen         float32 `json:"coins_bets_open"`
	CoinsBetsClosed       float32 `json:"coins_bets_closed"`
	CoinsWon              float32 `json:"coins_won"`
	PotentialWinnings     float32 `json:"potential_winnings"`
	TournamentCoinsOpen   float32 `json:"tournament_coins_open"`
	TournamentCoinsClosed float32 `json:"tournament_coins_closed"`
	TournamentCoinsWon    float32 `json:"tournament_coins_won"`
	StartCoins            float32 `json:"start_coins"`
	CoinsAvailable        float32 `json:"coins_available"`
}

type SortBalanceByCoinsWon []*UserTournamentBalance
type SortBalanceByCoinsAvailable []*UserTournamentBalance
type SortBalanceByPotentialCoins []*UserTournamentBalance

func (a SortBalanceByPotentialCoins) Len() int { return len(a) }
func (a SortBalanceByPotentialCoins) Less(i, j int) bool {
	return (a[i].CoinsAvailable + a[i].CoinsBetsOpen + a[i].PotentialWinnings) >
		(a[j].CoinsAvailable + a[j].CoinsBetsOpen + a[j].PotentialWinnings)
}
func (a SortBalanceByPotentialCoins) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

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
