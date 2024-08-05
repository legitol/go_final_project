package steps

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func AddTaskWM(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "scheduler.db")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Task{Error: err.Error()})
		return
	}
	defer db.Close()

	var task Task
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Task{Error: err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Task{Error: err.Error()})
		return
	}

	if task.Title == "" || task.Title == " " {
		w.WriteHeader(http.StatusBadRequest)
		wOut(w, Task{Error: "Не указан заголовок задачи"})
		return
	}

	if task.Date == "" || task.Date == " " {
		task.Date = time.Now().Format("20060102")
	}

	if task.Date == "today" {
		task.Date = time.Now().Format("20060102")
	}

	parseDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Task{Error: err.Error()})
		return
	}

	if parseDate.Before(time.Now()) {
		if task.Repeat == "" || task.Repeat == " " || task.Date == time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		} else {
			task.Date, err = NextDateWM(time.Now(), task.Date, task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				wOut(w, Task{Error: "Oшибка функции вычисления даты выполнения задачи"})
				return
			}
		}
	} else {
		task.Date = task.Date
	}

	insertId, err := Insert(db, task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		wOut(w, Task{Error: "Oшибка функции добавления записи в БД"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	wOut(w, Task{ID: insertId})

}
