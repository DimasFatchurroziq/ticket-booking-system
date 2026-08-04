package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// Impor layer internal aplikasi kamu
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/delivery/http"
	// "github.com/DimasFatchurroziq/ticket-booking-system/internal/user/delivery/http/middleware"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/delivery/http/route"
	// "github.com/DimasFatchurroziq/ticket-booking-system/internal/user/messaging"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/repository"
	"github.com/DimasFatchurroziq/ticket-booking-system/internal/user/usecase"
	// "github.com/DimasFatchurroziq/ticket-booking-system/internal/user/util"
)

type BootstrapConfig struct {
	DB          *pgxpool.Pool
	App         *fiber.App
	Log         *zap.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	Producer    *kafka.Writer
	RedisClient *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// 1. Setup Utilities & Helpers
	// Gunakan RedisClient dari config, JANGAN buat baru di sini!
	// secretKey := config.Config.GetString("JWT_SECRET")
	// tokenUtil := util.NewTokenUtil(secretKey, config.RedisClient)

	// 2. Setup Repositories (Data Access Layer)
	userRepository := repository.NewPostgresUserRepository(config.DB, config.Log)
	// contactRepository := repository.NewContactRepository(config.Log)
	// addressRepository := repository.NewAddressRepository(config.Log)

	// 3. Setup Messaging / Event Producers
	// userProducer := messaging.NewUserProducer(config.Producer, config.Log)
	// contactProducer := messaging.NewContactProducer(config.Producer, config.Log)
	// addressProducer := messaging.NewAddressProducer(config.Producer, config.Log)

	// 4. Setup Use Cases / Business Logic
	userUseCase := usecase.NewUserUsecase(userRepository, config.Log)
	// contactUseCase := usecase.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository, contactProducer)
	// addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, contactRepository, addressRepository, addressProducer)

	// 5. Setup Controllers / HTTP Handlers
	userHandler := http.NewUserHandler(config.App, userUseCase, config.Log, config.Validate)
	// contactController := http.NewContactController(contactUseCase, config.Log)
	// addressController := http.NewAddressController(addressUseCase, config.Log)
	// helloController := http.NewHelloController()

	// 6. Setup Middlewares
	// authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	// 7. Setup Routes & Delivery Layer
	routeConfig := route.RouteConfig{
		App:         config.App,
		UserHandler: userHandler,
		// ContactController: contactController,
		// AddressController: addressController,
		// AuthMiddleware:    authMiddleware,
		// HelloController:   helloController,
	}
	routeConfig.Setup()
}
