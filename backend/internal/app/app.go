package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"pigeon/internal/config"
	httptransport "pigeon/internal/transport/http"
)

type App struct {
	Config     *config.Config
	DB         *pgxpool.Pool
	Redis      *redis.Client
	Router     http.Handler
}

func New(cfg *config.Config) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	rdb, err := initRedis(ctx, cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	application := &App{
		Config: cfg,
		DB:     dbpool,
		Redis:  rdb,
	}

	application.Router = httptransport.NewRouter()

	return application, nil
}

func initRedis(ctx context.Context, redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}

func (a *App) Close() {
	a.DB.Close()
	a.Redis.Close()
}