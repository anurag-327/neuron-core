# Architecture & Request Flow

Neuron is built using a **Microservices Architecture**. This design decouples the "Business Logic" (User handling, Payments) from the "Execution Logic" (Docker, Security).

---

## 1. High-Level Overview

The system consists of two primary services:

1.  **Main Backend (`:8080`)**: The public-facing API Gateway.
    -   Handles User Authentication (JWT/OAuth).
    -   Manages Billing & Credits (Stripe/Database).
    -   Rate Limits users to prevent abuse.
    -   Validates input payload size.

2.  **Neuron Core (`:9000`)**: The internal execution engine.
    -   Runs on a private network (not exposed to internet).
    -   Manages Docker container pools.
    -   Executes code and returns raw stdout/stderr.
    -   Has NO concept of "Users" or "Credits"—it just runs code.

---

## 2. Request Lifecycle

Here is the step-by-step flow of a Code Execution Request:

### Step 1: User Request
The user (frontend/client) sends a POST request to the Public API.
```http
POST https://api.neuron-labs.xyz/api/v1/runner/submit
Authorization: Bearer <token>
Body: { "language": "cpp", "code": "..." }
```

### Step 2: Main Backend Processing
The **Backend Service** receives the request and performs checks:
1.  **Authentication**: Is the token valid?
2.  **Rate Limit**: Has this user exceeded 30 req/min?
3.  **Billing**: Does the user have enough credits?
4.  **Deduction**: Deduct 1 credit from the user's balance in the DB.

### Step 3: Forwarding to Core
If all checks pass, the Backend sends an HTTP request to the **Internal Core Service**.
```http
POST http://localhost:9000/api/v1/execute
Body: { "language": "cpp", "code": "..." }
```

### Step 4: Execution (Neuron Core)
The **Core Engine** takes over:
1.  **Pool**: Grabs a pre-warmed container for C++.
2.  **Setup**: Writes the code to a temp folder mounted to the container.
3.  **Compile**: Runs `g++` (if applicable).
4.  **Run**: Executes the binary with `timeout` and input.
5.  **Metrics**: Calculates execution time and memory usage.
6.  **Cleanup**: Returns the container to the pool (or destroys if dirty).

### Step 5: Response
1.  **Core** returns the JSON result (Stdout, Stderr, Time) to the **Backend**.
2.  **Backend** logs the execution for analytics.
3.  **Backend** sends the final response to the **User**.

---

## 3. Why this approach?

-   **Security**: The dangerous work (running code) happens in a service that doesn't have access to user databases or payment keys.
-   **Scalability**: We can scale the `Backend` for high traffic (1000s of requests/sec) and `Core` for high load (CPU usage) independently.
-   **Simplicity**: The Core engine code is clean and focused. It doesn't need to know about OAuth or SQL.
