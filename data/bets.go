package data

import (
	"errors"
	"github.com/kcapp/odds-api/models"
	"sort"
)

func GetUserTournamentCoinsOpen(userId, tournamentId, skipGameId int) (*models.CoinBalance, error) {
	var s string
	cb := new(models.CoinBalance)

	if skipGameId != 0 {
		s = `select COALESCE(sum(bg.bet1+bg.betx+bg.bet2), 0) as betCoins
			from bets_games bg
			where bg.user_id = ? and bg.tournament_id = ? and bg.outcome is null AND bg.id != ?`
		err := models.DB.QueryRow(s, userId, tournamentId, skipGameId).Scan(&cb.Coins)
		if err != nil {
			return nil, err
		}
	} else {
		s = `select COALESCE(sum(bg.bet1+bg.betx+bg.bet2), 0) as betCoins
			from bets_games bg
			where bg.user_id = ? and bg.tournament_id = ? and bg.outcome is null`
		err := models.DB.QueryRow(s, userId, tournamentId).Scan(&cb.Coins)
		if err != nil {
			return nil, err
		}
	}

	return cb, nil
}

func GetUserTournamentCoinsClosed(userId, tournamentId int) (*models.CoinBalance, error) {
	s := `select COALESCE(sum(bg.bet1+bg.betx+bg.bet2), 0) as betCoins
			from bets_games bg
			where bg.user_id = ? and bg.tournament_id = ? and bg.outcome is not null`

	cb := new(models.CoinBalance)
	err := models.DB.QueryRow(s, userId, tournamentId).
		Scan(&cb.Coins)
	if err != nil {
		return nil, err
	}

	return cb, nil
}

func GetUserTournamentCoinsWon(userId, tournamentId int) (*models.CoinBalance, error) {
	s := `select COALESCE(sum(ROUND((if(bgf.player1 = bgf.outcome, bet1 * odds1, 0) + 
                if(bgf.player2 = bgf.outcome, bet2 * odds2, 0)), 2)), 0) as coins
				from bets_games bgf
				where bgf.user_id = ? and bgf.tournament_id = ? and bgf.outcome is not null`

	cb := new(models.CoinBalance)
	err := models.DB.QueryRow(s, userId, tournamentId).
		Scan(&cb.Coins)
	if err != nil {
		return nil, err
	}

	return cb, nil
}

func GetTournamentRanking(tournamentId int) ([]*models.UserTournamentBalance, error) {
	rows, err := models.DB.Query(`
		select bgo.user_id, u.first_name, u.last_name, bgo.tournament_id, 
			(coalesce(count(bgo.user_id), 0) + coalesce(numBetsClosed, 0)) as numBets,
	   		sum(bgo.bet1+bgo.betx+bgo.bet2) as openBets,
	   		COALESCE(coins, 0) as closedBets,
	   		COALESCE(coinsWon, 0) as coinsWon,
		       1000, 1000, 1000
		from bets_games bgo
		left join (
			select user_id as uid, tournament_id as tid, sum(bet1+betx+bet2) as coins, count(user_id) as numBetsClosed
			from bets_games bgc where outcome is not null
			group by user_id
			) bgc on bgc.uid = bgo.user_id and bgc.tid = bgo.tournament_id
		left join (
			select user_id as uid, tournament_id as tid,
				   ROUND(SUM(if(bgf.player1 = bgf.outcome, bet1 * odds1, if(bgf.player2 = bgf.outcome, bet2 * odds2, 0))), 2) as coinsWon
			from bets_games bgf where outcome is not null
			group by user_id
			) bgf on bgf.uid = bgo.user_id and bgf.tid = bgo.tournament_id
		join users u on bgo.user_id = u.id
		where bgo.tournament_id = ? and bgo.outcome is null
		group by bgo.user_id`, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]*models.UserTournamentBalance, 0)
	for rows.Next() {
		b := new(models.UserTournamentBalance)
		err := rows.Scan(&b.UserId, &b.FirstName, &b.LastName, &b.TournamentId,
			&b.BetsPlaced, &b.CoinsBetsOpen, &b.CoinsBetsClosed, &b.CoinsWon,
			&b.TournamentCoinsOpen, &b.TournamentCoinsClosed, &b.StartCoins)
		b.CoinsAvailable = b.StartCoins - b.CoinsBetsOpen - b.CoinsBetsClosed + b.CoinsWon
		if err != nil {
			return nil, err
		}

		balances = append(balances, b)
	}

	sort.Sort(models.SortBalanceByCoinsAvailable(balances))

	if err != nil {
		return nil, err
	}

	return balances, nil
}

func GetUserTournamentGamesBets(userId, tournamentId int) ([]*models.BetMatch, error) {
	rows, err := models.DB.Query(`
			SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2, bm.outcome, 
			       COALESCE(gm.bets_off, 0) as bets_off
			FROM bets_games bm
			LEFT JOIN games_metadata gm on bm.match_id = gm.match_id
			WHERE bm.user_id = ? and bm.tournament_id = ?`, userId, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]*models.BetMatch, 0)
	for rows.Next() {
		b := new(models.BetMatch)
		err := rows.Scan(&b.ID, &b.UserId, &b.MatchId, &b.TournamentId, &b.Bet1, &b.BetX, &b.Bet2, &b.Outcome, &b.BetsOff)
		if err != nil {
			return nil, err
		}

		bets = append(bets, b)
	}

	if err != nil {
		return nil, err
	}

	return bets, nil
}

