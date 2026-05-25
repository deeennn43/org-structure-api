package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danil/org-structure-api/internal/config"
	"github.com/danil/org-structure-api/internal/database"
	"github.com/danil/org-structure-api/internal/handler"
	"github.com/danil/org-structure-api/internal/migrate"
	gormrepo "github.com/danil/org-structure-api/internal/repository/gorm"
	"github.com/danil/org-structure-api/internal/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("database", "err", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("sql db", "err", err)
		os.Exit(1)
	}
	if err := waitForDB(sqlDB); err != nil {
		slog.Error("wait db", "err", err)
		os.Exit(1)
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}
	if err := migrate.Up(sqlDB, migrationsDir); err != nil {
		slog.Error("migrate", "err", err)
		os.Exit(1)
	}

	deptRepo := gormrepo.NewDepartmentRepo(db)
	empRepo := gormrepo.NewEmployeeRepo(db)
	deptSvc := service.NewDepartmentService(deptRepo, empRepo)
	empSvc := service.NewEmployeeService(deptRepo, empRepo)
	deptHandler := handler.NewDepartmentHandler(deptSvc, empSvc)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler.NewRouter(deptHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		slog.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	slog.Info("server stopped")
}

func waitForDB(db *sql.DB) error {
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return db.Ping()
}
