/**
 * Project qualityManagerApi created by exluap
 * Date: 25.10.2018 23:57
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/tools"
	"strings"
	"time"
)

const (
	PORT   = "1337"
	SECRET = "Ra7G9XgMxwD8cehHp9Netf5EBpXMDCL3EBMX"
)

type Resultation struct {
	Result string
}

type Info struct {
	Note             string
	InfoInstr        string
	AdditionalAction string
}

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

const version = "1.3.4"

func main() {
	mode := flag.Bool("makeNewUser", false, "make new user")
	firstName := flag.String("firstname", "Nikita", "First Name")
	lastName := flag.String("lastname", "Zaycev", "Last Name")
	middleName := flag.String("middlename", "Alekseevich", "Middle Name")

	flag.Parse()

	/** v := semver.MustParse(version)

		latest, err := selfupdate.UpdateSelf(v,"exluap/forawalk")



		if err != nil {
			log.Printf("Update failed: %s",err)
		}


		if latest.Version.Equals(v) {
			log.Printf("Current version %s is latest", latest.Version)
		} else {
			log.Printf("The server is successfully updated to version: %s",latest.Version)
			log.Printf("Release note: \n %s", latest.ReleaseNotes)
		}
	 **/

	if *mode {
		makeNewUser(*firstName, *lastName, *middleName)
	} else {
		startWebServer()
	}

}

func startWebServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/get_query", getQuery)
	mux.HandleFunc("/add_query", addQuery)
	mux.HandleFunc("/get_queries", getQueries)
	mux.HandleFunc("/generate_note_and_instruction", generateNote)
	mux.HandleFunc("/delete_sr", deleteSR)
	mux.HandleFunc("/add_log", addingLog)

	handler := cors.AllowAll().Handler(mux)

	log.Println("Listening for connections on port: ", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, please auth")
}

func login(w http.ResponseWriter, r *http.Request) {
	selfupdate.EnableLog()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		tools.SaveLog("backend", "Login Failed, need logs from server side", "system")
		http.Error(w, "Login Failed", http.StatusUnauthorized)
	}

	var userData map[string]string

	json.Unmarshal(body, &userData)

	log.Println("try auth user: " + userData["login"])
	tools.SaveLog("backend", "try auth user: "+userData["login"], "system")

	if tools.CheckIfUserExist(userData["login"], userData["pass"]) {
		claims := JWTData{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(11 * time.Hour).Unix(),
			},

			CustomClaims: map[string]string{
				"userid": userData["login"],
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(SECRET))
		if err != nil {
			log.Println(err)
			tools.SaveLog("backend", "Login Failed, need logs from server side", "system")
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		json, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			tokenString,
		})

		if err != nil {
			log.Println(err)
			tools.SaveLog("backend", "Login Failed, need logs from server side", "system")
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		w.Write(json)
	} else {
		http.Error(w, "Login failed!", http.StatusUnauthorized)
		tools.SaveLog("backend", "Login Failed, need logs from server side", "system")
	}

}

