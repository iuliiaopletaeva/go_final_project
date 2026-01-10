package api

import (
	"main/pkg/db"
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	var now time.Time

	nowParam := req.FormValue("now")

	if nowParam == "" {
		now = time.Now()
	} else {
		parsedNow, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(res, "invalid now parameter", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))
}

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(res, req)
	case http.MethodGet:
		getTaskInfoHandler(res, req)
	case http.MethodPut:
		updateTaskHandler(res, req)
	case http.MethodDelete:
		deleteTaskHandler(res, req)
	default:
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskDoneHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		taskId := req.URL.Query().Get("id")

		if len(taskId) == 0 {
			writeJSON(res, http.StatusBadRequest, map[string]string{"error": "cannot find task without id"})
			return
		}

		taskInfo, err := db.GetTask(taskId)
		if err != nil {
			writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if taskInfo.Repeat == "" {
			err = db.DeleteTask(taskId)
			if err != nil {
				writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		} else {
			next, err := NextDate(time.Now(), taskInfo.Date, taskInfo.Repeat)
			if err != nil {
				writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}

			err = db.UpdateDate(next, taskId)
			if err != nil {
				writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		}

		writeJSON(res, http.StatusOK, map[string]interface{}{})
	} else {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
	}
}
