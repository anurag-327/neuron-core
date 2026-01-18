<h1 align="center">
  <br>
  <img src="https://svg-banners.vercel.app/api?type=glitch&text1=Neuron&width=800&height=200" alt="Neuron">
  <br>
  Neuron Core
  <br>
</h1>

<h4 align="center">A high-performance codebase execution engine for the modern web.</h4>


<p align="center">
  <img src="https://img.shields.io/badge/languages-4+-blue?style=for-the-badge" alt="Languages">
  <img src="https://img.shields.io/badge/latency-300--400ms-green?style=for-the-badge" alt="Latency">
  <img src="https://img.shields.io/badge/sandbox-Docker-2496ED?style=for-the-badge&logo=docker" alt="Docker">
  <img src="https://img.shields.io/badge/status-production%20ready-success?style=for-the-badge" alt="Status">
</p>


---

## Overview

Neuron Core is the engine behind **code execution at scale**. Designed for EdTech platforms, coding assessment tools, and competitive programming sites, it allows you to run untrusted user code **securely** and **instantly**.

Neuron Core handles the complexity of isolation, resource limits, and compilation so you can focus on building your platform.

👉 **Read the [System Design & Philosophy](./SYSTEM_DESIGN.md) guide to understand the engineering choices behind this engine.**


---
## Features

### ⚡ Blazing Fast
We pre-start containers ("Pooling") so code runs instantly. No 3-second Docker startup time.

### 💰 Cost Optimal
We built custom, tiny runtimes (Alpine Linux). They are 10x smaller than standard images, saving huge RAM and Disk space.

### 📊 Metric-Rich
We don't just say "Error". We tell you if it was Compilation Time, Run Time, or Memory Limit, down to the millisecond.

### 🛡️ Secure
- **No Internet**: User code cannot fetch external URLs.
- **Read-Only**: User code cannot delete system files.
- **Strict Limits**: CPU and RAM are capped at the kernel level.

---

## 🏗️ Deployment

Neuron Core runs as an **Internal Microservice** (port `9000`).

It sits safely behind your Main Backend (port `8080`), which handles User Auth and billing.

```
[ User ] -> [ Your Backend (Auth) ] -> [ Neuron Core (Runs Code) ]
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
}
```

**Response:**
```json
{
    "stdout": "Hello, World!",
    "stderr": "",
    "err_type": "",
    "err_msg": "",
    "exit_code": 0,
    "metrics": {
        "total_ms": 296,
        "compile_ms": 255,
        "run_ms": 40
    }
}
```

---

## 📚 Internal Documentation

Neuron Core is built on a layered architecture separating transport, engine, and runner logic. For a deep dive into the internal design, container pooling strategy, and runner implementation, please see:

👉 [**Internal Architecture**](./ARCHITECTURE.md)

---

## 📚 Documentation

- [**Setup Guide**](./SETUP.md) - Run locally or in production.
- [**Add New Language**](./ADD_LANGUAGE.md) - How to add Rust, Go, etc.
- [**Contributing**](./CONTRIBUTING.md) - Help us build better tools.
- [**System Design**](./SYSTEM_DESIGN.md) - Deep dive into architecture.

---

<div align="center">
  <b>Powering the next generation of coders</b>
</div>
