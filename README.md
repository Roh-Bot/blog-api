# Blog API – Clean Architecture in Go

A **production-ready Blog API** built using **Golang** and **Clean Architecture principles**, focusing on scalability, observability, and maintainability.  
The project demonstrates real-world backend patterns such as async structured logging, Redis caching, hot-reloadable configuration, and database migrations.

---

## ✨ Features

- **Clean Architecture**
  - Clear separation of concerns
  - Domain-driven design
  - Framework-agnostic core

- **High-performance API Framework**
  - **Fiber** (main branch)
  - **Echo** (migration-echo branch)

- **Asynchronous Structured Logging**
  - Powered by **Uber Zap**
  - Non-blocking logging pipeline
  - JSON logs ready for aggregation systems

- **PostgreSQL Integration**
  - Uses **pgx** for high-performance DB access
  - Explicit control over connection pooling

- **Database Migrations**
  - Managed using **Goose**
  - Versioned, repeatable migrations

- **Redis Caching**
  - Optimized read performance
  - Cache-aside strategy

- **Configuration Management**
  - **Koanf v2**
  - Supports multiple config sources (file/env)
  - **Hot reload using atomic pointer swaps** (lock-free reads)

- **Validation**
  - **Validator v10** for request validation
  - Centralized validation rules

- **API Documentation**
  - **Swagger / OpenAPI**
  - Auto-generated API specs

- **DevOps Ready**
  - `Makefile` for common workflows
  - `Dockerfile` for containerized deployments

---

## 🌿 Branch Strategy

| Branch | Description |
|------|------------|
| `main` | API implemented using **Fiber** |
| `migration-echo` | Same API implemented using **Echo** |

> The domain and usecase layers remain unchanged across branches, showcasing framework portability.

---

## 🚀 Getting Started

### Prerequisites

- Go **1.24+**
- PostgreSQL
- Redis
- Docker (optional)

---

### Clone the Repository

```bash
git clone https://github.com/your-username/blog-api.git
cd blog-api
```

---

⚙️ **Configuration**

Configuration is managed using Koanf v2 with support for:

YAML / JSON files

Environment variables

Runtime hot reloads

Configuration updates are applied using atomic pointer swaps, ensuring:

Lock-free reads

Zero downtime reloads

---

🗄️ **Database Migrations**

Run migrations using Goose:

```bash
make migrate-up
```

Rollback:

```bash
make migrate-down
```

---

▶️ **Run the Application**

Local

```bash
make run
```

Docker

```bash
docker build -t blog-api .
docker run -p 8080:8080 blog-api
```
---

📘 Swagger API Docs

Once the server is running:

```bash
http://localhost:8080/swagger/index.html
```
---

🧪 **Validation**

Request payloads validated using validator v10

Validation errors are standardized and user-friendly

---

🧵 **Logging**

Uber Zap with async log writes

Structured JSON logs

Correlation-friendly fields (request ID, method, latency)

---

🧠 **Architectural Goals**

Framework independence

Explicit dependencies

Predictable behavior under load

Observability-first design

Minimal runtime overhead

---

📌 **Future Improvements**

Distributed tracing

Rate limiting

Background workers

Metrics via Prometheus
