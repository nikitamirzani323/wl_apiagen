package entities

type Model_log struct {
	Log_id       int    `json:"log_id"`
	Log_datetime string `json:"log_datetime"`
	Log_user     string `json:"log_user"`
	Log_page     string `json:"log_page"`
	Log_tipe     string `json:"log_tipe"`
	Log_note     string `json:"log_note"`
}

type Controller_log struct {
	Master string `json:"master" validate:"required"`
}
