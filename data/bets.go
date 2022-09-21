package data

import (
	"database/sql"
	"errors"
	"sort"

	"github.com/kcapp/odds-api/models"
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

func GetUserTournamentTournamentCoinsOpen(userId, tournamentId, skipOutcomeId int) (*models.CoinBalance, error) {
	var s string
	cb := new(models.CoinBalance)

	if skipOutcomeId != 0 {
		s = `select COALESCE(sum(bg.bet1+bg.betx+bg.bet2), 0) as betCoins
			from bets_tournament bg
			where bg.user_id = ? and bg.tournament_id = ? and bg.outcome is null AND bg.id != ?`
		err := models.DB.QueryRow(s, userId, tournamentId, skipOutcomeId).Scan(&cb.Coins)
		if err != nil {
			return nil, err
		}
	} else {
		s = `select COALESCE(sum(bg.bet1+bg.betx+bg.bet2), 0) as betCoins
			from bets_tournament bg
			where bg.user_id = ? and bg.tournament_id = ? and bg.outcome is null`
		err := models.DB.QueryRow(s, userId, tournamentId).Scan(&cb.Coins)
		if err != nil {
			return nil, err
		}
	}

	return cb, nil
}

func GetUserTournamentTournamentCoinsClosed(userId, tournamentId int) (*models.CoinBalance, error) {
	s := `select COALESCE(sum(bg.bet1+bg.betx+bg.bet2), 0) as betCoins
			from bets_tournament bg
			where bg.user_id = ? and bg.tournament_id = ? and bg.outcome is not null`

	cb := new(models.CoinBalance)
	err := models.DB.QueryRow(s, userId, tournamentId).
		Scan(&cb.Coins)
	if err != nil {
		return nil, err
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

func GetUserTournamentTournamentCoinsWon(userId, tournamentId int) (*models.CoinBalance, error) {
	s := `select
				   coalesce(sum(coalesce(if (bt.outcome > ou.value, ou.odds1 * bt.bet1, ou.odds2 * bet2),
				   if (bt.outcome = op.value, op.oddsx * bt.betx, 0))), 0) as win
			from bets_tournament bt
			left join (select ou.id, ou.tournament_id,
						 ou.odds1, ou.oddsx, ou.odds2, m.type_id, ou.value
				  from outcomes ou
				  join markets m on ou.market_id = m.id
				  where m.type_id = 1
				  ) ou on ou.id = bt.outcome_id and ou.tournament_id = bt.tournament_id
			left join (select op.id, op.tournament_id,
						 op.odds1, op.oddsx, op.odds2, m.type_id, op.value
				  from outcomes op
				  join markets m on op.market_id = m.id
				  where m.type_id IN (2, 3)
				  ) op on op.id = bt.outcome_id and op.tournament_id = bt.tournament_id
			where bt.user_id = ? and bt.tournament_id = ? and bt.outcome is not null`

	cb := new(models.CoinBalance)
	err := models.DB.QueryRow(s, userId, tournamentId).
		Scan(&cb.Coins)
	if err != nil {
		return nil, err
	}

	return cb, nil
}

func GetTournamentGameRanking(tournamentId int) ([]*models.UserTournamentBalance, error) {
	rows, err := models.DB.Query(`
		select bg.user_id,
			   u.first_name,
			   u.last_name,
			   u.is_cheater,
			   bg.tournament_id,
			   (coalesce(bgo.numBetsOpen, 0) + coalesce(bgc.numBetsClosed, 0)) as numBets,
			   coalesce(bgc.numBetsClosed, 0)                                  as numBetsClosed,
			   coalesce(bgo.coinsOpenBets, 0)                                  as coinsOpenBets,
			   coalesce(bgc.coinsClosedBets, 0)                                as coinsClosedBets,
			   coalesce(bgc.coinsWon, 0)                                       as coinsWon,
			   coalesce(bgo.potentialWinnings, 0)                              as potentialWinnings,
			   1000,
			   1000,
			   1000
		from bets_games bg
				 left join (select bgo.user_id,
								   count(bgo.user_id)                  as numBetsOpen,
								   bgo.tournament_id,
								   bgo.bet1,
								   bgo.betx,
								   bgo.bet2,
								   sum(bgo.bet1 + bgo.betx + bgo.bet2) as coinsOpenBets,
								   ROUND(SUM(GREATEST(if(bgo.bet1 > 0, bgo.bet1 * bgo.odds1 - (bgo.bet1 + bgo.bet2 + bgo.betx), 0),
													  if(bgo.betx > 0, bgo.betx * bgo.oddsx - (bgo.bet1 + bgo.bet2 + bgo.betx), 0),
													  if(bgo.bet2 > 0, bgo.bet2 * bgo.odds2 - (bgo.bet1 + bgo.bet2 + bgo.betx), 0))),
										 2)                            as potentialWinnings
							from bets_games bgo
							where bgo.outcome IS NULL
							  and tournament_id = ?
							group by bgo.user_id) bgo
						   on bg.user_id = bgo.user_id and bg.tournament_id = bgo.tournament_id
				 left join (select bgc.user_id,
								   count(bgc.user_id)                                                       as numBetsClosed,
								   bgc.tournament_id,
								   sum(bgc.bet1 + bgc.betx + bgc.bet2)                                      as coinsClosedBets,
								   ROUND(SUM(if(bgc.player1 = bgc.outcome, bet1 * odds1 - bet1,
											 if(bgc.player2 = bgc.outcome, bet2 * odds2 - bet2,
											 if(bgc.outcome = 0, betx * oddsx - betx, 0)))), 2) as rawCoinsWon,
								   ROUND(SUM(if(bgc.player1 = bgc.outcome, bet1 * odds1,
											 if(bgc.player2 = bgc.outcome, bet2 * odds2,
											 if(bgc.outcome = 0, betx * oddsx, 0)))), 2)        as coinsWon
							from bets_games bgc
							where bgc.outcome IS NOT NULL
							  and tournament_id = ?
							group by bgc.user_id) bgc
						   on bg.user_id = bgc.user_id and bg.tournament_id = bgc.tournament_id
				 join users u on bg.user_id = u.id
		where bg.tournament_id = ?
		group by bg.user_id`, tournamentId, tournamentId, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]*models.UserTournamentBalance, 0)
	for rows.Next() {
		b := new(models.UserTournamentBalance)
		err := rows.Scan(&b.UserId, &b.FirstName, &b.LastName, &b.IsCheater, &b.TournamentId,
			&b.BetsPlaced, &b.BetsClosed, &b.CoinsBetsOpen, &b.CoinsBetsClosed, &b.CoinsWon, &b.PotentialWinnings,
			&b.TournamentCoinsOpen, &b.TournamentCoinsClosed, &b.StartCoins)
		b.CoinsAvailable = b.StartCoins - b.CoinsBetsOpen - b.CoinsBetsClosed + b.CoinsWon
		if err != nil {
			return nil, err
		}

		balances = append(balances, b)
	}

	sort.Sort(models.SortBalanceByPotentialCoins(balances))

	if err != nil {
		return nil, err
	}

	return balances, nil
}

func GetTournamentRanking(tournamentId int) ([]*models.UserTournamentBalance, error) {
	rows, err := models.DB.Query(`
		select bt.user_id, u.first_name, u.last_name, bt.tournament_id,
			(coalesce(bto1.numBetsOpen, 0) + coalesce(bto2.numBetsOpen, 0)) + (coalesce(btc1.numBetsClosed, 0) + coalesce(btc2.numBetsClosed, 0)) as numBets,
			(coalesce(btc1.numBetsClosed, 0) + coalesce(btc2.numBetsClosed, 0)) as numBetsClosed,
			0,0,0,
			coalesce(bto1.totalPotential, 0) + coalesce (bto2.futuresPropsPotential, 0) as potentialWinnings,
			coalesce(bto1.coinsOpenBets, 0) + coalesce(bto2.coinsOpenBets, 0) as coinsOpenBets,
			coalesce(btc1.coinsClosedBets, 0) + coalesce(btc2.coinsClosedBets, 0) as coinsClosedBets,
			coalesce(btc1.coinsWon, 0) + coalesce(btc2.coinsWon, 0) as coinsWon,
			1000
		from bets_tournament bt
            left join
            (select bto.user_id, count(bto.user_id) as numBetsOpen, bto.tournament_id,
                   bto.bet1, bto.betx, bto.bet2,
                   sum(bto.bet1 + bto.betx + bto.bet2) as coinsOpenBets,
                   ROUND(SUM(GREATEST(if(bto.bet1 > 0, bto.bet1 * o.odds1 - (bto.bet1 + bto.bet2), 0),
                                                          if(bto.bet2 > 0, bto.bet2 * o.odds2 - (bto.bet1 + bto.bet2), 0))),
                                            2)
                        as totalPotential
            from bets_tournament bto
            join outcomes o on bto.outcome_id = o.id
            join markets m on o.market_id = m.id
            where bto.outcome is null
            and m.type_id = 1
            group by bto.user_id
            ) bto1 on bt.user_id = bto1.user_id and bt.tournament_id = bto1.tournament_id
            left join
            (select bto.user_id, count(bto.user_id) as numBetsOpen, bto.tournament_id,
                   bto.bet1, bto.betx, bto.bet2,
                   sum(bto.bet1 + bto.betx + bto.bet2) as coinsOpenBets,
                   ROUND(SUM(if(bto.betx > 0, bto.betx * o.oddsx - bto.betx, 0)), 2)
                        as futuresPropsPotential
            from bets_tournament bto
            join outcomes o on bto.outcome_id = o.id
            join markets m on o.market_id = m.id
            where bto.outcome is null
            and m.type_id in (2,3)
            group by bto.user_id) bto2 on bt.user_id = bto2.user_id and bt.tournament_id = bto2.tournament_id
		    # closed bets for futures
            left join (select bto.user_id, count(bto.user_id) as numBetsClosed, bto.tournament_id,
                   bto.bet1, bto.betx, bto.bet2,
                   sum(bto.bet1 + bto.betx + bto.bet2) as coinsClosedBets,
                   sum(if(bto.outcome > o.value, bto.bet1 * o.odds1 - (bto.bet1 + bto.bet2), if(bto.outcome < o.value, bto.bet2 * o.odds2 - (bto.bet1 + bto.bet2), 0)))
                        as coinsWon
            from bets_tournament bto
            join outcomes o on bto.outcome_id = o.id
            join markets m on o.market_id = m.id
            where bto.outcome is not null
            and m.type_id = 1
            group by bto.user_id
            ) btc1 on bt.user_id = btc1.user_id and bt.tournament_id = btc1.tournament_id
            left join (select bto.user_id, count(bto.user_id) as numBetsClosed, bto.tournament_id,
                   bto.bet1, bto.betx, bto.bet2,
                   sum(bto.bet1 + bto.betx + bto.bet2) as coinsClosedBets,
                   if(bto.outcome = o.value, bto.betx * o.oddsx - (bto.betx), 0) as coinsWon
            from bets_tournament bto
            join outcomes o on bto.outcome_id = o.id
            join markets m on o.market_id = m.id
            where bto.outcome is not null
            and m.type_id IN (2,3)
            group by bto.user_id
            ) btc2 on bt.user_id = btc2.user_id and bt.tournament_id = btc2.tournament_id
						 join users u on bt.user_id = u.id
		where bt.tournament_id = ?
		group by bt.user_id`, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]*models.UserTournamentBalance, 0)
	for rows.Next() {
		b := new(models.UserTournamentBalance)
		err := rows.Scan(&b.UserId, &b.FirstName, &b.LastName, &b.TournamentId,
			&b.BetsPlaced, &b.BetsClosed, &b.CoinsBetsOpen, &b.CoinsBetsClosed, &b.CoinsWon, &b.PotentialWinnings,
			&b.TournamentCoinsOpen, &b.TournamentCoinsClosed, &b.TournamentCoinsWon, &b.StartCoins)
		b.CoinsAvailable = b.StartCoins - b.TournamentCoinsOpen + b.TournamentCoinsWon
		if err != nil {
			return nil, err
		}

		balances = append(balances, b)
	}

	sort.Sort(models.SortBalanceByPotentialCoins(balances))

	if err != nil {
		return nil, err
	}

	return balances, nil
}

func GetUserTournamentRanking(tournamentId int, userId int) (*models.UserTournamentBalance, error) {
	s := `
		select bg.user_id, u.first_name, u.last_name, bg.tournament_id,
			   (coalesce(bgo.numBetsOpen, 0) + coalesce(bgc.numBetsClosed, 0)) as numBets,
			   coalesce(bgc.numBetsClosed, 0) as numBetsClosed,
			   coalesce(bgo.coinsOpenBets, 0) as coinsOpenBets,
			   coalesce(bgc.coinsClosedBets, 0 ) as coinsClosedBets,
			   coalesce(bgc.coinsWon, 0) as coinsWon,
			   coalesce(bgo.potentialWinnings, 0) as potentialWinnings,
			   1000, 1000, 1000
		from bets_games bg
				 left join (select bgo.user_id, count(bgo.user_id) as numBetsOpen, bgo.tournament_id,
							bgo.bet1, bgo.betx, bgo.bet2,
							sum(bgo.bet1 + bgo.betx + bgo.bet2) as coinsOpenBets,
							ROUND(SUM(GREATEST(if(bgo.bet1 > 0, bgo.bet1 * bgo.odds1 - (bgo.bet1 + bgo.bet2 + bgo.betx), 0),
											   if(bgo.bet2 > 0, bgo.bet2 * bgo.odds2 - (bgo.bet1 + bgo.bet2 + bgo.betx), 0),
							    			   if(bgo.betx > 0, bgo.betx * bgo.oddsx - (bgo.bet1 + bgo.bet2 + bgo.betx), 0)
							    )),
								2) as potentialWinnings
							from bets_games bgo
							where bgo.outcome IS NULL
							group by bgo.user_id, bgo.tournament_id) bgo
						   on bg.user_id = bgo.user_id and bg.tournament_id = bgo.tournament_id
				 left join (select bgc.user_id, count(bgc.user_id) as numBetsClosed, bgc.tournament_id,
							sum(bgc.bet1 + bgc.betx + bgc.bet2) as coinsClosedBets,
							ROUND(SUM(
							    if(bgc.player1 = bgc.outcome, bet1 * odds1 - bet1, 
							        if(bgc.player2 = bgc.outcome, bet2 * odds2 - bet2, 
							            if(bgc.outcome = 0, betx * oddsx - betx, 0)
							            ))), 2) as rawCoinsWon,
							ROUND(SUM(
							    if(bgc.player1 = bgc.outcome, bet1 * odds1, 
							        if(bgc.player2 = bgc.outcome, bet2 * odds2, 
							          if(bgc.outcome = 0, betx * oddsx, 0)  
							            ))), 2) as coinsWon
							from bets_games bgc
							where bgc.outcome IS NOT NULL
							group by bgc.user_id, bgc.tournament_id) bgc
						   on bg.user_id = bgc.user_id and bg.tournament_id = bgc.tournament_id
				 join users u on bg.user_id = u.id
		where bg.tournament_id = ? and bg.user_id = ?
		group by bg.user_id`

	b := models.UserTournamentBalance{
		UserId:                userId,
		FirstName:             "",
		LastName:              "",
		TournamentId:          tournamentId,
		BetsPlaced:            0,
		BetsClosed:            0,
		CoinsBetsOpen:         0,
		CoinsBetsClosed:       0,
		CoinsWon:              0,
		PotentialWinnings:     0,
		TournamentCoinsOpen:   0,
		TournamentCoinsClosed: 0,
		TournamentCoinsWon:    0,
		StartCoins:            0,
		CoinsAvailable:        0,
	}

	err := models.DB.QueryRow(s, tournamentId, userId).
		Scan(&b.UserId, &b.FirstName, &b.LastName, &b.TournamentId,
			&b.BetsPlaced, &b.BetsClosed, &b.CoinsBetsOpen, &b.CoinsBetsClosed, &b.CoinsWon, &b.PotentialWinnings,
			&b.TournamentCoinsOpen, &b.TournamentCoinsClosed, &b.StartCoins)
	b.CoinsAvailable = b.StartCoins - b.CoinsBetsOpen - b.CoinsBetsClosed + b.CoinsWon

	if err == sql.ErrNoRows {
		u, err := GetUserById(userId)

		if err != nil {
			return nil, err
		}

		b.FirstName = u.FirstName
		b.LastName = u.LastName
		b.StartCoins = 1000
		b.CoinsAvailable = 1000
		return &b, err
	}

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func GetTournamentStatistics(tournamentId int) (*models.TournamentStatistics, error) {
	s := `select bg.tournament_id, count(bg.id) + bt.numTournamentBets as numBets
										from bets_games bg
										join (select bt.tournament_id, count(bt.id) as numTournamentBets
											  from bets_tournament bt
											  where bt.tournament_id = ?
											  ) bt on bt.tournament_id = bg.tournament_id
										where bg.tournament_id = ?
										group by bg.tournament_id`

	b := models.TournamentStatistics{
		TournamentId:        tournamentId,
		TournamentBetsCount: 0,
		TournamentUserCount: 0,
		BankWinnings:        0,
		UserWinnings:        0,
		BiggestWins:         nil,
		BiggestLosses:       nil,
	}

	err := models.DB.QueryRow(s, tournamentId, tournamentId).Scan(&b.TournamentId, &b.TournamentBetsCount)

	if err == sql.ErrNoRows {
		return &b, err
	}

	s = `select count(t.user_id) numUsers from
		(select bg.user_id
		from bets_games bg
		where bg.tournament_id = ?
		union
		select bt.user_id
		from bets_tournament bt
		where bt.tournament_id = ?) t`
	err = models.DB.QueryRow(s, tournamentId, tournamentId).Scan(&b.TournamentUserCount)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func GetTournamentOutcomes(tournamentId int) ([]*models.TournamentOutcome, error) {
	rows, err := models.DB.Query(`select o.id as outcomeId,
			   o.tournament_id as tournamentId,
			   m.id as marketId,
			   m.name as marketName,
			   m.type_id as marketTypeId,
			   mt.name as marketTypeName,
			   o.value as outcomeValue, o.odds1, o.odds2, o.oddsx,
			   concat(u.first_name, ' ', u.last_name) as playerName
		from outcomes o
		join markets m on o.market_id = m.id
		join market_types mt on m.type_id = mt.id
		left join users u on u.id = o.value and m.type_id IN (2, 3)
		where o.tournament_id = ?
		order by o.market_id, o.oddsx`, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outcomes := make([]*models.TournamentOutcome, 0)
	for rows.Next() {
		b := new(models.TournamentOutcome)
		err := rows.Scan(&b.ID, &b.TournamentId, &b.MarketId, &b.MarketName, &b.MarketTypeId, &b.MarketTypeName,
			&b.OutcomeValue, &b.Odds1, &b.Odds2, &b.OddsX, &b.PlayerName)
		if err != nil {
			return nil, err
		}

		outcomes = append(outcomes, b)
	}

	if err != nil {
		return nil, err
	}

	return outcomes, nil
}

func GetUserGamesBets(userId int) ([]*models.BetMatch, error) {
	rows, err := models.DB.Query(`
			SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2, bm.outcome,
			       bm.odds1, bm.oddsx, bm.odds2, bm.player1, bm.player2,
			       COALESCE(gm.bets_off, 0) as bets_off
			FROM bets_games bm
			LEFT JOIN games_metadata gm on bm.match_id = gm.match_id
			WHERE bm.user_id = ? ORDER BY bm.id DESC`, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]*models.BetMatch, 0)
	for rows.Next() {
		b := new(models.BetMatch)
		err := rows.Scan(&b.ID, &b.UserId, &b.MatchId, &b.TournamentId,
			&b.Bet1, &b.BetX, &b.Bet2, &b.Outcome,
			&b.Odds1, &b.OddsX, &b.Odds2, &b.Player1, &b.Player2, &b.BetsOff)
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

func GetUserTournamentGamesBets(userId, tournamentId int) ([]*models.BetMatch, error) {
	rows, err := models.DB.Query(`
			SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2, bm.outcome,
			       bm.odds1, bm.oddsx, bm.odds2, bm.player1, bm.player2,
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
		err := rows.Scan(&b.ID, &b.UserId, &b.MatchId, &b.TournamentId,
			&b.Bet1, &b.BetX, &b.Bet2, &b.Outcome,
			&b.Odds1, &b.OddsX, &b.Odds2, &b.Player1, &b.Player2, &b.BetsOff)
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

func GetUserTournamentTournamentsBets(userId, tournamentId int) ([]*models.BetTournament, error) {
	rows, err := models.DB.Query(`
			SELECT
			bt.id, bt.user_id, bt.tournament_id, bt.outcome_id, bt.bet1, bt.betx, bt.bet2, bt.outcome,
			o.value, o.odds1, o.odds2, o.oddsx, m.id as market_id, m.type_id as market_type_id
			FROM bets_tournament bt
			join outcomes o on bt.outcome_id = o.id
			join markets m on o.market_id = m.id
			WHERE bt.user_id = ? and bt.tournament_id = ?`, userId, tournamentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]*models.BetTournament, 0)
	for rows.Next() {
		b := new(models.BetTournament)
		err := rows.Scan(&b.ID, &b.UserId, &b.TournamentId, &b.OutcomeId,
			&b.Bet1, &b.BetX, &b.Bet2, &b.Outcome,
			&b.OutcomeValue, &b.Odds1, &b.Odds2, &b.OddsX, &b.MarketId, &b.MarketTypeId)
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

func GetGameBets(gameId int) ([]*models.BetMatch, error) {
	rows, err := models.DB.Query(`
			SELECT
			bm.id, bm.user_id, bm.match_id, bm.tournament_id, bm.bet1, bm.betx, bm.bet2, bm.outcome,
			       bm.odds1, bm.oddsx, bm.odds2, bm.player1, bm.player2,
			       COALESCE(gm.bets_off, 0) as bets_off
			FROM bets_games bm
			LEFT JOIN games_metadata gm on bm.match_id = gm.match_id
			WHERE bm.match_id = ?`, gameId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]*models.BetMatch, 0)
	for rows.Next() {
		b := new(models.BetMatch)
		err := rows.Scan(&b.ID, &b.UserId, &b.MatchId, &b.TournamentId,
			&b.Bet1, &b.BetX, &b.Bet2, &b.Outcome,
			&b.Odds1, &b.OddsX, &b.Odds2, &b.Player1, &b.Player2,
			&b.BetsOff)
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

func deleteTournamentBet(bet models.BetOutcome) (int64, error) {
	s := `DELETE FROM bets_tournament
				WHERE outcome_id = ? and user_id = ? AND outcome IS NULL
				AND tournament_id NOT IN (SELECT tm.tournament_id from tournaments_metadata tm WHERE tm.bets_off = 1)`
	args := make([]interface{}, 0)
	args = append(args, bet.OutcomeId, bet.UserId)

	lid, err := RunTransaction(s, args...)

	return lid, err
}

func insertTournamentBet(bet models.BetOutcome) (int64, error) {
	sq := `INSERT INTO bets_tournament (user_id, outcome_id, tournament_id, bet1, betx, bet2)
			VALUES (?, ?, ?, ?, ?, ?)`

	args := make([]interface{}, 0)
	args = append(args, bet.UserId, bet.OutcomeId, bet.TournamentId, bet.Bet1, bet.BetX, bet.Bet2)
	lid, err := RunTransaction(sq, args...)

	return lid, err
}

func AddTournamentBet(bet models.BetOutcome) (int64, error) {
	var err error

	err = ValidateTournamentBetInput(bet)
	if err != nil {
		return 0, err
	}

	// We can have an existing bet with 0s passed, then we should remove the row from the db
	if bet.Bet1+bet.Bet2+bet.BetX == 0 {
		lid, err := deleteTournamentBet(bet)
		return lid, err
	}

	if bet.ID == 0 {
		// We are placing new bet, but we might be passing 0s only
		if bet.Bet1+bet.Bet2+bet.BetX == 0 {
			return 0, errors.New("can't place an empty bet")
		}
		lid, err := insertTournamentBet(bet)

		return lid, err
	} else {

		sq := `UPDATE bets_tournament SET bet1 = ?, betx = ?, bet2 = ?, outcome_id = ?
				WHERE user_id = ? AND tournament_id = ? AND id = ?`

		args := make([]interface{}, 0)
		args = append(args, bet.Bet1, bet.BetX, bet.Bet2, bet.OutcomeId, bet.UserId, bet.TournamentId, bet.ID)
		_, err := RunTransaction(sq, args...)

		return int64(bet.ID), err
	}
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

func GetTournamentsMetadata() ([]*models.TournamentMetadata, error) {
	rows, err := models.DB.Query(`SELECT tm.tournament_id, tm.bets_off FROM tournaments_metadata tm`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	metadata := make([]*models.TournamentMetadata, 0)
	for rows.Next() {
		tm := new(models.TournamentMetadata)
		err := rows.Scan(&tm.TournamentId, &tm.BetsOff)
		if err != nil {
			return nil, err
		}

		metadata = append(metadata, tm)
	}

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func ValidateTournamentBetInput(bt models.BetOutcome) error {
	betsOff := CheckBetOff(bt.OutcomeId)
	if betsOff == 1 {
		return errors.New("bet are off for this match")
	}

	utb, err := GetUserTournamentBalance(bt.UserId, bt.TournamentId, bt.ID)
	if err != nil {
		return errors.New("can't fetch coin data")
	}

	newBets := bt.Bet1 + bt.BetX + bt.Bet2

	// this is a new bet, we just need to check the current balance
	coinsAvailable := utb.StartCoins - utb.TournamentCoinsClosed - utb.TournamentCoinsOpen + utb.TournamentCoinsWon
	if float32(newBets) > coinsAvailable {
		return errors.New("not enough coins")
	} else {
		return nil
	}

	return nil
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
