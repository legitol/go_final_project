package steps

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func CreateBD() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install == true {
		db, err := sql.Open("sqlite", "scheduler.db")

		if err != nil {
			log.Fatal(err)
			return
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler 
		(id INTEGER PRIMARY KEY AUTOINCREMENT, 
		date CHAR(8) NOT NULL DEFAULT '', 
		title VARCHAR(128) NOT NULL DEFAULT '', 
		comment VARCHAR(256) NOT NULL DEFAULT '', 
		repeat VARCHAR(128) NOT NULL DEFAULT '')`,
			`CREATE INDEX date_index ON scheduler (date)`)

		if err != nil {
			log.Fatal(err)
		}
	}
}
