/**
 * Project qualityManagerApi created by exluap
 * Date: 22.01.2019 00:22
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/auth"
	"qualityManagerApi/models"
	"qualityManagerApi/tools"
	"strings"
)

/**
@api {get} /api/user/overtime Checking overtime or set it
@apiName PostOver
@apiVersion 1.0.0
@apiGroup User
@apiHeader token Auth Token of user with information about him

@apiDescription Getting info about overtime or set it as u need


@apiSuccessExample {json} Success-Response (set overtime)
	{
		"Result": "Overtime changed"
	}

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"
*/

func CheckOver(w http.ResponseWriter, r *http.Request) {

	data := auth.CheckToken(w, r)

	userID := data.CustomClaims["userid"]

	var id string

	if tools.CheckIfUserInOver(userID) {
		id = "1"
	} else {
		id = "0"
	}

	tools.IneedMoreMoney(userID, id)

	res := &models.Resultation{
		Result: "Overtime changed",
	}

	jsonResult, _ := json.Marshal(res)

	w.Write(jsonResult)

}

func GetUserLogin(w http.ResponseWriter, r *http.Request) string {

	body, _ := ioutil.ReadAll(r.Body)

	var requestBody map[string]string

	json.Unmarshal(body, &requestBody)

	data := auth.CheckToken(w, r)

	userID := data.CustomClaims["userid"]

	return userID
}

/**
@api {get} /api/user/info Getting info about user
@apiName GetUserInfo
@apiVersion 1.0.0
@apiGroup User
@apiHeader token Auth Token of user with information about him


@apiDescription Getting info about user

@apiSuccess {String} Login Return User Login
@apiSuccess {Boolean} Overtime Return user status

@apiSuccessExample {json} Success-Response
	{
		"Login": "USERTEST",
		"Over": true
	}

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"

*/

func GetInfoAboutUser(w http.ResponseWriter, r *http.Request) {
	data := GetUserLogin(w, r)

	groups := tools.GetUserGroups(data)

	resultOfGroups := strings.Split(groups, ",")

	res := &models.User{
		Login:   data,
		Over:    tools.CheckIfUserInOver(data),
		IsAdmin: tools.CheckAdminMode(data),
		Groups:  resultOfGroups,
	}

	showReq, _ := json.Marshal(res)

	w.Write(showReq)
}

/**
@api {post} /api/user/changelogin Change User login
@apiVersion 1.0.0
@apiGroup User
@apiName PostNewLogin

@apiDescription Set new login for user

@apiParam {String} newLogin New login for user

@apiSuccess {json} Success-Response
	{
		"Result": "Login changed"
	}

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"
*/

func ChangeLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusInternalServerError)
	}
	data := auth.CheckToken(w, r)

	var changePass map[string]string

	json.Unmarshal(body, &changePass)

	userId := data.CustomClaims["userid"]

	if !tools.CheckIfExistRegister(userId) {
		err = errors.New("User is exist")
	} else {
		err = tools.ChangeUserLogin(userId, changePass["newLogin"])
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Bad credentials! User is exist", http.StatusBadRequest)
	} else {
		res := &models.Resultation{
			Result: "Login changed",
		}

		result, _ := json.Marshal(res)

		w.Write(result)
	}

}

func Logging(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("error with save log: ", err)
	}

	var b map[string]string

	err = json.Unmarshal(body, &b)

	if err != nil {
		log.Println("error with parse body: ", err)
	}

	err = tools.SaveLog("front", b["text"], b["user"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error with save log: ", err)
		_, err = w.Write([]byte("can not save log"))
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("ok"))
	}
}
