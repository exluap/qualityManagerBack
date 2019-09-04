/**
 * Project qualityManagerApi created by exluap
 * Date: 22.01.2019 00:23
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package constants

import (
	"github.com/tkanos/gonfig"
	"log"
)

const (
	PORT   = "1337"
	SECRET = "Ra7G9XgMxwD8cehHp9Netf5EBpXMDCL3EBMX"
)

type Config struct {
	Port               string `json:"Port"`
	ReportSaveFileSave string `json:"ReportSaveFileSave"`
	UploadFilePath     string `json:"UploadFilePath"`
	ReportFileName     string `json:"ReportFileName"`
}

func GetConfig() Config {
	var config Config

	err := gonfig.GetConf("config.json", &config)

	if err != nil {
		log.Printf("Can not open config file: %s", err)
	}

	return config
}
