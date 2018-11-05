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
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/tools"
	"strings"
	"time"
)

const (
	PORT = "1337"
	SECRET = "Ra7G9XgMxwD8cehHp9Netf5EBpXMDCL3EBMX"
)

type Info struct {
	Note string
	InfoInstr string
	AdditionalAction string
}

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

func main() {

	mode := flag.Bool("makeNewUser", false, "make new user")
	firstName := flag.String("firstname", "Nikita", "First Name")
	lastName := flag.String("lastname", "Zaycev", "Last Name")
	middleName := flag.String("middlename", "Alekseevich", "Middle Name")



	flag.Parse()

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
	mux.HandleFunc("/add_query",addQuery)
	mux.HandleFunc("/get_queries", getQueries)
	mux.HandleFunc("/generate_note_and_instruction", generateNote)

	handler := cors.AllowAll().Handler(mux)

	log.Println("Listening for connections on port: ", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Hello, please auth")
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, "Login Failed", http.StatusUnauthorized)
	}

	var userData map[string]string


	json.Unmarshal(body, &userData)

	log.Println(userData["login"])

	if tools.CheckIfUserExist(userData["login"], userData["pass"]) {
		claims := JWTData{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},

			CustomClaims: map[string]string{
				"userid": userData["login"],
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(SECRET))
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		json, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			tokenString,
		})

		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		w.Write(json)
	} else {
		http.Error(w, "Login failed!", http.StatusUnauthorized)
	}

}

func getQueries(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken," ")

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

	jsonData, err := tools.UserQueries(userID)

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	w.Write(jsonData)
}

