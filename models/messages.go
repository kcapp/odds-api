package models

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type Success struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
