# 📄 Frontend Product Requirement Document (PRD)
## C2C Marketplace Web Client (Tokopedia-Style)

---

## 1. Executive Summary & Goals
* **Project Name:** TokoMarket Web Client (Tokopedia-Style E-Commerce)
* **Objective:** Build a responsive, fast, and visually stunning Single Page Application (SPA) in **React.js** that connects seamlessly with the `go-marketplace-backend` REST API.
* **Design Identity:** **Tokopedia-Style Aesthetic** featuring:
  * Brand Primary Color: **Tokopedia Emerald Green** (`#03AC0E` / `#00AA5B`).
  * Secondary Accent: Warm Orange/Coral (`#FA591D` for promos/discounts).
  * Surface: Clean whites, subtle card borders (`#E5E7E9`), elegant shadows, crisp typography (Inter / Open Sans), and smooth micro-animations.
  * Mobile-friendly responsive layout with desktop navbar and mobile bottom navigation.

---

## 2. Tech Stack & Frontend Architecture

### 2.1 Tech Stack
* **Core:** React.js 18+ (Hooks, Context API / Zustand for state)
* **Build Tool:** Vite (Ultra-fast HMR and bundling)
* **Routing:** React Router v6
* **Icons:** `lucide-react` (Crisp modern SVG icons)
* **HTTP Client:** Axios with Request/Response interceptors for JWT token management
* **Styling:** Modular CSS / Modern Vanilla CSS with CSS Variables design tokens

### 2.2 Directory Structure
```text
frontend/
├── public/
├── src/
│   ├── assets/                 # Brand logos, banners, category icons
│   ├── components/             # Reusable UI components
│   │   ├── common/             # Button, Modal, Input, Badge, Toast
│   │   ├── layout/             # Header/Navbar, Footer, MobileNav, Sidebar
│   │   ├── product/            # ProductCard, VariantSelector, RatingStars
│   │   ├── cart/               # CartStoreGroup, CartItemRow, CartSummary
│   │   └── seller/             # SellerSidebar, StatCard, OrderManageCard
│   ├── context/                # AuthContext, CartContext, WishlistContext
│   ├── hooks/                  # useAuth, useCart, useWishlist, useDebounce
│   ├── pages/                  # Page Views
│   │   ├── HomePage.jsx        # Banner carousel, categories, flash sale, product grid
│   │   ├── ProductDetailPage.jsx# Gallery, variant picker, store info, review list
│   │   ├── CartPage.jsx        # Multi-store cart grouping, item quantity adjuster
│   │   ├── CheckoutPage.jsx    # Address selector, courier & voucher picker, order summary
│   │   ├── WishlistPage.jsx    # User saved/favorite items
│   │   ├── OrderHistoryPage.jsx# Order status tabs (All, Processing, Shipped, Completed), review modal
│   │   ├── AuthModal.jsx       # Login & Register tabs with validation
│   │   └── seller/             # Seller Center
│   │       ├── SellerDashboard.jsx # Balance, revenue stats, mutation ledger
│   │       ├── SellerProducts.jsx  # Add/edit product & variants
│   │       ├── SellerOrders.jsx    # Ship orders & input tracking number
│   │       └── SellerVouchers.jsx  # Create store vouchers
│   ├── services/               # API Service Clients
│   │   ├── api.js              # Axios instance with base URL & JWT interceptor
│   │   ├── authService.js      # Register, login, profile, addresses
│   │   ├── productService.js   # List, detail, categories, store products
│   │   ├── cartService.js      # Cart CRUD, multi-store grouping
│   │   ├── orderService.js     # Checkout, user orders, ship, complete
│   │   ├── voucherService.js   # List vouchers, apply voucher preview
│   │   ├── walletService.js    # Withdrawals & balance mutations
│   │   └── wishlistService.js  # Add/remove/list wishlist
│   ├── styles/                 # Global styles & design system tokens
│   │   ├── variables.css       # Tokopedia green colors, spacing, shadows
│   │   └── global.css          # Reset, typography, utility classes
│   ├── App.jsx                 # Route definitions & global providers
│   └── main.jsx                # Application entrypoint
├── index.html
├── package.json
└── vite.config.js
```

---

## 3. UI/UX Design System (Tokopedia Theme Tokens)

```css
:root {
  /* Tokopedia Color Palette */
  --color-primary: #03AC0E;
  --color-primary-dark: #028A0B;
  --color-primary-light: #E8F5E9;
  --color-accent-orange: #FA591D;
  --color-accent-orange-light: #FFF0EB;
  --color-accent-red: #E02954;
  --color-accent-yellow: #FFC400;

  /* Neutrals */
  --bg-main: #F3F4F5;
  --bg-white: #FFFFFF;
  --text-primary: #212121;
  --text-secondary: #6D7588;
  --text-disabled: #A6ABB5;
  --border-light: #E5E7E9;
  --border-focus: #03AC0E;

  /* Elevations */
  --shadow-sm: 0 1px 4px rgba(141, 150, 170, 0.1);
  --shadow-md: 0 4px 12px rgba(141, 150, 170, 0.18);
  --shadow-lg: 0 8px 24px rgba(141, 150, 170, 0.22);

  /* Radius */
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 16px;
  --radius-full: 9999px;
}
```

---

## 4. Key Functional Modules & Screens

