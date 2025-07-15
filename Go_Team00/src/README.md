# How to Start the Project

This guide explains how to start the project in two terminal windows using Docker Compose.

## Prerequisites

- Ensure you have Docker and Docker Compose installed on your system.
- Ensure the `.env` file is properly configured with the correct database credentials and server address.

## Steps to Start the Project


### Terminal 1: Start the gRPC Server

1. Open a new terminal window.
2. Navigate to the `src/server` directory:
   ```bash
   cd src/server
   ```
3. Run the gRPC server:
   ```bash
   go run main.go
   ```

   This will start the gRPC server and make it accessible on `localhost:50051`.

### Terminal 2: Start the PostgreSQL Database

1. Navigate to the `src` directory:
   ```bash
   cd src
   ```
2. Start the PostgreSQL database using Docker Compose:
   ```bash
   docker-compose up -d
   ```

   This will start the PostgreSQL container and make it accessible on `localhost:55432`.

3. Navigate to the `src/receiver` directory and run the receiver:
    ```bash
    cd ./receiver
    go run main.go -k <anomaly-coefficient>
    ```
    
    This will start the receiver service that listens for messages from the gRPC server and interacts with the PostgreSQL database.

