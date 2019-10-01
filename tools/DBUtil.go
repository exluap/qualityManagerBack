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
	"errors"
	"github.com/getsentry/raven-go"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func connectToDb() (*sql.DB, error) {
	conn, err := sql.Open("mysql", "quality_manager:1KTeMi7ZTKQ3LBSy@/quality_manager")

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Print(err)
	}

	return conn, err

}

func CheckIfUserExist(SIEBEL, PASS string) bool {

	db, err := connectToDb()

	sqlstmt := `SELECT SIEBEL, PASS FROM users WHERE SIEBEL = ? AND PASS = ?`

	err = db.QueryRow(sqlstmt, SIEBEL, PASS).Scan(&SIEBEL, &PASS)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}

		defer db.Close()

		return false
	}

	defer db.Close()
	return true
}

func UserQueries(userId, time_start, time_end string) ([]byte, error) {

	db, err := connectToDb()

	sqlString := `SELECT queries_list.siebel_login, queries_list.sr_number, queries_list.sr_type, queries_list.time_create, queries.sr_result,queries.overtime FROM (queries_list INNER JOIN queries ON queries_list.sr_number = queries.sr_number) WHERE queries_list.siebel_login = ? AND time_create between ? AND ? ORDER BY time_create DESC`

	rows, err := db.Query(sqlString, userId, time_start, time_end)
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
	defer db.Close()
	return []byte(jsonData), nil
}

func GetTaskInfo(taskId string) ([]byte, error) {

	db, err := connectToDb()

	sqlString := `SELECT * FROM tasks WHERE id = ?`

	rows, err := db.Query(sqlString, taskId)
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
	defer db.Close()
	return jsonData, nil
}

func UpdateTaskStatus(taskId, status, owner, assegnee string) bool {
	db, err := connectToDb()

	var sqlString string

	if status != "Closed" && status != "Canceled" {
		sqlString = `UPDATE tasks SET status = ?, assegnee = ? WHERE id = ? AND owner = ?`
	} else {
		sqlString = `UPDATE tasks SET status = ?, assegnee = ?, close_time = CURRENT_TIMESTAMP WHERE id = ? AND owner = ?`
	}

	_, err = db.Exec(sqlString, status, assegnee, taskId, owner)

	if err != nil {
		log.Print("Update task status error: ", err)
		defer db.Close()
		return false
	} else {
		defer db.Close()
		return true
	}
}

func UpdateTaskInfo(taskId string, taskInfo map[string]string) bool {
	db, err := connectToDb()

	var sqlString string

	sqlString = "UPDATE tasks SET type = ?, assegnee = ?, parent_sr = ?, contact_id = ?, account_id = ?, phone_number = ?, info = ? WHERE id = ?"

	_, err = db.Exec(sqlString, taskInfo["type"], taskInfo["assegnee"], taskInfo["parent_sr"], taskInfo["contact_id"], taskInfo["account_id"], taskInfo["phone_number"], taskInfo["info"], taskId)

	if err != nil {
		log.Print("Update task info error: ", err)
		defer db.Close()
		return false
	} else {
		defer db.Close()
		return true
	}
}

func PostNewTask(taskInfo map[string]string, owner string) bool {
	db, err := connectToDb()

	sqlString := `INSERT INTO tasks (type, owner, parent_sr, contact_id, account_id,phone_number,info, status) VALUES (?,?,?,?,?,?,?, 'Not assegnee')`

	_, err = db.Exec(sqlString, taskInfo["type"], owner, taskInfo["parent_sr"], taskInfo["contact_id"], taskInfo["account_id"], taskInfo["phone_number"], taskInfo["info"])

	if err != nil {
		log.Print("Task insert error: ", err)
		defer db.Close()
		return false
	} else {
		defer db.Close()
		return true
	}
}

func GetTasksByUserOwner(user string) ([]byte, error) {

	db, err := connectToDb()

	sqlString := `SELECT * FROM tasks WHERE owner = ?`

	rows, err := db.Query(sqlString, user)
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
	defer db.Close()
	return []byte(jsonData), nil
}

func GetTasksByAssegneeToUser(user string) ([]byte, error) {

	db, err := connectToDb()

	sqlString := `SELECT * FROM tasks WHERE assegnee = ?`

	rows, err := db.Query(sqlString, user)
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
	defer db.Close()
	return []byte(jsonData), nil
}

