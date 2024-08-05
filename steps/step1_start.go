package steps

import (
	"fmt"
	"log"
	"net/http"
)

func StartServ() {
	fmt.Println("Запуск сервера")
	mux := http.NewServeMux()
	
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/nextdate", nextdate)
	mux.HandleFunc("/api/task", reqSelect)
	mux.HandleFunc("/api/tasks", getList)
	mux.HandleFunc("/api/task/done", done)
	
	err := http.ListenAndServe(":7540", mux)
	fmt.Println("Слушается порт: 7540")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Завершение работы")
}

func reqSelect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTask(w, r)
		return
	case http.MethodGet:
		GetTaskId(w, r)
		return
	case http.MethodPut:
		EditTask(w, r)
		return
	case http.MethodDelete:
		DeleteTask(w, r)
		return
	}
}
