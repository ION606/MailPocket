# MailPocket

*MailPocket* is inspired by the idea of a pocket—a small, handy container designed to securely hold valuable items. In this case, it’s an extremely lightweight “pocket” for safely collecting email form submissions written in Go.

This project provides two independent server implementations:

1. **Batched Write Server**: A dependency-free server that stores emails in a CSV file.
2. **SQLite Server**: A server that uses SQLite for reliable and structured email storage.

Both servers are designed to be simple, efficient, and easy to deploy.

---

## Features

- **Batched Write Server**:
  - No external dependencies.
  - Stores emails in a CSV file (`emails.csv`).
  - Batches writes to reduce disk I/O.
  - Lightweight and minimalistic.

- **SQLite Server**:
  - Uses SQLite for structured and reliable email storage.
  - Prevents duplicate emails with a unique constraint.
  - Easy to query and manage stored emails.
  - Slightly more robust and feature-rich.

---

## Getting Started

### Prerequisites

- Go 1.20 or higher.
- For the SQLite Server, the `modernc.org/sqlite` dependency will be installed automatically.

---

### Running the Servers
1.0 Clone the repository:
   ```bash
   git clone https://github.com/your-username/forms-server.git
   cd forms-server
   ```

2.0 Using the Makefile:  
   - ```bash
   	 make run-batched # for CSV-based
   	 ```
   - ```bash
	 make run-sqlite # for SQLite-based
	 ```

2.1 **Batched Write Server**:
   - Navigate to the `batched-server` directory:
     ```bash
     cd batched-server
     ```
   - Run the server manually:
     ```bash
     go run main.go
     ```
   - The server will start on `http://localhost:3000`.

2.2 **SQLite Server**:
   - Navigate to the `sqlite-server` directory:
     ```bash
     cd sqlite-server
     ```
   - Install dependencies (if not already installed):
     ```bash
     go mod tidy
     ```
   - Run the server manually:
     ```bash
     go run main.go
     ```
   - The server will start on `http://localhost:3000`.

---

### API Endpoints

Both servers expose the following endpoints:

- **`GET /`**:
  - Returns a `200 OK` status with a message indicating the server is running.
  - Example response:
    ```json
    "Batched Write Server is running"
    ```

- **`POST /submit`**:
  - Accepts a form submission with an `email` field.
  - Example request:
    ```bash
    curl -X POST -d "email=user@example.com" http://localhost:3000/submit
    ```
  - Example response:
    ```json
    {"message": "data received"}
    ```

---

## Server Implementations

### 1. Batched Write Server (CSV-based)

- **Purpose**: A dependency-free implementation for users who want minimalism and simplicity.
- **Storage**: Emails are appended to a CSV file (`emails.csv`) in the `batched-server` directory.
- **Performance**: Batches writes to reduce disk I/O, flushing every 5 seconds or after 100 emails.
- **Format**: Each row in the CSV file contains an email and a timestamp.
- **Use Case**: Ideal for lightweight deployments where external dependencies are not desired.

### 2. SQLite Server

- **Purpose**: A more robust implementation using SQLite for structured storage.
- **Storage**: Emails are stored in an SQLite database (`emails.db`) in the `sqlite-server` directory.
- **Features**:
  - Prevents duplicate emails with a unique constraint.
  - Tracks the creation timestamp of each email.
- **Use Case**: Ideal for deployments where data integrity and querying capabilities are important.

---

## Configuration

- **Port**: Both servers run on port `3000` by default. To change the port, modify the `PORT` constant in the respective `main.go` file.
- **Storage Location**:
  - Batched Write Server: Emails are stored in `batched-server/emails.csv`.
  - SQLite Server: Emails are stored in `sqlite-server/emails.db`.

---

## Trade-offs

| Feature                | Batched Write Server (CSV)  | SQLite Server               |
|------------------------|------------------------------|-----------------------------|
| **Dependencies**       | None                        | Requires SQLite dependency  |
| **Storage**            | CSV file                    | SQLite database             |
| **Data Integrity**     | Basic                       | High (prevents duplicates)  |
| **Querying**           | Not supported               | Supported                   |
| **Performance**        | High (batched writes)       | High (SQLite optimized)     |
| **Use Case**           | Minimalist, dependency-free | Robust, structured storage  |

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
