package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
	"todoList/internal/models"
)

var (
	db sync.Map
)

func ValidateTask(task *models.TaskDTO) error {
	if task.Title == "" {
		return errors.New("title is required")
	}
	if len(task.Title) > 200 {
		return errors.New("title must not be more 200 characters")
	}
	if task.ActiveAt == "" {
		return errors.New("activeAt is required")
	}
	if !isValidDate(task.ActiveAt) {
		return errors.New("invalid date format")
	}
	if DuplicateExist(task.Title, task.ActiveAt) {
		return errors.New("task with the same title and activeAt already exists")
	}
	return nil
}

func isValidDate(activeAt string) bool {
	_, err := time.Parse("2006-01-02", activeAt)
	return err == nil
}

func DuplicateExist(title, activeAt string) bool {
	var foundDuplicate bool
	db.Range(func(_, value interface{}) bool {
		Task := value.(*models.Task)
		if Task.Title == title && Task.ActiveAt == activeAt {
			foundDuplicate = true
			return false
		}
		return true
	})
	return foundDuplicate
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO models.TaskDTO
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidateTask(&taskDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := models.NewTask(taskDTO.Title, taskDTO.ActiveAt)

	db.Store(task.ID, task)

	//db.Range(func(key, value interface{}) bool {
	//	t := value.(*models.Task)
	//	fmt.Printf("Task ID: %s\n", t.ID)
	//	fmt.Printf("Title: %s\n", t.Title)
	//	fmt.Printf("ActiveAt: %s\n", t.ActiveAt)
	//	fmt.Printf("Status: %s\n", t.Status)
	//	fmt.Println("-----")
	//})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]uuid.UUID{"id ": task.ID}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	taskID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, exists := db.Load(taskID)
	if !exists {
		log.Printf("Task not found in database for ID: %s", id)
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	var taskDTO models.TaskDTO

	err = json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidateTask(&taskDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldTask := task.(*models.Task)

	oldTask.Title = taskDTO.Title
	oldTask.ActiveAt = taskDTO.ActiveAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	taskID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, exists := db.Load(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	db.Delete(taskID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	taskID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	t, exists := db.Load(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task := t.(*models.Task)

	task.Status = "done"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	now := time.Now()

	var tasks []*models.Task

	db.Range(func(_, value interface{}) bool {
		task := value.(*models.Task)

		if task.Status != status {
			return true
		}

		activeAt, err := time.Parse("2006-01-02", task.ActiveAt)
		if err != nil {
			return true
		}

		if task.Status == "active" && activeAt.After(now) {
			return true
		}

		tasks = append(tasks, task)
		return true
	})

	for _, task := range tasks {
		activeAt, _ := time.Parse("2006-01-02", task.ActiveAt)
		if activeAt.Weekday() == time.Saturday || activeAt.Weekday() == time.Sunday {
			task.Title = "ВЫХОДНОЙ - " + task.Title
		}
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return byActiveAt(tasks, i, j)
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func byActiveAt(t []*models.Task, i int, j int) bool {
	date1, _ := time.Parse("2006-01-02", t[i].ActiveAt)
	date2, _ := time.Parse("2006-01-02", t[j].ActiveAt)

	return date1.Before(date2)
}
