package data

import (
	"errors"
	"github.com/kcapp/odds-api/models"
)

func RunTransaction(sql string, args ...interface{}) (int64, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return 0, errors.New("error creating transaction")
	}

	res, err := tx.Exec(sql, args...)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return lid, err
}
