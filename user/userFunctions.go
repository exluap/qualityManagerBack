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
	"flag"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/constants"
	"qualityManagerApi/models"
	"qualityManagerApi/tools"
	"strings"
)

/**
@api {post} /in_over Checking overtime or set it
@apiName PostOver
@apiVersion 1.0.0
@apiGroup User
@apiHeader token Auth Token of user with information about him

@apiDescription Getting info about overtime or set it as u need

@apiParam {String} overtime Result of overtime. 1 - is active, 0 - is inactive
@apiParam {String} action What r u need? "check" - checking if user in overtime or "update" - set overtime as !overtime


@apiSuccess {String} Result result of action

@apiSuccessExample {json} Success-Response (check action)
	true
@apiSuccessExample {json} Success-Response (set overtime)
	Ok

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"
*/

func CheckOver(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")
	body, err := ioutil.ReadAll(r.Body)

	if len(authArr) != 2 {
		tools.SaveLog("backend", "Authentication header is invalid: "+authToken, "system")
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
		tools.SaveLog("backend", "Failed Request! Need logs from user", "system")
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*models.JWTData)

	userID := data.CustomClaims["userid"]

	var userData map[string]string

	json.Unmarshal(body, &userData)

	if userData["action"] == "check" {
		jsonData := tools.CheckIfUserInOver(userID)

		jsonResult, _ := json.Marshal(jsonData)

		w.Write(jsonResult)
	} else if userData["action"] == "update" {
		tools.IneedMoreMoney(userID, userData["overtime"])
		jsonResult, _ := json.Marshal("Ok")
		w.Write(jsonResult)
	}

}

func GetUserLogin(w http.ResponseWriter, r *http.Request) string {
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

	return userID
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
