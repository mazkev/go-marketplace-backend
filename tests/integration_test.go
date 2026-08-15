package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-market/config"
	deliveryHttp "go-market/internal/delivery/http"
	"go-market/internal/delivery/http/handler"
	"go-market/internal/repository"
	"go-market/internal/usecase"
	"go-market/pkg/jwt"
)

func setupTestApp(t *testing.T) http.Handler {
	dbFile := fmt.Sprintf("test_%d.db", time.Now().UnixNano())
	t.Cleanup(func() {
		_ = os.Remove(dbFile)
	})

	cfg := &config.Config{
		DatabaseDriver:   "sqlite",
		DatabaseDSN:      dbFile,
		JWTSecret:        "test-secret-key-12345",
		JWTAccessExpiry:  1 * time.Hour,
		JWTRefreshExpiry: 24 * time.Hour,
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize test db: %v", err)
	}

	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)

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

	return deliveryHttp.SetupRouter(deliveryHttp.RouterParams{
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
}

func doRequest(handler http.Handler, method, url, token string, body interface{}) (*httptest.ResponseRecorder, map[string]interface{}) {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	return w, res
}

func TestCompleteC2CMarketplaceFlow(t *testing.T) {
	app := setupTestApp(t)

	// 1. Register Seller 1
	w, res := doRequest(app, "POST", "/api/v1/auth/register", "", map[string]interface{}{
		"name":     "Seller One",
		"email":    "seller1@example.com",
		"password": "password123",
		"phone":    "08123456789",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Seller 1 registration failed: %v", res)
	}
	seller1Data := res["data"].(map[string]interface{})
	seller1Token := seller1Data["access_token"].(string)

	// 2. Register Seller 2
	w, res = doRequest(app, "POST", "/api/v1/auth/register", "", map[string]interface{}{
		"name":     "Seller Two",
		"email":    "seller2@example.com",
		"password": "password123",
		"phone":    "08123456780",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Seller 2 registration failed: %v", res)
	}
	seller2Data := res["data"].(map[string]interface{})
	seller2Token := seller2Data["access_token"].(string)

	// 3. Register Buyer
	w, res = doRequest(app, "POST", "/api/v1/auth/register", "", map[string]interface{}{
		"name":     "Buyer Budi",
		"email":    "buyer@example.com",
		"password": "password123",
		"phone":    "08123456781",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Buyer registration failed: %v", res)
	}
	buyerData := res["data"].(map[string]interface{})
	buyerToken := buyerData["access_token"].(string)

	// 4. Create Store for Seller 1
	w, res = doRequest(app, "POST", "/api/v1/stores", seller1Token, map[string]interface{}{
		"store_name":  "Toko Elektronik Jaya",
		"domain_slug": "tokoelektronik",
		"city_id":     1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Seller 1 store creation failed: %v", res)
	}

	// 5. Create Store for Seller 2
	w, res = doRequest(app, "POST", "/api/v1/stores", seller2Token, map[string]interface{}{
		"store_name":  "Toko Fashion Mantap",
		"domain_slug": "tokofashion",
		"city_id":     2,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Seller 2 store creation failed: %v", res)
	}

	// 6. Seller 1 creates Product 1 (Laptop)
	w, res = doRequest(app, "POST", "/api/v1/seller/products", seller1Token, map[string]interface{}{
		"category_id": 5,
		"name":        "Laptop Gaming Pro",
		"description": "Laptop kencang untuk coding & gaming",
		"price":       15000000,
		"stock":       10,
		"weight":      2000,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Product 1 creation failed: %v", res)
	}
	prod1Data := res["data"].(map[string]interface{})
	prod1ID := uint(prod1Data["id"].(float64))

	// 7. Seller 2 creates Product 2 with Variants (T-Shirt)
	w, res = doRequest(app, "POST", "/api/v1/seller/products", seller2Token, map[string]interface{}{
		"category_id": 2,
		"name":        "Kaos Polos Cotton",
		"description": "Bahan katun adem combed 30s",
		"price":       75000,
		"stock":       50,
		"weight":      200,
		"variants": []map[string]interface{}{
			{"variant_name": "Size M", "stock": 25},
			{"variant_name": "Size L", "stock": 25},
		},
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Product 2 creation failed: %v", res)
	}
	prod2Data := res["data"].(map[string]interface{})
	prod2ID := uint(prod2Data["id"].(float64))
	variants := prod2Data["variants"].([]interface{})
	variant1ID := uint(variants[0].(map[string]interface{})["id"].(float64))

	// 8. Wishlist Testing: Buyer adds Laptop to Wishlist
	w, res = doRequest(app, "POST", fmt.Sprintf("/api/v1/wishlist/%d", prod1ID), buyerToken, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Add to wishlist failed: %v", res)
	}
	w, res = doRequest(app, "GET", "/api/v1/wishlist", buyerToken, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Get wishlist failed: %v", res)
	}
	wishlistItems := res["data"].([]interface{})
	if len(wishlistItems) != 1 {
		t.Fatalf("Expected 1 wishlist item, got %d", len(wishlistItems))
	}

	// 9. Voucher Testing: Seller 1 creates Store Voucher
	w, res = doRequest(app, "POST", "/api/v1/seller/vouchers", seller1Token, map[string]interface{}{
		"code":             "ELEKTRONIKHEBAT",
		"voucher_type":     "percentage",
		"discount_percent": 5.0,
		"max_discount":     100000,
		"min_spend":        500000,
		"quota":            50,
		"start_date":       time.Now().Add(-1 * time.Hour),
		"end_date":         time.Now().Add(24 * time.Hour),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Store voucher creation failed: %v", res)
	}

	// 10. Buyer adds Address
	w, res = doRequest(app, "POST", "/api/v1/user/addresses", buyerToken, map[string]interface{}{
		"receiver_name": "Budi Santoso",
		"phone":         "08123456781",
		"full_address":  "Jl. Sudirman No. 45, Jakarta",
		"city_id":       1,
		"is_primary":    true,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Address creation failed: %v", res)
	}
	addrData := res["data"].(map[string]interface{})
	addrID := uint(addrData["id"].(float64))

	// 11. Buyer adds items from both stores to Cart
	// Add Laptop (Store 1)
	w, res = doRequest(app, "POST", "/api/v1/cart/items", buyerToken, map[string]interface{}{
		"product_id": prod1ID,
		"quantity":   1,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Add product 1 to cart failed: %v", res)
	}

	// Add T-Shirt Size M (Store 2)
	w, res = doRequest(app, "POST", "/api/v1/cart/items", buyerToken, map[string]interface{}{
		"product_id": prod2ID,
		"variant_id": variant1ID,
		"quantity":   2,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Add product 2 with variant to cart failed: %v", res)
	}

	// 12. Buyer gets Cart grouped by Store
	w, res = doRequest(app, "GET", "/api/v1/cart", buyerToken, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Get cart failed: %v", res)
	}
	cartData := res["data"].(map[string]interface{})
	storesInCart := cartData["stores"].([]interface{})
	if len(storesInCart) != 2 {
		t.Fatalf("Expected cart grouped into 2 stores, got %d", len(storesInCart))
	}

	// 13. Buyer Checkout with Platform Voucher "DISKON10"
	w, res = doRequest(app, "POST", "/api/v1/orders/checkout", buyerToken, map[string]interface{}{
		"address_id":   addrID,
		"voucher_code": "DISKON10",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Checkout with voucher failed: %v", res)
	}
	orderData := res["data"].(map[string]interface{})
	invoiceNo := orderData["invoice_number"].(string)
	discountAmount := orderData["discount_amount"].(float64)
	if discountAmount <= 0 {
		t.Fatalf("Expected voucher discount to be applied, got %v", discountAmount)
	}
	orderItems := orderData["order_items"].([]interface{})
	if len(orderItems) != 2 {
		t.Fatalf("Expected 2 order items in checkout result, got %d", len(orderItems))
	}

	var store1OrderItemID uint
	for _, item := range orderItems {
		itemMap := item.(map[string]interface{})
		if uint(itemMap["product_id"].(float64)) == prod1ID {
			store1OrderItemID = uint(itemMap["id"].(float64))
		}
	}

	// 14. Payment Webhook Simulation (Simulate Gateway Callback)
	w, res = doRequest(app, "POST", "/api/v1/payments/webhook", "", map[string]interface{}{
		"invoice_number": invoiceNo,
		"amount":         orderData["final_amount"].(float64),
		"payment_status": "SETTLEMENT",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Payment webhook processing failed: %v", res)
	}

	// 15. Seller 1 ships order item
	w, res = doRequest(app, "PATCH", fmt.Sprintf("/api/v1/seller/orders/%d/ship", store1OrderItemID), seller1Token, map[string]interface{}{
		"tracking_number": "JNE-TRACK-998877",
		"courier_name":    "JNE YES",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Ship order item failed: %v", res)
	}

	// 16. Buyer confirms order completion -> releases escrow to seller 1 balance & records mutation
	w, res = doRequest(app, "POST", fmt.Sprintf("/api/v1/orders/items/%d/complete", store1OrderItemID), buyerToken, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Complete order item failed: %v", res)
	}

	// 17. Seller 1 checks balance & mutation history
	w, res = doRequest(app, "GET", "/api/v1/seller/mutations", seller1Token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Get mutations failed: %v", res)
	}
	mutations := res["data"].([]interface{})
	if len(mutations) == 0 {
		t.Fatalf("Expected at least 1 mutation record for escrow credit")
	}

	// 18. Seller 1 requests Withdrawal to Bank Account
	w, res = doRequest(app, "POST", "/api/v1/seller/withdrawals", seller1Token, map[string]interface{}{
		"bank_name":      "BCA",
		"account_number": "1234567890",
		"account_holder": "Seller One",
		"amount":         5000000,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Withdrawal request failed: %v", res)
	}

	// 19. Buyer creates Review for completed order item
	w, res = doRequest(app, "POST", "/api/v1/reviews", buyerToken, map[string]interface{}{
		"order_item_id": store1OrderItemID,
		"rating":        5,
		"comment":       "Barang sampai cepat dan laptopnya sangat kencang!",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Review creation failed: %v", res)
	}

	// 20. Verify Product Rating updated
	w, res = doRequest(app, "GET", fmt.Sprintf("/api/v1/products/%d", prod1ID), "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Get product failed: %v", res)
	}
	prodDetail := res["data"].(map[string]interface{})
	if prodDetail["rating_avg"].(float64) != 5.0 || prodDetail["rating_count"].(float64) != 1 {
		t.Fatalf("Expected product rating to be 5.0 with count 1, got %v (count %v)", prodDetail["rating_avg"], prodDetail["rating_count"])
	}
}
