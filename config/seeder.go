package config

import (
	"log"
	"time"

	"go-market/internal/domain"
	"go-market/pkg/hash"
	"gorm.io/gorm"
)

func SeedComprehensiveData(db *gorm.DB) {
	// Always sync and update image URLs for seeded products if missing
	productImages := map[string]string{
		"Laptop Gaming Legion Pro 5 Gen 8 AMD Ryzen 7":                   "https://images.unsplash.com/photo-1603302576837-37561b2e2302?auto=format&fit=crop&w=800&q=80",
		"Smartphone Ultra Pro Max 5G 256GB Snapdragon 8 Gen 3":           "https://images.unsplash.com/photo-1598327105666-5b89351aff97?auto=format&fit=crop&w=800&q=80",
		"TWS Wireless Earbuds ANC Active Noise Cancelling Bluetooth 5.3": "https://images.unsplash.com/photo-1590658268037-6bf12165a8df?auto=format&fit=crop&w=800&q=80",
		"Hoodie Oversize Unisex Cotton Fleece 330gsm Premium":             "https://images.unsplash.com/photo-1556905055-8f358a7a47b2?auto=format&fit=crop&w=800&q=80",
		"Kemeja Flannel Lengan Panjang Slim Fit Motif Kotak":              "https://images.unsplash.com/photo-1596755094514-f87e34085b2c?auto=format&fit=crop&w=800&q=80",
		"Mechanical Gaming Keyboard 75% RGB Hot-Swappable":               "https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=800&q=80",
		"Wireless Ultralight Gaming Mouse 49g Sensor 26000 DPI":           "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?auto=format&fit=crop&w=800&q=80",
		"Kopi Arabika Gayo Aceh Specialty Roast Beans 250g":               "https://images.unsplash.com/photo-1559056199-641a0ac8b55e?auto=format&fit=crop&w=800&q=80",
	}

	for name, imgURL := range productImages {
		db.Model(&domain.Product{}).Where("name = ? AND (image_url = '' OR image_url IS NULL)", name).Update("image_url", imgURL)
	}

	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount > 0 {
		return // Users already seeded
	}

	log.Println("Seeding comprehensive dummy data for marketplace...")

	// 1. Password Hash for all dummy accounts (Password123!)
	hashedPassword, _ := hash.HashPassword("Password123!")

	// 2. Users
	users := []domain.User{
		{
			Name:         "Budi Santoso (Buyer)",
			Email:        "buyer@market.com",
			PasswordHash: hashedPassword,
			Phone:        "081234567890",
			Role:         domain.RoleBuyer,
		},
		{
			Name:         "Kevin Gadget (Seller)",
			Email:        "gadget.seller@market.com",
			PasswordHash: hashedPassword,
			Phone:        "081234567891",
			Role:         domain.RoleSeller,
		},
		{
			Name:         "Siti Fashion (Seller)",
			Email:        "fashion.seller@market.com",
			PasswordHash: hashedPassword,
			Phone:        "081234567892",
			Role:         domain.RoleSeller,
		},
		{
			Name:         "Rian Gaming (Seller)",
			Email:        "gaming.seller@market.com",
			PasswordHash: hashedPassword,
			Phone:        "081234567893",
			Role:         domain.RoleSeller,
		},
	}
	for i := range users {
		db.Create(&users[i])
	}

	// 3. User Addresses for Buyer
	addresses := []domain.UserAddress{
		{
			UserID:       users[0].ID,
			ReceiverName: "Budi Santoso",
			Phone:        "081234567890",
			FullAddress:  "Jl. Sudirman Kav. 21 No. 45, RT 01 / RW 05, Kebayoran Baru, Jakarta Selatan",
			CityID:       1,
			IsPrimary:    true,
		},
		{
			UserID:       users[0].ID,
			ReceiverName: "Budi (Kantor)",
			Phone:        "081234567890",
			FullAddress:  "Gedung Bursa Efek Tower 2 Lantai 12, SCBD, Jakarta Selatan",
			CityID:       1,
			IsPrimary:    false,
		},
	}
	for i := range addresses {
		db.Create(&addresses[i])
	}

	// 4. Stores
	stores := []domain.Store{
		{
			UserID:     users[1].ID,
			StoreName:  "Official Gadget Store",
			DomainSlug: "gadget-store",
			CityID:     1,
			Balance:    24500000.0,
		},
		{
			UserID:     users[2].ID,
			StoreName:  "Urban Fashion Corner",
			DomainSlug: "urban-fashion",
			CityID:     2,
			Balance:    11800000.0,
		},
		{
			UserID:     users[3].ID,
			StoreName:  "GameZone Official Store",
			DomainSlug: "gamezone-official",
			CityID:     3,
			Balance:    16200000.0,
		},
	}
	for i := range stores {
		db.Create(&stores[i])
	}

	// 5. Categories
	categories := []domain.Category{
		{Name: "Elektronik", Slug: "elektronik"},
		{Name: "Fashion Pria", Slug: "fashion-pria"},
		{Name: "Fashion Wanita", Slug: "fashion-wanita"},
		{Name: "Handphone & Tablet", Slug: "handphone-tablet"},
		{Name: "Komputer & Laptop", Slug: "komputer-laptop"},
		{Name: "Makanan & Minuman", Slug: "makanan-minuman"},
	}
	for i := range categories {
		db.Create(&categories[i])
	}

	// 6. Products & Variants
	// Product 1: Laptop Gaming
	var1Price := 17500000.0
	var2Price := 19800000.0
	prod1 := domain.Product{
		StoreID:     stores[0].ID,
		CategoryID:  categories[4].ID, // Komputer & Laptop
		Name:        "Laptop Gaming Legion Pro 5 Gen 8 AMD Ryzen 7",
		Description: "Performa kencang untuk gaming AAA dan software berat.\n- AMD Ryzen 7 7745HX\n- NVIDIA GeForce RTX 4060 8GB GDDR6\n- Layar 16 inci WQXGA 240Hz 100% sRGB\n- Garansi Resmi 3 Tahun Lenovo Indonesia + ADP.",
		ImageURL:    productImages["Laptop Gaming Legion Pro 5 Gen 8 AMD Ryzen 7"],
		Price:       18500000,
		Stock:       15,
		Weight:      2500,
		RatingAvg:   5.0,
		RatingCount: 2,
		Variants: []domain.ProductVariant{
			{VariantName: "RAM 16GB / SSD 512GB", PriceOverride: &var1Price, Stock: 8},
			{VariantName: "RAM 32GB / SSD 1TB", PriceOverride: &var2Price, Stock: 7},
		},
	}
	db.Create(&prod1)

	// Product 2: Smartphone Flagship
	prod2 := domain.Product{
		StoreID:     stores[0].ID,
		CategoryID:  categories[3].ID, // Handphone & Tablet
		Name:        "Smartphone Ultra Pro Max 5G 256GB Snapdragon 8 Gen 3",
		Description: "Flagship dengan kamera 200MP Leica optics dan layar AMOLED 120Hz LTPO.\n- Baterai 5000mAh 120W HyperCharge\n- IP68 Water & Dust Resistant\n- Garansi Resmi 1 Tahun.",
		ImageURL:    productImages["Smartphone Ultra Pro Max 5G 256GB Snapdragon 8 Gen 3"],
		Price:       12999000,
		Stock:       25,
		Weight:      220,
		RatingAvg:   4.9,
		RatingCount: 5,
		Variants: []domain.ProductVariant{
			{VariantName: "Titanium Gray", Stock: 15},
			{VariantName: "Midnight Black", Stock: 10},
		},
	}
	db.Create(&prod2)

	// Product 3: TWS Earbuds
	prod3 := domain.Product{
		StoreID:     stores[0].ID,
		CategoryID:  categories[0].ID, // Elektronik
		Name:        "TWS Wireless Earbuds ANC Active Noise Cancelling Bluetooth 5.3",
		Description: "Audio jernih dengan bass bertenaga dan fitur peredam bising aktif hingga 45dB. Daya tahan baterai hingga 32 jam dengan casing.",
		ImageURL:    productImages["TWS Wireless Earbuds ANC Active Noise Cancelling Bluetooth 5.3"],
		Price:       899000,
		Stock:       40,
		Weight:      150,
		RatingAvg:   4.8,
		RatingCount: 3,
	}
	db.Create(&prod3)

	// Product 4: Hoodie Oversize
	prod4 := domain.Product{
		StoreID:     stores[1].ID,
		CategoryID:  categories[1].ID, // Fashion Pria
		Name:        "Hoodie Oversize Unisex Cotton Fleece 330gsm Premium",
		Description: "Bahan katun fleece tebal, halus di kulit, dan tidak mudah berbulu. Potongan boxy modern cocok untuk pria maupun wanita.",
		ImageURL:    productImages["Hoodie Oversize Unisex Cotton Fleece 330gsm Premium"],
		Price:       149000,
		Stock:       100,
		Weight:      600,
		RatingAvg:   5.0,
		RatingCount: 4,
		Variants: []domain.ProductVariant{
			{VariantName: "Hitam - Size L", Stock: 50},
			{VariantName: "Sage Green - Size XL", Stock: 50},
		},
	}
	db.Create(&prod4)

	// Product 5: Kemeja Flannel
	prod5 := domain.Product{
		StoreID:     stores[1].ID,
		CategoryID:  categories[1].ID, // Fashion Pria
		Name:        "Kemeja Flannel Lengan Panjang Slim Fit Motif Kotak",
		Description: "Kemeja bahan katun wol lembut, sangat nyaman dipakai casual maupun semi-formal.",
		ImageURL:    productImages["Kemeja Flannel Lengan Panjang Slim Fit Motif Kotak"],
		Price:       189000,
		Stock:       60,
		Weight:      350,
		RatingAvg:   4.7,
		RatingCount: 2,
		Variants: []domain.ProductVariant{
			{VariantName: "Maroon Plaid - Size M", Stock: 30},
			{VariantName: "Navy Plaid - Size L", Stock: 30},
		},
	}
	db.Create(&prod5)

	// Product 6: Mechanical Keyboard
	prod6 := domain.Product{
		StoreID:     stores[2].ID,
		CategoryID:  categories[4].ID, // Komputer & Laptop
		Name:        "Mechanical Gaming Keyboard 75% RGB Hot-Swappable",
		Description: "Keyboard mekanikal layout 75% compact dengan gasket mount, knob volume, RGB per-key, dan koneksi 3-mode (Wireless, BT, Type-C).",
		ImageURL:    productImages["Mechanical Gaming Keyboard 75% RGB Hot-Swappable"],
		Price:       650000,
		Stock:       35,
		Weight:      900,
		RatingAvg:   5.0,
		RatingCount: 6,
		Variants: []domain.ProductVariant{
			{VariantName: "Red Linear Switch (Silent)", Stock: 20},
			{VariantName: "Blue Clicky Switch (Tactile)", Stock: 15},
		},
	}
	db.Create(&prod6)

	// Product 7: Gaming Mouse
	prod7 := domain.Product{
		StoreID:     stores[2].ID,
		CategoryID:  categories[4].ID, // Komputer & Laptop
		Name:        "Wireless Ultralight Gaming Mouse 49g Sensor 26000 DPI",
		Description: "Mouse gaming super ringan 49 gram dengan optical sensor PAW3395 dan baterai tahan hingga 80 jam pemakaian.",
		ImageURL:    productImages["Wireless Ultralight Gaming Mouse 49g Sensor 26000 DPI"],
		Price:       350000,
		Stock:       50,
		Weight:      200,
		RatingAvg:   4.8,
		RatingCount: 3,
	}
	db.Create(&prod7)

	// Product 8: Kopi Arabika
	prod8 := domain.Product{
		StoreID:     stores[1].ID,
		CategoryID:  categories[5].ID, // Makanan & Minuman
		Name:        "Kopi Arabika Gayo Aceh Specialty Roast Beans 250g",
		Description: "Biji kopi arabika pilihan dari dataran tinggi Gayo Aceh dengan notes fruity, floral, dan aftertaste manis seimbang.",
		ImageURL:    productImages["Kopi Arabika Gayo Aceh Specialty Roast Beans 250g"],
		Price:       75000,
		Stock:       80,
		Weight:      260,
		RatingAvg:   5.0,
		RatingCount: 8,
	}
	db.Create(&prod8)

	// 7. Vouchers
	percent10 := 10.0
	percent15 := 15.0
	percent5 := 5.0
	max50k := 50000.0
	max30k := 30000.0
	max100k := 100000.0
	potongan50k := 50000.0
	now := time.Now()
	nextYear := now.AddDate(1, 0, 0)

	vouchers := []domain.Voucher{
		{
			Code:            "DISKON10",
			VoucherType:     domain.VoucherTypePercentage,
			DiscountPercent: &percent10,
			MaxDiscount:     &max50k,
			MinSpend:        100000.0,
			Quota:           200,
			StartDate:       now.AddDate(0, -1, 0),
			EndDate:         nextYear,
			IsActive:        true,
		},
		{
			Code:           "HEMAT50RB",
			VoucherType:    domain.VoucherTypeFixed,
			DiscountAmount: &potongan50k,
			MinSpend:       300000.0,
			Quota:          100,
			StartDate:      now.AddDate(0, -1, 0),
			EndDate:        nextYear,
			IsActive:       true,
		},
		{
			Code:            "GADGETHEMAT",
			StoreID:         &stores[0].ID,
			VoucherType:     domain.VoucherTypePercentage,
			DiscountPercent: &percent5,
			MaxDiscount:     &max100k,
			MinSpend:        1000000.0,
			Quota:           50,
			StartDate:       now.AddDate(0, -1, 0),
			EndDate:         nextYear,
			IsActive:        true,
		},
		{
			Code:            "FASHIONKEREN",
			StoreID:         &stores[1].ID,
			VoucherType:     domain.VoucherTypePercentage,
			DiscountPercent: &percent15,
			MaxDiscount:     &max30k,
			MinSpend:        100000.0,
			Quota:           50,
			StartDate:       now.AddDate(0, -1, 0),
			EndDate:         nextYear,
			IsActive:        true,
		},
	}
	for i := range vouchers {
		db.Create(&vouchers[i])
	}

	// 8. Completed Sample Order & Reviews
	expiredDate := now.Add(24 * time.Hour)
	sampleOrder := domain.Order{
		UserID:           users[0].ID,
		InvoiceNumber:    "INV/20260815/99810001",
		TotalAmount:      18510000.0,
		DiscountAmount:   50000.0,
		FinalAmount:      18460000.0,
		VoucherCode:      "DISKON10",
		PaymentMethod:    "BCA_VA",
		VANumber:         "88081234567890",
		PaymentStatus:    domain.PaymentStatusPaid,
		PaymentExpiredAt: &expiredDate,
	}
	db.Create(&sampleOrder)

	sampleOrderItem := domain.OrderItem{
		OrderID:        sampleOrder.ID,
		StoreID:        stores[0].ID,
		ProductID:      prod1.ID,
		VariantID:      &prod1.Variants[0].ID,
		Quantity:       1,
		Price:          17500000.0,
		ShippingCost:   10000.0,
		CourierName:    "JNE Reguler",
		TrackingNumber: "JNE-SAMPLE-778899",
		Status:         domain.OrderItemStatusCompleted,
	}
	db.Create(&sampleOrderItem)

	sampleReview := domain.Review{
		OrderItemID: sampleOrderItem.ID,
		UserID:      users[0].ID,
		ProductID:   prod1.ID,
		Rating:      5,
		Comment:     "Laptop original, performa sangat kencang, packing kayu dan bubble wrap sangat tebal dan aman!",
	}
	db.Create(&sampleReview)

	// 9. Store Balance Mutations
	mutations := []domain.BalanceMutation{
		{
			StoreID:      stores[0].ID,
			Amount:       17510000.0,
			Type:         domain.MutationTypeCredit,
			Description:  "Escrow released for completed order item #1 (Order INV/20260815/99810001)",
			BalanceAfter: 24500000.0,
		},
		{
			StoreID:      stores[1].ID,
			Amount:       5000000.0,
			Type:         domain.MutationTypeCredit,
			Description:  "Penjualan 30x Hoodie Oversize Cotton Fleece",
			BalanceAfter: 11800000.0,
		},
		{
			StoreID:      stores[2].ID,
			Amount:       6500000.0,
			Type:         domain.MutationTypeCredit,
			Description:  "Penjualan 10x Mechanical Gaming Keyboard",
			BalanceAfter: 16200000.0,
		},
	}
	for i := range mutations {
		db.Create(&mutations[i])
	}

	log.Println("Dummy data successfully seeded: Users, Stores, Products, Variants, Vouchers, Orders & Reviews.")
}
