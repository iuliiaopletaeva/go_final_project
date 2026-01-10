package api

import (
	"encoding/json"
	"io"
	"main/pkg/db"
	"net/http"
)

func updateTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task db.Task

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(body, &task); err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
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

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(res, http.StatusOK, map[string]interface{}{})
}
