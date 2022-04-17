package models

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apiagen/configs"
	"github.com/nikitamirzani323/wl_apiagen/db"
	"github.com/nikitamirzani323/wl_apiagen/entities"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
	"github.com/nleeper/goment"
)

func Login_Model(username, password string) (helpers.Response, bool, string, string, int, error) {
	var obj entities.Model_login
	var arraobj []entities.Model_login
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	render := time.Now()
	flag := false
	flag_login := false
	temp_password := ""
	idcompany := ""
	typeadmin := ""
	ruleadmin := 0

	sql_select := `SELECT 
			username_comp , password_comp, idcompany, typeadmin, idruleadmin  
			FROM ` + configs.DB_tbl_mst_Company_admin + ` 
			WHERE status_comp = 'ACTIVE'   
		`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_comp_db, password_comp_db, idcompany_db, typeadmin_db string
			idruleadmin_db                                                 int
		)
		err = row.Scan(&username_comp_db, &password_comp_db, &idcompany_db, &typeadmin_db, &idruleadmin_db)
		temp_username := strings.Trim(username_comp_db, "")
		temp_username = strings.ToLower(temp_username)
		if username_comp_db == temp_username {
			flag_login = true
			temp_password = password_comp_db
			idcompany = idcompany_db
			typeadmin = typeadmin_db
			ruleadmin = idruleadmin_db
		}

		obj.Login_username = username_comp_db
		obj.Login_password = password_comp_db
		obj.Login_typeadmin = typeadmin_db
		obj.Login_idcompany = idcompany
		obj.Login_idrule = idruleadmin_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	if flag_login {
		hashpass := helpers.HashPasswordMD5(password)
		log.Println("Pwd : ", password)
		log.Println("Pwd Hash :", hashpass)
		log.Println("DB Hash :", temp_password)
		if hashpass == temp_password {
			flag = true
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(render).String()

	return res, flag, idcompany, typeadmin, ruleadmin, nil
}
func Update_login(username, ipaddress, timezone, idcompany string) {
	tglnow, _ := goment.New()
	sql_update := `
			UPDATE ` + configs.DB_tbl_mst_Company_admin + ` 
			SET lastlogin_comp=$1, lastipaddres_comp=$2 , lasttimezone_comp=$3, 
			updatecomp_admin=$4,  updatedatecomp_admin=$5  
			WHERE username_comp  = $6 
			AND idcompany = $7  
			AND status_comp = 'ACTIVE' 
		`
	lastlogin := tglnow.Format("YYYY-MM-DD HH:mm:ss")
	flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_Company_admin, "UPDATE",
		lastlogin, ipaddress, timezone, username,
		tglnow.Format("YYYY-MM-DD HH:mm:ss"), username, idcompany)

	if flag_update {
		log.Println(msg_update)
	}

	notelog := ""
	notelog += "DATE : " + lastlogin + " <br>"
	notelog += "IP : " + ipaddress
	Insert_log("MASTER", idcompany, username, "LOGIN", "UPDATE", notelog)
}
