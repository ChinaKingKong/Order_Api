package app

import (
	"fmt"
	"log"
	"order_api/auth"
	"order_api/cache"
	"order_api/config"
	"order_api/database"
	"order_api/handler"
	"order_api/repository"
	"order_api/router"
	"order_api/service"

	"github.com/gin-gonic/gin"
)

type App struct {
	config      *config.Config
	db          *database.Database
	cache       *cache.Cache
	router      *gin.Engine
	authService *auth.AuthService
}

func NewApp() *App {
	return &App{
		config: config.NewConfig(),
	}
}

func (a *App) Initialize() error {
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := a.initCache(); err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	a.initAuth()

	if err := a.initRouter(); err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	return nil
}

func (a *App) initDatabase() error {
	db, err := database.NewDatabase(&a.config.Database)
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *App) initCache() error {
	cache, err := cache.NewCache(&a.config.Redis)
	if err != nil {
		return err
	}
	a.cache = cache
	return nil
}

func (a *App) initAuth() {
	a.authService = auth.NewAuthService(a.config)
}

func (a *App) initRouter() error {
	orderRepo := repository.NewOrderRepository(a.db.DB, a.cache)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	authHandler := handler.NewAuthHandler(a.authService)

	a.router = router.SetupRouter(orderHandler, authHandler, a.authService)
	return nil
}

func (a *App) Run() error {
	log.Printf("Server starting on port %s", a.config.Server.Port)
	return a.router.Run(":" + a.config.Server.Port)
}

func (a *App) Shutdown() error {
	if err := a.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	if err := a.cache.Close(); err != nil {
		log.Printf("Error closing cache connection: %v", err)
	}

	return nil
}
