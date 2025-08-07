package main

import (
	"context"
	"github.com/NeilPanic/SkillsRock_test/internal/handler"
	"github.com/NeilPanic/SkillsRock_test/internal/repo"
	"github.com/NeilPanic/SkillsRock_test/internal/service"
	"log/slog"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		slog.Error("DATABASE_DSN env is empty")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("pgxpool connect", "err", err)
		return
	}
	defer pool.Close()

	taskRepo := repo.NewTaskRepo(pool)
	taskSvc := service.New(taskRepo)

	app := fiber.New(fiber.Config{AppName: "todo-app"})
	handler.Register(app.Group("/tasks"), taskSvc)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			slog.Error("fiber stopped", "err", err)
			stop() // останавливаем контекст, чтобы перейти к Shutdown
		}
	}()

	<-ctx.Done() // ждём Ctrl-C

	if err := app.Shutdown(); err != nil {
		slog.Error("shutdown", "err", err)
	}

}
