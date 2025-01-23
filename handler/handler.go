package handler

import (
	"encoding/json"
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
	AddDB(DBFile string) error
	CheckTable(DBFile string) error
	CreateTable(DBFile string) error
	AddTask(t structurs.Task) (int, error)
	GetTasks() ([]structurs.Tasks, error)
	GetTaskId(id string) (structurs.Tasks, error)
	SearchTask(input string) ([]structurs.Tasks, error)
	PutTaskId(t structurs.Tasks) error
	DeleteTaskId(id string) error
	NextDate(now string, data string, repeat string) (string, error)
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
	resp := make(map[string]string)
	var TaskAdd structurs.Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&TaskAdd)
	if err != nil {
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json")
		return
	}
	date, _ := function.DataCheck(TaskAdd.Date)
	repeat, _ := function.RepeatChek(TaskAdd.Repeat)
	TaskAdd.Repeat = repeat
	TaskAdd.Date = date
	if TaskAdd.Title != "" && TaskAdd.Date != "" && TaskAdd.Repeat != "0" {
		id, errAdd := h.Repo.AddTask(TaskAdd)
		respId := strconv.Itoa(id)
		resp["id"] = respId
		json.NewEncoder(w).Encode(resp)
		if errAdd != nil {
			resp["error"] = errAdd.Error()
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json")
			return
		}
	} else {
		resp["error"] = "Не указан заголовок задачи"
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json")
		return
	}
}
func (h Handler) GetTasksSearch(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	search := r.URL.Query().Get("search")
	if search != "" {
		input := function.SearcCheck(search)
		searchData, err := h.Repo.SearchTask(input)
		if err != nil {
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json")
			return
		}
		respSearch := make(map[string][]structurs.Tasks)
		respSearch["tasks"] = searchData
		err = json.NewEncoder(w).Encode(respSearch)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json")
			return
		}
	} else {
		tasks, err := h.Repo.GetTasks()
		if tasks != nil {
			respTasks := make(map[string][]structurs.Tasks)
			respTasks["tasks"] = tasks
			json.NewEncoder(w).Encode(respTasks)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err != nil {
				resp["error"] = err.Error()
				json.NewEncoder(w).Encode(resp)
				w.Header().Set("Content-Type", "application/json")
				return
			}

		} else {
			nulSlais := make([]string, 0)
			respNil := make(map[string][]string)
			respNil["tasks"] = nulSlais
			err := json.NewEncoder(w).Encode(respNil)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err != nil {
				resp["error"] = err.Error()
				json.NewEncoder(w).Encode(resp)
				w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
		return
	}
	id, errA := function.IdCheck(TaskChange.Id)
	date, _ := function.DataCheck(TaskChange.Date)
	repeat, _ := function.RepeatChek(TaskChange.Repeat)
	TaskChange.Id = id
	TaskChange.Repeat = repeat
	TaskChange.Date = date
	if TaskChange.Title != "" && TaskChange.Date != "" && TaskChange.Repeat != "0" {
		err = h.Repo.PutTaskId(TaskChange)
		if err != nil || errA != nil {
			resp["error"] = "Задача не найдена"
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json")
			return
		}
		out := &structurs.Empty{}
		json.NewEncoder(w).Encode(out)
		w.Header().Set("Content-Type", "application/json")
	} else {
		resp["error"] = "Задача не найдена"
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json")
		return
	}
}

func (h Handler) DoneTaskId(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	nowNow := time.Now()
	format := "20060102"
	now := nowNow.Format(format)
	id := r.URL.Query().Get("id")
	task, err := h.Repo.GetTaskId(id)
	if err != nil {
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		return
	}
	date := task.Date
	repeat := task.Repeat
	if task.Repeat == "" {
		err := h.Repo.DeleteTaskId(id)
		if err != nil {
			resp["error"] = "Не завершилась"
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return
		} else {
			out := &structurs.Empty{}
			json.NewEncoder(w).Encode(out)
			w.Header().Set("Content-Type", "application/json")
		}
	} else {
		data, err := h.Repo.NextDate(now, date, repeat)
		if err != nil {
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return
		}
		task.Date = data
		err = h.Repo.PutTaskId(task)
		if err != nil {
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			return

		} else {
			out := &structurs.Empty{}
			json.NewEncoder(w).Encode(out)
			w.Header().Set("Content-Type", "application/json")
		}

	}
}
func (h Handler) DeleteTaskID(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	id := r.URL.Query().Get("id")
	idCheck, err := function.IdCheck(id)
	if err != nil {
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		return
	}
	errDel := h.Repo.DeleteTaskId(idCheck)
	if errDel != nil {
		resp["error"] = errDel.Error()
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		return
	} else {
		out := &structurs.Empty{}
		json.NewEncoder(w).Encode(out)
		w.Header().Set("Content-Type", "application/json")
	}
}
func (h Handler) NextData(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	nextdate, err := h.Repo.NextDate(now, date, repeat)
	answer, _ := strconv.Atoi(nextdate)
	if err != nil {
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		return
	}
	json.NewEncoder(w).Encode(answer)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}
