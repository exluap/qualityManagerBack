/**
 * Project qualityManagerApi created by exluap
 * Date: 04.03.2019 14:31
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package adminTools

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

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	data := auth.CheckToken(w, r)

	if tools.CheckAdminMode(data.CustomClaims["userid"]) {
		var userMap map[string]string

		err := json.Unmarshal(body, &userMap)

		if err != nil {
			log.Print("Admin error: ", err)
			raven.CaptureErrorAndWait(err, nil)
			http.Error(w, "Request is bad", http.StatusBadRequest)
		} else {
			res := tools.AddNewUser(userMap["firstName"], userMap["lastName"], userMap["middleName"], userMap["login"], userMap["password"], userMap["groups"], userMap["isAdmin"], userMap["winlogin"])

			if res {
				resultation := &models.Resultation{
					Result: "User created",
				}

				result, _ := json.Marshal(resultation)

				w.Write(result)
			} else {
				resultation := &models.Resultation{
					Result: "User exist. Cant create",
				}

				result, _ := json.Marshal(resultation)

				w.Write(result)
			}
		}
	} else {
		res := models.Resultation{
			Result: "User is not admin",
		}

		result, _ := json.Marshal(res)

		http.Error(w, string(result), http.StatusBadRequest)
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	data := auth.CheckToken(w, r)

	if tools.CheckAdminMode(data.CustomClaims["userid"]) {
		var userMap map[string]string

		err := json.Unmarshal(body, &userMap)

		if err != nil {
			log.Print("Admin error: ", err)
			raven.CaptureErrorAndWait(err, nil)
			http.Error(w, "Request is bad", http.StatusBadRequest)
		} else {
			res := tools.UpdateUserInfo(userMap["groups"], userMap["isAdmin"], userMap["login"])

			if res {
				resultation := &models.Resultation{
					Result: "User update",
				}

				result, _ := json.Marshal(resultation)

				w.Write(result)
			} else {
				resultation := &models.Resultation{
					Result: "User not updated. Not found",
				}

				result, _ := json.Marshal(resultation)

				w.Write(result)
			}
		}
	} else {
		res := models.Resultation{
			Result: "User is not admin",
		}

		result, _ := json.Marshal(res)

		http.Error(w, string(result), http.StatusBadRequest)
	}
}
