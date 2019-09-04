/**
 * Project IntelliJ IDEA created by exluap
 * Date: 04.09.2019 11:35
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */
package tools

import (
	"encoding/json"
	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"qualityManagerApi/constants"
)

type Report struct {
	Data []struct {
		Assegnee string `json:"assegnee"`
		DateWork string `json:"date_work"`
		ID       string `json:"id"`
	} `json:"data"`
}

func GenerateReport(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Can not read body: %s", err)
	}

	var report Report

	err = json.Unmarshal(body, &report)

	if err != nil {
		log.Printf("Can not parse json body: %s", err)
	}

	conf := constants.GetConfig()

	file, err := os.Create(conf.ReportSaveFileSave + conf.ReportFileName)

	if err != nil {
		log.Printf("Error with create file: %s", err)
	}

	log.Printf("Json: %s", body)
	log.Printf("reports: %s", &report.Data)

	defer file.Close()

	err = gocsv.MarshalFile(&report.Data, file)

	if err != nil {
		log.Printf("Can not write to report file: %s", err)
	}
}

func EncodeWindows1251(ba []uint8) []uint8 {
	enc := charmap.Windows1251.NewEncoder()
	out, _ := enc.String(string(ba))
	return []uint8(out)
}
