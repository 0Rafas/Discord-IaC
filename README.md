<div align="center">


# Discord IaC (Infrastructure-as-Code) CLI 🚀
*An Enterprise-Grade Go CLI for Generating High-Performance, Scalable Discord Bot Architectures.*

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io/)
[![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)](https://prometheus.io/)

</div>

## 🌟 Overview
**Discord IaC** is a powerful Command Line Interface designed to save senior engineers weeks of setup time. By writing a simple `bot.yaml` file, this tool instantly generates a production-ready **Clean Architecture** codebase in Go. 

It drops the traditional, slow direct-handling method and implements an advanced **Concurrent Worker Pool Router**, Prometheus Metrics, Redis caching, PostgreSQL databases, Dockerfiles, and Kubernetes Deployments out-of-the-box.

## ⚡ Features
- **Concurrent Event Routing:** Handles millions of Discord events seamlessly using Go Worker Pools & Channels, protecting your bot from gateway blocks and memory spikes.
- **Enterprise Observability:** Generates a built-in Prometheus metrics server (`:2112/metrics`) ready for Grafana dashboards.
- **Clean Architecture Principles:** Enforces rigorous separation of concerns (`cmd/`, `internal/handlers/`, etc.).
- **Deploy Anywhere:** Automatically generates `Dockerfile`, `docker-compose.yml`, GitHub Actions (`.github/workflows/ci.yml`), and Kubernetes manifests (`k8s/deployment.yaml`).
- **Data Layers:** Auto-scaffolds Redis and PostgreSQL integration endpoints.

---

## 🚀 Quick Start
### 1. Installation
Install the tool via Go (or download the binary from the releases page):
```bash
go install github.com/0Rafas/discord-iac@latest
```

### 2. Initialization
Run the init command inside an empty folder to generate the base configuration file:
```bash
discord-iac init
```
This generates `bot.yaml`:
```yaml
name: "TitanBot"
prefix: "!"
intents: ["Guilds", "GuildMessages", "MessageContent"]
database: "postgres"
cache: "redis"

# Enterprise Features
observability:
  prometheus: true
  tracing: true

concurrency:
  worker_pool_size: 100

deployment:
  kubernetes: true
  github_actions: true
```

### 3. Generation
Once your `bot.yaml` is configured, run:
```bash
discord-iac generate
```
The CLI will parse your configuration and generate a heavily-optimized, production-ready Go codebase inside the `out/` folder.

### 4. Build and Run
```bash
cd out
go mod tidy
DISCORD_TOKEN=your_token_here go run ./cmd/bot
```

---

## 🛠️ Architecture Generated
The tool generates the following structural pattern:
```text
out/
├── .github/workflows/ci.yml     # Automated testing
├── cmd/bot/main.go              # Bot Entrypoint
├── internal/
│   ├── bot/client.go            # Discord session & intents
│   ├── cache/redis.go           # Redis integration
│   ├── database/db.go           # Postgres integration
│   ├── metrics/metrics.go       # Prometheus server
│   └── handlers/handlers.go     # Concurrent Worker Pool Router
├── k8s/deployment.yaml          # Kubernetes container specs
├── Dockerfile          
├── docker-compose.yml 
└── go.mod
```

## 🤝 Contributing
Contributions are always welcome! Feel free to open a Pull Request to add more features like advanced Slash Command middleware generation, automated database migrations, or support for databases like MongoDB.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
