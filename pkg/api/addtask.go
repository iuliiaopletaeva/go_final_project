package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"main/pkg/db"
)

func addTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(body, &task); err != nil {
		writeJSON(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": "task has no title"})
		return
	}

	if err = checkDate(&task); err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	taskId, err := db.AddTask(&task)
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(res, http.StatusOK, map[string]interface{}{"id": taskId})

}

func writeJSON(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(data)
}

func checkDate(task *db.Task) error {
	var next string
	now := time.Now()
	today := now.Format(dateFormat)

	if task.Date == "today" || task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("cannot parse date")
	}

	if afterNow(now, t) {
		if task.Repeat != "" {
			next, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		} else {
			task.Date = today
		}
	}

	return nil
}
