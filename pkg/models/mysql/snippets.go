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
	// открытие транзакции
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO snippets (title, content, created_at, expires_at)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, resultErr := tx.Exec(query, title, content, expires)
	if resultErr != nil {
		tx.Rollback()
		return 0, resultErr
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return int(id), nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `SELECT id, title, content, created_at, expires_at FROM snippets 
	WHERE id = ? AND expires_at > UTC_TIMESTAMP()`

	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(query, id)

	// Инициализируем указатель на новую структуру OneEntry.
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
	query := `SELECT id, title, content, created_at, expires_at FROM snippets
	WHERE expires_at > UTC_TIMESTAMP() ORDER BY created_at LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится ПОСЛЕ проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Snippets.
	var snippets []*models.Snippet

	// Используем rows.Next() для перебора результата. Этот метод предоставляет
	// первую, а затем каждую следующую запись из базы данных, для обработки
	// методом rows.Scan().
	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// не возникла ли какая либо ошибка в ходе работы.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
