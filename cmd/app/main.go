package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

	handlers "auth_service/api/http"
	"auth_service/config"
	_ "auth_service/docs"
	"auth_service/repository/postgres"
	"auth_service/usecases/services"
)

const addr = ":8000"

// @title Auth Service
// @version 1.0
// @description Test task for medods.

// @host localhost:8000
// @BasePath /
func main() {
	flags := config.ParseFlags()
	var cfg config.HTTPConfig
	config.MustLoad(flags.ConfigPath, &cfg)

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.PG.Host, cfg.PG.Port, cfg.PG.User, cfg.PG.Password, cfg.PG.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	emailService := services.NewEmailService()

	authRepo := postgres.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, emailService, cfg.Tokens.AccessExp, cfg.Tokens.RefreshExp)
	authHandler := handlers.NewAuthHandler(authService)

	r := chi.NewRouter()

	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Post("/generate_tokens", authHandler.GenerateTokens)
	r.Post("/refresh_tokens", authHandler.RefreshTokens)

	httpServer := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	httpServer.ListenAndServe()
}
