#!/bin/bash

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Change to the backend directory
cd backend

# Install Go packages in the vendor directory
echo "Installing Go packages in the vendor directory..."
go mod vendor

# Check if the backend server is already running
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "Backend server is already running. Stopping it..."
    kill $(lsof -t -i:8080)
fi

# Start the backend server
echo "Starting backend server..."
go run main.go &

# Change to the frontend directory
cd ../frontend

# Check if Bun is installed
if ! command_exists bun ; then
    echo "Bun is not installed. Please install Bun and try again."
    exit 1
fi

# Install frontend packages using Bun
echo "Installing frontend packages..."
bun install

# Check if the frontend server is already running
if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null ; then
    echo "Frontend server is already running. Stopping it..."
    kill $(lsof -t -i:3000)
fi

# Start the frontend development server in the background
echo "Starting frontend development server..."
bun run start &

# Wait for the backend and frontend servers to finish
wait