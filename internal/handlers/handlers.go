package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/models"
	"strconv"
	"strings"
)

type Handler struct {
	store *database.TaskStore
}

func NewHandler(store *database.TaskStore) *Handler {
	return &Handler{store: store}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// GetAllTasks godoc
// @Summary      Получить список всех задач
// @Description  Возвращает массив всех задач из базы данных
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   models.Task
// @Failure      500  {object}  map[string]string  "Unable to fetch tasks"
// @Router       /tasks [get]
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)

}

// GetTask godoc
// @Summary      Получить задачу по ID
// @Description  Возвращает одну задачу по её идентификатору
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "ID задачи"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string  "Invalid task ID"
// @Router       /tasks/{id} [get]
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]

	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		log.Printf(err.Error())
		return
	}

	task, err := h.store.GetByID(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

// CreateTask godoc
// @Summary      Создать новую задачу
// @Description  Создаёт задачу на основе переданного JSON
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        input  body      models.CreateTaskInput  true  "Данные для создания задачи"
// @Success      201    {object}  models.Task
// @Failure      400    {object}  map[string]string  "Invalid request payload"
// @Failure      500    {object}  map[string]string  "Internal error"
// @Router       /tasks/create [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		log.Printf(err.Error())
		return
	}
	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Empty title")
		return
	}

	task, err := h.store.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary      Обновить задачу
// @Description  Обновляет поля задачи по её ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id     path      int                     true  "ID задачи"
// @Param        input  body      models.UpdateTaskInput  true  "Поля для обновления"
// @Success      200    {object}  models.Task
// @Failure      400    {object}  map[string]string  "Invalid request payload"
// @Failure      404    {object}  map[string]string  "Task not found"
// @Failure      500    {object}  map[string]string  "Internal error"
// @Router       /tasks/{id} [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]

	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		log.Printf(err.Error())
		return
	}
	var input models.UpdateTaskInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	task, err := h.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return

	}

	respondWithJSON(w, http.StatusOK, task)
}

// DeleteTask godoc
// @Summary      Удалить задачу
// @Description  Удаляет задачу по её ID
// @Tags         tasks
// @Produce      json
// @Param        id   path  int  true  "ID задачи"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string  "Invalid task ID"
// @Failure      404  {object}  map[string]string  "Task not found"
// @Failure      500  {object}  map[string]string  "Internal error"
// @Router       /tasks/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
