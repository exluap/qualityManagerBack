/**
 * Project qualityManagerApi created by exluap
 * Date: 22.01.2019 00:42
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package auth

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/raven-go"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/constants"
	"qualityManagerApi/models"
	"qualityManagerApi/tools"
	"strings"
	"time"
)

/**
@api {post} /api/auth/login Getting auth token
@apiVersion 1.0.0
@apiName GetToken
@apiGroup Authentication
@apiHeader token Auth Token of user with information about him

@apiDescription Get auth Token for use some methods for working

@apiParam {String} login Login of User
@apiParam {String} pass Password of User

@apiSuccess {String} token Bearer Auth token for use

@apiSuccessExample {json} Success-Response
	HTTP/1.1 200 OK
     {
       "token": "AUTH_TOKEN"
     }


@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"
*/

func Login(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//selfupdate.EnableLog()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		raven.CaptureErrorAndWait(err, nil)
		http.Error(w, "Login Failed", http.StatusUnauthorized)
	}

	var userData map[string]string

	json.Unmarshal(body, &userData)

	log.Println("try auth user: " + userData["login"])

	if tools.CheckIfUserExist(userData["login"], userData["pass"]) {
		claims := models.JWTData{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(11 * time.Hour).Unix(),
			},

			CustomClaims: map[string]string{
				"userid": userData["login"],
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(constants.SECRET))
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
			raven.CaptureErrorAndWait(err, nil)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		w.Write(json)
	} else {
		http.Error(w, "Login failed!", http.StatusUnauthorized)
	}

}

func CheckToken(w http.ResponseWriter, r *http.Request) *models.JWTData {
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
		raven.CaptureErrorAndWait(err, nil)
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	result := claims.Claims.(*models.JWTData)

	return result
}

/**
@api {post} /api/auth/changepassword Change User password
@apiVersion 1.0.0
@apiGroup Authentication
@apiName PostNewPassword

@apiDescription Set new password for user

@apiParam {String} passwordOld Old password of user
@apiParam {String} passwordNew New User password

@apiSuccessExample {json} Success-Response
	{
		"Result": "Password changed"
	}

@apiError Unauthorized Getting Bad Credentials

@apiErrorExample Error-Response
		"Request failed!"
*/

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusInternalServerError)
	}
	data := CheckToken(w, r)

	var changePass map[string]string

	json.Unmarshal(body, &changePass)

	userId := data.CustomClaims["userid"]

	if changePass["passwordOld"] != "" && changePass["passwordNew"] != "" {
		err = tools.ChangeUserPassword(userId, changePass["passwordOld"], changePass["passwordNew"])
	} else {
		err = errors.New("Not enough info")
	}

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Println(err)
		http.Error(w, "Bad credentials!", http.StatusUnauthorized)
	} else {
		res := &models.Resultation{
			Result: "Password changed",
		}

		result, _ := json.Marshal(res)

		w.Write(result)
	}

}
