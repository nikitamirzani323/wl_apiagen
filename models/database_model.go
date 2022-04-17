package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/nikitamirzani323/wl_apiagen/configs"
	"github.com/nikitamirzani323/wl_apiagen/db"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
	"github.com/nleeper/goment"
)

func Get_counter(field_column string) int {
	con := db.CreateCon()
	ctx := context.Background()
	idrecord_counter := 0
	sqlcounter := `SELECT 
					counter 
					FROM ` + configs.DB_tbl_counter + ` 
					WHERE nmcounter = $1 
				`
	var counter int = 0
	row := con.QueryRowContext(ctx, sqlcounter, field_column)
	switch e := row.Scan(&counter); e {
	case sql.ErrNoRows:
		log.Println("COUNTER - No rows were returned!")
	case nil:
		log.Println(counter)
	default:
		panic(e)
	}
	if counter > 0 {
		idrecord_counter = int(counter) + 1
		stmt, e := con.PrepareContext(ctx, "UPDATE "+configs.DB_tbl_counter+" SET counter=$1 WHERE nmcounter=$2 ")
		helpers.ErrorCheck(e)
		res, e := stmt.ExecContext(ctx, idrecord_counter, field_column)
		helpers.ErrorCheck(e)
		a, e := res.RowsAffected()
		helpers.ErrorCheck(e)
		if a > 0 {
			log.Println("COUNTER - UPDATE")
		}
	} else {
		stmt, e := con.PrepareContext(ctx, "insert into "+configs.DB_tbl_counter+" (nmcounter, counter) values ($1, $2)")
		helpers.ErrorCheck(e)
		res, e := stmt.ExecContext(ctx, field_column, 1)
		helpers.ErrorCheck(e)
		id, e := res.RowsAffected()
		helpers.ErrorCheck(e)
		log.Println("Insert id", id)
		log.Println("NEW")
		idrecord_counter = 1
	}
	return idrecord_counter
}
func Get_listitemsearch(data, pemisah, search string) bool {
	flag := false
	temp := strings.Split(data, pemisah)
	for i := 0; i < len(temp); i++ {
		if temp[i] == search {
			flag = true
			break
		}
	}
	return flag
}
func CheckDB(table, field, value string) bool {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	sql_db := `SELECT 
					` + field + ` 
					FROM ` + table + ` 
					WHERE ` + field + ` = $1 
				`
	row := con.QueryRowContext(ctx, sql_db, value)
	switch e := row.Scan(&field); e {
	case sql.ErrNoRows:
		log.Println("CHECK DB - No rows were returned!")
		flag = false
	case nil:
		flag = true
	default:
		flag = false
	}
	return flag
}
func CheckDBTwoField(table, field_1, value_1, field_2, value_2 string) bool {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	sql_db := `SELECT 
					` + field_1 + ` 
					FROM ` + table + ` 
					WHERE ` + field_1 + ` = $1 
					AND ` + field_2 + ` = $2
				`
	log.Println(sql_db)
	row := con.QueryRowContext(ctx, sql_db, value_1, value_2)
	switch e := row.Scan(&field_1); e {
	case sql.ErrNoRows:
		log.Println("CHECKDBTWOFIELD - No rows were returned!")
		flag = false
	case nil:
		flag = true
	default:
		flag = false
	}
	return flag
}
func Get_AdminRule(tipe, idcompany string, idrule int) string {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	result := ""
	rulecomp := ""

	sql_select := `SELECT
		rulecomp  
		FROM ` + configs.DB_tbl_mst_Company_adminrule + `  
		WHERE idcomprule = $1 
		AND idcompany = $2 
	`
	row := con.QueryRowContext(ctx, sql_select, idrule, idcompany)
	switch e := row.Scan(&rulecomp); e {
	case sql.ErrNoRows:
		flag = false
	case nil:
		flag = true

	default:
		panic(e)
	}
	if flag {
		switch tipe {
		case "rulecomp":
			result = rulecomp
		}
	}
	return result
}
func Delete_SQL(sql, table string, args ...interface{}) bool {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	stmt_delete, e_delete := con.PrepareContext(ctx, sql)
	helpers.ErrorCheck(e_delete)
	defer stmt_delete.Close()
	rec_delete, e_delete := stmt_delete.ExecContext(ctx, args...)

	helpers.ErrorCheck(e_delete)
	deletesource, e := rec_delete.RowsAffected()
	helpers.ErrorCheck(e)
	if deletesource > 0 {
		flag = true
		log.Printf("Data %s Berhasil di delete", table)
	} else {
		log.Printf("Data %s Failed di delete", table)
	}
	return flag
}
func Exec_SQL(sql, table, action string, args ...interface{}) (bool, string) {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	msg := ""
	stmt_exec, e_exec := con.PrepareContext(ctx, sql)
	helpers.ErrorCheck(e_exec)
	defer stmt_exec.Close()
	rec_exec, e_exec := stmt_exec.ExecContext(ctx, args...)

	helpers.ErrorCheck(e_exec)
	exec, e := rec_exec.RowsAffected()
	helpers.ErrorCheck(e)
	if exec > 0 {
		flag = true
		msg = "Data " + table + " Berhasil di " + action
	} else {
		msg = "Data " + table + " Failed di " + action
	}
	return flag, msg
}
func Insert_log(typeuser, idcompany, username, page, tipe, note string) {
	tglnow, _ := goment.New()
	sql_insert := `
		INSERT INTO 
		` + configs.DB_tbl_trx_log + ` (
			idlog, yearlog, 
			typeuser, company, userlog, pagelog, tipelog, notelog
		) VALUES (
			$1, $2, 
			$3, $4, $5, $6, $7, $8 
		)
	`

	year := tglnow.Format("YYYY")
	month := tglnow.Format("MM")
	field_col := configs.DB_tbl_trx_log + year + month
	idlog_counter := Get_counter(field_col)
	idlog := tglnow.Format("YY") + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idlog_counter)
	flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_log, "INSERT",
		idlog, year,
		typeuser, idcompany, username, page, tipe, note)
	if flag_insert {
		log.Println(msg_insert)
	} else {
		log.Println(msg_insert)
	}

}
