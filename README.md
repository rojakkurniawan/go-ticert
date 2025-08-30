# Ticert - Event Ticket Management System

Sistem manajemen tiket event yang komprehensif dibangun dengan Go, menampilkan arsitektur yang bersih dan memberikan solusi lengkap untuk mengelola event, kategori tiket, pemesanan, dan tugas administratif terkait tiket.

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Installation](#installation)
- [Running the Project](#running-the-project)
- [Using Docker](#using-docker)
- [API Documentation](#api-documentation)
- [Features](#features)
- [Configuration](#configuration)

## Overview

Ticert adalah sistem manajemen tiket event yang yang mengimplementasikan prinsip-prinsip clean architecture. Sistem ini mendukung multiple user roles (user dan admin) dan memberikan fitur untuk manajemen event termasuk pembuatan event, manajemen kategori tiket, sistem pemesanan, dan pelaporan administratif.

## Tech Stack

- **Go 1.23.3** - [Installation Guide](https://golang.org/doc/install)
- **Gin Framework** - [Documentation](https://gin-gonic.com/)
- **GORM** - [Documentation](https://gorm.io/)
- **MySQL 8.4** - [Installation Guide](https://dev.mysql.com/doc/mysql-installation-excerpt/8.0/en/)
- **Redis 8.2.1** - [Documentation](https://redis.io/documentation)
- **JWT v5** - [Documentation](https://github.com/golang-jwt/jwt)
- **Docker & Docker Compose** - [Installation Guide](https://docs.docker.com/get-docker/)
- **UUID** - Unique identifier generation
- **Bcrypt** - Password hashing
- **Godotenv** - Environment variable management
- **Gin Validator** - Request validation

## Getting Started

### Prerequisites

Pastikan Anda telah menginstall hal-hal berikut di sistem Anda:

- Go 1.23+
- MySQL 8.4+
- Redis 8.2+
- Docker dan Docker Compose (opsional, untuk deployment yang di-containerize)

### Installation

1. **Clone repository:**

   ```bash
   git clone https://github.com/rojakkurniawan/go-ticert
   cd ticert
   ```

2. **Copy environment file:**

   ```bash
   cp env.example .env
   ```

3. **Configure environment variables:**
   Edit file `.env` sesuai dengan environment lokal Anda:

   ```env
   # Server Configuration
   PORT=8080

   # Database Configuration (MySQL)
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=ticert_user
   DB_PASSWORD=ticert_password
   DB_NAME=ticert

   # Redis Configuration
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0

   # JWT Configuration
   JWT_ACCESS_SECRET=your_super_secret_access_key_here
   JWT_REFRESH_SECRET=your_super_secret_refresh_key_here
   JWT_ACCESS_EXPIRY=1
   JWT_REFRESH_EXPIRY=24

   # Gin Mode
   GIN_MODE=release
   ```

4. **Install dependencies:**

   ```bash
   go mod tidy
   ```

5. **Set up database:**
   - Buat database MySQL bernama `ticert`
   - Aplikasi akan otomatis menjalankan migrasi saat startup

## Running the Project

### Local Development

1. **Start the application:**

   ```bash
   go run main.go
   ```

2. **Server akan start di** `http://localhost:8080`

### Using Docker

#### Quick Start dengan Docker Compose

1. **Build dan jalankan semua services:**

   ```bash
   docker-compose up --build
   ```

   Ini akan:

   - Start container database MySQL
   - Start container Redis
   - Build dan start aplikasi Ticert
   - Setup networking antar container
   - Apply health checks

2. **Jalankan di background:**

   ```bash
   docker-compose up -d --build
   ```

3. **Stop services:**
   ```bash
   docker-compose down
   ```

#### Manual Docker Build

1. **Build Docker image:**

   ```bash
   docker build -t ticert:latest .
   ```

2. **Jalankan dengan database eksternal:**
   ```bash
   docker run -p 8080:8080 --env-file .env ticert:latest
   ```

## API Documentation

📖 **[Ticert API Documentation](https://documenter.getpostman.com/view/31887101/2sB3Hhs26W)**

API documentation mencakup:

- Authentication endpoints
- User management
- Event management
- Category management
- Order management
- Reporting system

## Features

### Core Event Features

- **Multi-role Authentication** - Support untuk user dan admin
- **Event Management** - Pembuatan dan manajemen event dengan detail lengkap
- **Category Management** - Manajemen kategori tiket dengan harga dan kuantitas
- **Order System** - Sistem pemesanan tiket dengan status tracking
- **Ticket Generation** - Generate tiket unik dengan kode tiket
- **User Management** - Manajemen profil user dan role

### Security & Authentication

- **JWT Authentication** - Secure token-based authentication
- **Role-based Access Control** - Different permissions untuk setiap user type
- **Password Encryption** - Secure password hashing dengan bcrypt

### Management Features

- **Event Scheduling** - Manajemen tanggal dan waktu event
- **Order Status Tracking** - Pending, paid, cancelled statuses
- **Revenue Tracking** - Financial reporting dan price management
- **System Analytics** - Comprehensive reporting dashboard
- **Redis Caching** - Performance optimization dengan Redis

## Configuration

### Database Configuration

Aplikasi menggunakan MySQL sebagai database utama dan Redis untuk caching. Database schema akan otomatis dibuat dan di-migrate saat aplikasi start.

**Supported entities:**

- Users (user, admin)
- Events dengan detail lengkap
- Categories untuk tiket event
- Orders dan OrderDetails
- System reports