func getQueries(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	if len(authArr) != 2 {
		tools.SaveLog("backend", "Authentication header is invalid: "+authToken, "system")
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		log.Println(err)
		tools.SaveLog("backend", "Failed Request! Need logs from user", "system")
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims["userid"]

	jsonData, err := tools.UserQueries(userID)

	if err != nil {
		log.Println(err)
		tools.SaveLog("backend", "Failed Request! Need logs from user", "system")
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	tools.SaveLog("backend", "User request queries list: "+string(jsonData), getUserLogin(w, r))

	w.Write(jsonData)
}

func addQuery(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	authToken := r.Header.Get("Authorization")
	authArray := strings.Split(authToken, " ")

	if len(authArray) != 2 {
		log.Println("Auth header is invalid: " + authToken)
		tools.SaveLog("backend", "Failed Request! Need logs from user", "system")
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	if err != nil {
		log.Println(err)
		tools.SaveLog("backend", "Failed Request! Need logs from user", "system")
		http.Error(w, "Save Query Failed", http.StatusUnauthorized)
	}

	jwtToken := authArray[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(SECRET), nil
	})

	data := claims.Claims.(*JWTData)

	var queryData map[string]string

	json.Unmarshal(body, &queryData)

	tools.AddQueryToDB(data.CustomClaims["userid"], queryData["sr_number"], queryData["sr_type"], queryData["sr_result"], queryData["sr_repeat_result"], queryData["inform"], queryData["no_records"], queryData["no_records_only"], queryData["expenditure"], queryData["more_thing"], queryData["exp_claim"], queryData["fin_korr"], queryData["close_account"], queryData["unblock_needed"], queryData["loyatly_needed"], queryData["phone_denied"], queryData["du_date_action"], queryData["due_date_zero"], queryData["due_date_move"], queryData["need_other"])

	res := &Resultation{
		Result: "ok",
	}

	logi, err := json.Marshal(queryData)

	tools.SaveLog("backend", "Saved query: "+string(logi), getUserLogin(w, r))

	showResponse, _ := json.Marshal(res)

	w.Write(showResponse)
}

func getQuery(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	var queryInfo map[string]string

	json.Unmarshal(body, &queryInfo)

	jsonData, err := tools.GetQueryInfo(queryInfo["sr_number"])

	if err != nil {
		log.Println(err)
		tools.SaveLog("backend", "Failed Request! Need logs from user", "system")
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	w.Write(jsonData)

	tools.SaveLog("backend", "User response for getQuery: "+string(jsonData), getUserLogin(w, r))

}

func generateNote(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	var queryInfo map[string]string

	json.Unmarshal(body, &queryInfo)

	var note string

	//Заголовок

	if queryInfo["sr_type"] == "ko_several_multi" {
		note += tools.KO_SEVERAL_MULTI + "\n"
		note += tools.DETAIL_CALL + "\n"
		note += tools.KO_BODY + "\n"
	} else if queryInfo["sr_type"] == "ko_repeat" {
		note = tools.CLAIM_TITLE
		note += tools.KOLLEGI_REPEAT
	} else {
		note = tools.CLAIM_TITLE
		note += tools.KOLLEGI
	}

	//Тело

	if (queryInfo["sr_result"] == "confirm" || queryInfo["sr_result"] == "partial") && queryInfo["sr_type"] != "ko_several_multi" {
		if queryInfo["expenditure"] == "1" { //Есть расходники

			if queryInfo["sr_result"] == "confirm" { //Подтверждена полностью
				if queryInfo["no_records"] == "1" || queryInfo["no_records_only"] == "1" {
					if queryInfo["no_records_only"] == "1" { //Только по отсутствию звонка
						note += tools.FULL_EXP_NO_CALL
					} else if queryInfo["no_records_only"] == "0" { //Не по отсутствию звонка

						if queryInfo["more_thing"] == "1" { //Несколько сутей
							if queryInfo["exp_claim"] == "1" {
								note += tools.FULL_MORE_THING
							} else if queryInfo["exp_claim"] == "0" {
								note += tools.FULL_MORE_CALL
							}
						} else if queryInfo["more_thing"] == "0" { //Одна суть
							note += tools.FULL_EXP_EXIST_ALL_CALL
						}
					}
				} else if queryInfo["no_records"] == "0" || queryInfo["no_records_only"] == "0" {
					note += tools.FULL_EXP_EXIST_ALL_CALL
				}
			} else if queryInfo["sr_result"] == "partial" { //Подтверждена частично
				if queryInfo["no_records"] == "1" || queryInfo["no_records_only"] == "1" {
					if queryInfo["more_thing"] == "1" {
						if queryInfo["exp_claim"] == "1" {
							note += tools.PARTIAL_EXP_MORE_THING
						} else if queryInfo["exp_claim"] == "0" {
							note += tools.PARTIAL_EXP
						}
					} else if queryInfo["more_thing"] == "0" {
						note += tools.PARTIAL_EXP_NO_CALL
					}
				} else if queryInfo["no_records"] == "0" || queryInfo["no_records_only"] == "0" {
					if queryInfo["more_thing"] == "1" {
						note += tools.PARTIAL_EXP_MORE_THING_ALL
					} else if queryInfo["more_thing"] == "0" {
						note += tools.PARTIAL_EXP_ONE_THING
					}
				}
			}

		} else if queryInfo["expenditure"] == "0" { //Нет расходников
			if queryInfo["sr_result"] == "confirm" { //Подтвержена полностью

				if queryInfo["no_records_only"] == "0" { //Не только по отсутствию звонк

					if queryInfo["no_records"] == "1" { //Нет одной из записей
						if queryInfo["more_thing"] == "1" { //Несколько сутей
							note += tools.CONFIRM_WIHOUT_EXP_MORE_THING
						} else if queryInfo["more_thing"] == "0" { //Одна суть
							note += tools.CONFIRM_WIHOUT_EXP_ALL_CALLS
						}
					} else if queryInfo["no_records"] == "0" {
						note += tools.CONFIRM_WIHOUT_EXP_ALL_CALLS
					}

				} else if queryInfo["no_records_only"] == "1" { //Только по отсутствию звонк
					note += tools.CONFIRM_WIHOUT_EXP_NO_CALLS
				}
			} else if queryInfo["sr_result"] == "partial" { //Частично
				if queryInfo["no_records"] == "0" && queryInfo["no_records_only"] == "0" {
					note += tools.PARTIAL_CHANGED
				} else {
					note += tools.PARTIAL_WITHOUT_EXP
					note += tools.ALL_AS_WE_CAN
				}
			}
		}
	} else if queryInfo["sr_result"] == "non_confirm" { //Претензия клиента не подтвердилась. Результат претензии non_confirm
		note += tools.NON_CONFIRM
	}

	if queryInfo["fin_korr"] == "0" && queryInfo["close_account"] == "0" && queryInfo["unblock_needed"] == "0" && queryInfo["loyatly_needed"] == "0" && queryInfo["phone_denied"] == "0" && queryInfo["due_date_action"] == "0" && queryInfo["need_other"] == "0" {

	} else {
		note += "\n "
		note += tools.ADDITIONAL_INFO + generateAdditionalAction(queryInfo)

	}

	if queryInfo["sr_result"] == "non_confirm" && queryInfo["sr_type"] != "ko_several_multi" {
		note += tools.FOOTER_TEXT
	} else if queryInfo["sr_result"] == "confirm" && queryInfo["expenditure"] != "1" {
		if queryInfo["no_records_only"] == "1" {
			note += tools.FOOTER_TEXT
		} else {
			note += tools.FOOTER_TEXT
			note += tools.RESULT_OF_CHECK
		}
	} else if queryInfo["sr_result"] == "partial" && queryInfo["sr_type"] != "ko_several_multi" {
		note += tools.FOOTER_TEXT

		if queryInfo["expenditure"] == "1" {
			note += tools.EXP_FOOTER
		}

	} else if queryInfo["sr_result"] == "confirm" && queryInfo["expenditure"] == "1" && queryInfo["sr_type"] != "ko_several_multi" {
		note += tools.FOOTER_TEXT

		if queryInfo["no_records_only"] != "1" {
			note += tools.EXP_FOOTER
		} else {
			note += "Результат проверки (для БЭК ОФИСА): нет записи разговора (-ов)"
		}

	}

	result := &Info{
		Note:             note,
		InfoInstr:        generateInfoInstr(queryInfo),
		AdditionalAction: generateAdditionalInstr(queryInfo),
	}

	showResponse, err := json.Marshal(result)

	if err != nil {
		log.Println(err)
	}

	tools.SaveLog("backend", "User response for generate Note Full: "+string(showResponse), getUserLogin(w, r))

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
		result += "Передан запрос на расторжение. Запрос на корректировки при необходимости передаст сотрудник ТМ "

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

func makeNewUser(firstName, lastName, middleName string) {

	flag.Parse()

	var login string

	login = string(firstName[0]) + string(middleName[0]) + lastName

	pass, result := tools.AddNewUser(firstName, lastName, middleName, login)

	if result {
		log.Println("User Login: " + login + "\n Password: " + pass)
	} else {
		log.Println("User is exist")
	}

}

func deleteSR(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	body, err := ioutil.ReadAll(r.Body)

	var requestBody map[string]string

	json.Unmarshal(body, &requestBody)

	if len(authArr) != 2 {
		log.Println("Authentication header is invalid: " + authToken)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*JWTData)

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

	tools.SaveLog("backend", "User delete SR: "+requestBody["sr_number"], userID)

}

func addingLog(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	body, err := ioutil.ReadAll(r.Body)

	var requestBody map[string]string

	json.Unmarshal(body, &requestBody)

	if len(authArr) != 2 {
		log.Println("Authentication header is invalid: " + authToken)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims["userid"]

	tools.SaveLog(requestBody["inter"], requestBody["infoText"], userID)

	resultShow, err := json.Marshal("Log is saved")

	w.Write(resultShow)

}

func getUserLogin(w http.ResponseWriter, r *http.Request) string {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	body, err := ioutil.ReadAll(r.Body)

	var requestBody map[string]string

	json.Unmarshal(body, &requestBody)

	if len(authArr) != 2 {
		log.Println("Authentication header is invalid: " + authToken)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims["userid"]

	return userID
}
