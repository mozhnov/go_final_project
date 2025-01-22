package handler

import (
	"encoding/json"
	"fmt"
	"go_final_project/function"
	"go_final_project/structurs"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Repo HandlerRepository
}
type HandlerRepository interface {
	AddTask(t structurs.Task) (int, error)
	AddDB(DBFile string)
	CheckTable(DBFile string) error
	CreateTable(DBFile string)
	GetTasks() ([]structurs.Tasks, error)
	GetTaskId(id string) (structurs.Tasks, error)
	SearchTask(search string) []structurs.Tasks
	PutTaskId(t structurs.Tasks) error
	DeleteTaskId(id string) error
	NextDate(d structurs.DataValid) (string, error)
}

func NewHandler(repo HandlerRepository) Handler {
	return Handler{repo}
}
func (h Handler) PostGetPutDeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.PostTask(w, r)
	} else if r.Method == http.MethodGet {
		h.GetTaskId(w, r)
	} else if r.Method == http.MethodDelete {
		h.DeleteTaskID(w, r)
	} else if r.Method == http.MethodPut {
		h.PutTask(w, r)
	}

}
func (h Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	var TaskAdd structurs.Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&TaskAdd)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	date, _ := function.DataCheck(TaskAdd.Date)
	repeat, _ := function.RepeatChek(TaskAdd.Repeat)
	TaskAdd.Repeat = repeat
	TaskAdd.Date = date
	if TaskAdd.Title != "" && TaskAdd.Date != "" && TaskAdd.Repeat != "0" {
		id, errAdd := h.Repo.AddTask(TaskAdd)
		respId := strconv.Itoa(id)
		json.NewEncoder(w).Encode(map[string]string{"id": respId})
		if errAdd != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": errAdd.Error()})
			return
		}
	} else {
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
}
func (h Handler) GetTasksSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search != "" {
		input := function.SearcCheck(search)
		searchData := h.Repo.SearchTask(input)
		respSearch := make(map[string][]structurs.Tasks)
		respSearch["tasks"] = searchData
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err := json.NewEncoder(w).Encode(respSearch)
		if err != nil {
			respErr := make(map[string]string)
			respErr["error"] = err.Error()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(respErr)
			return
		}
	} else {
		tasks, err := h.Repo.GetTasks()
		if tasks != nil {
			respTasks := make(map[string][]structurs.Tasks)
			respTasks["tasks"] = tasks
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(respTasks)
			if err != nil {
				respErr := make(map[string]string)
				respErr["error"] = err.Error()
				json.NewEncoder(w).Encode(respErr)
				return
			}

		} else {
			nulSlais := make([]string, 0)
			respNil := make(map[string][]string)
			respNil["tasks"] = nulSlais
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			err := json.NewEncoder(w).Encode(respNil)
			if err != nil {
				respErr := make(map[string]string)
				respErr["error"] = err.Error()
				json.NewEncoder(w).Encode(respErr)
				return
			}
		}
	}
}

func (h Handler) GetTaskId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := h.Repo.GetTaskId(id)
	resp := make(map[string]string)
	if err != nil {
		resp["error"] = "Задача не найдена"
		json.NewEncoder(w).Encode(resp)
		return
	}
	json.NewEncoder(w).Encode(task)
	w.Header().Set("Content-Type", "application/json")
}

func (h Handler) PutTask(w http.ResponseWriter, r *http.Request) {
	var TaskChange structurs.Tasks
	resp := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&TaskChange)
	if err != nil {
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = h.Repo.PutTaskId(TaskChange)
	if err != nil {
		resp["error"] = "Задача не найдена"
		json.NewEncoder(w).Encode(resp)
		return
	}
	out := &structurs.Empty{}
	json.NewEncoder(w).Encode(out)
	w.Header().Set("Content-Type", "application/json")
}

func (h Handler) DoneTaskId(w http.ResponseWriter, r *http.Request) {
	nowNow := time.Now()
	format := "20060102"
	now := nowNow.Format(format)
	id := r.URL.Query().Get("id")
	task, err := h.Repo.GetTaskId(id)
	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = err.Error()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	}
	var nextDate structurs.DataValid
	nextDate.Data = task.Date
	nextDate.Now = now
	nextDate.Repeat = task.Repeat
	if task.Repeat == "" {
		err := h.Repo.DeleteTaskId(id)
		if err != nil {
			respErr := make(map[string]string)
			respErr["error"] = "Не завершилась"
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(respErr)
			fmt.Println("err", respErr)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			var out structurs.Empty
			json.NewEncoder(w).Encode(out)
		}
	} else {
		fmt.Println("eeeeeeeee")
		data, err := h.Repo.NextDate(nextDate)
		fmt.Println("tttt", data)

		if err != nil {
			respErr := make(map[string]string)
			respErr["error"] = err.Error()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(respErr)
		}
		task.Date = data
		err = h.Repo.PutTaskId(task)
		if err != nil {
			respErr := make(map[string]string)
			respErr["error"] = err.Error()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(respErr)
		}
	}
}
func (h Handler) DeleteTaskID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.Repo.DeleteTaskId(id)
	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = "Не удалилась"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
		fmt.Println("err", respErr)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var out structurs.Empty
		json.NewEncoder(w).Encode(out)
	}
}
func (h Handler) NextData(w http.ResponseWriter, r *http.Request) {
	var nextDate structurs.DataValid
	nextDate.Now = r.URL.Query().Get("now")
	nextDate.Data = r.URL.Query().Get("date")
	nextDate.Repeat = r.URL.Query().Get("repeat")
	nextdate, err := h.Repo.NextDate(nextDate)
	answer, _ := strconv.Atoi(nextdate)
	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = err.Error()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(answer)

}
