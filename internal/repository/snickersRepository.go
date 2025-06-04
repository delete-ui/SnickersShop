package repository

import (
	"SnickersShopPet1.0/internal/models"
	"database/sql"
	"go.uber.org/zap"
)

type DBExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

type SnickersRepository struct {
	db   DBExecutor
	Zlog *zap.Logger
}

func NewSnickersRepository(db *sql.DB, Zlog *zap.Logger) *SnickersRepository {
	return &SnickersRepository{db: db, Zlog: Zlog}
}

func (r *SnickersRepository) Add(m *models.SnickerInput) (*models.Snickers, error) {
	const op = "internal.repository.SnickersRepository.Add"

	query := `INSERT INTO products (title, description, cost)
	VALUES ($1, $2,$3)
	RETURNING id, title, description, cost;`

	row := r.db.QueryRow(query, m.Title, m.Description, m.Cost)

	var snickers models.Snickers

	if err := row.Scan(&snickers.ID, &snickers.Title, &snickers.Description, &snickers.Cost); err != nil {
		r.Zlog.Error(op, zap.Error(err))
		return nil, err
	}

	return &snickers, nil

}

func (r *SnickersRepository) GetAll() (*[]models.Snickers, error) {
	const op = "internal.repository.SnickersRepository.GetAll"

	query := `SELECT * FROM products;`

	rows, err := r.db.Query(query)
	if err != nil {
		r.Zlog.Error(op, zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var snickersSlice []models.Snickers

	for rows.Next() {
		var snickers models.Snickers
		if err := rows.Scan(&snickers.ID, &snickers.Title, &snickers.Description, &snickers.Cost); err != nil {
			r.Zlog.Error(op, zap.Error(err))
			return nil, err
		}
		snickersSlice = append(snickersSlice, snickers)
	}

	if err := rows.Err(); err != nil {
		r.Zlog.Error(op, zap.Error(err))
		return nil, err
	}

	return &snickersSlice, nil

}
