package function

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func DataCheck(input string) (string, error) {
	var out string
	format := "20060102"
	date, err := time.Parse(format, input)
	if err != nil {
		log.Println(err)
		return "", err
	}
	a := date.Format("20060102")
	if a == "00010101" {
		log.Println(err)
		return "", nil
	} else {
		timeNow := time.Now()
		dateNow := timeNow.Format("20060102")
		if a <= dateNow || input == "" {
			out = dateNow
		} else if a > dateNow {
			out = a
		}
	}
	return out, nil
}
func RepeatChek(repeat string) (string, error) {
	var out string
	if (repeat == "y" && len(repeat) == 1) || repeat == "" {
		out = repeat
	} else if strings.Contains(repeat, "d") && len(repeat) > 2 {
		t := strings.Split(repeat, " ")
		if len(t) >= 2 {
			day := t[1]
			i, err := strconv.Atoi(day)
			if err != nil {
				return "0", err
			}
			if i < 401 && i > 0 {
				return repeat, nil
			}
		} else {
			return "0", nil
		}
		return "0", nil
	} else {
		return "0", nil
	}
	return out, nil
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
