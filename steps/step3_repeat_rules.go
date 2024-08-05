package steps

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextdate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" || repeat == " " {
		err := fmt.Errorf("Правило содержит пустую строку: \"%s\"", repeat)
		return "", err
	}

	parseDate, err := time.Parse("20060102", date)
	if err != nil {
		return fmt.Sprintf("Неверный формат даты: \"%s\"", date), err
	}

	split := strings.Split(repeat, " ")

	switch split[0] {
	case "d":
		if len(split) > 1 {
			add, err := strconv.Atoi(split[1])
			if err != nil {
				return fmt.Sprintf("\"%s\" не может быть преобразовано в корректный интервал дней", repeat), err
			}

			if add > 400 {
				err := fmt.Errorf("Превышен максимально допустимый интервал дней: %d", add)
				return "", err
			}

			if add < 1 {
				err := fmt.Errorf("Недопустимое значение интервала дней: %d", add)
				return "", err
				//return fmt.Sprintf("Недопустимое значение интервала дней: %d", add), err
			}

			addDate := parseDate.AddDate(0, 0, add)

			for {
				if addDate.After(now) {
					break
				}
				addDate = addDate.AddDate(0, 0, add)
			}

			return addDate.Format("20060102"), nil
		}
		err := fmt.Errorf("Не указан интервал в днях: \"%s\"", repeat)
		return "", err

	case "y":
		addDate := parseDate.AddDate(1, 0, 0)

		for {
			if addDate.After(now) {
				break
			}
			addDate = addDate.AddDate(1, 0, 0)
		}
		return addDate.Format("20060102"), nil

	default:
		err := fmt.Errorf("Недопустимый символ: \"%s\"", split[0])
		return "", err
	}
}
