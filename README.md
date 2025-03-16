# Modern Authentication System in Go

A robust authentication system built with Go, demonstrating secure user management and modern web practices from a backend engineer's perspective.

## 🚀 Features

- **Secure Authentication**
  - JWT-based session management
  - Password hashing and validation
  - Protected routes
  - Email verification

- **Real-time Validation**
  - Username availability checking
  - Email format validation
  - Password strength meter
  - Instant feedback using Alpine.js

- **Clean Architecture**
  - Modular Go backend design
  - Separation of concerns
  - Dependency injection
  - Configuration management

## 🛠️ Technology Stack

- **Backend**
  - Go 1.24.1
  - Chi Router (REST API)
  - PostgreSQL (Database)
  - JWT for authentication
  - Goose (Database migrations)
  - SQLC (Type-safe SQL)

- **Frontend**
  - Alpine.js (Lightweight reactivity)
  - Tailwind CSS (Styling)
  - Go html/template (Server-side rendering)

## 📋 Prerequisites

- Go 1.24.1 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## 🚀 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/authentication.git
   cd authentication
   ```

2. **Set up the configuration**
   ```bash
   cp config/config.example.yaml config/config.yaml
   # Edit config/config.yaml with your settings
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Set up the database**
   ```bash
   # Create PostgreSQL database
   createdb auth-db

   # Run migrations
   make migrate-up
   ```

5. **Run the application**
   ```bash
   make run
   # Or without Make:
   go run cmd/server/main.go
   ```

6. **Access the application**
   ```
   http://localhost:8080
   ```

## 📁 Project Structure

```
.
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── config/         # Configuration management
│   ├── db/             # Database operations and migrations
│   ├── middleware/     # HTTP middleware
│   └── service/        # Business logic
├── web/
│   └── templates/      # HTML templates
└── config/             # Configuration files
```

## 🔒 Security Features

- Secure password hashing using bcrypt
- JWT token-based authentication
- Protected routes with middleware
- SQL injection prevention with prepared statements
- XSS protection
- CSRF protection
- Rate limiting

## 🛠️ Development Commands

```bash
# Run the application
make run

# Run migrations
make migrate-up
make migrate-down

# Generate SQL code (using SQLC)
make sqlc

# Run tests
make test

# Build the application
make build
```

## 📝 API Endpoints

### Authentication
- `POST /api/register` - User registration
- `POST /api/login` - User login
- `GET /api/check-username` - Check username availability
- `GET /api/check-email` - Validate email
- `POST /api/check-password` - Check password strength

### Web Routes
- `GET /` - Home page
- `GET /login` - Login page
- `GET /register` - Registration page
- `GET /dashboard` - Protected dashboard

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

Dominic Yeboah
- LinkedIn: https://www.linkedin.com/in/dominic-yeboah/
- GitHub: https://github.com/yeboahd24
