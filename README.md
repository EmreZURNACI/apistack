
# 🎭 DVD Rental Service

This project is a **simple service** written in Go that performs CRUD operations on the `actor` table of the PostgreSQL **`dvdrental`** database. It has been designed with modern principles such as observability, structured logging, containerization, and graceful shutdown handling.





## 🚀 Technologies Used

- **Golang** – Backend development
- **PostgreSQL** – Relational database (`dvdrental`)
- **Docker & Docker Compose** – Containerization
- **OpenTelemetry** – Observability & tracing
- **Jaeger** – Distributed tracing UI
- **Prometheus** – Metrics collection
- **Grafana** – Metrics visualization dashboards
- **Postman** – API testing
## ✅ Functional Features

- Full CRUD support for the `actor` table:
  - ✔️ Create a new actor
  - 📖 Retrieve actor(s)
  - ✏️ Update actor details
  - ❌ Delete an actor
- Basic input validation (e.g., non-empty fields, string length checks)


## 🛠️ Non-Functional Features

- 📋 Structured logging (JSON logs)
- 🩺 Health check endpoint (`/health`)
- 📊 Prometheus metrics endpoint (`/metrics`)
- ☠️ Graceful shutdown with signal handling
- 🔍 Tracing with OpenTelemetry and Jaeger
- 📈 Grafana dashboards for real-time monitoring

## ⚙️ Installation & Running

### 1. Clone the Repository
```bash
git clone https://github.com/EmreZURNACI/apistack
```

### 2. Change Directory
```bash
cd apistack
```

### 3. Build Compose File
```bash
docker-compose up -d --build
```
## Screenshots

![Grafana](https://www.dropbox.com/scl/fi/t5zny9648905sori16mr6/Grafana.png?rlkey=6p8lxjnv97s7n35e83jcxc0kw&st=4p9sdoru&raw=1)

![Tracing-2](https://www.dropbox.com/scl/fi/0d42ctcvtom7hctuxtfeh/Tracing-2.png?rlkey=i4q4ss0ep8no7ojnyj47bvfcb&st=u7fbpouw&raw=1)

![Prometheus](https://www.dropbox.com/scl/fi/c7iseigegnncm8q2g8br9/Prometheus.png?rlkey=uzdtnvsokpdpnci215i8hqqjc&st=j8klh52s&raw=1)

![Tracing](https://www.dropbox.com/scl/fi/k4hplwhwxl001dyegbfh7/Tracing.png?rlkey=76lqon154vvwfsixgfj8mcqgb&st=ft12p4wg&raw=1)








## API Reference


| Method | Endpoint        | Description          |
| ------ | --------------- | -------------------- |
| GET    | /v1/actors      | Get all actors       |
| POST   | /v1/actors      | Create a new actor   |
| GET    | /v1/actor/{id}  | Get a specific actor |
| PUT    | /v1/actor/{id}  | Update actor details |
| DELETE | /v1/actor/{id}  | Delete an actor      |

## 📊 Monitoring & Tracing

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (Default login: admin / admin)
- **Jaeger UI**: http://localhost:16686


### 🔧 Variables

- **HOST**=postgres
- **PORT**=5432
- **USER**=postgres
- **DB**=dvdrental
- **PASSWORD**=123
- **SERVER_PORT**=:8080

### ⚠️ Limitations / Known Issues

- 🔐 **No authentication or authorization is implemented.**
- 🔁 **Only the `actor` table is implemented**; other entities in the `dvdrental` database are not yet supported.

### ⚠️ Note

This is currently a **single service** and **not a microservice architecture**. However, the codebase follows practices that allow for future scalability into microservices.
