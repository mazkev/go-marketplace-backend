package main

import (
	"log"

	"go-market/config"
	deliveryHttp "go-market/internal/delivery/http"
	"go-market/internal/delivery/http/handler"
	"go-market/internal/repository"
	"go-market/internal/usecase"
	"go-market/pkg/jwt"
)

func main() {
	// 1. Load Configurations
	cfg := config.LoadConfig()

	// 2. Initialize Database & Auto-Migrations
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 3. Initialize Shared Services
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)

	// 4. Initialize Repositories (Data Access Layer)
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	voucherRepo := repository.NewVoucherRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	wishlistRepo := repository.NewWishlistRepository(db)

	// 5. Initialize Usecases (Business Logic Layer)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)
	storeUsecase := usecase.NewStoreUsecase(storeRepo, userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo, storeRepo)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, cartRepo, userRepo, storeRepo, productRepo, voucherRepo, walletRepo)
	reviewUsecase := usecase.NewReviewUsecase(reviewRepo, orderRepo, productRepo)
	voucherUsecase := usecase.NewVoucherUsecase(voucherRepo, storeRepo)
	paymentUsecase := usecase.NewPaymentUsecase(orderRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, storeRepo)
	wishlistUsecase := usecase.NewWishlistUsecase(wishlistRepo, productRepo)

	// 6. Initialize HTTP Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	storeHandler := handler.NewStoreHandler(storeUsecase)
	productHandler := handler.NewProductHandler(productUsecase, categoryRepo)
	cartHandler := handler.NewCartHandler(cartUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	reviewHandler := handler.NewReviewHandler(reviewUsecase)
	voucherHandler := handler.NewVoucherHandler(voucherUsecase)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	wishlistHandler := handler.NewWishlistHandler(wishlistUsecase)

	// 7. Setup Router & Routes
	router := deliveryHttp.SetupRouter(deliveryHttp.RouterParams{
		JWTService:      jwtService,
		AuthHandler:     authHandler,
		StoreHandler:    storeHandler,
		ProductHandler:  productHandler,
		CartHandler:     cartHandler,
		OrderHandler:    orderHandler,
		ReviewHandler:   reviewHandler,
		VoucherHandler:  voucherHandler,
		PaymentHandler:  paymentHandler,
		WalletHandler:   walletHandler,
		WishlistHandler: wishlistHandler,
	})

	// 8. Start Server
	log.Printf("Starting C2C Marketplace Server on port %s...", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
