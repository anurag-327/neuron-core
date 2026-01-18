# System Design & architecture

> "The true complexity of a Code Execution Engine lies not in running the code, but in doing so securely, instantly, and reliably at scale."

This document explains the **Why** and **How** behind Neuron Core, specifically written for engineers and evaluators exploring this codebase.

---

## 1. The Challenge

Building a platform like LeetCode or Remote Interview tools requires running code submitted by untrusted users. This presents three massive challenges:
1.  **Security**: A malicious user could run `rm -rf /` or mine crypto.
2.  **Latency**: Users expect "Runs..." -> "Output" in milliseconds. Docker containers typically take 2-3 seconds to start.
3.  **Stability**: One user's infinite loop shouldn't crash the server for everyone else.

## 2. The Solution: Dedicated Microservice

We moved the execution logic into its own service (`Neuron Core`). This wasn't just for clean code—it solves real infrastructure problems.

### ❓ Why separate it?

#### 1. Scaling Code is Different
- **Backend (API)**: Handles users and database. Needs little CPU.
- **Engine (Core)**: Compiles code. Needs huge CPU.
- **Win**: We can run the Engine on powerful servers and the API on cheap ones. If traffic spikes, we just add more Engine nodes.

#### 2. Safety First ("Blast Radius")
- Running user code is dangerous. If a malicious script crashes the server, it **only** crashes the Engine.
- **Win**: Your main website (Login, Payments, Dashboard) stays online no matter what.

#### 3. Faster Updates
- Want to add Rust support? Just update the Engine.
- **Win**: You don't have to redeploy your entire website just to upgrade a compiler.

#### 4. Reusability
- This Engine is now a generic tool.
- **Win**: You can use the same Engine for a Coding Interview app, a homework tool, or a contest platform simultaneously.

## 3. Key Design Decisions

### A. The "Cold Start" Problem (Solved by Pooling)
**Problem**: Docker `run` is slow (1-3s).
**Our Solution**: We implemented a **Container Pool Manager**.
- The system keeps `N` containers "warm" (running `sleep infinity`) for each language.
- When a request comes in, we grab a hot container from the channel (`< 5ms` latency).
- After execution, if the container is "dirty" (e.g., file system modified), we trash it. If "clean", we recycle it.

### B. Security via Defense-in-Depth
1.  **Network**: Containers run with `--network none`. No data exfiltration possible.
2.  **Filesystem**: Root FS is Read-Only. Temp writes only allowed in `/tmp`.
3.  **Kernel Limits**: We use Docker's `pids-limit` to prevent fork bombs.

### C. Lightweight Runtimes
We rejected standard images (`gcc:latest` is >1GB).
We built custom **Alpine-based images**:
- **C++**: 300MB (vs 1.2GB)
- **Node**: 200MB (vs 900MB)
- **Python**: 250MB (vs 900MB)

This allows us to run **hundreds of runners** on a single cheap VPS.

---

## 4. Why Go?

We chose Go for the core engine because:
1.  **Concurrency**: The `Pool Manager` relies heavily on Channels and Goroutines to manage container lifecycles without blocking.
2.  **Docker SDK**: Go is the native language of Docker. usage is first-class.
3.  **Performance**: Near C-level performance with memory safety.

---

## 5. Integrating this Engine

While this repository stands alone, in a production environment it sits behind your API Gateway.

**Flow:**
1.  User posts code to `api.neuron-labs.xyz/api/v1/runner/submit`.
2.  Your Backend checks if user has credits.
3.  Your Backend forwards request to `internal-neuron-core:9000`.
4.  Neuron Core executes and returns metrics.
5.  Your Backend saves result to DB and responding to user.

---

## 6. Future Roadmap

- [ ] **Firecracker MicroVMs**: Moving from Docker to AWS Firecracker for even harder isolation.
- [ ] **WASM**: Explore WebAssembly for sandboxing.