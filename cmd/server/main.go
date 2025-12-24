package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portal-data-backend/infrastructure/config"
	"portal-data-backend/infrastructure/db"
	"portal-data-backend/infrastructure/http/middleware"
	"portal-data-backend/infrastructure/http/response"
	"portal-data-backend/infrastructure/logger"
	"portal-data-backend/infrastructure/security"
	"portal-data-backend/infrastructure/storage"

	// Auth module
	authDelivery "portal-data-backend/internal/auth/delivery/http"
	authRepo "portal-data-backend/internal/auth/repository"
	authUsecase "portal-data-backend/internal/auth/usecase"

	// User module
	userDelivery "portal-data-backend/internal/user/delivery/http"
	userRepo "portal-data-backend/internal/user/repository"
	userUsecase "portal-data-backend/internal/user/usecase"

	// Organization module
	orgDelivery "portal-data-backend/internal/organization/delivery/http"
	orgRepo "portal-data-backend/internal/organization/repository"
	orgUsecase "portal-data-backend/internal/organization/usecase"

	// Dataset module
	datasetDelivery "portal-data-backend/internal/dataset/delivery/http"
	datasetRepo "portal-data-backend/internal/dataset/repository"
	datasetUsecase "portal-data-backend/internal/dataset/usecase"

	// Tag module
	tagDelivery "portal-data-backend/internal/tag/delivery/http"
	tagRepo "portal-data-backend/internal/tag/repository"
	tagUsecase "portal-data-backend/internal/tag/usecase"

	// BusinessField module
	bfDelivery "portal-data-backend/internal/business_field/delivery/http"
	bfRepo "portal-data-backend/internal/business_field/repository"
	bfUsecase "portal-data-backend/internal/business_field/usecase"

	// Topic module
	topicDelivery "portal-data-backend/internal/topic/delivery/http"
	topicRepo "portal-data-backend/internal/topic/repository"
	topicUsecase "portal-data-backend/internal/topic/usecase"

	// Unit module
	unitDelivery "portal-data-backend/internal/unit/delivery/http"
	unitRepo "portal-data-backend/internal/unit/repository"
	unitUsecase "portal-data-backend/internal/unit/usecase"

	// Feedback module
	fbDelivery "portal-data-backend/internal/feedback/delivery/http"
	fbRepo "portal-data-backend/internal/feedback/repository"
	fbUsecase "portal-data-backend/internal/feedback/usecase"

	// File module
	fileDelivery "portal-data-backend/internal/file/delivery/http"
	fileRepo "portal-data-backend/internal/file/repository"
	fileUsecase "portal-data-backend/internal/file/usecase"

	// Analytics module
	analyticsDelivery "portal-data-backend/internal/analytics/delivery/http"
	analyticsRepo "portal-data-backend/internal/analytics/repository"
	analyticsUsecase "portal-data-backend/internal/analytics/usecase"

	// Visualization module
	vizDelivery "portal-data-backend/internal/visualization/delivery/http"
	vizRepo "portal-data-backend/internal/visualization/repository"
	vizUsecase "portal-data-backend/internal/visualization/usecase"

	// Publication module
	pubDelivery "portal-data-backend/internal/publication/delivery/http"
	pubRepo "portal-data-backend/internal/publication/repository"
	pubUsecase "portal-data-backend/internal/publication/usecase"

	// Settings module
	settingsDelivery "portal-data-backend/internal/settings/delivery/http"
	settingsRepo "portal-data-backend/internal/settings/repository"
	settingsUsecase "portal-data-backend/internal/settings/usecase"

	// Notification module
	notifDelivery "portal-data-backend/internal/notification/delivery/http"
	notifRepo "portal-data-backend/internal/notification/repository"
	notifUsecase "portal-data-backend/internal/notification/usecase"

	// DataRow module
	dataRowDelivery "portal-data-backend/internal/data_row/delivery/http"
	dataRowRepo "portal-data-backend/internal/data_row/repository"
	dataRowUsecase "portal-data-backend/internal/data_row/usecase"

	// Desk module
	deskDelivery "portal-data-backend/internal/desk/delivery/http"
	deskRepo "portal-data-backend/internal/desk/repository"
	deskUsecase "portal-data-backend/internal/desk/usecase"

	// Integration module
	integrationDelivery "portal-data-backend/internal/integration/delivery/http"
	integrationRepo "portal-data-backend/internal/integration/repository"
	integrationUsecase "portal-data-backend/internal/integration/usecase"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.App.Debug, cfg.App.Environment)

	logger.Info("Starting %s v%s", cfg.App.Name, cfg.App.Version)
	logger.Info("Environment: %s", cfg.App.Environment)

	// Initialize database
	postgres, err := db.NewPostgres(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer postgres.Close()

	logger.Info("Database connected successfully")

	// Initialize infrastructure components
	jwtManager := security.NewJWTManager(&cfg.JWT)
	passwordHasher := security.NewPasswordHandler()

	// Initialize Auth module
	userRepository := authRepo.NewUserPostgresRepository(postgres.DB)
	tokenRepository := authRepo.NewTokenPostgresRepository(postgres.DB)

	authUsecaseInstance := authUsecase.NewAuthUsecase(
		userRepository,
		tokenRepository,
		jwtManager,
		passwordHasher,
	)

	authHandler := authDelivery.NewHandler(authUsecaseInstance)

	// Initialize User module
	userRepositoryInstance := userRepo.NewUserPostgresRepository(postgres.DB)
	userUsecaseInstance := userUsecase.NewUserUsecase(userRepositoryInstance)
	userHandler := userDelivery.NewHandler(userUsecaseInstance)

	// Initialize Organization module
	orgRepository := orgRepo.NewOrgPostgresRepository(postgres.DB)
	orgUsecaseInstance := orgUsecase.NewOrgUsecase(orgRepository)
	orgHandler := orgDelivery.NewHandler(orgUsecaseInstance)

	// Initialize Dataset module
	datasetRepository := datasetRepo.NewDatasetPostgresRepository(postgres.DB)
	datasetUsecaseInstance := datasetUsecase.NewDatasetUsecase(datasetRepository)
	datasetHandler := datasetDelivery.NewHandler(datasetUsecaseInstance)

	// Initialize Tag module
	tagRepository := tagRepo.NewTagPostgresRepository(postgres.DB)
	tagUsecaseInstance := tagUsecase.NewTagUsecase(tagRepository)
	tagHandler := tagDelivery.NewHandler(tagUsecaseInstance)

	// Initialize BusinessField module
	bfRepository := bfRepo.NewBusinessFieldPostgresRepository(postgres.DB)
	bfUsecaseInstance := bfUsecase.NewBusinessFieldUsecase(bfRepository)
	bfHandler := bfDelivery.NewHandler(bfUsecaseInstance)

	// Initialize Topic module
	topicRepository := topicRepo.NewTopicPostgresRepository(postgres.DB)
	topicUsecaseInstance := topicUsecase.NewTopicUsecase(topicRepository)
	topicHandler := topicDelivery.NewHandler(topicUsecaseInstance)

	// Initialize Unit module
	unitRepository := unitRepo.NewUnitPostgresRepository(postgres.DB)
	unitUsecaseInstance := unitUsecase.NewUnitUsecase(unitRepository)
	unitHandler := unitDelivery.NewHandler(unitUsecaseInstance)

	// Initialize Feedback module
	fbRepository := fbRepo.NewFeedbackPostgresRepository(postgres.DB)
	fbUsecaseInstance := fbUsecase.NewFeedbackUsecase(fbRepository)
	fbHandler := fbDelivery.NewHandler(fbUsecaseInstance)

	// Initialize File module with MinIO storage
	minioStorage, err := storage.NewMinIOStorage(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.Bucket,
		cfg.MinIO.UseSSL,
	)
	if err != nil {
		logger.Fatal("Failed to connect to MinIO: %v", err)
	}
	logger.Info("MinIO connected successfully")

	fileRepository := fileRepo.NewFilePostgresRepository(postgres.DB)
	fileUsecaseInstance := fileUsecase.NewFileUsecase(fileRepository, minioStorage, "files")
	fileHandler := fileDelivery.NewHandler(fileUsecaseInstance)

	// Initialize Analytics module
	analyticsRepository := analyticsRepo.NewAnalyticsPostgresRepository(postgres.DB)
	analyticsUsecaseInstance := analyticsUsecase.NewAnalyticsUsecase(analyticsRepository)
	analyticsHandler := analyticsDelivery.NewHandler(analyticsUsecaseInstance)

	// Initialize Visualization module
	vizRepository := vizRepo.NewVisualizationPostgresRepository(postgres.DB)
	vizUsecaseInstance := vizUsecase.NewVisualizationUsecase(vizRepository)
	vizHandler := vizDelivery.NewHandler(vizUsecaseInstance)

	// Initialize Publication module
	pubRepository := pubRepo.NewPublicationPostgresRepository(postgres.DB)
	pubUsecaseInstance := pubUsecase.NewPublicationUsecase(pubRepository)
	pubHandler := pubDelivery.NewHandler(pubUsecaseInstance)

	// Initialize Settings module
	settingsRepository := settingsRepo.NewSettingsPostgresRepository(postgres.DB)
	settingsUsecaseInstance := settingsUsecase.NewSettingsUsecase(settingsRepository)
	settingsHandler := settingsDelivery.NewHandler(settingsUsecaseInstance)

	// Initialize Notification module
	notifRepository := notifRepo.NewNotificationPostgresRepository(postgres.DB)
	notifUsecaseInstance := notifUsecase.NewNotificationUsecase(notifRepository)
	notifHandler := notifDelivery.NewHandler(notifUsecaseInstance)

	// Initialize DataRow module
	dataRowRepository := dataRowRepo.NewDataRowPostgresRepository(postgres.DB)
	dataRowUsecaseInstance := dataRowUsecase.NewDataRowUsecase(dataRowRepository)
	dataRowHandler := dataRowDelivery.NewHandler(dataRowUsecaseInstance)

	// Initialize Desk module
	deskRepository := deskRepo.NewDeskPostgresRepository(postgres.DB)
	deskUsecaseInstance := deskUsecase.NewDeskUsecase(deskRepository)
	deskHandler := deskDelivery.NewHandler(deskUsecaseInstance)

	// Initialize Integration module
	integrationRepository := integrationRepo.NewIntegrationPostgresRepository(postgres.DB)
	integrationUsecaseInstance := integrationUsecase.NewIntegrationUsecase(integrationRepository)
	integrationHandler := integrationDelivery.NewHandler(integrationUsecaseInstance)

	// Setup HTTP router
	router := setupRouter(
		cfg,
		authHandler,
		userHandler,
		orgHandler,
		datasetHandler,
		tagHandler,
		bfHandler,
		topicHandler,
		unitHandler,
		fbHandler,
		fileHandler,
		analyticsHandler,
		vizHandler,
		pubHandler,
		settingsHandler,
		notifHandler,
		dataRowHandler,
		deskHandler,
		integrationHandler,
		jwtManager,
	)

	// Setup HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server listening on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited successfully")
}