### 4.1 Header & Navigation (Tokopedia Desktop & Mobile)
* **Desktop Navbar**:
  * Top Bar: Download App, Mitra Tokopedia, Seller Center link, Help Center.
  * Main Bar: Brand Logo, Category Dropdown Mega Menu, Global Search Input with Auto-Search, Wishlist Icon (with badge counter), Cart Icon (with badge counter), Notification Icon, and User Profile / Login Register buttons.
* **Mobile Bottom Navigation**:
  * Home, Feed/Wishlist, Official Store, Cart, and Account Profile.

### 4.2 Home Page (`/`)
* **Hero Banner Carousel**: Promo highlights and discount events.
* **Category Icon Grid**: 6+ visual categories (Elektronik, Fashion, Laptop, etc.).
* **Special Promo / Voucher Slider**: Tokopedia-style voucher cards with "Klaim" button.
* **Infinite / Paginated Product Feed**:
  * Product card: Product image, official store badge, product title (2 lines max), price (bold black), discount tag (e.g. `10%`), city location badge, rating star with total sold counter.

### 4.3 Product Detail Page (`/products/:id`)
* **Left Gallery**: Multi-image preview with hover zoom.
* **Center Details**: Title, rating stars, sold count, price, variant selector (Size, RAM, Color with real-time price & stock updates), product description tabs.
* **Store Card**: Store name, city location, store badge, "Follow" & "Chat Toko" buttons.
* **Right Sticky Purchase Card**: Quantity selector, total subtotal, "Beli Langsung" & "+ Keranjang" buttons, Wishlist toggle.
* **Bottom Reviews Section**: Average rating summary (e.g. 5.0 / 5.0) and buyer review comments.

### 4.4 Multi-Store Shopping Cart Page (`/cart`)
* **Grouped Per Store**:
  * Store checkbox & Store name header with location badge.
  * Product rows: Thumbnail, Variant label, Unit price, Quantity counter (`-` `qty` `+`), and Delete button.
* **Sticky Summary Sidebar**:
  * Total items count, Subtotal price, Voucher input button, and "Beli (Checkout)" CTA button.

### 4.5 Checkout & Payment Page (`/checkout`)
* **Delivery Address Section**: Displays primary address with "Pilih Alamat Lain" modal.
* **Order Item Breakdown (Per Store)**:
  * Courier selection dropdown per store (JNE Reguler, SiCepat, dll).
  * Auto-calculated shipping fee per store.
* **Voucher Selection**: Platform voucher & Store voucher modal picker with real-time discount preview.
* **Payment Method**: BCA Virtual Account, Mandiri VA, QRIS.
* **Place Order CTA**: Triggers atomic multi-store checkout API and presents payment instruction dialog.

### 4.6 Seller Center Dashboard (`/seller`)
* **Store Header & Balance Card**: Total balance, "Tarik Dana" (Withdrawal) button, mutation history.
* **Product Catalog Management**: Add new product with multiple variants, edit prices & stock, delete products.
* **Order Fulfillment Center**:
  * Incoming store orders list.
  * "Input Resi Pengiriman" action with courier tracking number.
* **Store Voucher Creator**: Form to create custom store discount vouchers.

### 4.7 Order History & Reviews (`/orders`)
* Status tabs: *Semua, Menunggu Pembayaran, Diproses, Dikirim, Selesai, Dibatalkan*.
* Order card with tracking number and "Selesaikan Pesanan" (Confirm receipt) button.
* Interactive Review Modal: 1–5 Star interactive rating + text comment submission.

---

## 5. API Integration Mapping

| Screen / Feature | Backend API Endpoint | HTTP Method |
| :--- | :--- | :--- |
| **Auth** | `/api/v1/auth/login`, `/api/v1/auth/register` | `POST` |
| **User Profile & Addresses** | `/api/v1/auth/profile`, `/api/v1/user/addresses` | `GET`, `POST` |
| **Categories & Products** | `/api/v1/categories`, `/api/v1/products`, `/api/v1/products/:id` | `GET` |
| **Wishlist** | `/api/v1/wishlist`, `/api/v1/wishlist/:product_id` | `GET`, `POST`, `DELETE` |
| **Vouchers** | `/api/v1/vouchers`, `/api/v1/vouchers/apply` | `GET`, `POST` |
| **Cart** | `/api/v1/cart`, `/api/v1/cart/items`, `/api/v1/cart/items/:id` | `GET`, `POST`, `PUT`, `DELETE` |
| **Checkout & Orders** | `/api/v1/orders/checkout`, `/api/v1/orders`, `/api/v1/orders/:id` | `POST`, `GET` |
| **Order Completion** | `/api/v1/orders/items/:id/complete` | `POST` |
| **Product Reviews** | `/api/v1/reviews`, `/api/v1/products/:id/reviews` | `POST`, `GET` |
| **Seller Store & Balance** | `/api/v1/seller/store`, `/api/v1/seller/balance`, `/api/v1/seller/mutations` | `GET` |
| **Seller Products** | `/api/v1/seller/products`, `/api/v1/seller/products/:id` | `POST`, `PUT`, `DELETE` |
| **Seller Orders & Shipping**| `/api/v1/seller/orders`, `/api/v1/seller/orders/:id/ship` | `GET`, `PATCH` |
| **Seller Withdrawals** | `/api/v1/seller/withdrawals` | `POST`, `GET` |
