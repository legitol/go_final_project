package steps

import (
	"database/sql"
	"net/http"
	"time"
)

func done(w http.ResponseWriter, r *http.Request) {
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

	r.Method = http.MethodPost
	id := r.FormValue("id")
	task, err := ScanId(db, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		wOut(w, Err{Error: "Задача не найдена"})
		return
	}

	if task.Repeat == "" || task.Repeat == " " {
		err = DeleteId(db, id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			wOut(w, Err{Error: "Oшибка функции удаления записи БД"})
			return
		}

		w.WriteHeader(http.StatusOK)
		wOut(w, Task{})
		return
	}

	nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("id", task.ID),
		sql.Named("date", nextDate))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Err{Error: "Oшибка функции обновления даты задачи в записи таблицы БД"})
		return
	}

	w.WriteHeader(http.StatusOK)
	wOut(w, Task{})
}

func DeleteId(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
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

	id := r.FormValue("id")
	_, err = ScanId(db, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		wOut(w, Err{Error: "Задача не найдена"})
		return
	}

	err = DeleteId(db, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		wOut(w, Err{Error: "Ошибка функции удаления записи БД"})
		return
	}

	w.WriteHeader(http.StatusOK)
	wOut(w, Task{})

}
