package models

type User struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	RequiresChange bool   `json:"requires_change"`
}

type Authentication struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	UserId         int    `json:"user_id"`
	Login          string `json:"login"`
	TokenString    string `json:"token"`
	RequiresChange bool   `json:"requires_change"`
}
