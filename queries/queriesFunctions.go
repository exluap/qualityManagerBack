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
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/raven-go"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/auth"
	"qualityManagerApi/constants"
	"qualityManagerApi/models"
	"qualityManagerApi/tools"
	"strings"
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
@api {post} /api/query/helper Generate Note text and get instructions
@apiName GetInstructions
@apiVersion 1.0.0
@apiGroup Queries
@apiHeader token Auth Token of user with information about him

@apiDescription Generate Note text and instruction for query with special param

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


@apiSuccess {String} Note Note text for Query
@apiSuccess {String} InfoInstr Information instruction
@apiSuccess {String} AdditionalAction Additional action instruction

@apiSuccessExample {json} Success-Response
{
  "Note": "ОТВЕТ ПО ПРЕТЕНЗИИ: \nКоллеги все проверили и выяснили, что действительно была ошибка. Мы сделаем все возможное,чтобы такое больше не повторилось, приносим извинения.\n Доп. инфо: \n<УКАЖИ результат решения вопроса с расторжением и корректировками>\n\nДЕТАЛИ РАЗБОРА: \nДата, время, тип коммуникации: \n < > \nСуть претензии:\n < > \nРезультат проверки: \n  < ___ >",
  "InfoInstr": "Выбери подстатус - Обоснована. Фин. кор-ки.\n \n \n ПОРЯДОК ИНФОРМИРОВАНИЯ \n \n Зарегистрируй SR Исходящий e-mail - Информирование \n -В ПО запроса напиши текст ответа клиенту (придумывай сам)\n-Вложения (если надо) вкладывай до выполнения\n-Для срочной отправки перед выполнением поставь 1-й приоритет",
  "AdditionalAction": "Действуй по инструкции для расторжения договора. \n \n"
}

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"

*/

