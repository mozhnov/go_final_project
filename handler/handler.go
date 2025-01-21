package handler

import (
	"encoding/json"
	"fmt"
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
	err := json.NewDecoder(r.Body).Decode(&TaskAdd)
	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = err.Error()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	}

	id, errAdd := h.Repo.AddTask(TaskAdd)
	strId := strconv.Itoa(id)
	respID := make(map[string]string)
	respID["id"] = strId
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(respID)
	if errAdd != nil {
		respErr := make(map[string]string)
		respErr["error"] = "Не указан заголовок задачи"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	}

}
func (h Handler) GetTasksSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search != "" {
		searchData := h.Repo.SearchTask(search)
		respSearch := make(map[string][]structurs.Tasks)
		respSearch["tasks"] = searchData
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err := json.NewEncoder(w).Encode(respSearch)
		if err != nil {
			respErr := make(map[string]string)
			respErr["error"] = err.Error()
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(respErr)
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
			}
		}
	}
}

func (h Handler) GetTaskId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := h.Repo.GetTaskId(id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = err.Error()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	}
}

func (h Handler) PutTask(w http.ResponseWriter, r *http.Request) {
	var TaskChange structurs.Tasks
	err := json.NewDecoder(r.Body).Decode(&TaskChange)
	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = err.Error()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	}
	err = h.Repo.PutTaskId(TaskChange)

	if err != nil {
		respErr := make(map[string]string)
		respErr["error"] = "Не изменилась"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(respErr)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var out structurs.Empty
		json.NewEncoder(w).Encode(out)
	}
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
