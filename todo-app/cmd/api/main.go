package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-app/internal/handler"
	"todo-app/internal/repo"
	"todo-app/internal/service"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		slog.Error("DATABASE_DSN env is empty")
		return
	}

	// ── graceful shutdown ────────────────────────────────────────────────
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// ── Postgres ────────────────────────────────────────────────────────
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("pgxpool connect", "err", err)
		return
	}
	defer pool.Close()

	// ── DI-цепочка: repo → service → handler ────────────────────────────
	taskRepo := repo.NewTaskRepo(pool)
	taskSvc := service.New(taskRepo)

	app := fiber.New(fiber.Config{AppName: "todo-app"})
	handler.Register(app.Group("/tasks"), taskSvc)

	// ── run server ──────────────────────────────────────────────────────
	go func() {
		if err := app.Listen(":8080"); err != nil {
			slog.Error("fiber stopped", "err", err)
			stop() // останавливаем контекст, чтобы перейти к Shutdown
		}
	}()

	<-ctx.Done() // ждём Ctrl-C

	// ── graceful stop (без контекста в v2<2.54) ─────────────────────────
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fiber ≤ v2.53:
	if err := app.Shutdown(); err != nil {
		slog.Error("shutdown", "err", err)
	}

	// если обновишься до v2.54+:
	// if err := app.ShutdownWithContext(shutdownCtx); err != nil { … }
}
