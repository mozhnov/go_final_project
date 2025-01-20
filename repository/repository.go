package repository

import (
	"database/sql"
	"go_final_project/function"
	"go_final_project/structurs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func NewDBwork(db *sql.DB) DB {
	return DB{db: db}
}
func (s DB) AddDB(DBFile string) {
	_, err := os.Create(DBFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL CHECK(date !=''), title TEXT NOT NULL CHECK(title !=''), comment TEXT, repeat TEXT);")
	if err != nil {
		log.Println(err, "create_tablr err")
	}
	_, err = s.db.Exec("CREATE INDEX ID_Date ON scheduler (date);")
	if err != nil {
		log.Println(err, "create index err")
	}
}
func (s DB) CheckTable(DBFile string) error {
	_, err := s.db.Query("SELECT * FROM scheduler;")
	if err != nil {
		return err
	}
	return nil
}
func (s DB) CreateTable(DBFile string) {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL CHECK(date !=''), title TEXT NOT NULL CHECK(title !=''), comment TEXT, repeat TEXT);")
	if err != nil {
		log.Println(err, "create_tablr err")
	}
	_, err = s.db.Exec("CREATE INDEX ID_Date ON scheduler (date);")
	if err != nil {
		log.Println(err, "create index err")
	}
}
func (s DB) AddTask(t structurs.Task) (int, error) {
	input := function.DataCheck(t.Date)
	layout := "20060102"
	t.Date = input
	timeNow := time.Now()
	dateNow := timeNow.Format(layout)
	if t.Date < dateNow {
		t.Date = dateNow
	}
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}
func (s DB) GetTasks() ([]structurs.Tasks, error) {
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC;")
	if err != nil {
		return []structurs.Tasks{}, err
	}
	defer rows.Close()

	var res []structurs.Tasks
	for rows.Next() {
		p := structurs.Tasks{}
		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return []structurs.Tasks{}, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return []structurs.Tasks{}, err
	}

	return res, nil
}
func (s DB) GetTaskId(id string) (structurs.Tasks, error) {
	p := structurs.Tasks{}
	rows := s.db.QueryRow("SELECT * FROM scheduler WHERE id = :id;",
		sql.Named("id", id))
	err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
	return p, err
}
func (s DB) SearchTask(search string) []structurs.Tasks {
	input := function.DataCheck(search)
	rows, err := s.db.Query("SELECT * FROM scheduler WHERE id LIKE CONCAT('%', :input, '%') OR date LIKE CONCAT('%', :input, '%') OR title LIKE CONCAT('%', :input, '%') OR comment LIKE CONCAT('%', :input, '%') ORDER BY date DESC;",
		sql.Named("input", input))
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var res []structurs.Tasks
	for rows.Next() {
		p := structurs.Tasks{}

		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			log.Println(err)
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return res

}
func (s DB) PutTaskId(t structurs.Tasks) error {
	input := function.DataCheck(t.Date)
	layout := "20060102"
	t.Date = input
	timeNow := time.Now()
	dateNow := timeNow.Format(layout)
	if t.Date < dateNow {
		t.Date = dateNow
	}
	_, err := s.db.Exec("UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id;",
		sql.Named("id", t.Id),
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return err
	}
	return nil
}
func (s DB) DeleteTaskId(id string) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id;",
		sql.Named("id", id))

	if err != nil {
		return err
	}
	return nil
}
func (s DB) NextDate(now time.Time, date string, repeat string) (string, error) {
	format := "20060102"
	dateTimeFormat, err := time.Parse(format, date)
	if err != nil {
		log.Println(err)
	}
	var yearAdd time.Time
	if dateTimeFormat.After(now) && repeat == "y" {
		yearAdd = dateTimeFormat.AddDate(1, 0, 0)
	} else if dateTimeFormat.After(now) && strings.Contains(repeat, "d") {
		repeatSplit := strings.Split(repeat, " ")
		day := repeatSplit[1]
		i, err := strconv.Atoi(day)
		if err != nil {
			log.Println(err)
		}
		Add := dateTimeFormat.AddDate(0, 0, i)
		dayAdd := Add.Format(format)
		return dayAdd, err
	} else if dateTimeFormat.Before(now) && repeat == "y" {
		var yearAddbefore string
		for dateTimeFormat.Before(now) {
			Add := dateTimeFormat.AddDate(1, 0, 0)
			dateTimeFormat = Add
			yearAddbefore = Add.Format(format)
		}
		return yearAddbefore, err
	} else if dateTimeFormat.Before(now) && strings.Contains(repeat, "d") {
		repeatSplit := strings.Split(repeat, " ")
		day := repeatSplit[1]
		i, err := strconv.Atoi(day)
		if err != nil {
			log.Println(err)
		}
		var dayAddbefore string
		for dateTimeFormat.Before(now) {
			Add := dateTimeFormat.AddDate(0, 0, i)
			dateTimeFormat = Add
			dayAddbefore = Add.Format(format)
		}
		return dayAddbefore, err
	} else {
		return "", err
	}
	return yearAdd.Format(format), err
}
