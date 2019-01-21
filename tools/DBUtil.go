/**
 * Project qualityManagerApi created by exluap
 * Date: 26.10.2018 00:28
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package tools

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"time"
)

func CheckIfUserExist(SIEBEL, PASS string) bool {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	sqlstmt := `SELECT SIEBEL, PASS FROM users WHERE SIEBEL = ? AND PASS = ?`

	err = db.QueryRow(sqlstmt, SIEBEL, PASS).Scan(&SIEBEL, &PASS)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}

		return false
	}

	return true
}

func UserQueries(userId string) ([]byte, error) {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	sqlString := `SELECT queries_list.siebel_login, queries_list.sr_number, queries_list.sr_type, queries_list.time_create, queries.sr_result,queries.overtime FROM (queries_list INNER JOIN queries ON queries_list.sr_number = queries.sr_number) WHERE queries_list.siebel_login = ? AND time_create between strftime('%d.%m.%Y %H:%M',date('now')) AND strftime('%d.%m.%Y %H:%M',datetime('now'), '+3 hours') ORDER BY time_create DESC`

	rows, err := db.Query(sqlString, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v

			if col == "overtime" {
				switch v {
				case "1":
					entry[col] = "Да"
					break
				case "0":
					entry[col] = "Нет"
					break
				}
			}
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return nil, err
	}
	//log.Println("user: " + userId + " response: " + string(jsonData))
	return []byte(jsonData), nil
}

func AddQueryToDB(userId, sr_number, sr_type, sr_result, sr_repeat_result, inform, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, need_other, note string) {
	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	var sqlQuery string

	if checkIfQueryExist(sr_number) {
		sqlQuery = `UPDATE queries SET sr_type = ?, sr_result = ?, sr_repeat_result = ?, no_records = ?, no_records_only = ?, expenditure = ?, more_thing = ?, exp_claim = ?, fin_korr = ?, close_account = ?, unblock_needed = ?, loyatly_needed = ?, phone_denied = ?, due_date_action = ?, due_date_zero = ?, due_date_move = ?, inform = ?, need_other = ?, note = ?, additional_actions="", how_inform = "" WHERE sr_number = ?`
		_, err = db.Exec(sqlQuery, sr_type, sr_result, sr_repeat_result, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, inform, need_other, sr_number, note)
	} else {
		sqlQuery = `INSERT INTO queries (sr_number, sr_type, sr_result, sr_repeat_result, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, inform, need_other, note, additional_actions, how_inform) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?, ?, "", "")`
		_, err = db.Exec(sqlQuery, sr_number, sr_type, sr_result, sr_repeat_result, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, inform, need_other, note)

		if CheckIfUserInOver(userId) {
			sqlQuery = "UPDATE queries SET overtime = 1 WHERE sr_number = ?"
			_, err = db.Exec(sqlQuery, sr_number)
		}
	}

	if err != nil {
		log.Print(err)
	}

	var sr_type_rus string

	switch sr_type {
	case "ko_normal":
		sr_type_rus = "Сотрудник СС / Бизнес"
		break
	case "ko_repeat":
		sr_type_rus = "Рассмотрение КО"
		break
	case "ko_several":
		sr_type_rus = "КО на несколько подразделений"
		break
	case "ko_several_multi":
		sr_type_rus = "Проверка в рамках НП"
		break
	}

	if checkIfQueryExist(sr_number) {
		sqlQuery = `UPDATE queries_list SET sr_type = ? WHERE sr_number = ?`
		_, err = db.Exec(sqlQuery, sr_type_rus, sr_number)
	} else {
		sqlQuery = `INSERT INTO queries_list (time_create, siebel_login, sr_number, sr_type) VALUES (strftime('%d.%m.%Y %H:%M',datetime('now'), '+3 hours'),?,?,?)`
		_, err = db.Exec(sqlQuery, userId, sr_number, sr_type_rus)
	}

	if err != nil {
		log.Print(err)
	}

}

func checkIfQueryExist(sr_number string) bool {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	sqlstmt := `SELECT sr_number FROM queries_list WHERE sr_number = ?`

	err = db.QueryRow(sqlstmt, sr_number).Scan(&sr_number)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}

		return false
	}

	return true
}

func CheckIfUserInOver(user string) bool {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	sqlstmt := `SELECT overtime FROM users WHERE siebel = ? AND overtime = "1"`

	err = db.QueryRow(sqlstmt, user).Scan(&user)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}

		return false
	}

	return true
}

func IneedMoreMoney(user, action string) {
	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	var sqlstmt string

	if action == "1" {
		sqlstmt = `UPDATE users SET overtime = 0 WHERE siebel = ?`
	} else if action == "0" {
		sqlstmt = `UPDATE users SET overtime = 1 WHERE siebel = ?`
	}

	db.Exec(sqlstmt, user)
}

func GetQueryInfo(sr_number string) ([]byte, error) {
	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	sqlString := `SELECT  * FROM queries WHERE sr_number = ?`

	rows, err := db.Query(sqlString, sr_number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return nil, err
	}
	//log.Println("sr_number: " + sr_number + " response: " + string(jsonData))
	return []byte(jsonData), nil
}

func AddNewUser(firstName, lastName, middleName, login string) (string, bool) {
	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	var pass string

	pass = randStringRunes(8)

	if checkIfExistRegister(login) {
		return "", false
	} else {
		sqlQuery := `INSERT INTO users (SIEBEL, PASS, firstName, lastName, middleName) VALUES (?,?,?,?,?)`

		_, err = db.Exec(sqlQuery, login, pass, firstName, lastName, middleName)

		return pass, true
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func checkIfExistRegister(SIEBEL string) bool {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	sqlstmt := `SELECT SIEBEL FROM users WHERE SIEBEL = ?`

	err = db.QueryRow(sqlstmt, SIEBEL).Scan(&SIEBEL)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}

		return false
	}

	return true
}

func DeleteQuery(sr_number, user string) bool {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	querySQL := `DELETE FROM queries_list WHERE siebel_login = ? AND sr_number = ?`

	_, err = db.Exec(querySQL, user, sr_number)

	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

func SaveLog(inter, logText, userName string) {

	db, err := sql.Open("sqlite3", "./goqualityBD.db")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	querySQL := `INSERT INTO logs (inter, logText, user, logTime) VALUES (?,?,?,?)`

	timeCreate := time.Now().Format("02.01.2006 15:04")

	_, err = db.Exec(querySQL, inter, logText, userName, timeCreate)

}
