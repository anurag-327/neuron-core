# 🌐 Adding a New Language

So you want to add support for Rust, Go, or Ruby? Follow this 3-step guide.

---

## Step 1: Create the Dockerfile

1. Create a new folder in `runtime/images/`:
   ```bash
   mkdir runtime/images/rust
   ```

2. Create a `Dockerfile` inside it. 
   **Rule:** Keep it under 300MB if possible (use Alpine).

   ```dockerfile
   # runtime/images/rust/Dockerfile
   FROM alpine:3.19

   # Install compiler & busybox (required for timeout)
   RUN apk add --no-cache rust cargo busybox

   # Setup sandbox
   RUN mkdir -p /sandbox && chmod 777 /sandbox
   WORKDIR /sandbox
   CMD ["sh"]
   ```

3. Build and test it manually:
   ```bash
   docker build -t neuron-rust runtime/images/rust
   ```

---

## Step 2: Register in Runtime Config

Open `runtime/runtime.go` and add an entry to `LanguageRegistry`.

```go
"rust": {
    Language:       "rust",
    Image:          "neuron-rust",  // Must match your docker tag
    InitSize:       1,
    MaxSize:        5,
    HealthCmd:      []string{"echo", "ok"},
    ResourceLimits: commonResouceLimits,

    // Define the Entry File
    EntryFile: EntryFile{
        FileName:  "main",
        Extension: "rs",
    },

    // Compilation Command
    CompileCmd: func(n FileNames) string {
        // rustc main.rs -o main
        return fmt.Sprintf("rustc %s -o %s", n.FullName, n.FileName)
    },

    // Execution Command
    RunCmd: func(n FileNames) string {
        // ./main < input.txt
        return fmt.Sprintf("./%s < input.txt", n.FileName)
    },
},
```

---

## Step 3: Update Build Script

Add your new image to `scripts/setup.sh` so it builds automatically for others.

```bash
# scripts/setup.sh
docker build -t neuron-rust runtime/images/rust  
```

---

## Step 4: Restart & Verify

1. Restart the server:
   ```bash
   go run cmd/api/main.go
   ```
   *You should see logs: `Pre-warming container pool for rust...`*

2. Test execution:
   ```bash
   curl -X POST http://localhost:9000/api/v1/execute \
     -d '{
       "language": "rust",
       "code": "fn main() { println!(\"Hello Rust\"); }",
       "input": ""
     }'
   ```
