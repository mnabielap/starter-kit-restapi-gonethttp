# 🚀 Starter Kit REST API - Go (net/http)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html#go1.22)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Architecture](https://img.shields.io/badge/Architecture-Standard_Go_Layout-purple)](https://github.com/golang-standards/project-layout)

A production-ready, feature-rich RESTful API starter kit built with **Go (Golang)** using the standard `net/http` library.

This project is designed as a high-performance alternative to Express.js boilerplates, featuring a **Standard Go Project Layout**, **GORM** for database interactions (compatible with SQLite and PostgreSQL), and fully automated **Python-based API testing**.

---

## ✨ Features

- **🏗 Standard Go Layout**: Clean separation of concerns (`cmd`, `internal`, `pkg`).
- **💾 Dual Database Support**: Seamlessly switch between **SQLite** (Pure Go, CGO-free) and **PostgreSQL**.
- **🔐 Authentication**: Robust JWT implementation (Access & Refresh Tokens).
- **👮 Authorization (RBAC)**: Role-Based Access Control ensuring only Admins can manage users.
- **🛡 Security**: Password hashing (Bcrypt) and API Rate Limiting.
- **📝 Logging**: Structured logging using Go's `log/slog`.
- **🐳 Docker Ready**: Multi-stage builds with Alpine Linux for tiny images.
- **📧 Email Service**: Built-in SMTP support for verification and password resets.
- **⚡ Pagination & Filtering**: Built-in utilities for data queries.
- **🧪 Automated Testing**: No Postman needed! Includes a suite of Python scripts for API testing.

---

## 📂 Project Structure

```text
starter-kit-restapi-gonethttp/
├── cmd/api/               # Application entry point
├── config/                # Environment and DB configuration
├── internal/
│   ├── handlers/          # HTTP Controllers (Transport Layer)
│   ├── middleware/        # Auth, Logger, Rate Limiter, Roles
│   ├── models/            # GORM Structs (Database Models)
│   ├── repository/        # Data Access Layer (SQL Queries)
│   ├── routes/            # Route definitions
│   └── services/          # Business Logic Layer
├── pkg/                   # Public Utilities (Logger, Response, Validator)
├── migrations/            # SQL Migrations (if needed manually)
├── api_tests/             # Python scripts for API Testing
├── .env                   # Environment variables
├── Dockerfile             # Docker configuration
├── entrypoint.sh          # Docker entry script
└── README.md              # Documentation
```

---

## ⚙️ Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

**Key Variables (`.env`):**
```properties
PORT=8080
GO_ENV=development

# Choose: 'sqlite' or 'postgres'
DB_DRIVER=sqlite
# For SQLite, DB_NAME is the filename (e.g., starter_kit_db)
DB_NAME=starter_kit_db

# For Postgres (Fill these if DB_DRIVER=postgres)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password

# JWT Secrets
JWT_SECRET=change_this_to_something_secure
```

---

## 🏃‍♂️ Getting Started (Local Development)

We recommend running the project locally first to understand how it works.

### Prerequisites
- **Go** (version 1.22 or higher)
- **Git**
- **Python 3.x** (for testing scripts)

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Run the Application
```bash
go run cmd/api/main.go
```
*The server will start at `http://localhost:8080`.*
*If using SQLite, a `.db` file will be created automatically.*

---

## 🐳 Docker Deployment

This project includes a production-ready `Dockerfile`. We use individual Docker commands (without Compose) to simulate a real-world orchestrated environment.

**Note:** Commands below use `\` for line breaks. If you are using Windows CMD, replace `\` with `^`.

### 1. Create Network 🌐
Create a dedicated network so the App and Database can communicate.
```bash
docker network create restapi_gonethttp_network
```

### 2. Create Volumes 📦
Create persistent storage for the Database and Media files.
```bash
# Volume for Postgres Data
docker volume create restapi_gonethttp_postgres_data

# Volume for Application Media/Uploads
docker volume create restapi_gonethttp_media_volume

# Volume for Application DB (if using SQLite)
docker volume create restapi_gonethttp_db_volume
```

### 3. Start Database (PostgreSQL) 🐘
Run the Postgres container attached to our network and volume.
```bash
docker run -d \
  --name restapi-gonethttp-postgres \
  --network restapi_gonethttp_network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=starter_kit_db \
  -v restapi_gonethttp_postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine
```

### 4. Setup Environment for Docker 📝
Create a `.env.docker` file. **Crucial:** `DB_HOST` must match the Postgres container name (`restapi-gonethttp-postgres`).

```properties
PORT=5005
GO_ENV=production
DB_DRIVER=postgres
DB_HOST=restapi-gonethttp-postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=mysecretpassword
DB_NAME=starter_kit_db
```

### 5. Build the App Image 🏗
```bash
docker build -t restapi-gonethttp-app .
```

### 6. Run the App Container 🚀
Run the application container. We inject the environment file and mount volumes.

```bash
docker run -d -p 5005:5005 \
  --env-file .env.docker \
  --network restapi_gonethttp_network \
  -v restapi_gonethttp_db_volume:/usr/src/app/db \
  -v restapi_gonethttp_media_volume:/usr/src/app/media \
  --name restapi-gonethttp-container \
  restapi-gonethttp-app
```

The API is now accessible at `http://localhost:5005`.

---

## 🛠 Docker Management Commands

Here are useful commands to manage your containers:

#### 📜 View Logs
See what's happening inside the container in real-time.
```bash
docker logs -f restapi-gonethttp-container
```

#### 🛑 Stop Container
```bash
docker stop restapi-gonethttp-container
```

#### ▶️ Start Container
Start the container again (data persists).
```bash
docker start restapi-gonethttp-container
```

#### 🗑 Remove Container
Removes the container instance (requires stopping first).
```bash
docker rm restapi-gonethttp-container
```

#### 📦 List Volumes
```bash
docker volume ls
```

#### ⚠️ Remove Volume
**WARNING:** This deletes your database data permanently!
```bash
docker volume rm restapi_gonethttp_postgres_data
```

---

## 🧪 API Testing (Automated)

Forget Postman! This project comes with a suite of **Python scripts** to test every endpoint. These scripts automatically handle token storage (`secrets.json`) and chaining requests.

### Setup
1. Ensure Python 3 is installed.
2. Navigate to the `api_tests` directory (or wherever you placed the scripts).
3. If running via Docker, edit `api_tests/utils.py` and set `BASE_URL = "http://localhost:5005/v1"`.

### How to Run
Simply run the script file. No arguments needed.

**1. Register a User (Saves token automatically):**
```bash
python api_tests/A1.auth_register.py
```

**2. Login (Refreshes token in secrets.json):**
```bash
python api_tests/A2.auth_login.py
```
*Note: Ensure you login as an **Admin** to perform User Management tests (B series).*

**3. Get Users List (Pagination):**
```bash
python api_tests/B2.user_get_list.py
```

---

## 📝 License

This project is licensed under the MIT License.