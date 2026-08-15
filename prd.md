# 📄 Product Requirement Document (PRD)
## C2C Marketplace Back-End Service (Tokopedia-Style)

---

## 1. Executive Summary & Goals
* **Project Name:** C2C Marketplace Back-End Service
* **Objective:** Build a reliable, secure, and modular *RESTful API back-end* for a large-scale *C2C (Consumer-to-Consumer) Marketplace* platform.
* **Core Philosophy:** Utilize **Clean Architecture** to completely decouple business logic from data source implementations. This approach allows rapid MVP development using **SQLite**, while ensuring seamless migration to **PostgreSQL** when the application is ready for production.

---

## 2. Tech Stack & Architecture

### 2.1 Tech Stack (MVP Phase)
* **Language:** Go (Golang) v1.21+
* **Framework / Router:** Gin / Fiber / Echo
* **Database (MVP):** SQLite via `gorm.io/driver/sqlite`
* **ORM:** GORM (`gorm.io/gorm`)
* **Authentication:** JWT (JSON Web Token) with *Access Token* & *Refresh Token*
* **Password Hashing:** `bcrypt`

### 2.2 System Architecture (Clean Architecture Layout)
```text
marketplace-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point & dependency injection
├── config/
│   └── database.go                 # Database connection & Auto-Migration setup
├── internal/
│   ├── domain/                     # Entities & Data Structures (DB Models)
│   ├── repository/                 # Data Access Layer (GORM queries)
│   ├── usecase/                    # Core Business Logic Layer
│   └── delivery/
│       └── http/
│           ├── handler/            # HTTP Request/Response Controllers
│           ├── middleware/         # Auth, Logger, Rate-limiting
│           └── router.go           # Route declarations
└── pkg/                            # Shared utilities (JWT, Response Wrappers)
```

---

## 3. User Roles & Permissions

* **Buyer:**
  * Browse and search products/stores.
  * Manage cart items (Multi-Store Cart).
  * Proceed to checkout, perform payments, and submit product reviews.
* **Seller:**
  * Open and manage store profiles.
  * Manage product catalogs (CRUD, variants, and stock).
  * Process incoming orders and attach shipping tracking numbers.
  * Withdraw store balance (Store Balance Withdrawal).
* **Admin Platform:**
  * Moderate stores and products (suspend/verify).
  * Manage global product categories.
  * Audit transaction records and balance payout requests.

---

## 4. Functional Requirements & Core Modules

### 4.1 Auth & User Profile Module
* **Features:**
  * User Registration & Login (Email/Password).
  * Profile management and Multi-Address system (Primary Address vs Additional Addresses).
  * Role-Based Access Control (Standard users can seamlessly register a store without creating a new account).

### 4.2 Merchant / Store Management Module
* **Features:**
  * New store registration (Store name, domain slug, and store location/city).
  * Internal Store Showcase (Etalase) management.
  * Store Balance tracking (Automatically credited upon order completion).

### 4.3 Product Management Module
* **Features:**
  * Product CRUD + Support for Product Variants (Size, Color, etc.).
  * Inventory & Price Management (Safe Stock Deduction during concurrent transactions).
  * Hierarchical product categories and sub-categories.

### 4.4 Multi-Store Cart & Checkout System
* **Features:**
  * **Multi-Store Cart:** Shopping cart capable of holding items from multiple distinct sellers simultaneously.
  * **Grouped Checkout:** Automatic splitting into Sub-Orders grouped by Store/Seller.
  * Per-store shipping fee calculation (Courier API / Mock Integration).
  * **Escrow mechanism:** Funds held by the platform until the Buyer confirms receipt of the order.

### 4.5 Review & Rating System
* **Features:**
  * Product reviews can only be created after the order status reaches `COMPLETED`.
  * 1–5 Star Rating + Text Comments.
  * Automatic aggregation for calculating average product and store ratings.

---

## 5. Database Schema (SQLite Entities)

```sql
-- Users & Store
users             (id, name, email, password_hash, phone, role, created_at, updated_at)
stores            (id, user_id, store_name, domain_slug, city_id, balance, created_at, updated_at)
user_addresses    (id, user_id, receiver_name, phone, full_address, city_id, is_primary)

-- Products & Categories
categories        (id, parent_id, name, slug)
products          (id, store_id, category_id, name, description, price, stock, weight, created_at)
product_variants (id, product_id, variant_name, price_override, stock)

-- Shopping & Orders
carts             (id, user_id, product_id, variant_id, quantity)
orders            (id, user_id, invoice_number, total_amount, payment_status, created_at)
order_items       (id, order_id, store_id, product_id, quantity, price, shipping_cost, courier_name, tracking_number, status)
reviews           (id, order_item_id, user_id, product_id, rating, comment, created_at)
```

---

## 6. Core RESTful API Specifications

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | User Registration | ❌ No |
| `POST` | `/api/v1/auth/login` | Login & Obtain JWT Token Pair | ❌ No |
| `GET` | `/api/v1/products` | Get Product List (Search, Filter, Pagination) | ❌ No |
| `POST` | `/api/v1/stores` | Create / Register New Store | ✅ Yes (Buyer) |
| `POST` | `/api/v1/cart/items` | Add Item to Shopping Cart | ✅ Yes (Buyer) |
| `GET` | `/api/v1/cart` | Get Cart Items (Grouped by Store) | ✅ Yes (Buyer) |
| `POST` | `/api/v1/orders/checkout` | Create Order (Multi-store checkout) | ✅ Yes (Buyer) |
| `PATCH` | `/api/v1/seller/orders/:id/ship` | Input Shipping Tracking Number | ✅ Yes (Seller) |

---

## 7. Non-Functional Requirements

* **Database Portability:** Repository logic must rely purely on GORM / ANSI SQL abstractions so migrating to PostgreSQL later only requires updating the database driver without modifying Usecase logic.
* **Performance:** API latency $< 200\text{ ms}$ for read operations.
* **Data Consistency:** Order creation and stock deduction processes must strictly use Database Transactions (`db.Begin()`, `tx.Commit()`, `tx.Rollback()`).
* **Security:** Secure password hashing via `bcrypt`, JWT-based route authentication middleware, and environment variable configuration via `.env` files.
