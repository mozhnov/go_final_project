package function

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func DataCheck(search string) string {
	out := search
	dateLay := "02.01.2006"
	layout := "20060102"
	date, err := time.Parse(dateLay, search)
	search = date.Format(layout)

	if err != nil {
		return out
	}
	return search
}
func DBconnect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	DBFile := os.Getenv("TODO_DBFILE")
	db, err := sql.Open("sqlite", DBFile)
	if err != nil {
		log.Println(err)
	}
	return db
}
