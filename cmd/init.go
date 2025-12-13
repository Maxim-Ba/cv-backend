package main

import (
	"context"

	"github.com/Maxim-Ba/cv-backend/config"
	"github.com/Maxim-Ba/cv-backend/internal/dbconn"
	"github.com/Maxim-Ba/cv-backend/internal/router"
)

func initApplication(ctx context.Context,db *dbconn.DB, cfg *config.Config) (*router.Router, error) {
	//TODO init other services
r:= router.New() 
return  r, nil
}
