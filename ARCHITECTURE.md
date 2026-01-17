<h1 align="center">
  <br>
  Neuron Core
  <br>
</h1>

<h4 align="center">The high-performance execution engine powering the Neuron platform.</h4>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/docker-2496ED?style=for-the-badge&logo=docker" alt="Docker">
  <img src="https://img.shields.io/badge/architecture-clean-success?style=for-the-badge" alt="Architecture">
  <img src="https://img.shields.io/badge/status-production%20ready-success?style=for-the-badge" alt="Status">
</p>

---

## 🎯 What is Neuron Core?

Neuron Core is the dedicated **code execution microservice** responsible for running untrusted code in secure, isolated environments. It decouples the complex execution logic from the main application/backend, ensuring:

- **Isolation**: Heavy execution load doesn't impact the main API or database.
- **Security**: Dedicated sandbox environment with strict limits.
- **Scalability**: Can be deployed independently on worker nodes.
- **Speed**: Optimized container pooling and custom lightweight runtimes.

---

## ✨ Key Features

### ⚡ Blazing Fast Implementation
- **Container Pooling**: Pre-warmed containers eliminate cold start latency (< 5ms acquisition).
- **Custom Runtimes**: Minimal Alpine-based images (< 200MB) for C++, Python, Node.js, and Java.
- **Direct Execution**: Optimized via `Docker Exec` API with attached standard streams.

### 🛡️ Sandboxed & Secure
- **Network Isolation**: No internet access (`--network none`).
- **Read-Only RootFS**: Prevents system modifications.
- **Resource Quotas**: Strict CPU, Memory, and PIDs limits.
- **Secure Handling**: Input/Output sanitization to prevent path leakage.

### 📊 Rich Metrics
- **Time Breakdown**: Precise tracking of **Compilation** vs **Execution** time.
- **Status Reporting**: Granular classification of TLE (Time Limit Exceeded), MLE (Memory Limit Exceeded), and Runtime Errors.

---

## 🛠️ Technology Stack

- **Language**: Go (Golang)
- **Containerization**: Docker SDK
- **Architecture**: Modular (Transport -> Engine -> Runner)
- **Logging**: Structured JSON Logging

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22+
- Docker Engine (Running)

### Running Locally

```bash
# Clone the repository
git clone https://github.com/anurag-327/neuron-core.git
cd neuron-core

# Install dependencies
go mod download

# Create .env file
echo "PORT=9000" > .env

# Run the service
go run cmd/api/main.go
```

The service will start on port `9000` and automatically warm up container pools.

---

## 📡 API Reference

### `POST /api/v1/execute`

Execute a snippet of code.

**Request:**
```json
{
  "language": "cpp",
  "code": "#include <iostream>\nint main() { ... }",
  "input": "test input",
  "limit": {
    "memory_kb": 256000,
    "time_ms": 2000
  }
}
```

**Response:**
```json
{
  "stdout": "Hello World",
  "stderr": "",
  "exit_code": 0,
  "err_type": "",
  "err_msg": "",
  "metrics": {
    "total_ms": 150,
    "compile_ms": 120,
    "run_ms": 30
  },
  "container_dirty": false
}
```

---

## 🏗️ Architecture

```
User Request
    │
    ▼
[ HTTP Handler ] -> Validates Request
    │
    ▼
[ Execution Engine ] -> Orchestrates Flow
    │
    ▼
[ Docker Runner ]
    │
    ├── 1. Acquire Container (Pool)
    ├── 2. Prepare Workspace (Host Bind Mount)
    ├── 3. Compile Code (if needed)
    ├── 4. Execute Code (with timeout)
    └── 5. Process & Sanitize Output
    │
    ▼
[ Response ]
```

---

## 📝 Configuration

Configuration is managed via `runtime/runtime.go`. You can customize:
- **Images**: Docker images used for each language.
- **Commands**: Compile and Run commands (`g++`, `python3`, etc.).
- **Limits**: Default CPU/Memory quotas.

---

## 🤝 Contributing

1. Fork the Project
2. Create your Feature Branch
3. Commit your Changes
4. Push to the Branch
5. Open a Pull Request

---

<div align="center">
  <b>Powered by Neuron Engine</b>
</div>