func GenerateNote(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	var queryInfo map[string]string

	json.Unmarshal(body, &queryInfo)

	var note string

	//Заголовок

	if queryInfo["sr_type"] == "ko_several_multi" {
		note += constants.KO_SEVERAL_MULTI + "\n"
		note += constants.DETAIL_CALL + "\n"
		note += constants.KO_BODY + "\n"
	} else if queryInfo["sr_type"] == "ko_repeat" {
		note = constants.CLAIM_TITLE
		note += constants.KOLLEGI_REPEAT
	} else {
		note = constants.CLAIM_TITLE
		note += constants.KOLLEGI
	}

	//Тело

	if (queryInfo["sr_result"] == "confirm" || queryInfo["sr_result"] == "partial") && queryInfo["sr_type"] != "ko_several_multi" {
		if queryInfo["expenditure"] == "1" { //Есть расходники

			if queryInfo["sr_result"] == "confirm" { //Подтверждена полностью
				if queryInfo["no_records"] == "1" || queryInfo["no_records_only"] == "1" {
					if queryInfo["no_records_only"] == "1" { //Только по отсутствию звонка
						note += constants.FULL_EXP_NO_CALL
					} else if queryInfo["no_records_only"] == "0" { //Не по отсутствию звонка

						if queryInfo["more_thing"] == "1" { //Несколько сутей
							if queryInfo["exp_claim"] == "1" {
								note += constants.FULL_MORE_THING
							} else if queryInfo["exp_claim"] == "0" {
								note += constants.FULL_MORE_CALL
							}
						} else if queryInfo["more_thing"] == "0" { //Одна суть
							note += constants.FULL_EXP_EXIST_ALL_CALL
						}
					}
				} else if queryInfo["no_records"] == "0" || queryInfo["no_records_only"] == "0" {
					note += constants.FULL_EXP_EXIST_ALL_CALL
				}
			} else if queryInfo["sr_result"] == "partial" { //Подтверждена частично
				if queryInfo["no_records"] == "1" || queryInfo["no_records_only"] == "1" {
					if queryInfo["more_thing"] == "1" {
						if queryInfo["exp_claim"] == "1" {
							note += constants.PARTIAL_EXP_MORE_THING
						} else if queryInfo["exp_claim"] == "0" {
							note += constants.PARTIAL_EXP
						}
					} else if queryInfo["more_thing"] == "0" {
						note += constants.PARTIAL_EXP_NO_CALL
					}
				} else if queryInfo["no_records"] == "0" || queryInfo["no_records_only"] == "0" {
					if queryInfo["more_thing"] == "1" {
						note += constants.PARTIAL_EXP_MORE_THING_ALL
					} else if queryInfo["more_thing"] == "0" {
						note += constants.PARTIAL_EXP_ONE_THING
					}
				}
			}

		} else if queryInfo["expenditure"] == "0" && queryInfo["sr_type"] != "ko_several_multi" { //Нет расходников
			if queryInfo["sr_result"] == "confirm" { //Подтвержена полностью

				if queryInfo["no_records_only"] == "0" { //Не только по отсутствию звонк

					if queryInfo["no_records"] == "1" { //Нет одной из записей
						if queryInfo["more_thing"] == "1" { //Несколько сутей
							note += constants.CONFIRM_WIHOUT_EXP_MORE_THING
						} else if queryInfo["more_thing"] == "0" { //Одна суть
							note += constants.CONFIRM_WIHOUT_EXP_ALL_CALLS
						}
					} else if queryInfo["no_records"] == "0" {
						note += constants.CONFIRM_WIHOUT_EXP_ALL_CALLS
					}

				} else if queryInfo["no_records_only"] == "1" { //Только по отсутствию звонк
					note += constants.CONFIRM_WIHOUT_EXP_NO_CALLS
				}
			} else if queryInfo["sr_result"] == "partial" { //Частично
				if queryInfo["no_records"] == "0" && queryInfo["no_records_only"] == "0" {
					note += constants.PARTIAL_CHANGED
				} else {
					note += constants.PARTIAL_WITHOUT_EXP
					note += constants.ALL_AS_WE_CAN
				}
			}
		}
	} else if queryInfo["sr_result"] == "non_confirm" && queryInfo["sr_type"] != "ko_several_multi" { //Претензия клиента не подтвердилась. Результат претензии non_confirm
		note += constants.NON_CONFIRM
	}

	if queryInfo["fin_korr"] == "0" && queryInfo["close_account"] == "0" && queryInfo["unblock_needed"] == "0" && queryInfo["loyatly_needed"] == "0" && queryInfo["phone_denied"] == "0" && queryInfo["due_date_action"] == "0" && queryInfo["need_other"] == "0" {

	} else {
		note += "\n "
		note += constants.ADDITIONAL_INFO + generateAdditionalAction(queryInfo)

	}

	if queryInfo["sr_result"] == "non_confirm" && queryInfo["sr_type"] != "ko_several_multi" {
		note += constants.FOOTER_TEXT
	} else if queryInfo["sr_result"] == "confirm" && queryInfo["expenditure"] != "1" && queryInfo["sr_type"] != "ko_several_multi" {
		if queryInfo["no_records_only"] == "1" {
			note += constants.FOOTER_TEXT
		} else {
			note += constants.FOOTER_TEXT
			//note += constants.RESULT_OF_CHECK
		}
	} else if queryInfo["sr_result"] == "partial" && queryInfo["sr_type"] != "ko_several_multi" {
		note += constants.FOOTER_TEXT

		if queryInfo["expenditure"] == "1" {
			note += constants.EXP_FOOTER
		}

	} else if queryInfo["sr_result"] == "confirm" && queryInfo["expenditure"] == "1" && queryInfo["sr_type"] != "ko_several_multi" {
		note += constants.FOOTER_TEXT

		if queryInfo["no_records_only"] != "1" {
			note += constants.EXP_FOOTER
		} else {
			note += "Результат проверки (для БЭК ОФИСА): нет записи разговора (-ов)"
		}

	}

	//Тянем инфу о пользователе

	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	if len(authArr) != 2 {
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &models.JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(constants.SECRET), nil
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*models.JWTData)

	userID := data.CustomClaims["userid"]

	//Проверяем, не в овертайме ли пользователь

	checkOver := tools.CheckIfUserInOver(userID)

	if checkOver {
		//note += "\n \n ДЛЯ УКК: OVR16$"
	}

	result := &models.Info{
		Note:             note,
		InfoInstr:        generateInfoInstr(queryInfo),
		AdditionalAction: generateAdditionalInstr(queryInfo),
	}

	showResponse, err := json.Marshal(result)

	if err != nil {
		log.Println(err)
	}

	w.Write([]byte(showResponse))
}

func generateInfoInstr(queryInfo map[string]string) string {
	var result string

	if queryInfo["fin_korr"] == "1" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Выбери подстатус - Необоснована. Фин. кор-ки."
		} else {

			result += "Выбери подстатус - Обоснована. Фин. кор-ки."

		}
	} else if queryInfo["fin_korr"] == "0" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Выбери подстатус - Необоснована"
		} else {

			result += "Выбери подстатус - Обоснована"

		}
	}

	if queryInfo["inform"] != "non_inform" && queryInfo["inform"] != "" {
		result += "\n \n \n ПОРЯДОК ИНФОРМИРОВАНИЯ \n \n "
	}

	if queryInfo["inform"] == "call_inform" {
		result += "Создай активность 11. Исходящий звонок, позвони клиенту и проинформируй о результате рассмотрения претензии"
	}
	if queryInfo["inform"] == "mail_inform" {
		result += "Зарегистрируй SR Исходящий e-mail - Информирование \n -В ПО запроса напиши текст ответа клиенту (придумывай сам)\n" +
			"-Вложения (если надо) вкладывай до выполнения\n" +
			"-Для срочной отправки перед выполнением поставь 1-й приоритет"
	}
	if queryInfo["inform"] == "chat_inform" {
		result += "Действуй согласно регламенту отправки исходящего чата клиентам"
	}
	if queryInfo["inform"] == "call_pm_inform" {
		result += "Передай запрос Информирование.. - Результат рассмотрения КО - Звонок (ПМ Бизнес). Обязательно выбери компанию на запросе"
	}
	if queryInfo["inform"] == "call_pm_head_inform" {
		result += "Передай запрос Информирование.. - Результат рассмотрения КО - Звонок (РГ ПМ Бизнес). Обязательно выбери компанию на запросе"
	}
	if queryInfo["inform"] == "sms_inform" {
		result += "Отправь СМС запросом: Информирование.. - Результат рассмотрения КО - SMS"
		if queryInfo["fin_korr"] == "1" {
			result += " (рассм. корректировк)"
		} else if queryInfo["fin_korr"] == "0" {
			if queryInfo["loyatly_needed"] == "0" {
				if queryInfo["sr_type"] == "ko_repeat" {
					if queryInfo["sr_repeat_result"] == "changed" {
						result += " (решение изм.)"
					} else {
						result += " (решение не изм.)"
					}
				} else {
					if queryInfo["sr_result"] == "non_confirm" {
						if queryInfo["unblock_needed"] == "1" {
							result += " (восстановл. - необоснов.)"
						} else {
							result += " (не обосн., без корр.)"
						}
					} else {
						if queryInfo["unblock_needed"] == "1" {
							result += " (восстановл. - обоснов.)"
						} else {

							if queryInfo["sr_result"] == "confirm" {
								result += " (обосн., без корр)"
							} else {
								result += " (частич.обосн., без корр.)"
							}
						}
					}
				}
			} else if queryInfo["loyatly_needed"] == "1" {
				result += " (начисления по ПЛ/акции)"
			}
		}
	}
	if queryInfo["inform"] == "non_standart_sms" {
		result += "Действуй согласно регламенту отправки нестандартных СМС"
	}

	return result
}

