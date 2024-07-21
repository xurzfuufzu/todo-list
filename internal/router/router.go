package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"todoList/internal/handlers"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Get("/healthCheck", handlers.HealthCheck)
		r.Post("/todo-list/tasks", handlers.CreateTask)
		r.Put("/todo-list/tasks/{id}", handlers.UpdateTask)
		r.Delete("/todo-list/tasks/{id}", handlers.DeleteTask)
		r.Put("/todo-list/tasks/{id}/done", handlers.MarkTaskDone)
		r.Get("/todo-list/tasks", handlers.GetTasks)
	})

	return router
}
