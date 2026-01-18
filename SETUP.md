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

### 2. Prepare Images
You need to build the custom runtime images once.

```bash
# Easy way: Run this one-liner
docker build -t neuron-cpp runtime/images/cpp && \
docker build -t neuron-python runtime/images/python && \
docker build -t neuron-node runtime/images/node && \
docker build -t neuron-java runtime/images/java
```

### 3. Run Server

```bash
# Create .env (Optional, defaults to 9000)
echo "PORT=9000" > .env

# Start the engine
go run cmd/api/main.go
```

You should see:
```
[info] Initializing sandbox container pools...
[info] Pre-warming container pool for cpp...
[info] Server starting {"port":"9000"}
```

---

## 🚀 Production Deployment

We recommend deploying Neuron Core on a dedicated worker node (e.g., AWS EC2, DigitalOcean Droplet).

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
