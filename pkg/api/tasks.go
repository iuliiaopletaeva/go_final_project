package api

import (
	"main/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(res http.ResponseWriter, req *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(res, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

func getTaskInfoHandler(res http.ResponseWriter, req *http.Request) {
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

	writeJSON(res, http.StatusOK, taskInfo)
}