func generateAdditionalInstr(queryInfo map[string]string) string {
	var result string

	if queryInfo["close_account"] == "1" {
		result += "Действуй по инструкции для расторжения договора. \n \n"
	}
	if queryInfo["fin_korr"] == "1" && queryInfo["close_account"] == "0" {
		result += "Передай запрос Претензия - Комиссии или Претензия - Проценты. Подтему выбери в зависимости от тематики запроса \n \n "
	}
	if queryInfo["loyatly_needed"] == "1" {
		result += "Передай тематический запрос по процедуре: Акции и программы лояльности \n \n"
	}
	if queryInfo["unblock_needed"] == "1" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Действуй согласно стандартному порядку восстановления обслуживания для КЦ \n \n"
		} else {
			result += "Если подтверждаешь ошибочную блокировку счета и счет ещё не закрыт, заполни опросник и выполни запрос Действие - Восстановление обслуживания \n \n"
		}
	}
	if queryInfo["due_date_move"] == "1" {
		result += "Отправь письмо по шаблону: Перенос МП на correction_expert \n Для correction_expert: " +
			"\n	Действуй по инструкции для переноса МП. \n \n "
	}
	if queryInfo["due_date_zero"] == "1" {
		result += "Отправь письмо по шаблону: Обнуление МП на correction_expert \n Для correction_expert: " +
			"\n	Действуй по инструкции для обнуления МП. \n \n "
	}

	return result
}

func generateAdditionalAction(queryInfo map[string]string) string {
	var result string

	if queryInfo["close_account"] == "1" && queryInfo["fin_korr"] == "0" {
		result += "< укажи результат решения вопроса с расторжением договора >"

	}
	if queryInfo["fin_korr"] == "1" && queryInfo["close_account"] == "0" {
		result += "Передан запрос на рассмотрение корректировок. "

	}

	if queryInfo["fin_korr"] == "1" && queryInfo["close_account"] == "1" {
		result += "<УКАЖИ результат решения вопроса с расторжением и корректировками>"

	}

	if queryInfo["loyatly_needed"] == "1" {
		result += "Передан запрос на рассмотрение начисления баллов/бонуса по акции либо участия в ней (УБЕРИ ЛИШНЕЕ)"

	}
	if queryInfo["unblock_needed"] == "1" {
		result += "< укажи результат решения вопроса с восстановлением обслуживания >"

	}
	if queryInfo["due_date_move"] == "1" {
		result += "Перенесли текущую дату минимального платежа. "

	}
	if queryInfo["due_date_zero"] == "1" {
		result += "Минимальный платеж в текущем периоде обнулен. "

	}
	if queryInfo["phone_denied"] == "1" {
		result += "В предоставлении записи звонка отказано. \n На вопрос о причине отказа ответить: \n" +
			"\"Согласно п. 3.4.3 УКБО банк вправе: \"При заключении Договоров, а также при ином обращении Клиента в банк осуществлять наблюдение, " +
			"фотографирование, аудио- и видеозапись, включая запись телефонных разговоров, без уведомления Клиента о такой записи. Клиент соглашается, " +
			"что Банк вправе хранить такие записи в течение 5 (пяти) лет с момента прекращения отношения с Клиентом, а также использовать их при " +
			"проведении любых расследований в связи с Универсальным Договором.\" При этом не указано, что банк обязан предоставлять такие записи по обращению клиента.\" "

	}
	if queryInfo["need_other"] == "1" {
		result += "<____________> "
	}

	return result
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
