# Parkin AI System - Backend

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![GoFrame](https://img.shields.io/badge/GoFrame-v2-blue?style=flat)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)
![Jenkins](https://img.shields.io/badge/Jenkins-CI/CD-D33833?style=flat&logo=jenkins)

A comprehensive parking management system backend built with GoFrame framework, featuring AI-powered parking solutions, real-time management, and integrated payment processing.

## 🚀 Features

### Core Features
- **🅿️ Parking Management**: Real-time parking slot monitoring and reservation
- **💳 Payment Integration**: PayOS payment gateway integration with webhook support
- **👤 User Management**: JWT-based authentication and authorization
- **📱 Mobile API**: RESTful API for mobile applications
- **🔔 Notification System**: Real-time notifications for users
- **⭐ Favorites**: User favorite parking lots management
- **📊 Dashboard**: Admin dashboard with analytics

### Technical Features
- **🏗️ Clean Architecture**: Domain-driven design with clear separation of concerns
- **🔒 Security**: JWT authentication, middleware protection, input validation
- **📝 API Documentation**: Auto-generated Swagger documentation
- **🐳 Containerized**: Docker support for easy deployment
- **🚀 CI/CD**: Jenkins pipeline for automated testing and deployment
- **🗄️ Database**: PostgreSQL with GORM integration
- **📈 Monitoring**: Health checks and logging

## 🛠️ Tech Stack

- **Framework**: [GoFrame v2](https://goframe.org/)
- **Database**: PostgreSQL (Neon DB)
- **Payment**: PayOS Gateway
- **Containerization**: Docker
- **CI/CD**: Jenkins
- **Cloud**: AWS EC2
- **Documentation**: Swagger/OpenAPI

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 13+
- Docker (optional)
- Git

## 🚀 Quick Start

### 1. Clone the repository
```bash
git clone https://github.com/Parking-AI-Systems/parkin_ai_systems_be.git
cd parkin_ai_systems_be
```

### 2. Setup environment
```bash
# Copy config template
cp manifest/config/config.yaml.example manifest/config/config.yaml

# Edit configuration
nano manifest/config/config.yaml
```

### 3. Install dependencies
```bash
go mod tidy
go mod download
```

### 4. Database setup
```bash
# Make sure PostgreSQL is running and create database
createdb parkin_ai_system

# Run migrations (if available)
# go run main.go migrate
```

### 5. Run the application
```bash
# Development mode
go run main.go

# Production build
go build -o parkin-ai-system main.go
./parkin-ai-system
```

The server will start on `http://localhost:8000`

## 📖 API Documentation

Once the server is running, you can access:

- **Swagger UI**: `http://localhost:8000/backend/parkin/v1/swagger`
- **OpenAPI JSON**: `http://localhost:8000/backend/parkin/v1/api.json`

## 🏗️ Project Structure

```
parkin-ai-system/
├── api/                    # API definitions and request/response models
│   ├── payment/           # Payment API definitions
│   ├── parking_order/     # Parking orders API
│   └── ...
├── internal/              # Private application code
│   ├── cmd/              # Command line interface
│   ├── controller/       # HTTP handlers
│   ├── logic/           # Business logic layer
│   ├── service/         # Service interfaces
│   ├── dao/             # Data access objects
│   ├── model/           # Data models and entities
│   └── config/          # Configuration management
├── manifest/             # Deployment and configuration files
│   └── config/          # Configuration files
├── Dockerfile           # Docker configuration
├── Jenkinsfile         # Jenkins CI/CD pipeline
├── go.mod              # Go module dependencies
└── main.go             # Application entry point
```

## 🔧 Configuration

Edit `manifest/config/config.yaml`:

```yaml
server:
  address: ":8000"
  openapiPath: "/backend/parkin/v1/api.json"
  swaggerPath: "/backend/parkin/v1/swagger"

database:
  default:
    link: "pgsql:user:password@tcp(localhost:5432)/dbname"

auth:
  secretKey: "your-secret-key"
  refreshTokenExpireMinute: 1440

payos:
  clientID: "your-payos-client-id"
  apiKey: "your-payos-api-key"
  checkSum: "your-payos-checksum"
```

## 🐳 Docker Deployment

### Build and run with Docker
```bash
# Build image
docker build -t parkin-ai-system .

# Run container
docker run -d \
  --name parkin-ai-system \
  -p 8000:8000 \
  -e GF_DATABASE_LINK="your-db-connection" \
  parkin-ai-system
```

### Docker Compose (recommended)
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8000:8000"
    environment:
      - GF_DATABASE_LINK=pgsql://user:pass@db:5432/parkin_ai_system
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: parkin_ai_system
      POSTGRES_USER: your_user
      POSTGRES_PASSWORD: your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 🚀 CI/CD Pipeline

The project includes a Jenkins pipeline (`Jenkinsfile`) that:

1. **Tests** the application
2. **Builds** the Go binary
3. **Creates** Docker image
4. **Pushes** to Docker Hub
5. **Deploys** to AWS EC2

### Pipeline Stages:
- ✅ Clean Workspace
- ✅ Checkout from SCM
- ✅ Test Application
- ✅ Build Application
- ✅ Build & Push Docker Image
- ✅ Deploy to AWS EC2

## 📊 API Endpoints

### Authentication
- `POST /backend/parkin/v1/auth/login` - User login
- `POST /backend/parkin/v1/auth/register` - User registration
- `POST /backend/parkin/v1/auth/refresh` - Refresh token

### Parking Management
- `GET /backend/parkin/v1/parking-lots` - List parking lots
- `GET /backend/parkin/v1/parking-slots` - List parking slots
- `POST /backend/parkin/v1/parking-orders` - Create parking order

### Payment
- `POST /backend/parkin/v1/create-payment-link` - Create payment link
- `POST /backend/parkin/v1/payment-requests` - Create payment request
- `POST /backend/parkin/v1/webhook` - Payment webhook

### Health Check
- `GET /backend/parkin/v1/health` - Application health status

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/logic/payment/...
```

## 📝 Development

### Adding new features
1. Define API in `api/` directory
2. Implement controller in `internal/controller/`
3. Add business logic in `internal/logic/`
4. Create service interface in `internal/service/`
5. Update database models in `internal/model/`

### Code Generation
```bash
# Generate API documentation
gf gen swagger

# Generate database models
gf gen dao
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Backend Development**: GoFrame + PostgreSQL
- **Payment Integration**: PayOS Gateway
- **DevOps**: Docker + Jenkins + AWS
- **API Documentation**: Swagger/OpenAPI

## 📞 Support

If you have any questions or need help, please:

1. Check the [API Documentation](http://localhost:8000/backend/parkin/v1/swagger)
2. Open an [Issue](https://github.com/Parking-AI-Systems/parkin_ai_systems_be/issues)
3. Contact the development team

---

**Made with ❤️ by Parking AI Systems Team**
