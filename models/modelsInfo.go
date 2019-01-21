/**
 * Project qualityManagerApi created by exluap
 * Date: 22.01.2019 00:25
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package models

import "github.com/dgrijalva/jwt-go"

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
