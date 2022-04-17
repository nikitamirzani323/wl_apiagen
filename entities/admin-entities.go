package entities

type Model_admin struct {
	Admin_username      string `json:"admin_username"`
	Admin_type          string `json:"admin_type"`
	Admin_idrule        int    `json:"admin_idrule"`
	Admin_rule          string `json:"admin_rule"`
	Admin_nama          string `json:"admin_nama"`
	Admin_phone         string `json:"admin_phone"`
	Admin_email         string `json:"admin_email"`
	Admin_joindate      string `json:"admin_joindate"`
	Admin_lastlogin     string `json:"admin_lastlogin"`
	Admin_lastIpaddress string `json:"admin_lastipaddres"`
	Admin_status        string `json:"admin_status"`
	Admin_create        string `json:"admin_create"`
	Admin_update        string `json:"admin_update"`
}
type Model_adminrule struct {
	Admin_idrule int    `json:"adminrule_idrule"`
	Admin_nmrule string `json:"adminrule_nmrule"`
}
type Model_adminsave struct {
	Username string `json:"admin_username"`
	Nama     string `json:"admin_nama"`
	Rule     string `json:"admin_rule"`
	Status   string `json:"admin_status"`
	Create   string `json:"admin_create"`
	Update   string `json:"admin_update"`
}
type Controller_admin struct {
	Master string `json:"master" validate:"required"`
}
type Controller_admindetail struct {
	Username string `json:"agen_username" validate:"required"`
}
type Controller_adminsave struct {
	Sdata    string `json:"sdata" validate:"required"`
	Page     string `json:"page" validate:"required"`
	Username string `json:"admin_username" validate:"required"`
	Password string `json:"admin_password"`
	Nama     string `json:"admin_nama" validate:"required"`
	Email    string `json:"admin_email" validate:"required"`
	Phone    string `json:"admin_phone" validate:"required"`
	Idrule   int    `json:"admin_idrule" validate:"required"`
	Status   string `json:"admin_status" validate:"required"`
}

type Responseredis_adminrule struct {
	Adminrule_idrule string `json:"adminrule_idrule"`
}
