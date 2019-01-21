/**
 * Project qualityManagerApi created by exluap
 * Date: 22.01.2019 00:31
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package tools

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/constants"
	"qualityManagerApi/models"
	"strings"
)

/**
@api {post} /add_log Adding log record
@apiName PostLogInfo
@apiVersion 1.0.0
@apiGroup Logging
@apiHeader token Auth Token of user with information about him

@apiDescription For looking log about of actions by user need saving it

@apiParam {String} inter Interface of getting log
@apiParam {String} infoText Text who's need add for this log record

@apiSuccess {String} Result Result of saving logs
@apiSuccessExample {json} Success-Response:
	"Log is saved"

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response:
		"Request failed!"

*/
func AddingLog(w http.ResponseWriter, r *http.Request) {
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

	SaveLog(requestBody["inter"], requestBody["infoText"], userID)

	resultShow, err := json.Marshal("Log is saved")

	w.Write(resultShow)

}
