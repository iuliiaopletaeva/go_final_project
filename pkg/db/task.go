package db

import "errors"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(task *Task) (int64, error) {
	if db == nil {
		panic("database not initialized")
	}

	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	if db == nil {
		return []*Task{}, errors.New("database not initialized")
	}

	if limit <= 0 {
		return []*Task{}, errors.New("limit has to be more than 0")
	}

	tasks := make([]*Task, 0, limit)

	query := "SELECT * FROM scheduler ORDER BY date ASC LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	task := &Task{}

	query := "SELECT * FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, errors.New("cannot find task with such id")
	}

	return task, err
}

func UpdateTask(task *Task) error {
	if db == nil {
		return errors.New("database not initialized")
	}

	query := `UPDATE scheduler SET 
		date = ?,
		title = ?,
		comment = ?,
		repeat = ?
	WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {
	if db == nil {
		return errors.New("database not initialized")
	}

	if id == "" {
		return errors.New("cannot delete task without id")
	}

	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect id for updating task")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	if db == nil {
		return errors.New("database not initialized")
	}

	query := `UPDATE scheduler SET 
		date = ?
	WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("incorrect id for updating task")
	}
	return nil
}
