<h1 align="center">
  <br>
  Neuron
  <br>
</h1>

<h4 align="center">A high-performance codebase execution engine for the modern web.</h4>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/docker-2496ED?style=for-the-badge&logo=docker" alt="Docker">
  <img src="https://img.shields.io/badge/status-production%20ready-success?style=for-the-badge" alt="Status">
</p>

---

## 🚀 Overview

Neuron Core is the engine behind **code execution at scale**. Designed for EdTech platforms, coding assessment tools, and competitive programming sites, it allows you to run untrusted user code **securely** and **instantly**.

Neuron Core handles the complexity of isolation, resource limits, and compilation so you can focus on building your platform.

---

## ✨ Features

### ⚡ Blazing Fast
By using **Container Pooling**, Neuron eliminates the typical 2-3 second "cold start" of Docker containers. Code execution starts in **milliseconds**, making it feel instant to the end-user.

### � Cost Optimal
Run thousands of executions on minimal hardware. Our lightweight, **Alpine-based custom runtimes** (< 200MB) share kernel resources efficiently, allowing high density on small VPS instances.

### 📊 Detailed Metrics
Neuron doesn't just say "Pass" or "Fail". It provides deep insights:
- **Compile Time vs Run Time**: Know exactly where the bottlenecks are.
- **Granular Errors**: Distinguish between Runtime Errors, Memory Limit Exceeded (MLE), and Time Limit Exceeded (TLE).
- **Exact Profiling**: Millisecond-level precision on execution duration.

### �️ Enterprise-Grade Security
Running user code is dangerous. Neuron handles this with:
- **Network Isolation**: Zero internet access for user code.
- **Read-Only Filesystems**: Prevents tampering with the runner.
- **Strict Quotas**: CPU, Memory, and Disk limits enforced at the kernel level.

---

## 🎯 Use Cases

- **Online Code Editors**: Power real-time "Run" buttons.
- **Assessment Platforms**: Auto-grade millions of student submissions.
- **Competitive Programming**: Host contests with high-reliability judgement.
- **Interview Tools**: Collaborative coding environments.

---

## 🏗️ Deployment Architecture

Neuron Core is built to run as an **Internal Microservice** (listening on port `9000`). It is **not** exposed to the public internet.

Instead, your main backend (e.g., API Gateway, SDK Server) running on port `8080` handles authentication, rate-limiting, and user management, and then forwards the sanitized execution requests to this core engine.

```
       [ Public User ]
              │
              ▼
   [ Main Backend / SDK ]  <-- Public (Port 8080)
      (Auth, Rate Limit)
              │
              │ (gRPC / HTTP)
              ▼
    [ Neuron Core Engine ] <-- Internal (Port 9000)
    (Execution, Sandboxing)
```

## 🔌 API Usage (Internal)

The core engine exposes a raw execution endpoint for your internal infrastructure.

**Endpoint:** `POST /api/v1/execute` (Internal Network Only)

**Request:**
```json
{
  "language": "cpp",
  "code": "#include <iostream>...",
  "input": "test input",
  "limit": { "time_ms": 2000, "memory_kb": 256000 }
}
```

**Response:**
```json
{
  "stdout": "Hello World",
  "exit_code": 0,
  "metrics": { "total_ms": 150, "compile_ms": 120, "run_ms": 30 }
}
```

---

## 📚 Internal Documentation

Neuron Core is built on a layered architecture separating transport, engine, and runner logic. For a deep dive into the internal design, container pooling strategy, and runner implementation, please see:

👉 [**Internal Architecture**](./ARCHITECTURE.md)

---

## 🤝 Contributing

We welcome community contributions to make code execution better for everyone. Please check our issues page to get started.

<div align="center">
  <b>Powering the next generation of coders</b>
</div>
