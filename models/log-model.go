package models

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_apiagen/configs"
	"github.com/nikitamirzani323/wl_apiagen/db"
	"github.com/nikitamirzani323/wl_apiagen/entities"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
)

func Fetch_logHome(typeuser, idcompany string) (helpers.Response, error) {
	var obj entities.Model_log
	var arraobj []entities.Model_log
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idlog, to_char(COALESCE(datetimelog,now()), 'YYYY-MM-DD HH24:MI:SS'), userlog, pagelog,  
			tipelog, notelog    
			FROM ` + configs.DB_tbl_trx_log + ` 
			WHERE typeuser=$1 AND company=$2 
			ORDER BY datetimelog DESC LIMIT 300 
		`

	row, err := con.QueryContext(ctx, sql_select, typeuser, idcompany)

	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idlog_db                                                       int
			datetimelog_db, userlog_db, pagelog_db, tipelog_db, notelog_db string
		)

		err = row.Scan(&idlog_db, &datetimelog_db, &userlog_db, &pagelog_db, &tipelog_db, &notelog_db)

		helpers.ErrorCheck(err)

		obj.Log_id = idlog_db
		obj.Log_datetime = datetimelog_db
		obj.Log_user = userlog_db
		obj.Log_page = pagelog_db
		obj.Log_tipe = tipelog_db
		obj.Log_note = notelog_db
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
