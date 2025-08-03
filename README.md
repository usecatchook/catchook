# 🎣 Catchook

![Catchook Logo](./app/public/og-image.png)

**The most powerful webhook platform for developers**  
Capture, debug, transform, and route webhooks with zero configuration.

> [!WARNING]  
> Catchook is under active development. APIs may change. Not recommended for production use yet.


[![Go Version](https://img.shields.io/badge/Go-1.24.4-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-15-black?style=flat&logo=next.js)](https://nextjs.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

---

## ✨ Why Catchook?

**Stop fighting with webhooks.** Catchook gives you superpowers:

- 🚀 **Zero Config** - Start capturing webhooks instantly
- 🔍 **Smart Debugging** - Real-time inspection with advanced filtering
- 🔄 **Intelligent Routing** - Route webhooks based on content, headers, or custom rules
- 📊 **Live Monitoring** - Beautiful dashboard with metrics and alerting
- 🛠️ **Transform & Replay** - Modify payloads and replay events
- ⚡ **High Performance** - Built with Go + Fiber for maximum throughput

## 🚀 Quick Start

Get Catchook running locally in under 2 minutes:

### Prerequisites

- [Go 1.24+](https://golang.org/doc/install)
- [Node.js 18+](https://nodejs.org/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)

### 1. Clone & Setup

```bash
git clone https://github.com/theotruvelot/catchook.git
cd catchook

# Start PostgreSQL & Redis
docker-compose -f docker-compose.dev.yml up -d

# Install dependencies
go mod tidy
cd app && npm install && cd ..
```

### 2. Environment Setup

```bash
# Copy example environment
cp .env.example .env

# The default config works with docker-compose setup!
# Edit .env if you need custom database credentials
```

### 3. Start Development

```bash
# Terminal 1: Start the API
make dev-api

# Terminal 2: Start the Frontend
make dev-app

# 🎉 Open http://localhost:3000
```

Your first webhook endpoint is ready at `http://localhost:8080/hooks/your-unique-id`

## 🏗️ Architecture

Catchook is built for **performance** and **developer experience**:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Next.js App   │───▶│    Go Fiber API  │───▶│   PostgreSQL    │
│   (Frontend)    │    │   (Backend)      │    │   (Storage)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │     Redis       │
                       │(Cache & Pub/Sub)│
                       └─────────────────┘
```

### 🛠️ Tech Stack

| Component | Technology | Why? |
|-----------|------------|------|
| **Backend** | Go + Fiber | Blazing fast HTTP performance |
| **Database** | PostgreSQL + SQLC | Type-safe SQL with zero ORM overhead |
| **Cache** | Redis | Real-time features & smart caching |
| **Frontend** | Next.js + TypeScript | Modern React with full-stack capabilities |
| **Styling** | Tailwind CSS + Shadcn UI | Rapid UI development |
| **Auth** | JWT | Stateless auth with performance |

## 📁 Project Structure

```
catchook/
├── api/                    # API-related configs
├── app/                    # Next.js frontend application
├── cmd/api/               # API entry point
├── internal/              # Private Go packages
│   ├── app/               # HTTP handlers & dependency injection
│   ├── config/            # Configuration management
│   ├── domain/            # Business logic & interfaces
│   ├── middleware/        # HTTP middleware (auth, logging, etc.)
│   ├── repository/        # Data access layer (SQLC)
│   └── service/           # Business orchestration
├── pkg/                   # Public Go packages
│   ├── cache/             # Redis abstraction
│   ├── jwt/               # JWT token management
│   ├── logger/            # Structured logging
│   └── validator/         # Request validation
├── storage/postgres/      # Database schemas & queries
└── bruno/                 # API testing collection
```

## 🤝 Contributing

We ❤️ contributions! Catchook is designed to be **contributor-friendly**.

### 🎯 Good First Issues

Look for issues labeled [`good first issue`](https://github.com/theotruvelot/catchook/labels/good%20first%20issue):

- 🐛 Bug fixes
- 📝 Documentation improvements
- ✨ Small feature additions
- 🧪 Test coverage improvements

### 🔧 Development Workflow

1. **Fork & Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/catchook.git
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **Make Changes & Test**
   ```bash
   make test
   make lint
   ```

4. **Submit PR**
   - Write clear commit messages
   - Add tests for new features
   - Update documentation

### 🧪 Testing (WIP)

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Test specific package
go test ./internal/service/...
```

## 📖 API Documentation (WIP)

Explore the API with our [Bruno collection](./bruno/).

### Manual Deployment

```bash
# Build API
make build-api

# Build Frontend
cd app && npm run build

# Run migrations
make migrate-up
```

## 🗺️ Roadmap

**🎯 v1.0 Goals:**
- [ ] Complete webhook capture & replay
- [ ] Advanced filtering & transformation
- [ ] Webhook routing rules
- [ ] Real-time dashboard
- [ ] API rate limiting


## 📄 License

Catchook is [MIT licensed](./LICENSE).

## 🙏 Acknowledgments

Built with amazing open-source tools:
- [Fiber](https://github.com/gofiber/fiber) - Express-inspired Go web framework
- [SQLC](https://github.com/sqlc-dev/sqlc) - Type-safe SQL in Go
- [Next.js](https://nextjs.org/) - React production framework
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS
- [Shadcn UI](https://ui.shadcn.com/) - React UI library

A big thank to [OpenSourceTogether](https://opensource-together.com/) for the amazing support and resources.

---

<div align="center">

**[⭐ Star this repo](https://github.com/theotruvelot/catchook)** | **[🐛 Report Bug](https://github.com/theotruvelot/catchook/issues)** | **[💡 Request Feature](https://github.com/theotruvelot/catchook/issues)**

*Made with ❤️ by [@theotruvelot](https://github.com/theotruvelot) and [contributors](https://github.com/theotruvelot/catchook/contributors)*

</div>


