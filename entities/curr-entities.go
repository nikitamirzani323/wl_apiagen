package entities

type Model_curr struct {
	Curr_idcurr string `json:"curr_idcur"`
	Curr_nama   string `json:"curr_nama"`
	Curr_create string `json:"curr_create"`
	Curr_update string `json:"curr_update"`
}

type Controller_curr struct {
	Master string `json:"master" validate:"required"`
}

type Controller_currsave struct {
	Sdata  string `json:"sdata" validate:"required"`
	Page   string `json:"page" validate:"required"`
	Idcurr string `json:"curr_idcurr" validate:"required"`
	Name   string `json:"curr_name" validate:"required"`
}