func addQuery(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	authToken := r.Header.Get("Authorization")
	authArray := strings.Split(authToken, " ")

	if len(authArray) != 2 {
		log.Println("Auth header is invalid: " + authToken)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	if err != nil {
		log.Println(err)
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

	w.Write([]byte("Query Added"))
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
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	fmt.Println(jsonData)

	w.Write(jsonData)

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
		if queryInfo["expenditure"] == "1" {
			log.Println("exp1 start")
			if queryInfo["sr_result"] == "confirm" {
				if queryInfo["no_records"] == "0" || queryInfo["no_records_only"] == "0" {
					note += tools.FULL_EXP_EXIST_ALL_CALL
					note += tools.DETAIL_CALL
					note += tools.CLAIM_INFO

				} else if queryInfo["no_records"] == "1" || queryInfo["no_records_only"] == "1" {
					if queryInfo["no_records_only"] == "1" {
						note += tools.FULL_EXP_NO_CALL
						note += tools.DETAIL_CALL
						note += tools.CLAIM_INFO
					} else if queryInfo["no_records_only"] == "0" {
						if queryInfo["no_records"] == "0" {
							note += tools.FULL_EXP_EXIST_ALL_CALL
							note += tools.DETAIL_CALL
							note += tools.CLAIM_INFO

						} else if queryInfo["no_records"] == "1" {
							if queryInfo["more_thing"] == "0" {
								note += tools.FULL_EXP_NO_CALL
								note += tools.DETAIL_CALL
								note += tools.CLAIM_INFO

							} else if queryInfo["more_thing"] == "1" {
								if queryInfo["exp_claim"] == "1" {
									note += tools.FULL_MORE_THING
								} else if queryInfo["exp_claim"] == "0" {
									note += tools.FULL_MORE_CALL
									note += tools.DETAIL_CALL
									note += tools.CLAIM_INFO
								}
							}

						}

					}
				}
			} else if queryInfo["sr_result"] == "partial" {
				if queryInfo["no_records"] == "0" || queryInfo["no_records_only"] == "0" {
					if queryInfo["more_thing"] == "1" {
						note += tools.PARTIAL_EXP_MORE_THING_ALL
						note += tools.DETAIL_CALL
						note += tools.CLAIM_INFO

					} else if queryInfo["more_thing"] == "0" {
							note += tools.FULL_EXP_EXIST_ALL_CALL
							note += tools.DETAIL_CALL
							note += tools.CLAIM_INFO
					}

				} else if queryInfo["no_records"] == "1" || queryInfo["no_records_only"] == "1" {
					if queryInfo["more_thing"] == "1" {
						if queryInfo["exp_claim"] == "1" {
							note += tools.PARTIAL_EXP_MORE_THING
							note += tools.DETAIL_CALL
							note += tools.CLAIM_INFO
						} else if queryInfo["exp_claim"] == "0" {
							note += tools.PARTIAL_EXP
							note += tools.DETAIL_CALL
							note += tools.CLAIM_INFO
						}


					} else if queryInfo["more_thing"] == "0" {
						note += tools.PARTIAL_EXP_NO_CALL
						note += tools.DETAIL_CALL
						note += tools.CLAIM_INFO
					}
				}
			}
		} else if queryInfo["expenditure"] == "0" {
			if queryInfo["sr_result"] == "confirm" {
				if queryInfo["no_records_only"] == "1" {
					note += tools.CONFIRM_WIHOUT_EXP_NO_CALLS
					note += tools.DETAIL_CALL
					note += tools.CLAIM_INFO

				} else if queryInfo["no_records_only"] == "0" {
					if queryInfo["no_records"] == "0" {
						note += tools.CONFIRM_WIHOUT_EXP_ALL_CALLS
						note += tools.DETAIL_CALL
						note += tools.CLAIM_INFO
					} else if queryInfo["no_records"] == "1" {
						if queryInfo["more_thing"] == "1" {
							note += tools.CONFIRM_WIHOUT_EXP_MORE_THING
							note += tools.DETAIL_CALL
							note += tools.CLAIM_INFO
						} else if queryInfo["more_thing"] == "0" {
							note += tools.CONFIRM_WIHOUT_EXP_ALL_CALLS
							note += tools.DETAIL_CALL
							note += tools.CLAIM_INFO
						}
					}
				}
			} else if queryInfo["sr_result"] == "partial" {
				if queryInfo["no_records"] == "1" && queryInfo["no_records_only"] == "1" {
					note += tools.PARTIAL_WITHOUT_EXP
					note += tools.DETAIL_CALL
					note += tools.CLAIM_INFO
				} else if queryInfo["no_records"] == "0" && queryInfo["no_records_only"] == "0" {
					note += tools.PARTIAL
					note += tools.DETAIL_CALL
					note += tools.CLAIM_INFO
				}
			}
		}


		if queryInfo["expenditure"] == "1" {
			if queryInfo["sr_result"] == "confirm" {
				note += tools.RESULT_OF_CHECK
			}
			note += tools.EXP_FOOTER
		} else if queryInfo["expenditure"] == "0" {

		}


	} else if queryInfo["sr_result"] == "non_confirm" && queryInfo["sr_type"] != "ko_several_multi" {

		note += tools.NON_CONFIRM
		note += tools.DETAIL_CALL
		note += tools.CLAIM_INFO


	}

	if queryInfo["fin_korr"] == "0" && queryInfo["close_account"] == "0" && queryInfo["unblock_needed"] == "0" && queryInfo["loyatly_needed"] == "0" && queryInfo["phone_denied"] == "0" && queryInfo["due_date_action"] == "0" && queryInfo["need_other"] == "0"{

	} else {
		note += tools.ADDITIONAL_INFO + generateAdditionalAction(queryInfo)
	}


	result := &Info{
		Note:note,
		InfoInstr: generateInfoInstr(queryInfo),
		AdditionalAction:generateAdditionalInstr(queryInfo),
	}

	showResponse, err := json.Marshal(result)

	if err != nil {
		log.Println(err)
	}

	w.Write([]byte(showResponse))
}

func generateInfoInstr(queryInfo map[string]string) string{
	var result string

	if queryInfo["fin_korr"] == "1" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Выбери подстатус - Необоснована. Фин. кор-ки."
		} else {
			if queryInfo["no_records_only"] == "1" {
				result += "Выбери подстатус - Рассмотрена. Фин. кор-ки."
			} else {
				result += "Выбери подстатус - Обоснована. Фин. кор-ки."
			}
		}
	} else if queryInfo["fin_korr"] == "0" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Выбери подстатус - Необоснована"
		} else {
			if queryInfo["no_records_only"] == "1" {
				result += "Выбери подстатус - Рассмотрена"
			} else if queryInfo["no_records_only"] == "0"{
				result += "Выбери подстатус - Обоснована"
			}
		}
	}

	result += "\n \n \n ПОРЯДОК ИНФОРМИРОВАНИЯ \n \n "

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
							if queryInfo["no_records_only"] == "1" {
								result += " (проверено, нет записи)"
							} else {
								if queryInfo["sr_result"] == "confirm" {
									result += " (обосн., без корр)"
								} else {
									result += " (частич.обосн., без корр.)"
								}
							}
						}
					}
				}
			} else if queryInfo["loyatly_needed"] == "1" {
				result += " (начисления по ПЛ/акции)"
			}
		}

		result += " /n /n"
	}

	return result
}

