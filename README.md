# Go PostgreSQL Project

A Go application that demonstrates basic PostgreSQL database operations using the `pgx` driver and Viper for configuration management.

## Overview

This project connects to a PostgreSQL database, creates a `users` table, and performs insert operations with duplicate-key handling. It serves as a practical example of:
- Database connectivity with PostgreSQL
- Configuration management using environment variables
- SQL table creation and data insertion
- Error handling and logging in Go

## Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database (local or remote)
- Environment variables configured (see [Configuration](#configuration) section)

## Project Structure

```
go-postgres/
├── main.go          # Main application code
├── go.mod           # Module definition and dependencies
├── go.sum           # Checksum file for dependencies
├── .env             # Environment variables (not included in repo)
└── README.md        # This file
```

## Dependencies

The project uses the following key dependencies:

- **github.com/jackc/pgx/v5** - PostgreSQL driver for Go with excellent performance
- **github.com/spf13/viper** - Configuration management library for handling environment variables
- **github.com/joho/godotenv** - Utility to load environment variables from `.env` file (optional)

## Configuration

Create a `.env` file in the project root with the following variable:

```env
CONN_STR=postgres://username:password@localhost:5432/database_name
```

### Connection String Format

The connection string follows the standard PostgreSQL URI format:
```
postgres://[username[:password]@][host][:port][/database][?param=value]
```

**Example:**
```env
CONN_STR=postgres://postgres:mypassword@localhost:5432/testdb
```

## Building and Running

### Run Directly

```bash
go run main.go
```

### Build Executable

```bash
go build -o go-postgres.exe
```

Then run:
```bash
.\go-postgres.exe
```

### Install Dependencies

```bash
go mod tidy
go mod download
```

## What the Application Does

1. **Loads Configuration**: Reads the database connection string from the `.env` file using Viper
2. **Connects to Database**: Establishes a connection to PostgreSQL
3. **Creates Table**: Creates a `users` table with the following schema:
   - `id` - Auto-incrementing primary key (SERIAL)
   - `username` - Unique username (VARCHAR 50)
   - `email` - Unique email (VARCHAR 100)
   - `created_at` - Timestamp with default value (CURRENT_TIMESTAMP)
4. **Inserts Data**: Attempts to insert three user records with duplicate-key conflict handling
5. **Displays Results**: Prints the current database time and configuration values

### Table Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Features

- **Duplicate Key Handling**: Uses `ON CONFLICT (username) DO NOTHING` to gracefully handle duplicate usernames
- **Error Logging**: Implements comprehensive error handling with detailed log messages
- **Context Management**: Uses Go's context for timeout and cancellation support
- **Configuration Flexibility**: Supports both `.env` files and system environment variables

## Output Example

```
Table 'users' created or already exists.
User alice inserted successfully
User bob inserted successfully
Failed to insert user alice: <error details>
Current time: 2025-12-09 15:30:45.123456 +0000 UTC
Developer: Hozana
```

## Error Handling

The application implements error handling for:
- Configuration file loading failures
- Database connection failures
- SQL execution failures
- Individual record insertion failures

Each error includes contextual information to aid debugging.

## Development Notes

- The application uses `pgx` for direct database access without an ORM
- Configuration is managed through Viper with automatic environment variable reading
- The code includes commented-out `godotenv` usage as an alternative configuration method
- Contexts are properly managed with deferred connection closing

## Future Enhancements

Potential improvements could include:
- Query result retrieval and display
- Transaction support
- Connection pooling configuration
- Prepared statements for better security
- Database schema versioning/migrations
- Structured logging

## License

This is a demonstration project. Feel free to use and modify as needed.

## Author

Hozana
#