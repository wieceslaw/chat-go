package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wieceslaw/chat-go/config"
	"github.com/wieceslaw/chat-go/internal/environment"
	"github.com/wieceslaw/chat-go/internal/server/auth"
	"github.com/wieceslaw/chat-go/internal/server/hello"
)

type Server struct {
	Config     *config.Config
	HttpServer *http.Server
	Db         *sql.DB
}

func New(ctx context.Context, configFile string) (*Server, error) {
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	// Setup router
	gin.SetMode(cfg.Server.Mode)

	// Connect to db
	db, err := getDb(ctx, &cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to db: %v", err)
	}

	// Setup router
	router := setupRouter(ctx, db, cfg)

	// Create http server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &Server{
		Config:     cfg,
		HttpServer: srv,
		Db:         db,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Server failed: %v", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	// TODO: shutdown multiple things with async and context
	if err := s.Db.Close(); err != nil {
		log.Println(fmt.Errorf("Failed to stop db connection %v", err).Error())
	}
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) String() string {
	return fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port)
}

func setupRouter(ctx context.Context, db *sql.DB, cfg *config.Config) *gin.Engine {
	router := gin.New()

	if cfg.Server.EnableLogger {
		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output:    os.Stdout,
			SkipPaths: []string{"/health"},
		}))
	}

	if cfg.Server.EnableRecovery {
		router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
			log.Printf("Panic recovered: %v", recovered)
			c.AbortWithStatus(http.StatusInternalServerError)
		}))
	}

	router.MaxMultipartMemory = cfg.Server.MaxMultipartMemory

	if len(cfg.Server.TrustedProxies) > 0 {
		router.SetTrustedProxies(cfg.Server.TrustedProxies)
	}

	// TODO: Add CORS middleware

	// TODO: Add rate limiting

	setupRoutes(ctx, db, router)

	return router
}

func setupRoutes(ctx context.Context, db *sql.DB, r *gin.Engine) {
	repository := auth.NewUserRepository(db)
	service, _ := auth.NewUserService(ctx, repository, auth.MockJwtProvider())

	authHandler := auth.NewAuthHanlder(service)
	authHandler.RegisterRoutes(r.Group(""))

	authMiddleware := auth.NewAuthMiddleware(service)

	api := r.Group("/api/v1")
	api.Use(authMiddleware.AuthRequired())
	{
		if environment.Get() == environment.DevelopmentEnv {
			helloHandler := hello.NewHelloHandler()
			helloHandler.RegisterRoutes(api.Group("/hello"))
		}
	}
}

func getDb(ctx context.Context, cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("Failed to open db driver %v", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping db %v", err)
	}

	log.Println("Database connected successfully")

	return db, nil
}
