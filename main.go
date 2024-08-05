package main

import (
	"github.com/legitol/go_final_project/steps"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}

func main() {
	//задание со звёздочкой _asterisk/Aster
	steps.CreateBDAster()
	steps.StartServAster()

	//задание без звёздочки
	// steps.CreateBD()
	// steps.StartServ()
}
