/**
 * Project qualityManagerApi created by exluap
 * Date: 04.03.2019 14:04
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */

package tasks

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"qualityManagerApi/auth"
	"qualityManagerApi/models"
	"qualityManagerApi/tools"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	_ = auth.CheckToken(w, r)

	result, err := tools.ListOfTasks()

	if err != nil {
		log.Print("Error with get tasks list")
		raven.CaptureErrorAndWait(err, nil)
		http.Error(w, "Some problems with query", http.StatusBadRequest)
	} else {
		w.Write(result)
	}
}

func GetTasksByUserOwner(w http.ResponseWriter, r *http.Request) {
	_ = auth.CheckToken(w, r)

	vars := mux.Vars(r)

	result, err := tools.GetTasksByUserOwner(vars["user"])

	if err != nil {
		log.Print("Cant get tasks where user is owner")
		raven.CaptureErrorAndWait(err, nil)
		http.Error(w, "Some problems with query", http.StatusBadRequest)
	} else {
		w.Write(result)
	}
}

func GetTasksByUserAssignee(w http.ResponseWriter, r *http.Request) {
	_ = auth.CheckToken(w, r)

	vars := mux.Vars(r)

	result, err := tools.GetTasksByAssegneeToUser(vars["user"])

	if err != nil {
		log.Print("Cant get tasks where user is owner")
		raven.CaptureErrorAndWait(err, nil)
		http.Error(w, "Some problems with query", http.StatusBadRequest)
	} else {
		w.Write(result)
	}
}

func PostNewTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	var taskInfo map[string]string

	json.Unmarshal(body, &taskInfo)

	data := auth.CheckToken(w, r)

	result := tools.PostNewTask(taskInfo, data.CustomClaims["userid"])

	if result {
		resultation := &models.Resultation{
			Result: "Task was created",
		}

		res, _ := json.Marshal(resultation)

		w.Write(res)

	} else {
		resultation := &models.Resultation{
			Result: "Task was not created",
		}

		res, _ := json.Marshal(resultation)

		w.Write(res)
	}

	if err != nil {
		http.Error(w, "Error with create task", http.StatusInternalServerError)
	}
}

func GetTaskInfo(w http.ResponseWriter, r *http.Request) {
	_ = auth.CheckToken(w, r)

	vars := mux.Vars(r)

	result, err := tools.GetTaskInfo(vars["taskId"])

	if err != nil {
		log.Print("Cant get tasks info")
		raven.CaptureErrorAndWait(err, nil)
		http.Error(w, "Some problems with query", http.StatusBadRequest)
	} else {
		w.Write(result)
	}

}

func PostNewTaskStatus(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	var taskInfo map[string]string

	json.Unmarshal(body, &taskInfo)

	data := auth.CheckToken(w, r)

	vars := mux.Vars(r)

	updateStatus := prepareToUpdate(vars["taskId"], taskInfo["status"], data.CustomClaims["userid"])

	if updateStatus {
		res := &models.Resultation{
			Result: "Status changed to " + taskInfo["status"],
		}

		result, _ := json.Marshal(res)

		w.Write(result)
	} else {
		res := &models.Resultation{
			Result: "Status not changed to " + taskInfo["status"],
		}

		result, _ := json.Marshal(res)

		w.Write(result)
	}

	if err != nil {
		log.Print("Cant update task status")
		raven.CaptureErrorAndWait(err, nil)
		http.Error(w, "Some problems with query:", http.StatusBadRequest)
	}
}

func prepareToUpdate(taskId, status, userId string) bool {
	var result bool

	owner := tools.GetTaskOwner(taskId)
	statusTask := tools.GetTaskStatus(taskId)

	switch status {
	case "In work":
		result = tools.UpdateTaskStatus(taskId, "In work", owner, userId)
		break
	case "In progress":
		result = tools.UpdateTaskStatus(taskId, "In progress", owner, userId)
	case "Not assegnee":
		result = tools.UpdateTaskStatus(taskId, "Not assegnee", owner, userId)
		break
	case "Closed":
		if statusTask == "Not assegnee" {
			result = false
		} else {
			result = tools.UpdateTaskStatus(taskId, "Closed", owner, userId)
		}
		break
	case "Canceled":
		result = tools.UpdateTaskStatus(taskId, "Canceled", owner, userId)
		break
	}

	return result
}
