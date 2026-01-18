# 🤝 Contributing to Neuron Core

First off, thank you for considering contributing to Neuron! It's people like you that make specialized tools great.

---

## 💡 How can I help?

### 1. Add a New Language
We want to support Go, Rust, Ruby, PHP, and more.
👉 Follow the **[Adding a Language Guide](./ADD_LANGUAGE.md)**.

### 2. Improve Security
Found a way to escape the sandbox?
- **Do NOT open a public issue.**
- Email me directly at `anuragsrivastav0027@gmail.com`.
- Building better isolation (e.g., gVisor, Firecracker support) is a high priority.

### 3. Reporting Bugs
- Open a GitHub Issue.
- Include the `docker logs` output.
- Include the exact code snippet that caused the failure.

---

## 🧑‍💻 Development Workflow

1. **Fork** the repository.
2. **Clone** your fork.
3. Create a branch:
   ```bash
   git checkout -b feat/add-rust-support
   ```
4. **Test** your changes:
   - Ensure `go run cmd/api/main.go` starts without errors.
   - Run a test execution via `curl`.
5. **Commit** follows [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat: add rust runtime`
   - `fix: resolve memory leak in pool`
   - `docs: update setup guide`
6. **Push** and open a **Pull Request**.

---

## 📐 Code Standards

- **Go**: We adhere strictly to `gofmt`.
- **Docker**: Keep images minimal (Alpine preferred).
- **Logging**: Use the `logger` package, do not use `fmt.Print` in production code.
- **Error Handling**: Return specific errors (TLE, MLE) rather than generic error strings.
