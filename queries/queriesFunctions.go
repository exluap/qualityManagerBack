/**
 * Project qualityManagerApi created by exluap
 * Date: 22.01.2019 00:27
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package queries

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/auth"
	"qualityManagerApi/models"
	"qualityManagerApi/tools"
)

/**
@api {post} /api/query/list Getting today records
@apiVersion 1.0.0
@apiName GetQueryOfUser
@apiGroup Queries
@apiHeader token Auth Token of user with information about him

@apiParam {String} start Start date with format YYYY-mm-dd hh:mm:ss
@apiParam {String} end End date with format as start

@apiDescription Getting today queries of user

@apiSuccess {Object[]} Queries Array of queries
@apiSuccess {String} Queries.overtime Type of query
@apiSuccess {String} Queries.siebel_login Login of user in DataBase
@apiSuccess {String} Queries.sr_number Special number, not ID
@apiSuccess {String} Queries.sr_result Result of SR
@apiSuccess {String} Queries.sr_type Type of SR
@apiSuccess {String} Queries.time_create Date and Time of create record in DataBase


@apiSuccessExample {json} Success-Response
	[
  {
    "overtime": "Нет",
    "siebel_login": "KOBolotova",
    "sr_number": "541658498",
    "sr_result": "confirm",
    "sr_type": "Обычное КО",
    "time_create": "13.01.2018 00:53"
  },
  {
    "overtime": "Нет",
    "siebel_login": "KOBolotova",
    "sr_number": "5-9653258",
    "sr_result": "confirm",
    "sr_type": "Обычное КО",
    "time_create": "13.01.2018 00:51"
  }
]

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"

*/

func GetQueries(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	data := auth.CheckToken(w, r)

	var queryData map[string]string

	json.Unmarshal(body, &queryData)

	userID := data.CustomClaims["userid"]

	jsonData, err := tools.UserQueries(userID, queryData["start"], queryData["end"])

	if queryData["start"] == "" && queryData["end"] == "" {
		res := &models.Resultation{
			Result: "Not set range. Cant get queries",
		}

		_, err = json.Marshal(res)

		if err != nil {
			log.Print("Error with get queries list")
			raven.CaptureErrorAndWait(err, nil)
			http.Error(w, "Not set range. Cant get queries", http.StatusBadRequest)
		}
	} else {
		_, err = w.Write(jsonData)

		if err != nil {
			log.Print("Auth error: ", err)
			raven.CaptureErrorAndWait(err, nil)
			http.Error(w, "Request is bad", http.StatusBadRequest)
		}
	}

}

/**
@api {post} /api/query/add Adding new record
@apiVersion 1.0.0
@apiName AddQuery
@apiGroup Queries
@apiHeader token Auth Token of user with information about him

@apiDescription Adding new record with info about query to DataBase

@apiParam {String} sr_number Number of SR
@apiParam {String} sr_result Result of SR
@apiParam {String} sr_repeat_result Result of SR if it is repeat SR
@apiParam {String} no_records
@apiParam {String} no_records_only
@apiParam {String} expenditure
@apiParam {String} more_thing
@apiParam {String} exp_claim
@apiParam {String} fin_korr
@apiParam {String} close_account
@apiParam {String} unblock_needed
@apiParam {String} loyatly_needed
@apiParam {String} phone_denied
@apiParam {String} due_date_action
@apiParam {String} due_date_zero
@apiParam {String} due_date_move
@apiParam {String} inform Type of Information
@apiParam {String} sr_type Type of SR
@apiParam {String} need_other
@apiParam {String} additional_actions
@apiParam {String} how_inform Inform instruction
@apiParam {String} note Generated note text

@apiSuccess {String} Result Result of creating
@apiSuccessExample {json} Success-Response
	{
		"Result": "ok"
	}

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"

*/