func AddBet(bet models.BetMatch) (int64, error) {
	var sq string
	var err error

	if bet.ID == 0 {
		sq = `INSERT INTO bets_games (id, user_id, match_id, tournament_id, player1, player2, bet1, betx, bet2, odds1, oddsx, odds2) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		// We are placing new bet, but we might be passing 0s only
		if bet.Bet1+bet.Bet2+bet.BetX == 0 {
			return 0, errors.New("can't place an empty bet")
		}
	} else {
		sq = `REPLACE INTO bets_games (id, user_id, match_id, tournament_id, player1, player2, bet1, betx, bet2, odds1, oddsx, odds2)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	err = ValidateInput(bet)
	if err != nil {
		return 0, err
	}

	// We can have an existing bet with 0s passed, then we should remove the row from the db
	if bet.Bet1+bet.Bet2+bet.BetX == 0 {
		s := `DELETE FROM bets_games
				WHERE match_id = ? and user_id = ? AND outcome IS NULL
				AND match_id NOT IN (SELECT gm.match_id from games_metadata gm WHERE gm.bets_off = 1)`
		args := make([]interface{}, 0)
		args = append(args, bet.MatchId, bet.UserId)

		lid, err := RunTransaction(s, args...)

		return lid, err
	}

	args := make([]interface{}, 0)
	args = append(args, bet.ID, bet.UserId, bet.MatchId, bet.TournamentId, bet.Player1, bet.Player2, bet.Bet1, bet.BetX, bet.Bet2, bet.Odds1, bet.OddsX, bet.Odds2)
	lid, err := RunTransaction(sq, args...)

	return lid, err
}

func CheckBetOff(matchId int) int {
	var bo int
	_ = models.DB.QueryRow(`
				select COALESCE(mg.bets_off, 0)
				from games_metadata mg
				where mg.match_id = ?`, matchId).Scan(&bo)

	return bo
}

func GetUserActiveBets(bm models.BetMatch) (*models.UserActiveBets, error) {
	var err error

	uab := new(models.UserActiveBets)
	if bm.ID != 0 {
		err = models.DB.QueryRow(`
						select uc.user_id, uc.tournament_id,
					   coalesce(sum(bm.bet1 + bm.betx + bm.bet2),0) as betsTotal,
					   uc.coins as availableCoins, coalesce(sum(bmc.bet1 + bmc.betx + bmc.bet2),0) currentSavedBet
				from user_coins uc
				left join bets_games bm on uc.user_id = bm.user_id and uc.tournament_id = bm.tournament_id
				left join bets_games bmc on uc.user_id = bmc.user_id and uc.tournament_id = bmc.tournament_id and bmc.id = ?
				and bm.id != ? and bm.outcome is null
				where uc.user_id = ? and uc.tournament_id = ?`, bm.ID, bm.ID, bm.UserId, bm.TournamentId).
			Scan(&uab.UserId, &uab.TournamentId, &uab.BetsTotal, &uab.AvailableCoins, &uab.CurrentSavedBet)
	} else {
		err = models.DB.QueryRow(`
				select uc.user_id, uc.tournament_id,
					   coalesce(sum(bm.bet1 + bm.betx + bm.bet2),0) as bets,
					   uc.coins
				from user_coins uc
				left join bets_games bm on uc.user_id = bm.user_id and uc.tournament_id = bm.tournament_id and bm.outcome is null
				where uc.user_id = ? and uc.tournament_id = ? 
		`, bm.UserId, bm.TournamentId).
			Scan(&uab.UserId, &uab.TournamentId, &uab.BetsTotal, &uab.AvailableCoins)
	}

	if err != nil {
		return nil, err
	}

	return uab, nil
}

func GetUserBetById(betId int) (*models.BetMatch, error) {
	var err error
	bm := new(models.BetMatch)
	{
		err = models.DB.QueryRow(`
		SELECT bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2 
		FROM bets_games bm
		WHERE bm.id = ?`, betId).
			Scan(&bm.ID, &bm.UserId, &bm.MatchId, &bm.TournamentId, &bm.Bet1, &bm.BetX, &bm.Bet2)
	}

	if err != nil {
		return nil, err
	}

	return bm, nil
}

func ValidateInput(bm models.BetMatch) error {
	betsOff := CheckBetOff(bm.MatchId)
	if betsOff == 1 {
		return errors.New("bet are off for this match")
	}

	utb, err := GetUserTournamentBalance(bm.UserId, bm.TournamentId, bm.ID)
	if err != nil {
		return errors.New("can't fetch coin data")
	}

	newBets := bm.Bet1 + bm.BetX + bm.Bet2

	// this is a new bet, we just need to check the current balance
	coinsAvailable := utb.StartCoins - utb.CoinsBetsClosed - utb.CoinsBetsOpen + utb.CoinsWon
	if float32(newBets) > coinsAvailable {
		return errors.New("not enough coins")
	} else {
		return nil
	}

	return nil
}
