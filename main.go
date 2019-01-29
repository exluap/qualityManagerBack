/**
 * Project qualityManagerApi created by exluap
 * Date: 25.10.2018 23:57
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"qualityManagerApi/auth"
	"qualityManagerApi/constants"
	"qualityManagerApi/queries"
	"qualityManagerApi/user"
)

func main() {

	startWebServer()

}

func startWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", hello)

	//User Group
	userGroup := r.PathPrefix("/api/user").Subrouter()
	userGroup.HandleFunc("/overtime", user.CheckOver)
	userGroup.HandleFunc("/info", user.GetInfoAboutUser).Methods("GET")

	//Query Group
	queryGroup := r.PathPrefix("/api/query").Subrouter()
	queryGroup.HandleFunc("/info", queries.GetQuery).Methods("GET")
	queryGroup.HandleFunc("/list", queries.GetQueries).Methods("GET")
	queryGroup.HandleFunc("/add", queries.AddQuery).Methods("POST")
	queryGroup.HandleFunc("/delete", queries.DeleteSR).Methods("POST")
	queryGroup.HandleFunc("/helper", queries.GenerateNote).Methods("POST")

	//Auth Group
	authGroup := r.PathPrefix("/api/auth").Subrouter()
	authGroup.HandleFunc("/login", auth.Login).Methods("GET")
	authGroup.HandleFunc("/changepassword", auth.ChangePassword).Methods("POST")

	//r.HandleFunc("/add_log", tools.AddingLog)

	log.Println("Listening for connections on port: ", constants.PORT)
	log.Fatal(http.ListenAndServe(":"+constants.PORT, r))

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, please auth")
}
