package api

import (
	"main/pkg/db"
	"net/http"
)

func deleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("id")

	if len(taskId) == 0 {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": "cannot find task without id"})
		return
	}

	err := db.DeleteTask(taskId)
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(res, http.StatusOK, map[string]interface{}{})
}
