package database

import (
	"restapi-tasks/internal/models"

	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func GetAll(s *TaskStore) ([]models.Task, error) {
	var tasks []models.Task

	query := `SELECT * FROM tasks order by created_at desc;`

	err := s.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
