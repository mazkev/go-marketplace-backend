# 🛍️ C2C Marketplace Back-End Service (Tokopedia-Style)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue)
![Database](https://img.shields.io/badge/Database-SQLite%20%7C%20PostgreSQL-003B57?style=flat&logo=sqlite)
![Framework](https://img.shields.io/badge/Framework-Gin-008ECF)
![License](https://img.shields.io/badge/License-MIT-green)

A modular, scalable, and secure **RESTful API Back-End** for a **C2C (Consumer-to-Consumer) Marketplace** built with **Golang**, **Gin**, and **GORM** following **Clean Architecture** principles.

---

## 📑 Table of Contents
- [Architecture & Tech Stack](#-architecture--tech-stack)
- [Project Structure](#-project-structure)
- [Key Features](#-key-features)
- [Getting Started](#-getting-started)
- [API Endpoints](#-api-endpoints)
- [Testing](#-testing)
- [License](#-license)

---

## 🛠 Architecture & Tech Stack

* **Language:** Go (Golang) v1.21+
* **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
* **ORM:** [GORM](https://gorm.io/)
* **Database Driver (Pure Go):** [`github.com/glebarez/sqlite`](https://github.com/glebarez/sqlite) (Zero-CGO SQLite driver with ANSI SQL compatibility, seamlessly portable to PostgreSQL)
* **Authentication:** JWT (JSON Web Tokens) with Access & Refresh Token pair
* **Security:** `bcrypt` password hashing & HMAC-SHA256 signature verification

---

## 📂 Project Structure

```text
marketplace-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point & dependency injection wire-up
├── config/
│   ├── config.go                   # Environment configuration loader
│   └── database.go                 # Database connection & Auto-Migration setup
├── internal/
│   ├── domain/                     # Entities, DTOs & Business Interfaces
│   │   ├── user.go                 # User, Address & Role definitions
│   │   ├── store.go                # Store entity & balance
│   │   ├── category.go             # Category hierarchy
│   │   ├── product.go              # Products & Variants
│   │   ├── cart.go                 # Multi-Store Cart
│   │   ├── order.go                # Orders, OrderItems & Payment Webhooks
│   │   ├── voucher.go              # Platform & Store Vouchers
│   │   ├── wallet.go               # Store Withdrawals & Mutation Ledger
│   │   ├── wishlist.go             # User Wishlist
│   │   └── review.go               # Ratings & Reviews
│   ├── repository/                 # Data Access Layer (GORM)
│   ├── usecase/                    # Business Logic Layer
│   └── delivery/
│       └── http/                   # HTTP Delivery (Handlers & Middlewares)
│           ├── handler/
│           ├── middleware/
│           └── router.go
├── pkg/
│   ├── hash/                       # Bcrypt password hashing
│   ├── jwt/                        # JWT generator & validator
│   └── response/                   # Standardized JSON response envelope
├── tests/
│   └── integration_test.go         # Complete automated integration tests
├── test.http                       # REST Client interactive request collection
├── prd.md                          # Product Requirement Document
└── README.md
```

---

## ✨ Key Features

1. **Clean Architecture Separation:** Independent layers (Domain, Repository, Usecase, Delivery) ensuring loose coupling and testability.
2. **User & Merchant Management:**
   - Multi-role support (`buyer`, `seller`, `admin`).
   - Seamless store onboarding for existing users without separate accounts.
   - Multi-address management with primary address selection.
3. **Product Catalog & Multi-Variants:**
   - Products with flexible variants (Size, Color, etc.) and stock overrides.
   - Search, category filtering, price range, and pagination.
4. **Multi-Store Shopping Cart:**
   - Holds products from multiple independent sellers concurrently.
   - Groups items by store for transparent checkout.
5. **Transactional Checkout & Escrow System:**
   - Atomic database transactions (`db.Transaction`) for concurrent stock deduction.
   - Escrow payment protection: Store balance credited only after buyer confirms receipt.
6. **Voucher & Promo Engine:**
   - Platform & Store vouchers (Percentage with cap or Fixed amount discount).
   - Validations: Minimum spend, quota, and validity period.
7. **Payment Gateway Simulator & Webhook:**
   - Virtual Account generation and idempotent webhook status processing (`PENDING` -> `PAID` -> `PROCESSING`).
8. **Store Withdrawal & Mutation Ledger:**
   - Seller payout requests with debit deductions.
   - Full credit/debit audit trail for merchant balances.
9. **Wishlist & Rating System:**
   - Bookmark favorite products.
   - 1–5 star reviews restricted to completed orders with dynamic rating recalculations.

---

## 🚀 Getting Started

### 1. Prerequisites
- [Go 1.21+](https://go.dev/dl/)

### 2. Clone and Setup
```bash
git clone https://github.com/<your-username>/go-market.git
cd go-market
```

### 3. Environment Configuration
Create a `.env` file from `.env.example`:
```bash
cp .env.example .env
```

Default `.env` contents:
```env
PORT=8080
DB_DRIVER=sqlite
DB_DSN=marketplace.db
JWT_SECRET=super-secret-jwt-key-change-in-production-12345
```

### 4. Run the Application
```bash
go run ./cmd/api
```
The server will start at `http://localhost:8080` and automatically run database migrations and seed initial categories and vouchers.

---

## 📡 API Endpoints Summary

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Register new account | Public |
| `POST` | `/api/v1/auth/login` | Login and get JWT token pair | Public |
| `GET` | `/api/v1/categories` | Get category tree | Public |
| `GET` | `/api/v1/products` | Search & filter products | Public |
| `GET` | `/api/v1/products/:id` | Get product details & reviews | Public |
| `POST` | `/api/v1/stores` | Register new store | Buyer |
| `POST` | `/api/v1/seller/products` | Create product with variants | Seller |
| `GET` | `/api/v1/cart` | Get cart items grouped by store | Buyer |
| `POST` | `/api/v1/cart/items` | Add item to cart | Buyer |
| `POST` | `/api/v1/orders/checkout` | Multi-store atomic checkout | Buyer |
| `POST` | `/api/v1/payments/webhook` | Payment Gateway webhook notification | Public |
| `PATCH` | `/api/v1/seller/orders/:id/ship` | Input tracking number | Seller |
| `POST` | `/api/v1/orders/items/:id/complete` | Complete order & release escrow | Buyer |
| `POST` | `/api/v1/seller/withdrawals` | Request store balance withdrawal | Seller |
| `GET` | `/api/v1/seller/mutations` | View balance transaction ledger | Seller |
| `POST` | `/api/v1/reviews` | Submit review on completed item | Buyer |
| `GET` | `/api/v1/wishlist` | Get user saved items | Buyer |

> 💡 **Tip:** Use [test.http](file:///d:/prog/go-market/test.http) with the REST Client extension for interactive API testing.

---

## 🧪 Testing

Run all automated unit and integration tests:
```bash
go test -v ./tests
```

---

## 📄 License
This project is licensed under the MIT License.
