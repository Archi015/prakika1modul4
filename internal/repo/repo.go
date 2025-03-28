package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"simple-service/internal/config"
)

// Слой репозитория, здесь должны быть все методы, связанные с базой данных

// SQL-запрос на вставку задачи (sql-запросы в константе)
const (
	checkUserExists = `"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"`
	createUser      = `"INSERT INTO users (username, password) VALUES ($1, $2)"`
)

type userRepository struct {
	db *sql.DB
}

// Userepository - интерфейс с методом создания задачи
type UserRepository interface {
	CheckUserExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, user config.User) (string, error)
}

type repository struct {
	pool *pgxpool.Pool
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewRepository(ctx context.Context, cfg config.PostgreSQL) (*repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)
	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CheckUserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := checkUserExists
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user config.User) (string, error) {
	var userID string
	query := createUser
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Password)
	return userID, err
}
