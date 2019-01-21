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
	"github.com/rs/cors"
	"log"
	"net/http"
	"qualityManagerApi/auth"
	"qualityManagerApi/constants"
	"qualityManagerApi/queries"
	"qualityManagerApi/tools"
	"qualityManagerApi/user"
)

func main() {

	startWebServer()

}

func startWebServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/get_query", queries.GetQuery)
	mux.HandleFunc("/add_query", queries.AddQuery)
	mux.HandleFunc("/get_queries", queries.GetQueries)
	mux.HandleFunc("/generate_note_and_instruction", queries.GenerateNote)
	mux.HandleFunc("/delete_sr", queries.DeleteSR)
	mux.HandleFunc("/add_log", tools.AddingLog)
	mux.HandleFunc("/in_over", user.CheckOver)

	handler := cors.AllowAll().Handler(mux)

	log.Println("Listening for connections on port: ", constants.PORT)
	log.Fatal(http.ListenAndServe(":"+constants.PORT, handler))

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, please auth")
}
