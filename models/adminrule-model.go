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

func Fetch_adminruleHome(idcompany string) (helpers.Response, error) {
	var obj entities.Model_adminruleall
	var arraobj []entities.Model_adminruleall
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcomprule , nmcomprule , rulecomp ,  
			createcomprule, to_char(COALESCE(createdatecomprule,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecomprule, 
			updatecomprule, to_char(COALESCE(updatedatecomprule,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecomprule 
			FROM ` + configs.DB_tbl_mst_Company_adminrule + ` 
			WHERE idcompany = $1 
			ORDER BY createcomprule DESC 
		`

	row, err := con.QueryContext(ctx, sql_select, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcomprule_db                                                                      int
			nmcomprule_db, rulecomp_db                                                         string
			createcomprule_db, createdatecomprule_db, updatecomprule_db, updatedatecomprule_db string
		)

		err = row.Scan(&idcomprule_db, &nmcomprule_db, &rulecomp_db,
			&createcomprule_db, &createdatecomprule_db, &updatecomprule_db, &updatedatecomprule_db)

		helpers.ErrorCheck(err)

		create := createcomprule_db + ", " + createdatecomprule_db
		update := ""
		if updatecomprule_db != "" {
			update = updatecomprule_db + ", " + updatedatecomprule_db
		}

		obj.Adminrule_idrule = idcomprule_db
		obj.Adminrule_nmrule = nmcomprule_db
		obj.Adminrule_rule = rulecomp_db
		obj.Adminrule_create = create
		obj.Adminrule_update = update
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
func Save_adminrule(admin, idcompany, name, rule, sData string, idrule int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_mst_Company_adminrule + ` (
					idcomprule, idcompany, nmcomprule, 
					createcomprule, createdatecomprule  
				) values (
					$1, $2, $3, 
					$4, $5  
				) 
			`
		idcounter := Get_counter(configs.DB_tbl_mst_Company_adminrule + "_" + idcompany)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_Company_adminrule, "INSERT",
			idcounter, idcompany, name, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)

			notelog := ""
			notelog += "NEW RULE <br>"
			notelog += "RULE : " + name
			Insert_log("MASTER", idcompany, admin, "ADMIN RULE", "INSERT", notelog)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_Company_adminrule + `   
				SET rulecomp=$1, updatecomprule=$2, updatedatecomprule=$3 
				WHERE idcomprule=$4 
				AND idcompany=$5  
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_Company_adminrule, "UPDATE",
			rule, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrule, idcompany)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)

			notelog := ""
			notelog += "UPDATE RULE <br>"
			notelog += "RULE : " + name
			Insert_log("MASTER", idcompany, admin, "ADMIN RULE", "UPDATE", notelog)
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
