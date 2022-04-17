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

func Fetch_currHome() (helpers.Response, error) {
	var obj entities.Model_curr
	var arraobj []entities.Model_curr
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcurr , nmcurr,   
			createcurr, to_char(COALESCE(createdatecurr,now()), 'YYYY-MM-DD HH24:MI:SS') as createdatecurr, 
			updatecurr, to_char(COALESCE(updatedatecurr,now()), 'YYYY-MM-DD HH24:MI:SS') as updatedatecurr 
			FROM ` + configs.DB_tbl_mst_currency + ` 
		`

	row, err := con.QueryContext(ctx, sql_select)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcurr_db, nmcurr_db                                               string
			createcurr_db, createdatecurr_db, updatecurr_db, updatedatecurr_db string
		)

		err = row.Scan(&idcurr_db, &nmcurr_db,
			&createcurr_db, &createdatecurr_db, &updatecurr_db, &updatedatecurr_db)

		helpers.ErrorCheck(err)

		create := createcurr_db + ", " + createdatecurr_db
		update := ""
		if updatecurr_db != "" {
			update = updatecurr_db + ", " + updatedatecurr_db
		}
		obj.Curr_idcurr = idcurr_db
		obj.Curr_nama = nmcurr_db
		obj.Curr_create = create
		obj.Curr_update = update
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
func Save_currHome(admin, idcurr, nama, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false
	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_currency, "idcurr", idcurr)
		if !flag {
			sql_insert := `
				insert into 
				` + configs.DB_tbl_mst_currency + ` (
					idcurr , nmcurr,  
					createcurr, createdatecurr
				) values (
					$1, $2, 
					$3, $4
				)
			`

			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_currency, "INSERT",
				idcurr, nama,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				flag = true
				msg = "Succes"
				log.Println(msg_insert)

				notelog := ""
				notelog += "NEW CURRENCY <br>"
				notelog += "CURRENCRY : " + idcurr
				Insert_log("MASTER", "", admin, "CURRENCY", "INSERT", notelog)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
