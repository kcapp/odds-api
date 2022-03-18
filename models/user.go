package models

type User struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Login          string `json:"login"`
	Pin            string `json:"pin"`
	RequiresChange bool   `json:"requires_change"`
}
