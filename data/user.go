package data

import "github.com/kcapp/odds-api/models"

func GetUser(id int) (*models.User, error) {
	u := new(models.User)
	err := models.DB.QueryRow(`
		SELECT
			u.id, u.first_name, u.last_name, u.login, u.pin, u.requires_change
		FROM users u
		WHERE u.id = ?`, id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Login, &u.Pin, &u.RequiresChange)
	if err != nil {
		return nil, err
	}

	return u, nil
}
