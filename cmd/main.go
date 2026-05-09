// @title           CV Backend API
// @version         1.0
// @description     REST API для CV-приложения: теги, технологии, история работы, образование.
// @host            localhost:3333
// @BasePath        /api

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Maxim-Ba/cv-backend/config"
	_ "github.com/Maxim-Ba/cv-backend/docs"
	"github.com/Maxim-Ba/cv-backend/internal/dbconn"
	"github.com/Maxim-Ba/cv-backend/internal/repository"
	"github.com/Maxim-Ba/cv-backend/internal/router"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/pkg/logger"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cfg := config.GetConfig()
	fmt.Printf("Config: %+v\n", cfg)
	logger.InitLogger(cfg)
	db, err := dbconn.New(*cfg)

	if err != nil {
		log.Panicf("%v", err)
	}
	router, err := initApplication(ctx, db, cfg)
	if err != nil {
		log.Panicf("%v", err)
	}
	var wg sync.WaitGroup
	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router.R,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			close(exit)
		}
	}()

	select {
	case <-exit:
	case <-ctx.Done():
	}
	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		if err := server.Close(); err != nil {
			slog.Error("forced server close failed", "error", err)
		}
	}
	wg.Wait()
	slog.Info("server stopped")
}

func initApplication(ctx context.Context, db *dbconn.DB, cfg *config.Config) (*router.Router, error) {
	// Инициализация репозиториев
	repos := defineRepositories(db)
	
	// Инициализация сервисов с использованием репозиториев
	tagSvc := services.NewTagService(repos.TagRepository)
	techSvc := services.NewTechService(repos.TechRepository)
	eduSvc := services.NewEducationService(repos.EducationRepository)
	whSvc := services.NewWorkHistoryService(repos.WorkHistoryRepository)
	pdfSvc := services.NewPDFService(repos.ProfileRepository, whSvc, techSvc, eduSvc)

	deps := &router.Dependencies{
		TagService:         tagSvc,
		TechService:        techSvc,
		EducationService:   eduSvc,
		WorkHistoryService: whSvc,
		PDFService:         pdfSvc,
	}
	
	// Инициализация роутера с зависимостями
	r := router.New(deps, db.GetConnection(), cfg.AllowedOrigin, cfg.AdminUser, cfg.AdminPassword, cfg.Secret)
	return r, nil
}

// Repositories структура для хранения всех репозиториев приложения
type Repositories struct {
	TagRepository         *repository.TagRepo
	TechRepository        *repository.TechnologyRepo
	EducationRepository   *repository.EducationRepo
	WorkHistoryRepository *repository.WorkHistoryRepo
	ProfileRepository     *repository.ProfileRepo
}

// defineRepositories создает экземпляры всех репозиториев
func defineRepositories(db *dbconn.DB) *Repositories {
	return &Repositories{
		TagRepository:         repository.NewTagRepo(db.GetConnection()),
		TechRepository:        repository.NewTechnologyRepo(db.GetConnection()),
		EducationRepository:   repository.NewEducationRepo(db.GetConnection()),
		WorkHistoryRepository: repository.NewWorkHistoryRepo(db.GetConnection()),
		ProfileRepository:     repository.NewProfileRepo(db.GetConnection()),
	}
}
