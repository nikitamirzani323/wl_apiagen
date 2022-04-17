package models

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apiagen/configs"
	"github.com/nikitamirzani323/wl_apiagen/db"
	"github.com/nikitamirzani323/wl_apiagen/entities"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
	"github.com/nleeper/goment"
)

func Fetch_agenHome(idcompany string) (helpers.ResponseAgen, error) {
	var obj entities.Model_agen
	var arraobj []entities.Model_agen
	var res helpers.ResponseAgen
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idagen , to_char(COALESCE(startjoinagen,now()), 'YYYY-MM-DD HH24:MI:SS'), to_char(COALESCE(endjoinagen,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			idcurr , nmagen, owneragen, owneragen_phone, owneragen_email, endpointagenurl, statusagen, 
			createagen, to_char(COALESCE(createdateagen,now()), 'YYYY-MM-DD HH24:MI:SS') as createdateagen, 
			updateagen, to_char(COALESCE(updatedateagen,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedateagen 
			FROM ` + configs.DB_tbl_mst_Agen + ` 
			WHERE idcompany=$1 
			ORDER BY createagen DESC 
		`

	row, err := con.QueryContext(ctx, sql_select, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idagen_db, startjoinagen_db, endjoinagen_db                                                                   string
			idcurr_db, nmagen_db, owneragen_db, owneragen_phone_db, owneragen_email_db, endpointagenurl_db, statusagen_db string
			createagen_db, createdateagen_db, updateagen_db, updatedateagen_db                                            string
		)

		err = row.Scan(
			&idagen_db, &startjoinagen_db, &endjoinagen_db,
			&idcurr_db, &nmagen_db, &owneragen_db, &owneragen_phone_db, &owneragen_email_db, &endpointagenurl_db, &statusagen_db,
			&createagen_db, &createdateagen_db, &updateagen_db, &updatedateagen_db)

		helpers.ErrorCheck(err)
		if startjoinagen_db == "0000-00-00 00:00:00" {
			startjoinagen_db = ""
		}
		if endjoinagen_db == "0000-00-00 00:00:00" {
			endjoinagen_db = ""
		}
		create := createagen_db + ", " + createdateagen_db
		update := ""
		if updateagen_db != "" {
			update = updateagen_db + ", " + updatedateagen_db
		}
		obj.Agen_idagen = idagen_db
		obj.Agen_startjoin = startjoinagen_db
		obj.Agen_endjoin = endjoinagen_db
		obj.Agen_idcurr = idcurr_db
		obj.Agen_nmagen = nmagen_db
		obj.Agen_ownername = owneragen_db
		obj.Agen_ownerphone = owneragen_phone_db
		obj.Agen_owneremail = owneragen_email_db
		obj.Agen_urlendpoint = endpointagenurl_db
		obj.Agen_status = statusagen_db
		obj.Agen_create = create
		obj.Agen_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	var objCurr entities.Model_agencurr
	var arraobjCurr []entities.Model_agencurr
	sql_listcurr := `SELECT 
		idcurr 	
		FROM ` + configs.DB_tbl_mst_currency + ` 
	`
	row_listcurr, err_listcurr := con.QueryContext(ctx, sql_listcurr)
	helpers.ErrorCheck(err_listcurr)
	for row_listcurr.Next() {
		var (
			idcurr_db string
		)

		err = row_listcurr.Scan(&idcurr_db)

		helpers.ErrorCheck(err)

		objCurr.Curr_idcurr = idcurr_db
		arraobjCurr = append(arraobjCurr, objCurr)
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Listcurr = arraobjCurr
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_agenHome(
	admin, idcompany, idagen, idcurr, nmagen,
	owneragen, ownerphone, owneremail, endpointurl,
	status, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_Agen, "idagen", idagen)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_Agen + ` (
					idagen, idcompany , startjoinagen, idcurr, 
					nmagen, owneragen, owneragen_phone, owneragen_email, endpointagenurl, statusagen, 
					createagen, createdateagen
				) values (
					$1, $2, $3, $4,   
					$5, $6, $7, $8, $9, $10, 
					$11, $12  
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_Agen, "INSERT",
				idagen, idcompany, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcurr,
				nmagen, owneragen, ownerphone, owneremail, endpointurl, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)

				notelog := ""
				notelog += "NEW AGEN <br>"
				notelog += "IDAGEN : " + idagen + "<br>"
				notelog += "CURRENCY : " + idcurr + "<br>"
				notelog += "AGEN : " + nmagen + "<br>"
				notelog += "OWNER : " + owneragen + "<br>"
				notelog += "PHONE : " + ownerphone + "<br>"
				notelog += "EMAIL : " + owneremail + "<br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "AGEN", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_Agen + `   
				SET nmagen=$1, owneragen=$2, owneragen_phone=$3, owneragen_email=$4,  
				endpointagenurl =$5, statusagen=$6,  
				updateagen=$7, updatedateagen=$8 
				WHERE idcompany=$9 
				AND idagen=$10  
			`
		flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_Agen, "UPDATE",
			nmagen, owneragen, ownerphone, owneremail, endpointurl, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idcompany, idagen)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE AGEN <br>"
			notelog += "CURRENCY : " + idcurr + "<br>"
			notelog += "AGEN : " + nmagen + "<br>"
			notelog += "OWNER : " + owneragen + "<br>"
			notelog += "PHONE : " + ownerphone + "<br>"
			notelog += "EMAIL : " + owneremail + "<br>"
			notelog += "STATUS : " + status
			Insert_log("MASTER", idcompany, admin, "AGEN", "UPDATE", notelog)
		} else {
			log.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Fetch_agenListAdmin(idcompany, idagen string) (helpers.Response, error) {
	var obj entities.Model_agenadmin
	var arraobj []entities.Model_agenadmin
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			username_agen , typeadmin_agen, nama_agen, phone_agen, email_agen,  
			status_agen , to_char(COALESCE(lastlogin_agen,now()), 'YYYY-MM-DD HH24:MI:SS'), lastipaddres_agen, 
			createadmin_agen, to_char(COALESCE(createdateadmin_agen,now()), 'YYYY-MM-DD HH24:MI:SS') as createdateadmin_agen, 
			updateadmin_agen, to_char(COALESCE(updatedateadmin_agen,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedateadmin_agen 
			FROM ` + configs.DB_tbl_mst_Agen_admin + ` 
			WHERE idagen = $1 
			AND idcompany = $2  
			ORDER BY createadmin_agen DESC 
		`

	row, err := con.QueryContext(ctx, sql_select, idagen, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_agen_db, typeadmin_agen_db, nama_agen_db, phone_agen_db, email_agen_db            string
			status_agen_db, lastlogin_agen_db, lastipaddres_agen_db                                    string
			createadmin_agen_db, createdateadmin_agen_db, updateadmin_agen_db, updatedateadmin_agen_db string
		)

		err = row.Scan(
			&username_agen_db, &typeadmin_agen_db, &nama_agen_db, &phone_agen_db, &email_agen_db,
			&status_agen_db, &lastlogin_agen_db, &lastipaddres_agen_db,
			&createadmin_agen_db, &createdateadmin_agen_db, &updateadmin_agen_db, &updatedateadmin_agen_db)

		helpers.ErrorCheck(err)
		if lastlogin_agen_db == "0000-00-00 00:00:00" {
			lastlogin_agen_db = ""
		}
		create := createadmin_agen_db + ", " + createdateadmin_agen_db
		update := ""
		if updateadmin_agen_db != "" {
			update = updateadmin_agen_db + ", " + updatedateadmin_agen_db
		}
		obj.Agenadmin_username = username_agen_db
		obj.Agenadmin_type = typeadmin_agen_db
		obj.Agenadmin_name = nama_agen_db
		obj.Agenadmin_phone = phone_agen_db
		obj.Agenadmin_email = email_agen_db
		obj.Agenadmin_status = status_agen_db
		obj.Agenadmin_lastlogin = lastlogin_agen_db
		obj.Agenadmin_lastipaddress = lastipaddres_agen_db
		obj.Agenadmin_create = create
		obj.Agenadmin_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_agenlistadmin(
	admin, idcompany, idagen, username, password,
	name, email, phone, status,
	sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_Agen_admin, "username_agen", username)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_Agen_admin + ` (
					username_agen, password_agen, idagen, idcompany, typeadmin_agen, idruleadmin_agen, 
					nama_agen, phone_agen, email_agen, status_agen,   
					createadmin_agen, createdateadmin_agen  
				) values (
					$1, $2, $3, $4, $5, $6,   
					$7, $8, $9, $10, 
					$11, $12 
					  
				)
			`
			hashpwd := helpers.HashPasswordMD5(password)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_Agen_admin, "INSERT",
				username, hashpwd, idagen, idcompany, "MASTER", "0",
				name, phone, email, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)

				notelog := ""
				notelog += "NEW ADMIN - AGEN <br>"
				notelog += "USERNAME : " + username + "<br>"
				notelog += "TYPE : MASTER<br>"
				notelog += "NAME : " + name + "<br>"
				notelog += "PHONE : " + phone + "<br>"
				notelog += "EMAIL : " + email + "<br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "AGEN ADMIN", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password != "" {
			haspwd := helpers.HashPasswordMD5(password)
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_Agen_admin + `   
				SET password_agen=$1, nama_agen=$2, phone_agen=$3, email_agen=$4,  
				status_agen=$5, updateadmin_agen=$6, updatedateadmin_agen=$7 
				WHERE username_agen=$8  
				AND idagen=$9   
				AND idcompany=$10   
			`
			flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_Agen_admin, "UPDATE",
				haspwd, name, phone, email, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), username, idagen, idcompany)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)

				notelog := ""
				notelog += "UPDATE ADMIN - AGEN - CHANGE PASSWORD<br>"
				notelog += "USERNAME : " + username + "<br>"
				notelog += "NAME : " + name + "<br>"
				notelog += "PHONE : " + phone + "<br>"
				notelog += "EMAIL : " + email + "<br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "AGEN ADMIN", "UPDATE", notelog)
			} else {
				log.Println(msg_update)
			}
		} else {
			sql_update2 := `
				UPDATE 
				` + configs.DB_tbl_mst_Agen_admin + `   
				SET nama_agen=$1, phone_agen=$2, email_agen=$3,  
				status_agen=$4, updateadmin_agen=$5, updatedateadmin_agen=$6  
				WHERE username_agen=$7  
				AND idagen=$8   
				AND idcompany=$9    
			`
			flag_update, msg_update := Exec_SQL(sql_update2, configs.DB_tbl_mst_Agen_admin, "UPDATE",
				name, phone, email, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), username, idagen, idcompany)

			if flag_update {
				flag = true
				msg = "Succes"
				log.Println(msg_update)

				notelog := ""
				notelog += "UPDATE ADMIN - AGEN - CHANGE PASSWORD<br>"
				notelog += "USERNAME : " + username + "<br>"
				notelog += "NAME : " + name + "<br>"
				notelog += "PHONE : " + phone + "<br>"
				notelog += "EMAIL : " + email + "<br>"
				notelog += "STATUS : " + status
				Insert_log("MASTER", idcompany, admin, "AGEN ADMIN", "UPDATE", notelog)
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