func AddQuery(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	data := auth.CheckToken(w, r)

	if err != nil {
		http.Error(w, "Request failed!", http.StatusInternalServerError)
	}

	var queryData map[string]string

	json.Unmarshal(body, &queryData)

	tools.AddQueryToDB(data.CustomClaims["userid"], queryData["sr_number"], queryData["sr_type"], queryData["sr_result"], queryData["sr_repeat_result"], queryData["inform"], queryData["no_records"], queryData["no_records_only"], queryData["expenditure"], queryData["more_thing"], queryData["exp_claim"], queryData["fin_korr"], queryData["close_account"], queryData["unblock_needed"], queryData["loyatly_needed"], queryData["phone_denied"], queryData["du_date_action"], queryData["due_date_zero"], queryData["due_date_move"], queryData["need_other"], queryData["note"], queryData["note_sub_1"], queryData["note_sub_2"], queryData["claim_info"], queryData["comm_chat"], queryData["comm_call"], queryData["comm_mail"], queryData["comm_meet"], queryData["comm_nothing"], queryData["communications"])

	res := &models.Resultation{
		Result: "ok",
	}

	_, err = json.Marshal(queryData)

	showResponse, _ := json.Marshal(res)

	w.Write(showResponse)
}

/**
@api {get} /api/query/info Getting record info
@apiName GetQuery
@apiVersion 1.0.0
@apiGroup Queries
@apiHeader token Auth Token of user with information about him

@apiDescription Getting record info by SR number

@apiParam {String} sr_number Number of SR

@apiSuccess {String} additional_actions
@apiSuccess {String} admin_check
@apiSuccess {String} close_account
@apiSuccess {String} due_date_action
@apiSuccess {String} due_date_move
@apiSuccess {String} due_date_zero
@apiSuccess {String} exp_claim
@apiSuccess {String} expenditure
@apiSuccess {String} fin_korr
@apiSuccess {String} how_inform
@apiSuccess {Integer} id
@apiSuccess {String} inform
@apiSuccess {String} loyatly_needed
@apiSuccess {String} more_thing
@apiSuccess {String} need_other
@apiSuccess {String} no_records
@apiSuccess {String} no_records_only
@apiSuccess {String} note
@apiSuccess {String} overtime
@apiSuccess {String} phone_denied
@apiSuccess {String} sr_number
@apiSuccess {String} sr_repeat_result
@apiSuccess {String} sr_result
@apiSuccess {String} sr_type
@apiSuccess {String} unblock_needed

@apiSuccessExample {json} Success-Response
[
  {
    "additional_actions": "",
    "admin_check": 0,
    "close_account": "0",
    "due_date_action": "",
    "due_date_move": "0",
    "due_date_zero": "0",
    "exp_claim": "0",
    "expenditure": "0",
    "fin_korr": "0",
    "how_inform": "",
    "id": 1344,
    "inform": "",
    "loyatly_needed": "0",
    "more_thing": "0",
    "need_other": "0",
    "no_records": "0",
    "no_records_only": "0",
    "note": "ОТВЕТ ПО ПРЕТЕНЗИИ: \nКоллеги все проверили и выяснили, что действительно была ошибка. Мы сделаем все возможное,чтобы такое больше не повторилось, приносим извинения.\n\nДЕТАЛИ РАЗБОРА: \nДата, время, тип коммуникации: \n < > \nСуть претензии:\n < > wsadsadsad\nРезультат проверки: \n  < ___ >",
    "overtime": "0",
    "phone_denied": "0",
    "sr_number": "123455667890",
    "sr_repeat_result": "",
    "sr_result": "confirm",
    "sr_type": "ko_normal",
    "unblock_needed": "0"
  }
]

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"


*/

func GetQuery(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	var queryInfo map[string]string

	json.Unmarshal(body, &queryInfo)

	jsonData, err := tools.GetQueryInfo(queryInfo["sr_number"])

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusInternalServerError)
	}

	w.Write(jsonData)

}

/**
@api {post} /api/query/delete Delete record from list
@apiName PostDeleteSR
@apiVersion 1.0.0
@apiGroup Queries

@apiHeader token Auth Token of user with information about him

@apiDescription Delete record about SR from list

@apiParam {String} sr_number SR value in list

@apiSuccessExample {json} Success-Response
	Delete

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"

@apiError Can't delete SR
@apiErrorExample Error-Response (delete)
	"Huston we have a problem. See a log!"
*/

func DeleteSR(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	var requestBody map[string]string

	json.Unmarshal(body, &requestBody)

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := auth.CheckToken(w, r)

	userID := data.CustomClaims["userid"]

	result := tools.DeleteQuery(requestBody["sr_number"], userID)

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	if result {
		log.Println("Deleted SR " + requestBody["sr_body"])
		w.Write([]byte("Delete"))
	} else {
		log.Println("NOT Deleted SR PANIC" + requestBody["sr_body"])
		w.Write([]byte("Huston we have a problem. See a log!"))
	}

}
