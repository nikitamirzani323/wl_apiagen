package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apiagen/configs"
	"github.com/nikitamirzani323/wl_apiagen/db"
	"github.com/nikitamirzani323/wl_apiagen/entities"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
	"github.com/nleeper/goment"
)

func Fetch_adminHome(idcompany string) (helpers.ResponseAdmin, error) {
	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var res helpers.ResponseAdmin
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			username_comp , typeadmin, idruleadmin, 
			nama_comp, email_comp, phone_comp, 
			status_comp, to_char(COALESCE(lastlogin_comp,now()), 'YYYY-MM-DD HH24:MI:SS') as lastlogin_comp, 
			lastipaddres_comp,   
			createcomp_admin, to_char(COALESCE(createdatecomp_admin,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecomp_admin, 
			updatecomp_admin, to_char(COALESCE(updatedatecomp_admin,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecomp_admin 
			FROM ` + configs.DB_tbl_mst_Company_admin + ` 
			WHERE idcompany = $1
			ORDER BY lastlogin_comp DESC 
		`

	row, err := con.QueryContext(ctx, sql_select, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_comp_db, typeadmin_db                                                             string
			idruleadmin_db                                                                             int
			nama_comp_db, email_comp_db, phone_comp_db                                                 string
			status_comp_db, lastlogin_comp_db, lastipaddres_comp_db                                    string
			createcomp_admin_db, createdatecomp_admin_db, updatecomp_admin_db, updatedatecomp_admin_db string
		)

		err = row.Scan(
			&username_comp_db, &typeadmin_db, &idruleadmin_db,
			&nama_comp_db, &email_comp_db, &phone_comp_db,
			&status_comp_db, &lastlogin_comp_db, &lastipaddres_comp_db,
			&createcomp_admin_db, &createdatecomp_admin_db, &updatecomp_admin_db, &updatedatecomp_admin_db)

		helpers.ErrorCheck(err)
		nmrule := _adminrule(idruleadmin_db, "nmcomprule", idcompany)
		create := createcomp_admin_db + ", " + createdatecomp_admin_db
		update := ""
		if updatecomp_admin_db != "" {
			update = updatecomp_admin_db + ", " + updatedatecomp_admin_db
		}
		obj.Admin_username = username_comp_db
		obj.Admin_type = typeadmin_db
		obj.Admin_idrule = idruleadmin_db
		obj.Admin_rule = nmrule
		obj.Admin_nama = nama_comp_db
		obj.Admin_phone = phone_comp_db
		obj.Admin_email = email_comp_db
		obj.Admin_joindate = createdatecomp_admin_db
		obj.Admin_lastlogin = lastlogin_comp_db
		obj.Admin_lastIpaddress = lastipaddres_comp_db
		obj.Admin_status = status_comp_db
		obj.Admin_create = create
		obj.Admin_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objRule entities.Model_adminrule
	var arraobjRule []entities.Model_adminrule
	sql_listrule := `SELECT 
		idcomprule, nmcomprule 	
		FROM ` + configs.DB_tbl_mst_Company_adminrule + ` 
		WHERE idcompany=$1 
	`
	row_listrule, err_listrule := con.QueryContext(ctx, sql_listrule, idcompany)

	helpers.ErrorCheck(err_listrule)
	for row_listrule.Next() {
		var (
			idcomprule_db int
			nmcomprule_db string
		)

		err = row_listrule.Scan(&idcomprule_db, &nmcomprule_db)

		helpers.ErrorCheck(err)

		objRule.Admin_idrule = idcomprule_db
		objRule.Admin_nmrule = nmcomprule_db
		arraobjRule = append(arraobjRule, objRule)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listrule = arraobjRule
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_adminHome(admin, idcompany, username, password, nama, email, phone, status, sData string, idrule int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false
	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_Company_admin, "username_comp", username)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_Company_admin + ` (
					username_comp , password_comp, idcompany, typeadmin, idruleadmin, 
					nama_comp, email_comp, phone_comp, status_comp, 
					createcomp_admin, createdatecomp_admin
				) values (
					$1, $2, $3, $4, $5,  
					$6, $7
				)
			`
			hashpass := helpers.HashPasswordMD5(password)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_Company_admin, "INSERT",
				username, hashpass, idcompany, "ADMIN", idrule,
				nama, email, phone, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)
				nmrule := _adminrule(idrule, "nmcomprule", idcompany)

				notelog := ""
				notelog += "NEW ADMIN <br>"
				notelog += "USERNAME : " + username + " <br>"
				notelog += "RULE : " + nmrule + " <br>"
				notelog += "TYPE : ADMIN <br>"
				notelog += "NAME : " + nama + " <br>"
				notelog += "PHONE : " + phone + " <br>"
				notelog += "EMAIL : " + email + " <br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "ADMIN", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password == "" {
			sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_Company_admin + `  
				SET idruleadmin=$1 ,nama_comp=$2, email_comp=$3, phone_comp=$4, status_comp=$5, 
				updatecomp_admin=$6, updatedatecomp_admin=$7  
				WHERE username_comp=$8  
				AND idcompany=$9   
			`

			flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_Company_admin, "UPDATE",
				idrule, nama, email, phone, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"),
				username, idcompany)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)
				nmrule := _adminrule(idrule, "nmcomprule", idcompany)
				notelog := ""
				notelog += "UPDATE ADMIN <br>"
				notelog += "USERNAME : " + username + " <br>"
				notelog += "RULE : " + nmrule + " <br>"
				notelog += "TYPE : ADMIN <br>"
				notelog += "NAME : " + nama + " <br>"
				notelog += "PHONE : " + phone + " <br>"
				notelog += "EMAIL : " + email + " <br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "ADMIN", "UPDATE", notelog)
			} else {
				log.Println(msg_update)
			}
		} else {
			hashpass := helpers.HashPasswordMD5(password)
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_Company_admin + `   
				SET password_comp=$1, idruleadmin=$2 ,nama_comp=$3, email_comp=$4, phone_comp=$5, status_comp=$6, 
				updatecomp_admin=$7, updatedatecomp_admin=$8   
				WHERE username_comp=$9   
				AND idcompany=$10    
			`
			flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_Company_admin, "UPDATE",
				hashpass, idrule, nama, email, phone, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"),
				username, idcompany)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)
				nmrule := _adminrule(idrule, "nmcomprule", idcompany)
				notelog := ""
				notelog += "UPDATE ADMIN - CHANGE PASSWORD <br>"
				notelog += "USERNAME : " + username + " <br>"
				notelog += "RULE : " + nmrule + " <br>"
				notelog += "TYPE : ADMIN <br>"
				notelog += "NAME : " + nama + " <br>"
				notelog += "PHONE : " + phone + " <br>"
				notelog += "EMAIL : " + email + " <br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "ADMIN", "UPDATE", notelog)
			} else {
				log.Println(msg_update)
			}
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func _adminrule(idrule int, tipe, idcompany string) string {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	result := ""
	nmcomprule := ""

	sql_select := `SELECT
		nmcomprule  
		FROM ` + configs.DB_tbl_mst_Company_adminrule + `  
		WHERE idcomprule = $1 
		AND idcompany = $2 
	`
	row := con.QueryRowContext(ctx, sql_select, idrule, idcompany)
	switch e := row.Scan(&nmcomprule); e {
	case sql.ErrNoRows:
		flag = false
	case nil:
		flag = true

	default:
		panic(e)
	}
	if flag {
		switch tipe {
		case "nmcomprule":
			result = nmcomprule
		}
	}
	return result
}
