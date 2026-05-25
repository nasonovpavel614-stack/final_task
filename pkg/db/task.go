package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrTaskNotFound = errors.New("Задача не найдена")

const (
	DateFormat        = "20060102"
	dateDisplayFormat = "02.01.2006"
	TasksLimitDefault = 50
)

// Task описывает задачу планировщика.
type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает id новой записи.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetTask возвращает задачу по идентификатору.
func GetTask(id string) (*Task, error) {
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	var taskID int64
	var task Task
	err = DB().QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		idNum,
	).Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	task.ID = strconv.FormatInt(taskID, 10)
	return &task, nil
}

// UpdateTask обновляет задачу в таблице scheduler.
func UpdateTask(task *Task) error {
	idNum, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return ErrTaskNotFound
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idNum)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// DeleteTask удаляет задачу по идентификатору.
func DeleteTask(id string) error {
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ErrTaskNotFound
	}

	res, err := DB().Exec(`DELETE FROM scheduler WHERE id = ?`, idNum)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// UpdateDate обновляет дату задачи.
func UpdateDate(next, id string) error {
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ErrTaskNotFound
	}

	res, err := DB().Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, idNum)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// Tasks возвращает ближайшие задачи, отсортированные по дате.
// Если search не пустой — фильтрует по подстроке в title/comment или по дате (02.01.2006).
func Tasks(limit int, search string) ([]*Task, error) {
	if limit <= 0 {
		limit = TasksLimitDefault
	}

	if search == "" {
		today := time.Now().Format(DateFormat)
		return queryTasks(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date LIMIT ?`,
			today, limit,
		)
	}

	if t, err := time.Parse(dateDisplayFormat, search); err == nil {
		return queryTasks(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`,
			t.Format(DateFormat), limit,
		)
	}

	pattern := strings.ToLower(search)
	allTasks, err := queryTasks(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date`,
	)
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, 0, limit)
	for _, task := range allTasks {
		if containsFold(task.Title, pattern) || containsFold(task.Comment, pattern) {
			tasks = append(tasks, task)
			if len(tasks) >= limit {
				break
			}
		}
	}
	return tasks, nil
}

func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}

func queryTasks(query string, args ...any) ([]*Task, error) {
	rows, err := DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func scanTask(rows *sql.Rows) (*Task, error) {
	var id int64
	var task Task
	if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return nil, err
	}
	task.ID = strconv.FormatInt(id, 10)
	return &task, nil
}
