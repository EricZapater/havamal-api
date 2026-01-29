package server

import (
	"havamal-api/config"
	"havamal-api/internal/auth"
	"havamal-api/internal/categories"
	"havamal-api/internal/images"
	"havamal-api/internal/navigation"
	"havamal-api/internal/posts"
	"havamal-api/internal/versions"

	"havamal-api/internal/users"

	"database/sql"
	"havamal-api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router *gin.Engine
	config config.Config
	db *sql.DB
}

func NewServer(config config.Config, db *sql.DB) *Server {
	return &Server{
		router: gin.New(),
		config: config,
		db:     db,
	}
}

func(s *Server)Setup()error{
	s.router.Use(middleware.SetupCORS())
	s.router.Use(middleware.ObservabilityMiddleware())
	
	authMiddleware, err := middleware.SetupJWT(s.config)
	if err != nil{
		return err
	}
	

	//Repositories
	userRepo := users.NewRepository(s.db)
	postRepo := posts.NewRepository(s.db)
	categoryRepo := categories.NewRepository(s.db)
	versionRepo := versions.NewRepository(s.db)
	navigationRepo := navigation.NewRepository(s.db)


	//Services
	userService := users.NewService(userRepo)
	authService := auth.NewAuthService(userService, authMiddleware)
	postService := posts.NewService(postRepo, userService)
	categoryService := categories.NewService(categoryRepo)
	versionService := versions.NewService(versionRepo)
	navigationService := navigation.NewService(navigationRepo)

	//Handlers
	userHandler := users.NewHandler(userService)
	authHandler := auth.NewAuthHandler(authService, authMiddleware)
	postHandler := posts.NewHandler(postService)
	categoryHandler := categories.NewHandler(categoryService)
	versionHandler := versions.NewHandler(versionService)
	navigationHandler := navigation.NewHandler(navigationService)
	imageHandler := images.NewHandler()

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// Prometheus metrics endpoint
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Serve static files from images directory
	s.router.Static("/images", "../images")

	// Public routes
	public := s.router.Group("/auth")
	auth.RegisterRoutes(public, authHandler, authMiddleware)	

	//Blog routes
	blog := s.router.Group("/blog")
	posts.RegisterPublicRoutes(blog, &postHandler)	
	categories.RegisterPublicRoutes(blog, &categoryHandler)
	versions.RegisterPublicRoutes(blog, &versionHandler)	
	navigation.RegisterPublicRoutes(blog, &navigationHandler)
	

	//protected routes
	protected := s.router.Group("/api")
	protected.Use(authMiddleware.MiddlewareFunc())
	protected.Use(middleware.ContextMiddleware()) // Inject context values	
	users.RegisterRoutes(protected, &userHandler) // internally has /users prefix		
	posts.RegisterRoutes(protected, &postHandler)		
	categories.RegisterRoutes(protected, &categoryHandler)		
	versions.RegisterRoutes(protected, &versionHandler)		
	navigation.RegisterRoutes(protected, &navigationHandler)
	images.RegisterRoutes(protected, &imageHandler)

	return nil
	
}

func(s *Server)Run()error{
	return s.router.Run(":" + s.config.App.Port)
}