// setupRouter configures and returns the HTTP router
func setupRouter(
	cfg *config.Config,
	authHandler *authDelivery.Handler,
	userHandler *userDelivery.Handler,
	orgHandler *orgDelivery.Handler,
	datasetHandler *datasetDelivery.Handler,
	tagHandler *tagDelivery.Handler,
	bfHandler *bfDelivery.Handler,
	topicHandler *topicDelivery.Handler,
	unitHandler *unitDelivery.Handler,
	fbHandler *fbDelivery.Handler,
	fileHandler *fileDelivery.Handler,
	analyticsHandler *analyticsDelivery.Handler,
	vizHandler *vizDelivery.Handler,
	pubHandler *pubDelivery.Handler,
	settingsHandler *settingsDelivery.Handler,
	notifHandler *notifDelivery.Handler,
	dataRowHandler *dataRowDelivery.Handler,
	deskHandler *deskDelivery.Handler,
	integrationHandler *integrationDelivery.Handler,
	jwtManager *security.JWTManager,
) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(middleware.Logger(cfg.App.Debug))
	r.Use(middleware.CORS())
	r.Use(middleware.ContentType)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, response.CodeSuccess, "Service is healthy", map[string]string{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// Public auth routes
	authDelivery.RegisterRoutes(r, authHandler)

	// Public routes (no authentication required)
	r.Group(func(r chi.Router) {
		// Organizations - public read access
		r.Route("/organizations", func(r chi.Router) {
			r.Get("/", orgHandler.List)
			r.Get("/code/{code}", orgHandler.GetByCode)
			r.Get("/{id}", orgHandler.GetByID)
		})

		// Datasets - public read access
		r.Route("/datasets", func(r chi.Router) {
			r.Get("/", datasetHandler.List)
			r.Get("/slug/{slug}", datasetHandler.GetBySlug)
			r.Get("/{id}", datasetHandler.GetByID)
		})

		// Tags - public read access
		r.Route("/tags", func(r chi.Router) {
			r.Get("/", tagHandler.List)
			r.Get("/{id}", tagHandler.GetByID)
		})

		// BusinessFields - public read access
		r.Route("/business-fields", func(r chi.Router) {
			r.Get("/", bfHandler.List)
			r.Get("/{id}", bfHandler.GetByID)
		})

		// Topics - public read access
		r.Route("/topics", func(r chi.Router) {
			r.Get("/", topicHandler.List)
			r.Get("/{id}", topicHandler.GetByID)
		})

		// Units - public read access
		r.Route("/units", func(r chi.Router) {
			r.Get("/", unitHandler.List)
			r.Get("/{id}", unitHandler.GetByID)
		})

		// Visualizations - public read access
		r.Route("/visualizations", func(r chi.Router) {
			r.Get("/", vizHandler.List)
			r.Get("/stats", vizHandler.GetStats)
			r.Get("/dataset/{datasetId}", vizHandler.GetByDatasetID)
			r.Get("/organization/{orgId}", vizHandler.GetByOrganizationID)
			r.Get("/{id}", vizHandler.GetByID)
		})

		// Publications - public read access
		r.Route("/publications", func(r chi.Router) {
			r.Get("/", pubHandler.List)
			r.Get("/dataset/{datasetId}", pubHandler.GetByDatasetID)
			r.Get("/organization/{orgId}", pubHandler.GetByOrganizationID)
			r.Get("/{id}", pubHandler.GetByID)
		})

		// Analytics - public read access
		r.Get("/analytics/dashboard", analyticsHandler.GetDashboard)
		r.Get("/analytics/stats/datasets", analyticsHandler.GetDatasetStats)
		r.Get("/analytics/stats/organizations", analyticsHandler.GetOrganizationStats)
		r.Get("/analytics/stats/users", analyticsHandler.GetUserStats)
		r.Get("/analytics/popular/datasets", analyticsHandler.GetPopularDatasets)
		r.Get("/analytics/popular/tags", analyticsHandler.GetPopularTags)
		r.Get("/analytics/trend/datasets", analyticsHandler.GetDatasetTrend)
	})

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))

		// Auth protected routes
		r.Post("/auth/revoke-all", authHandler.RevokeAllTokens)
		r.Get("/me", authHandler.GetCurrentUser)

		// User management
		userDelivery.RegisterRoutes(r, userHandler)

		// Organization management (write access)
		r.Route("/organizations", func(r chi.Router) {
			r.Post("/", orgHandler.Create)
			r.Put("/{id}", orgHandler.Update)
			r.Delete("/{id}", orgHandler.Delete)
			r.Patch("/{id}/status", orgHandler.UpdateStatus)
		})

		// Dataset management (write access)
		r.Route("/datasets", func(r chi.Router) {
			r.Post("/", datasetHandler.Create)
			r.Put("/{id}", datasetHandler.Update)
			r.Delete("/{id}", datasetHandler.Delete)
			r.Patch("/{id}/status", datasetHandler.UpdateStatus)
		})

		// Tag management (write access)
		r.Route("/tags", func(r chi.Router) {
			r.Post("/", tagHandler.Create)
			r.Put("/{id}", tagHandler.Update)
			r.Delete("/{id}", tagHandler.Delete)
		})

		// BusinessField management (write access)
		r.Route("/business-fields", func(r chi.Router) {
			r.Post("/", bfHandler.Create)
			r.Put("/{id}", bfHandler.Update)
			r.Delete("/{id}", bfHandler.Delete)
		})

		// Topic management (write access)
		r.Route("/topics", func(r chi.Router) {
			r.Post("/", topicHandler.Create)
			r.Put("/{id}", topicHandler.Update)
			r.Delete("/{id}", topicHandler.Delete)
		})

		// Unit management (write access)
		r.Route("/units", func(r chi.Router) {
			r.Post("/", unitHandler.Create)
			r.Put("/{id}", unitHandler.Update)
			r.Delete("/{id}", unitHandler.Delete)
		})

		// Feedback management
		fbDelivery.RegisterRoutes(r, fbHandler)

		// File management
		fileDelivery.RegisterRoutes(r, fileHandler)

		// Visualization management (write access)
		r.Route("/visualizations", func(r chi.Router) {
			r.Post("/", vizHandler.Create)
			r.Put("/{id}", vizHandler.Update)
			r.Delete("/{id}", vizHandler.Delete)
			r.Patch("/{id}/status", vizHandler.UpdateStatus)
		})

		// Publication management (write access)
		r.Route("/publications", func(r chi.Router) {
			r.Post("/", pubHandler.Create)
			r.Put("/{id}", pubHandler.Update)
			r.Delete("/{id}", pubHandler.Delete)
			r.Patch("/{id}/status", pubHandler.UpdateStatus)
			r.Post("/{id}/download", pubHandler.IncrementDownloadCount)
		})

		// Settings management
		settingsDelivery.RegisterRoutes(r, settingsHandler)

		// Notification management
		notifDelivery.RegisterRoutes(r, notifHandler)

		// DataRow management
		dataRowDelivery.RegisterRoutes(r, dataRowHandler)

		// Desk/Ticket management
		deskDelivery.RegisterRoutes(r, deskHandler)

		// Integration management
		integrationDelivery.RegisterRoutes(r, integrationHandler)
	})

	return r
}
