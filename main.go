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
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"qualityManagerApi/adminTools"
	"qualityManagerApi/auth"
	"qualityManagerApi/constants"
	"qualityManagerApi/queries"
	"qualityManagerApi/tasks"
	"qualityManagerApi/tools"
	"qualityManagerApi/user"
)

const version = "5.1.0"

func init() {
	err := raven.SetDSN("https://b65f1572d92948cfbd2c5a2bb3e4adc2:1ca4a46c2c7f408fbaf655de030a0e4f@log.exluap.com/2")

	raven.SetRelease("4.0.0")

	if err != nil {
		log.Print(err)
	}
}

func main() {
	startWebServer()

}

func startWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", hello)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                  // All origins
		AllowedMethods:   []string{"GET", "POST", "PUT"}, // Allowing only get, just an example
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	//User Group
	userGroup := r.PathPrefix("/api/user").Subrouter()
	userGroup.HandleFunc("/overtime", raven.RecoveryHandler(user.CheckOver)).Methods("POST")
	userGroup.HandleFunc("/info", raven.RecoveryHandler(user.GetInfoAboutUser)).Methods("GET")
	userGroup.HandleFunc("/changelogin", raven.RecoveryHandler(user.ChangeLogin)).Methods("POST")

	//Query Group
	queryGroup := r.PathPrefix("/api/query").Subrouter()
	queryGroup.HandleFunc("/info", raven.RecoveryHandler(queries.GetQuery)).Methods("POST")
	queryGroup.HandleFunc("/list", raven.RecoveryHandler(queries.GetQueries)).Methods("POST")
	queryGroup.HandleFunc("/add", raven.RecoveryHandler(queries.AddQuery)).Methods("POST")
	queryGroup.HandleFunc("/delete", raven.RecoveryHandler(queries.DeleteSR)).Methods("POST")
	queryGroup.HandleFunc("/helper", raven.RecoveryHandler(queries.GenerateNote)).Methods("POST")

	//Auth Group
	authGroup := r.PathPrefix("/api/auth").Subrouter()
	authGroup.HandleFunc("/login", raven.RecoveryHandler(auth.Login)).Methods("POST")
	authGroup.HandleFunc("/changepassword", raven.RecoveryHandler(auth.ChangePassword)).Methods("POST")

	//Admin Func

	adminGroup := r.PathPrefix("/api/admin").Subrouter()
	adminGroup.HandleFunc("/createNewUser", raven.RecoveryHandler(adminTools.CreateNewUser)).Methods("POST")
	adminGroup.HandleFunc("/updateUserInfo", raven.RecoveryHandler(adminTools.UpdateUser)).Methods("POST")

	//Tasks Group

	tasksGroup := r.PathPrefix("/api/tasks").Subrouter()
	tasksGroup.HandleFunc("/list", tasks.GetAllTasks).Methods("GET")
	tasksGroup.HandleFunc("/{user}/owner", tasks.GetTasksByUserOwner).Methods("GET")
	tasksGroup.HandleFunc("/{user}/assignee", tasks.GetTasksByUserAssignee).Methods("GET")
	tasksGroup.HandleFunc("/createNewTask", tasks.PostNewTask).Methods("POST")
	tasksGroup.HandleFunc("/{taskId}/info", tasks.GetTaskInfo).Methods("GET")
	tasksGroup.HandleFunc("/{taskId}/statusUpdate", tasks.PostNewTaskStatus).Methods("POST")
	tasksGroup.HandleFunc("/{taskId}/update", tasks.UpdateTaskInfo).Methods("POST")

	//File Utils
	r.HandleFunc("/api/upload", tools.GetFile).Methods("POST")

	//Report Generator
	r.HandleFunc("/api/report", tools.GenerateReport).Methods("POST")

	//r.HandleFunc("/add_log", tools.AddingLog)

	conf := constants.GetConfig()

	log.Println("Current version:", version)
	log.Println("Listening for connections on port: ", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, c.Handler(r)))

}

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "The request is not correct.  Please read the documentation")

	if err != nil {
		log.Printf("Error with handle /: %s", err)
	}
}
