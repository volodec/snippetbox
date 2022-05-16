package mysql

import (
	"database/sql"
	"errors"
	"github.com/volodec/snippetbox/pkg/models"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `INSERT INTO snippets (title, content, created_at, expires_at)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, resultErr := m.DB.Exec(query, title, content, expires)
	if resultErr != nil {
		return 0, resultErr
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `SELECT id, title, content, created_at, expires_at FROM snippets 
	WHERE id = ? AND expires_at > UTC_TIMESTAMP()`

	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(query, id)

	// Инициализируем указатель на новую структуру Snippet.
	s := &models.Snippet{}

	scanErr := row.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)

	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, scanErr
		}
	}

	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
