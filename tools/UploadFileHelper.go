/**
 * Project IntelliJ IDEA created by exluap
 * Date: 29.08.2019 17:35
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"qualityManagerApi/constants"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		log.Printf("Cant parse form with file: %s", err)
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Error read form from uploadFile: %s", err)
	}
	defer file.Close()

	conf := constants.GetConfig()

	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(conf.UploadFilePath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Printf("Error create file: %s", err)
	}

	defer f.Close()
	_, err = io.Copy(f, file)

	if err != nil {
		log.Printf("Can't copy file: %s", err)
	}
}
