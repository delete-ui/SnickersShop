package repository

import (
	"SnickersShopPet1.0/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type DBExecutor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

type SnickersRepository struct {
	db        DBExecutor
	Zlog      *zap.Logger
	redisRepo *redis.Client
}

func NewSnickersRepository(db *sql.DB, Zlog *zap.Logger, redisRepo *redis.Client) *SnickersRepository {
	return &SnickersRepository{db: db, Zlog: Zlog, redisRepo: redisRepo}
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

func (r *SnickersRepository) GetAll(pagination *models.Pagination) (*[]models.Snickers, error) {
	const op = "internal.repository.SnickersRepository.GetAll"

	offset := (pagination.Offset - 1) * pagination.Limit
	cacheKey := fmt.Sprintf("snickers:page:%d:size:%d", pagination.Offset, pagination.Limit)
	ctx := context.Background()

	cachedData, err := r.redisRepo.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var snickers []models.Snickers
		if err := json.Unmarshal(cachedData, &snickers); err == nil {
			r.Zlog.Debug("Served by redis")
			return &snickers, nil
		}
		r.Zlog.Warn("Failed to unmarshal cached data", zap.String("location", op), zap.Error(err))
	} else if err != redis.Nil {
		r.Zlog.Warn("Redis error", zap.String("location", op), zap.Error(err))
	}

	query := `SELECT * FROM products OFFSET $1 LIMIT $2;`
	rows, err := r.db.Query(query, offset, pagination.Limit)
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

	if len(snickersSlice) > 0 {
		serialized, err := json.Marshal(snickersSlice)
		if err != nil {
			r.Zlog.Warn("Failed to marshal data for cache", zap.String("location", op), zap.Error(err))
		} else {
			if err := r.redisRepo.Set(ctx, cacheKey, serialized, 5*time.Minute).Err(); err != nil {
				r.Zlog.Warn("Failed to set cache", zap.String("location", op), zap.Error(err))
			}
		}
	}

	return &snickersSlice, nil
}

func (r *SnickersRepository) GetByID(idUUID string) (*models.Snickers, error) {
	const op = "internal.repository.SnickersRepository.GetByID"

	query := `SELECT * FROM products WHERE id = $1;`

	row := r.db.QueryRow(query, idUUID)

	var snickers models.Snickers

	if err := row.Scan(&snickers.ID, &snickers.Title, &snickers.Description, &snickers.Cost); err != nil {
		r.Zlog.Error(op, zap.Error(err))
		return nil, err
	}

	return &snickers, nil

}

func (r *SnickersRepository) GetByCost(cost *models.CostRange) (*[]models.Snickers, error) {
	const op = "internal.repository.SnickersRepository.GetByCost"

	query := `SELECT * FROM products WHERE cost > $1 AND cost < $2;`

	rows, err := r.db.Query(query, cost.Min, cost.Max)
	if err != nil {
		r.Zlog.Error("Error while executing query", zap.String(" location:", op), zap.Error(err))
		return nil, err
	}

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
