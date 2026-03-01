# 🛠️ Setup Guide

This guide will help you set up **Neuron Core** on your local machine or a production server.

---

## 💻 Local Development

### Prerequisites
- **Go**: Version 1.22 or higher
- **Docker**: Desktop or Engine (Must be running)

### 1. Installation

```bash
# Clone the repo
git clone https://github.com/anurag-327/neuron-core.git
cd neuron-core

# Install Go modules
go mod tidy
```

### 2. Setup Environment
Run the setup script. This will build all necessary Docker images and compile the binary.

```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### 3. Run Server

**Option A: Standard Run**
```bash
go run cmd/api/main.go
```

**Option B: Hot Reload (Development)**
Neuron includes an `.air.toml` configuration for live reloading.
```bash
# Install Air
go install github.com/air-verse/air@latest

# Run
air
```

You should see:
```
[info] Initializing sandbox container pools...
[info] Pre-warming container pool for cpp...
[info] Server starting {"port":"9000"}
```

---

## 🚀 Production Deployment

Recommend deploying Neuron Core on a dedicated worker node

### Recommended Specs
- **CPU**: 2+ vCPUs (Compute Optimized preferred)
- **RAM**: 4GB+ (Depends on pool size)
- **OS**: Ubuntu 22.04 / Debian 12

### Deployment Steps

1. **Install Docker & Go** on the server.
2. **Clone & Build** the binary:
   ```bash
   go build -o neuron-core cmd/api/main.go
   ```
3. **Build Images**: Run the docker build commands from above.
