package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/admin"
	"github.com/yourusername/online-library/internal/auth"
	"github.com/yourusername/online-library/internal/book"
	"github.com/yourusername/online-library/internal/bookmark"
	"github.com/yourusername/online-library/internal/config"
	"github.com/yourusername/online-library/internal/donation"
	"github.com/yourusername/online-library/internal/handover"
	"github.com/yourusername/online-library/internal/idea"
	"github.com/yourusername/online-library/internal/infrastructure/db/postgres"
	"github.com/yourusername/online-library/internal/notification"
	"github.com/yourusername/online-library/internal/repository"
	"github.com/yourusername/online-library/internal/review"
	"github.com/yourusername/online-library/internal/successscore"
	"github.com/yourusername/online-library/internal/user"

	adminhandler "github.com/yourusername/online-library/internal/rest/handler/admin"
	authhandler "github.com/yourusername/online-library/internal/rest/handler/auth"
	bookhandler "github.com/yourusername/online-library/internal/rest/handler/book"
	bookmarkhandler "github.com/yourusername/online-library/internal/rest/handler/bookmark"
	donationhandler "github.com/yourusername/online-library/internal/rest/handler/donation"
	handoverhandler "github.com/yourusername/online-library/internal/rest/handler/handover"
	ideahandler "github.com/yourusername/online-library/internal/rest/handler/idea"
	notificationhandler "github.com/yourusername/online-library/internal/rest/handler/notification"
	reviewhandler "github.com/yourusername/online-library/internal/rest/handler/review"
	swaggerhandler "github.com/yourusername/online-library/internal/rest/handler/swagger"
	userhandler "github.com/yourusername/online-library/internal/rest/handler/user"
	"github.com/yourusername/online-library/internal/rest/middleware"

	"go.uber.org/zap"
)

func run(ctx context.Context, cfg *config.Config, log *zap.Logger) error {
	// Connect to database
	conn, err := postgres.NewConnection(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()
	log.Info("database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(conn.DB, log)
	bookRepo := repository.NewBookRepository(conn.DB, log)
	ideaRepo := repository.NewIdeaRepository(conn.DB, log)
	reviewRepo := repository.NewReviewRepository(conn.DB, log)
	donationRepo := repository.NewDonationRepository(conn.DB, log)
	bookmarkRepo := repository.NewBookmarkRepository(conn.DB, log)
	scoreRepo := repository.NewSuccessScoreRepository(conn.DB, log)
	notificationRepo := repository.NewNotificationRepository(conn.DB, log)
	adminRepo := repository.NewAdminRepository(conn.DB, log)
	handoverRepo := repository.NewHandoverRepository(conn.DB, log)

	// Initialize services
	successScoreSvc := successscore.NewService(scoreRepo, log)
	notificationSvc := notification.NewService(notificationRepo, log)
	authSvc := auth.NewService(userRepo, cfg.JWT.Secret, log)
	userSvc := user.NewService(userRepo, log)
	bookSvc := book.NewService(bookRepo, log)
	ideaSvc := idea.NewService(ideaRepo, successScoreSvc, notificationSvc, log)
	reviewSvc := review.NewService(reviewRepo, successScoreSvc, notificationSvc, log)
	donationSvc := donation.NewService(donationRepo, successScoreSvc, log)
	bookmarkSvc := bookmark.NewService(bookmarkRepo, log)
	handoverSvc := handover.NewService(handoverRepo, notificationSvc, log)
	adminSvc := admin.NewService(adminRepo, successScoreSvc, notificationSvc, handoverRepo, log)

	// Initialize handlers
	authHandler := authhandler.NewHandler(authSvc, log)
	userHandler := userhandler.NewHandler(userSvc, log)
	bookHandler := bookhandler.NewHandler(bookSvc, log)
	ideaHandler := ideahandler.NewHandler(ideaSvc, log)
	reviewHandler := reviewhandler.NewHandler(reviewSvc, log)
	donationHandler := donationhandler.NewHandler(donationSvc, log)
	bookmarkHandler := bookmarkhandler.NewHandler(bookmarkSvc, log)
	adminHandler := adminhandler.NewHandler(adminSvc, log)
	notificationHandler := notificationhandler.NewHandler(notificationSvc, log)
	handoverHandler := handoverhandler.NewHandler(handoverSvc, log)

	// Setup router
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/docs", swaggerhandler.ServeSwaggerUI)
	router.GET("/docs/swagger.yaml", swaggerhandler.ServeSwaggerYAML)
	log.Info("📚 API documentation available at http://localhost:" + cfg.Server.Port + "/docs")

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		authhandler.RegisterPublicRoutes(api, authHandler)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authSvc, log))
		{
			authhandler.RegisterProtectedRoutes(protected, authHandler)
			userhandler.RegisterRoutes(protected, userHandler)
			bookhandler.RegisterRoutes(protected, bookHandler)
			ideahandler.RegisterRoutes(protected, ideaHandler)
			reviewhandler.RegisterRoutes(protected, reviewHandler)
			donationhandler.RegisterRoutes(protected, donationHandler)
			bookmarkhandler.RegisterRoutes(protected, bookmarkHandler)
			notificationhandler.RegisterRoutes(protected, notificationHandler)
			handoverhandler.RegisterRoutes(protected, handoverHandler)
		}

		// Admin routes (requires admin role)
		adminRoutes := api.Group("")
		adminRoutes.Use(middleware.AuthMiddleware(authSvc, log))
		adminRoutes.Use(middleware.AdminMiddleware())
		{
			adminhandler.RegisterRoutes(adminRoutes, adminHandler)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.Info("starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info("shutting down server")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
		return err
	}

	return nil
}
