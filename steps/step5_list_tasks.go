package steps

import (
	"database/sql"
	"net/http"
)

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Err struct {
	Error string `json:"error,omitempty"`
}

func getList(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "scheduler.db")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Err{Error: err.Error()})
		return
	}

	err = db.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Err{Error: err.Error()})
		return
	}
	defer db.Close()

	tasks, err := Scan(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Err{Error: "Ошибка чтения БД"})
		return
	}

	w.WriteHeader(http.StatusOK)
	wOut(w, Tasks{Tasks: tasks})
}

func Scan(db *sql.DB) ([]Task, error) {
	tasks := make([]Task, 0, 20)
	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date ASC LIMIT ?", 20)

	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil

}
