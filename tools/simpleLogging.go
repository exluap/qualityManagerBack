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
