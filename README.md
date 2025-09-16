**Backend Test**

**Local development**

```bash
make up
```

**Clarifications**

First of all, I had almost no experience with Go before this trail (I had only used it for literally two days to make a mini prototype). Even so, I decided to do it with Go to learn and practice.
It is quite likely that there are things to improve, such as the unlock process.


- The `blocks` and `mempool` microservices have been separated into their corresponding folders, as they have different responsibilities
- Some improvements have been made to the original `blocks` microservice (see commits history)
- The `mempool` microservice:
  - Collects transactions, fetches the current gas price and fee for each transaction, and submits them in batches to the `blocks` microservice
  - Exposes a REST API to accept transactions (e.g., `POST /submitTransactions`). This endpoint try to be as fast as possible, so it just stores the transactions in memory and returns a response
  - Processes transactions sequentially, ensuring that no additional transactions are submitted while a batch is being processed
  - Prioritizes transactions to maximize total fees collected
- A `Makefile` file has been created to help local development

**TODOs (with more time)**
- Add unit tests
- Add an Interface in the `client` (for testing and for decoupling from the concrete implementation)
- Save unsubmitted transactions to avoid endlessly retrying (and apply an alternative strategy to submit them)
- The `TransactionsToSubmit` and `TransactionsToProcess` could store addresses instead of the whole transaction to save memory
- Investigate how to handle panics
- Delegate `Transactions` sorted logic to a database (e.g. Redis with sorted sets)

---

**Objective:**
The purpose of this test is to evaluate your ability to design and implement a backend service that efficiently manages and submits transactions to a simulated blockchain while maximizing profitability.

---

**Scenario:**
You have been given a microservice that simulates a blockchain called **Blocks**. This microservice provides two main endpoints:

1. **GET http://localhost:8080/getCurrentPrice**
   - Returns the current gas price and transaction fee.
   - **Important:** You **MUST** call this endpoint for **every transaction** to retrieve the current fee and gas price before processing it.
   - Example response:
   ```json
   {
       "fee": 13,
       "gasPrice": 105
   }
   ```

2. **POST http://localhost:8080/simulateBlock**
   - Accepts an array of transactions and submits them to Blocks.
   - Constraints:
      - The total gas of the submitted transactions cannot exceed **10,000 gas units**.
      - The service processes transactions sequentially. If a request is currently being processed (which takes between **1 to 3 seconds**), any additional requests will be rejected.
   - Example request:
   ```json
   {
       "transactions": [
           {
               "id": "tx1",
               "fee": 11,
               "gasPrice": 944
           },
           {
               "id": "tx2",
               "fee": 13,
               "gasPrice": 105
           }
       ]
   }
   ```
   - Example response:
   ```json
   {
       "gasLimit": 10000,
       "processingTimeSeconds": 1.058,
       "totalFees": 24,
       "totalGas": 1049,
       "transactions": {
           "transactions": [
               {
                   "id": "tx1",
                   "gasPrice": 944,
                   "fee": 11
               },
               {
                   "id": "tx2",
                   "gasPrice": 105,
                   "fee": 13
               }
           ]
       }
   }
   ```

---

**Your Task:**
You need to build a backend service that acts as a **mempool** or **transaction orchestrator**, which collects, organizes, and efficiently submits transactions to Blocks while maximizing profitability.

### **Requirements:**
1. **Transaction Management:**
   - Accept transactions from multiple sources in real-time.
   - Store and manage transactions in a queue or mempool.
   - **Fetch the gas price and fee for each transaction** from `GET /getCurrentPrice` before processing it.

2. **Batching and Optimization:**
   - Accumulate transactions to maximize block utilization (up to 10,000 gas units per batch).
   - Prioritize transactions to **maximize total fees** collected.

3. **Blockchain Submission:**
   - Submit transactions in an optimized manner to Blocks.
   - Ensure that no additional transactions are submitted while a batch is being processed.

4. **Resilience and Error Handling:**
   - Handle request failures and retry if necessary.
   - Avoid losing transactions due to transient errors.

5. **Concurrency Considerations:**
   - Ensure transactions are processed sequentially while allowing new transactions to be continuously collected.

---

### **Running the Project:**
To run the project, you only need to execute:

```sh
docker-compose up
```

You must **update the provided `docker-compose.yaml` file** to include your solution.

Example `docker-compose.yaml`:

```yaml
version: "3.8"

services:
  blocks:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
  your-service:
    build:
      context: .
      dockerfile: Dockerfile.your_service
    depends_on:
      - blocks
    ports:
      - "9090:9090"
```

Ensure that the container is running before testing your implementation.

### **Technology Choice:**
You are free to use **any programming language or framework** to implement this solution.

---

### **Deliverables & Expectations:**
- The solution does **not** need to be **production-ready**, as we understand that time is limited.
- What is **most important** to us is:
   - **Your ability to adapt to product changes** and iterate over your solution.
   - **Following best practices** in software development.
   - **Providing a working demo** of your solution.
- Start with a **basic version** and iterate on it, refining your approach as you go.
- If you use AI assistance, please share the converstations with us.

### **Repository Setup:**
- You **must fork the provided repository** to complete the test.
- Submit your solution via a pull request when finished.

---

### **Evaluation Criteria:**
Your solution will be assessed based on:
- **Correctness** – Does your service meet the given constraints and requirements?
- **Efficiency** – How well does your service optimize transaction submission?
- **Code Quality** – Is your code clean, well-structured, and maintainable?
- **Error Handling** – How well does your service handle failures and edge cases?
- **Adaptability** – How well can you iterate and improve your initial implementation?

---

**Setup and Testing Environment:**
You will have **four hours** to complete this test in our office environment. A pre-configured instance of the simulated blockchain microservice (**Blocks**) will be provided.

Good luck!