package entities

type Model_game struct {
	Game_id     int    `json:"game_id"`
	Game_name   string `json:"game_name"`
	Game_create string `json:"game_create"`
	Game_update string `json:"game_update"`
}

type Controller_game struct {
	Game_search string `json:"game_search"`
	Game_page   int    `json:"game_page"`
}
type Controller_gamesave struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	Game_search   string `json:"game_search"`
	Game_page     int    `json:"game_page"`
	Game_idrecord int    `json:"game_id"`
	Game_name     string `json:"game_name" validate:"required"`
}
