package entities

type Model_agen struct {
	Agen_idagen      string `json:"agen_idagen"`
	Agen_startjoin   string `json:"agen_startjoin"`
	Agen_endjoin     string `json:"agen_endjoin"`
	Agen_idcurr      string `json:"agen_idcurr"`
	Agen_nmagen      string `json:"agen_nmagen"`
	Agen_ownername   string `json:"agen_owneragen"`
	Agen_ownerphone  string `json:"agen_ownerphone"`
	Agen_owneremail  string `json:"agen_owneremail"`
	Agen_urlendpoint string `json:"agen_urlendpoint"`
	Agen_status      string `json:"agen_status"`
	Agen_create      string `json:"agen_create"`
	Agen_update      string `json:"agen_update"`
}
type Model_agencurr struct {
	Curr_idcurr string `json:"curr_idcurr"`
	Curr_nama   string `json:"curr_nama"`
}
type Model_agenadmin struct {
	Agenadmin_username      string `json:"agenadmin_username"`
	Agenadmin_type          string `json:"agenadmin_type"`
	Agenadmin_name          string `json:"agenadmin_name"`
	Agenadmin_phone         string `json:"agenadmin_phone"`
	Agenadmin_email         string `json:"agenadmin_email"`
	Agenadmin_lastlogin     string `json:"agenadmin_lastlogin"`
	Agenadmin_lastipaddress string `json:"agenadmin_lastipaddress"`
	Agenadmin_status        string `json:"agenadmin_status"`
	Agenadmin_create        string `json:"agenadmin_create"`
	Agenadmin_update        string `json:"agenadmin_update"`
}
type Controller_agen struct {
	Master string `json:"master" validate:"required"`
}
type Controller_agensave struct {
	Sdata            string `json:"sdata" validate:"required"`
	Page             string `json:"page" validate:"required"`
	Agen_idagen      string `json:"agen_idagen" validate:"required"`
	Agen_idcurr      string `json:"agen_idcurr" validate:"required"`
	Agen_nmagen      string `json:"agen_nmagen" validate:"required"`
	Agen_ownername   string `json:"agen_owneragen" validate:"required"`
	Agen_ownerphone  string `json:"agen_ownerphone" validate:"required"`
	Agen_owneremail  string `json:"agen_owneremail" validate:"required"`
	Agen_urlendpoint string `json:"agen_urlendpoint"`
	Agen_status      string `json:"agen_status" validate:"required"`
}
type Controller_agenadmin struct {
	Idagen string `json:"idagen" validate:"required"`
}
type Controller_agenadminsave struct {
	Sdata              string `json:"sdata" validate:"required"`
	Page               string `json:"page" validate:"required"`
	Agen_idagen        string `json:"agen_idagen" validate:"required"`
	Agenadmin_username string `json:"agenadmin_username" validate:"required"`
	Agenadmin_password string `json:"agenadmin_password" `
	Agenadmin_name     string `json:"agenadmin_name" validate:"required"`
	Agenadmin_phone    string `json:"agenadmin_phone" validate:"required"`
	Agenadmin_email    string `json:"agenadmin_email" validate:"required"`
	Agenadmin_status   string `json:"agenadmin_status" validate:"required"`
}
