TodoList Project

This project is a microservices-based application built with Go, Gin, and MongoDB, demonstrating a simple architecture for user management, task management, and auditing, all orchestrated with Docker Compose and Kubernetes. It includes four main services:

- **UserService**: Handles user registration and authentication (JWT).
- **TaskService**: Manages user tasks.
- **AuditService**: Logs/audits requests (placeholder for future expansion).
- **APIGateway**: Routes external requests to the appropriate service.

---

## Table of Contents

- Architecture Overview
- Project Structure
- Services
  - UserService
  - TaskService
  - AuditService
  - APIGateway
- API Endpoints
  - UserService Endpoints
  - TaskService Endpoints
  - AuditService Endpoints
- Running Locally with Docker Compose
- Kubernetes Deployment
- Environment Variables
- Building and Running Services Individually
- Extending the Project
- Troubleshooting
- License

---

## Architecture Overview

```
[Client] <---> [APIGateway] <---> [UserService]
                               \--> [TaskService]
                               \--> [AuditService]
                               \--> [MongoDB]
```

- **APIGateway** exposes a single entrypoint for clients and proxies requests to the appropriate backend service.
- **UserService** and **TaskService** use MongoDB for persistence.
- **AuditService** is designed to log requests for auditing purposes.
- All services are containerized and can be orchestrated with Docker Compose or Kubernetes.

---

## Project Structure

```
.
├── APIGateway/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── AuditService/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── TaskService/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── UserService/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── README
├── k8s/
│   ├── audit-deployment.yaml
│   ├── audit-service.yaml
│   ├── gateway-deployment.yaml
│   ├── gateway-service.yaml
│   ├── ingress.yaml
│   ├── task-deployment.yaml
│   ├── task-service.yaml
│   ├── user-deployment.yaml
│   ├── user-mongo.yaml
│   ├── user-service.yaml
│   └── (subfolders for user/task/gateway/audit)
├── docker-compose.yml
└── README.md
```

---

## Services

### UserService

- **Port**: 8081
- **Responsibilities**: User registration and authentication (JWT).
- **Database**: MongoDB (`userservice.users` collection).
- **Endpoints**:
  - `POST /users/signup`
  - `POST /users/login`

### TaskService

- **Port**: 8082
- **Responsibilities**: Task creation and retrieval for users.
- **Database**: MongoDB (`taskservice.tasks` collection).
- **Endpoints**:
  - `POST /tasks/`
  - `GET /tasks/?username=...`

### AuditService

- **Port**: 8083
- **Responsibilities**: Logs incoming requests (currently a placeholder).
- **Endpoints**:
  - `GET /audit/health`

### APIGateway

- **Port**: 8080
- **Responsibilities**: Reverse proxy for all services.
- **Routes**:
  - `/users/*` → UserService
  - `/tasks/*` → TaskService
  - `/audit/*` → AuditService

---

## API Endpoints

### UserService Endpoints

#### `POST /users/signup`

Register a new user.

**Request Body:**
```json
{
  "username": "your_username",
  "password": "your_password"
}
```

**Response:**
- `201 Created` on success
- `400 Bad Request` on validation error

#### `POST /users/login`

Authenticate a user and receive a JWT.

**Request Body:**
```json
{
  "username": "your_username",
  "password": "your_password"
}
```

**Response:**
- `200 OK` with `{ "token": "<jwt>" }` on success
- `401 Unauthorized` on invalid credentials

---

### TaskService Endpoints

#### `POST /tasks/`

Create a new task.

**Request Body:**
```json
{
  "title": "Task title",
  "status": "pending",
  "username": "user1"
}
```

**Response:**
- `201 Created` with the created task object

#### `GET /tasks/?username=...`

List all tasks for a user.

**Response:**
- `200 OK` with an array of tasks

---

### AuditService Endpoints

#### `GET /audit/health`

Health check endpoint.

**Response:**
- `200 OK` with `{ "status": "Audit Service is running" }`

---

## Running Locally with Docker Compose

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

### Steps

1. **Clone the repository:**
   ```sh
   git clone <repository-url>
   cd <repository-folder>
   ```

2. **Start all services:**
   ```sh
   docker-compose up --build
   ```

3. **Access the services:**
   - APIGateway: http://localhost:8080
   - UserService: http://localhost:8081
   - TaskService: http://localhost:8082
   - AuditService: http://localhost:8083
   - MongoDB: localhost:27017

4. **Example requests:**
   - Register user:
     ```sh
     curl -X POST http://localhost:8080/users/signup -H "Content-Type: application/json" -d '{"username":"user1","password":"pass"}'
     ```
   - Login:
     ```sh
     curl -X POST http://localhost:8080/users/login -H "Content-Type: application/json" -d '{"username":"user1","password":"pass"}'
     ```
   - Create task:
     ```sh
     curl -X POST http://localhost:8080/tasks/ -H "Content-Type: application/json" -d '{"title":"Test Task","status":"pending","username":"user1"}'
     ```
   - List tasks:
     ```sh
     curl "http://localhost:8080/tasks/?username=user1"
     ```

---

## Kubernetes Deployment

### Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- A running Kubernetes cluster (e.g., [minikube](https://minikube.sigs.k8s.io/docs/))

### Steps

1. **Apply MongoDB deployment and service:**
   ```sh
   kubectl apply -f k8s/user-mongo.yaml
   ```

2. **Deploy UserService, TaskService, AuditService, and APIGateway:**
   ```sh
   kubectl apply -f k8s/user-deployment.yaml
   kubectl apply -f k8s/user-service.yaml
   kubectl apply -f k8s/task-deployment.yaml
   kubectl apply -f k8s/task-service.yaml
   kubectl apply -f k8s/audit-deployment.yaml
   kubectl apply -f k8s/audit-service.yaml
   kubectl apply -f k8s/gateway-deployment.yaml
   kubectl apply -f k8s/gateway-service.yaml
   ```

3. **(Optional) Deploy Ingress:**
   ```sh
   kubectl apply -f k8s/ingress.yaml
   ```

4. **Access the API Gateway:**
   - Use `kubectl port-forward` or your ingress controller's external IP.

---

## Environment Variables

- `MONGO_URI`: MongoDB connection string (e.g., `mongodb://mongo:27017` for Docker Compose, or `mongodb://user-mongo:27017` for Kubernetes).

---

## Building and Running Services Individually

Each service can be built and run separately:

```sh
cd UserService
go build -o main .
MONGO_URI=mongodb://localhost:27017 ./main
```

Replace UserService with TaskService, AuditService, or APIGateway as needed.

---

## Extending the Project

- **Add more endpoints**: Expand UserService and TaskService for more features.
- **Implement AuditService**: Store logs in a database or external system.
- **Add authentication/authorization**: Use JWT in TaskService for protected endpoints.
- **Monitoring & Logging**: Integrate with Prometheus, Grafana, or ELK stack.
- **CI/CD**: Add GitHub Actions or other CI/CD pipelines.

---

## Troubleshooting

- **Ports already in use**: Make sure no other services are running on 8080-8083 or 27017.
- **MongoDB connection errors**: Check `MONGO_URI` and ensure MongoDB is running.
- **Kubernetes issues**: Use `kubectl get pods` and `kubectl logs` for debugging.

---

## License

This project is licensed under the MIT License.

---

**Contributions are welcome!** Please open issues or submit pull requests for improvements or bug fixes.