func generateAdditionalInstr(queryInfo map[string]string) string {
	var result string

	if queryInfo["close_account"] == "1" {
		result += "Передай запрос Закрытие договора - Клиентское обслуживание \n \n"
	}
	if queryInfo["fin_korr"] == "1" {
		result += "Передай запрос Претензия - Комиссии или Претензия - Проценты. Подтему выбери в зависимости от тематики запроса \n \n "
	}
	if queryInfo["loyatly_needed"] == "1" {
		result += "Передай тематический запрос по процедуре: Акции и программы лояльности \n \n"
	}
	if queryInfo["unblock_needed"] == "1" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Если подтверждаешь ошибочную блокировку счета и счет ещё не закрыт, передай запрос Действие - Восстановление обслуживания и заполни опросник \n \n"
		} else {
			result += "Действуй согласно стандартному порядку восстановления обслуживания для КЦ \n \n"
		}
	}
	if queryInfo["due_date_move"] == "1" {
		result += "Отправь письмо по шаблону: Перенос МП на correction_expert \n \n"
	}
	if queryInfo["due_date_zero"] == "1" {
		result += "Отправь письмо по шаблону: Обнуление МП на correction_expert \n \n"
	}

	return result
}

func generateAdditionalAction(queryInfo map[string]string) string {
	var result string

	if queryInfo["close_account"] == "1" {
		result += "Передан запрос на расторжение. \n"
		result += "\n"
	}
	if queryInfo["fin_korr"] == "1" {
		result += "Передан запрос на рассмотрение корректировок \n"
		result += "\n"
	}
	if queryInfo["loyatly_needed"] == "1" {
		result += "Передан запрос на рассмотрение начисления баллов/бонуса по акции либо участия в ней \n"
		result += "\n"
	}
	if queryInfo["unblock_needed"] == "1" {
		if queryInfo["sr_result"] == "non_confirm" {
			result += "Передан запрос на восстановление. \n"
		} else {
			result += "Счет обслуживания восстановлен. \n"
		}
		result += "\n"
	}
	if queryInfo["due_date_move"] == "1" {
		result += "Дата минимального платежа изменена только в текущем периоде. \n"
		result += "\n"
	}
	if queryInfo["due_date_zero"] == "1" {
		result += "Минимальный платеж в текущем периоде обнулен. \n"
		result += "\n"
	}
	if queryInfo["phone_denied"] == "1" {
		result += "В предоставлении записи звонка отказано. \n На вопрос о причине отказа ответить: \n" +
			"\"Согласно п. 3.4.3 УКБО банк вправе: \"При заключении Договоров, а также при ином обращении Клиента в банк осуществлять наблюдение, " +
			"фотографирование, аудио- и видеозапись, включая запись телефонных разговоров, без уведомления Клиента о такой записи. Клиент соглашается, " +
			"что Банк вправе хранить такие записи в течение 5 (пяти) лет с момента прекращения отношения с Клиентом, а также использовать их при " +
			"проведении любых расследований в связи с Универсальным Договором.\" При этом не указано, что банк обязан предоставлять такие записи по обращению клиента.\" \n"
		result += "\n"
		}
	if queryInfo["need_other"] == "1" {
		result += "<____________> \n"
	}


	return result
}

func makeNewUser(firstName, lastName, middleName string) {

	flag.Parse()

	var login string

	login = string(firstName[0]) +  string(middleName[0]) + string(lastName[0])

	pass,result := tools.AddNewUser(firstName, lastName, middleName, login)


	if result {
		log.Println("User Login: " + login + "\n Password: " + pass)
	} else {
		log.Println("User is exist")
	}


}