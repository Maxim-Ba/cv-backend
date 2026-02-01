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
		//TODO log
		if err := server.Close(); err != nil {
			slog.Error(err.Error())

		}
	}
	//TODO shutdown actions
	wg.Wait()
}

func initApplication(ctx context.Context, db *dbconn.DB, cfg *config.Config) (*router.Router, error) {
	// Инициализация репозиториев
	repos := defineRepositories(db)
	
	// Инициализация сервисов с использованием репозиториев
	deps := &router.Dependencies{
		TagService:         services.NewTagServise(repos.TagRepository),
		TechService:        services.NewTechService(repos.TechRepository),
		EducationService:   services.NewEducationService(repos.EducationRepository),
		WorkHistoryService: services.NewWorkHistoryService(repos.WorkHistoryRepository),
	}
	
	// Инициализация роутера с зависимостями
	r := router.New(deps)
	return r, nil
}

// Repositories структура для хранения всех репозиториев приложения
type Repositories struct {
	TagRepository         *repository.TagRepo
	TechRepository        *repository.TechnologyRepo
	EducationRepository   *repository.EducationRepo
	WorkHistoryRepository *repository.WorkHistoryRepo
}

// defineRepositories создает экземпляры всех репозиториев
func defineRepositories(db *dbconn.DB) *Repositories {
	return &Repositories{
		TagRepository:         repository.NewTagRepo(db.GetConnection()),
		TechRepository:        repository.NewTechnologyRepo(db.GetConnection()),
		EducationRepository:   repository.NewEducationRepo(db.GetConnection()),
		WorkHistoryRepository: repository.NewWorkHistoryRepo(db.GetConnection()),
	}
}