func ListOfTasks() ([]byte, error) {

	db, err := connectToDb()

	sqlString := `SELECT * FROM tasks`

	rows, err := db.Query(sqlString)
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
	defer db.Close()
	return []byte(jsonData), nil
}

func GetUserGroups(login string) string {
	db, err := connectToDb()

	sqlString := `SELECT user_group FROM users WHERE SIEBEL = ?`

	row := db.QueryRow(sqlString, login)

	var groups string

	err = row.Scan(&groups)

	if err != nil {
		log.Print("User Func error: ", err)
	}

	defer db.Close()

	return groups

}

func GetTaskOwner(taskID string) string {
	db, err := connectToDb()

	sqlString := `SELECT owner FROM tasks WHERE id = ?`

	row := db.QueryRow(sqlString, taskID)

	var owner string

	err = row.Scan(&owner)

	if err != nil {
		log.Print("Get owner Func error: ", err)
	}

	defer db.Close()

	return owner

}

func GetTaskStatus(taskID string) string {
	db, err := connectToDb()

	sqlString := `SELECT status FROM tasks WHERE id = ?`

	row := db.QueryRow(sqlString, taskID)

	var status string

	err = row.Scan(&status)

	if err != nil {
		log.Print("Get status Func error: ", err)
	}

	defer db.Close()

	return status

}

