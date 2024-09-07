package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rogue0026/task-tracker/internal/models"
	"github.com/rogue0026/task-tracker/internal/storage"
)

type TasksStorage struct {
	conn *pgxpool.Pool
}

func New(dsn string) (TasksStorage, error) {
	const fn = "internal.storage.postgres.New"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return TasksStorage{}, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return TasksStorage{}, err
	}

	s := TasksStorage{
		conn: pool,
	}

	return s, nil
}

func (s *TasksStorage) SaveTask(t models.Task) error {
	const fn = "internal.storage.postgres.SaveTask"
	q := "INSERT INTO t_tasks (name, deadline, user_id) VALUES ($1, $2, $3)"
	_, err := s.conn.Exec(context.Background(), q, t.Name, t.Deadline, t.UserID)
	if err != nil {
		return fmt.Errorf("func=%s error=%w", fn, err)
	}
	return nil
}

func (s *TasksStorage) DeleteTask(taskName string, userID int64) error {
	const fn = "internal.storage.postgres.DeleteTask"
	q := "DELETE FROM t_tasks WHERE name = $1 and user_id = $2"
	_, err := s.conn.Exec(context.Background(), q, taskName, userID)
	if err != nil {
		return fmt.Errorf("func=%s error=%w", fn, err)
	}

	return nil
}

func (s *TasksStorage) UserTasks(userID int64) ([]models.Task, error) {
	const fn = "internal.storage.postgres.UserTasks"
	q := "SELECT name, deadline FROM t_tasks WHERE user_id = $1 ORDER BY deadline"
	rows, err := s.conn.Query(context.Background(), q, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNoTasksForUser
		}
		return nil, fmt.Errorf("func=%s error=%w", fn, err)
	}
	defer rows.Close()

	userTasks := make([]models.Task, 0)

	for rows.Next() {
		var t models.Task
		if err = rows.Scan(&t.Name, &t.Deadline); err != nil {
			return nil, fmt.Errorf("func=%s error=%w", fn, err)
		}
		userTasks = append(userTasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("func=%s error=%w", fn, err)
	}

	if len(userTasks) == 0 {
		return nil, storage.ErrNoTasksForUser
	}

	return userTasks, nil
}
