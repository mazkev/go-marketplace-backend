package http

import (
	"github.com/gin-gonic/gin"
	"go-market/internal/delivery/http/handler"
	"go-market/internal/delivery/http/middleware"
	"go-market/pkg/jwt"
)

type RouterParams struct {
	JWTService      jwt.JWTService
	AuthHandler     *handler.AuthHandler
	StoreHandler    *handler.StoreHandler
	ProductHandler  *handler.ProductHandler
	CartHandler     *handler.CartHandler
	OrderHandler    *handler.OrderHandler
	ReviewHandler   *handler.ReviewHandler
	VoucherHandler  *handler.VoucherHandler
	PaymentHandler  *handler.PaymentHandler
	WalletHandler   *handler.WalletHandler
	WishlistHandler *handler.WishlistHandler
}

func SetupRouter(p RouterParams) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "c2c-marketplace-backend"})
	})

	api := r.Group("/api/v1")
	{
		// 1. Auth Endpoints
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", p.AuthHandler.Register)
			authGroup.POST("/login", p.AuthHandler.Login)
		}

		// 2. Public Catalog, Categories & Vouchers
		api.GET("/categories", p.ProductHandler.ListCategories)
		api.GET("/products", p.ProductHandler.ListProducts)
		api.GET("/products/:id", p.ProductHandler.GetProductByID)
		api.GET("/products/:id/reviews", p.ReviewHandler.GetProductReviews)
		api.GET("/stores/:id", p.StoreHandler.GetStoreProfile)
		api.GET("/vouchers", p.VoucherHandler.ListAvailable)
		api.POST("/vouchers/apply", p.VoucherHandler.ApplyVoucher)

		// 3. Payment Gateway Webhook (Public with signature/payload)
		api.POST("/payments/webhook", p.PaymentHandler.HandleWebhook)

		// Protected Routes
		authMiddleware := middleware.AuthMiddleware(p.JWTService)
		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			// User Profile & Addresses
			protected.GET("/auth/profile", p.AuthHandler.GetProfile)
			protected.GET("/user/addresses", p.AuthHandler.GetAddresses)
			protected.POST("/user/addresses", p.AuthHandler.AddAddress)

			// Store Registration (Buyer becomes Seller)
			protected.POST("/stores", p.StoreHandler.CreateStore)

			// Seller Management Endpoints
			seller := protected.Group("/seller")
			{
				seller.GET("/store", p.StoreHandler.GetMyStore)
				seller.GET("/balance", p.StoreHandler.GetStoreBalance)
				seller.POST("/products", p.ProductHandler.CreateProduct)
				seller.PUT("/products/:id", p.ProductHandler.UpdateProduct)
				seller.DELETE("/products/:id", p.ProductHandler.DeleteProduct)
				seller.GET("/orders", p.OrderHandler.GetStoreOrders)
				seller.PATCH("/orders/:id/ship", p.OrderHandler.ShipOrderItem)

				// Seller Vouchers
				seller.POST("/vouchers", p.VoucherHandler.CreateStoreVoucher)

				// Seller Wallet & Withdrawals
				seller.POST("/withdrawals", p.WalletHandler.RequestWithdrawal)
				seller.GET("/withdrawals", p.WalletHandler.GetStoreWithdrawals)
				seller.GET("/mutations", p.WalletHandler.GetStoreMutations)
			}

			// Cart Endpoints
			cart := protected.Group("/cart")
			{
				cart.GET("", p.CartHandler.GetCart)
				cart.POST("/items", p.CartHandler.AddToCart)
				cart.PUT("/items/:id", p.CartHandler.UpdateCartItem)
				cart.DELETE("/items/:id", p.CartHandler.DeleteCartItem)
			}

			// Orders & Checkout
			orders := protected.Group("/orders")
			{
				orders.GET("", p.OrderHandler.GetUserOrders)
				orders.GET("/:id", p.OrderHandler.GetOrderByID)
				orders.POST("/checkout", p.OrderHandler.Checkout)
				orders.POST("/items/:id/complete", p.OrderHandler.CompleteOrderItem)
			}

			// Wishlist Endpoints
			wishlist := protected.Group("/wishlist")
			{
				wishlist.GET("", p.WishlistHandler.GetUserWishlist)
				wishlist.POST("/:product_id", p.WishlistHandler.AddToWishlist)
				wishlist.DELETE("/:product_id", p.WishlistHandler.RemoveFromWishlist)
			}

			// Reviews
			protected.POST("/reviews", p.ReviewHandler.CreateReview)
		}
	}

	return r
}