func AddQueryToDB(userId, sr_number, sr_type, sr_result, sr_repeat_result, inform, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, need_other, note, note_sub_1, note_sub_2, claim_info, comm_chat, comm_call, comm_mail, comm_meet, comm_nothing, communications string) {
	db, err := connectToDb()

	var sqlQuery string

	if checkIfQueryExist(sr_number) {
		sqlQuery = `UPDATE queries SET sr_type = ?, sr_result = ?, sr_repeat_result = ?, no_records = ?, no_records_only = ?, expenditure = ?, more_thing = ?, exp_claim = ?, fin_korr = ?, close_account = ?, unblock_needed = ?, loyatly_needed = ?, phone_denied = ?, due_date_action = ?, due_date_zero = ?, due_date_move = ?, inform = ?, need_other = ?, note = ?, note_sub_1 = ?, note_sub_2 = ?, claim_info = ?, comm_chat = ?, comm_call = ?, comm_mail = ?, comm_meet = ?, comm_nothing = ?, communications = ?  WHERE sr_number = ?`
		_, err = db.Exec(sqlQuery, sr_type, sr_result, sr_repeat_result, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, inform, need_other, note, note_sub_1, note_sub_2, claim_info, comm_chat, comm_call, comm_mail, comm_meet, comm_nothing, communications, sr_number)
	} else {
		sqlQuery = `INSERT INTO queries (sr_number, sr_type, sr_result, sr_repeat_result, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, inform, need_other, note, note_sub_1, note_sub_2, claim_info,comm_chat, comm_call,comm_mail, comm_meet, comm_nothing, communications) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(sqlQuery, sr_number, sr_type, sr_result, sr_repeat_result, no_records, no_records_only, expenditure, more_thing, exp_claim, fin_korr, close_account, unblock_needed, loyatly_needed, phone_denied, due_date_action, due_date_zero, due_date_move, inform, need_other, note, note_sub_1, note_sub_2, claim_info, comm_chat, comm_call, comm_mail, comm_meet, comm_nothing, communications)

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
		sqlQuery = `INSERT INTO queries_list (time_create, siebel_login, sr_number, sr_type) VALUES (CURRENT_TIMESTAMP,?,?,?)`
		_, err = db.Exec(sqlQuery, userId, sr_number, sr_type_rus)
	}

	if err != nil {
		log.Print(err)

	}

	defer db.Close()

}

func checkIfQueryExist(sr_number string) bool {

	db, err := connectToDb()

	sqlstmt := `SELECT sr_number FROM queries_list WHERE sr_number = ?`

	err = db.QueryRow(sqlstmt, sr_number).Scan(&sr_number)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)

		}

		defer db.Close()

		return false
	}

	defer db.Close()

	return true
}

func CheckIfUserInOver(user string) bool {

	db, err := connectToDb()

	sqlstmt := `SELECT overtime FROM users WHERE siebel = ? AND overtime = "1"`

	err = db.QueryRow(sqlstmt, user).Scan(&user)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)

		}

		defer db.Close()

		return false
	}

	defer db.Close()

	return true
}

func CheckAdminMode(user string) bool {
	db, err := connectToDb()

	sqlstmt := `SELECT admin FROM users WHERE siebel = ? AND admin = "1"`

	err = db.QueryRow(sqlstmt, user).Scan(&user)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)

		}

		defer db.Close()

		return false
	}

	defer db.Close()

	return true
}

func IneedMoreMoney(user, action string) {
	db, err := connectToDb()

	var sqlstmt string

	if action == "1" {
		sqlstmt = `UPDATE users SET overtime = 0 WHERE siebel = ?`
	} else if action == "0" {
		sqlstmt = `UPDATE users SET overtime = 1 WHERE siebel = ?`
	}

	db.Exec(sqlstmt, user)

	if err != nil {

		log.Print(err)
	}

	defer db.Close()
}

func GetQueryInfo(sr_number string) ([]byte, error) {
	db, err := connectToDb()

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
	defer db.Close()
	return []byte(jsonData), nil
}

func AddNewUser(firstName, lastName, middleName, login, password, groups, isAdmin, winlogin string) bool {
	db, err := connectToDb()

	if err != nil {

		log.Print(err)
	}

	if CheckIfExistRegister(login) {
		defer db.Close()
		return false
	} else {
		sqlQuery := `INSERT INTO users (SIEBEL, PASS, firstName, lastName, middleName, user_group, admin, winlogin) VALUES (?,?,?,?,?,?,?,?)`

		_, err = db.Exec(sqlQuery, login, password, firstName, lastName, middleName, groups, isAdmin, winlogin)
		defer db.Close()
		return true
	}
}

func UpdateUserInfo(groups, isAdmin, login string) bool {
	db, err := connectToDb()

	if err != nil {

		log.Print(err)
	}

	if CheckIfExistRegister(login) {
		sqlQuery := `UPDATE users SET admin = ?, user_group = ? WHERE SIEBEL = ?`

		_, err = db.Exec(sqlQuery, isAdmin, groups, login)
		defer db.Close()
		return true

	} else {
		defer db.Close()
		return false
	}
}

func CheckIfExistRegister(SIEBEL string) bool {

	db, err := connectToDb()

	sqlstmt := `SELECT SIEBEL FROM users WHERE SIEBEL = ?`

	err = db.QueryRow(sqlstmt, SIEBEL).Scan(&SIEBEL)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)

		}
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func validUser(SIEBEL, PASS string) bool {

	db, err := connectToDb()

	sqlstmt := `SELECT SIEBEL FROM users WHERE SIEBEL = ? AND PASS = ?`

	err = db.QueryRow(sqlstmt, SIEBEL, PASS).Scan(&SIEBEL)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)

		}
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func DeleteQuery(sr_number, user string) bool {
	db, err := connectToDb()

	querySQL := `DELETE FROM queries_list WHERE siebel_login = ? AND sr_number = ?`

	_, err = db.Exec(querySQL, user, sr_number)

	if err != nil {
		log.Println(err)

		defer db.Close()
		return false
	}

	defer db.Close()
	return true

}

func SaveLog(inter, logText, userName string) error {

	db, err := connectToDb()

	querySQL := `INSERT INTO logs (inter, logText, user, logTime) VALUES (?,?,?,?)`

	timeCreate := time.Now().Format("02.01.2006 15:04")

	_, err = db.Exec(querySQL, inter, logText, userName, timeCreate)

	if err != nil {

		log.Print(err)
	}
	defer db.Close()

	return err

}

func ChangeUserPassword(user, passwordold, newpassword string) error {

	db, err := connectToDb()

	if validUser(user, passwordold) {
		querySql := `UPDATE users SET PASS = ? WHERE SIEBEL = ? AND PASS = ?`

		_, err = db.Exec(querySql, newpassword, user, passwordold)
	} else {
		err = errors.New("User not valid")
	}

	defer db.Close()

	return err

}

func ChangeUserLogin(oldLogin, newLogin string) error {

	db, err := connectToDb()

	querySql := `UPDATE users SET SIEBEL = ? WHERE SIEBEL = ?`

	_, err = db.Exec(querySql, newLogin, oldLogin)

	defer db.Close()
	return err

}
