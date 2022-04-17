package entities

type Model_adminruleall struct {
	Adminrule_idrule int    `json:"adminrule_idrule"`
	Adminrule_nmrule string `json:"adminrule_nmrule"`
	Adminrule_rule   string `json:"adminrule_rule"`
	Adminrule_create string `json:"adminrule_create"`
	Adminrule_update string `json:"adminrule_update"`
}

type Controller_adminrule struct {
	Master string `json:"master" validate:"required"`
}

type Controller_adminrulesave struct {
	Sdata  string `json:"sdata" validate:"required"`
	Page   string `json:"page" validate:"required"`
	Idrule int    `json:"adminrule_idrule" "`
	Nmrule string `json:"adminrule_nmrule" validate:"required"`
	Rule   string `json:"adminrule_rule" `
}

type Responseredis_adminruleall struct {
	Adminrule_idadmin string `json:"adminrule_idadmin"`
	Adminrule_rule    string `json:"adminrule_rule"`
}
