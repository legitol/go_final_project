package steps

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServAster() {
	fmt.Println("Запуск сервера")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	mux.HandleFunc("/api/nextdate", nextDateWeekMon)

	mux.HandleFunc("/api/signin", auth)
	mux.HandleFunc("/api/task", authTask(reqSelectAster))
	mux.HandleFunc("/api/tasks/", authTask(searchHandler))
	mux.HandleFunc("/api/task/done", authTask(done))

	portEnv, exists := os.LookupEnv("TODO_PORT")
	var port string
	if exists {
		port = portEnv
	} else {
		port = "7540"
	}

	err := http.ListenAndServe((":" + port), mux)
	fmt.Printf("Слушается порт: %s", port)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Завершение работы")
}

func reqSelectAster(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskWM(w, r)
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
