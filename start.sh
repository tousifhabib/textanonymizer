#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to stop a service if it is running
stop_service_if_running() {
    local port=$1
    local service_name=$2
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        echo "$service_name is already running. Stopping it..."
        kill $(lsof -t -i:$port)
    fi
}

# Function to start a service
start_service() {
    local service_command=$1
    local service_name=$2
    echo "Starting $service_name..."
    $service_command &
}

# Change to a specified directory and check for success
change_dir() {
    local dir=$1
    if cd $dir; then
        echo "Changed to directory: $dir"
    else
        echo "Failed to change to directory: $dir"
        exit 1
    fi
}

# Backend setup
change_dir "backend/GPT"

echo "Installing Go packages in the vendor directory..."
go mod vendor

stop_service_if_running 8080 "Backend server"
start_service "go run main.go" "backend server"

# spaCy setup
change_dir "../spacy"

if ! command_exists python3 ; then
    echo "Python is not installed. Please install Python and try again."
    exit 1
fi

echo "Creating a virtual environment..."
python3 -m venv venv

echo "Activating the virtual environment..."
source venv/bin/activate

echo "Installing Python packages from requirements.txt..."
pip install -r requirements.txt

echo "Downloading the spaCy English model..."
python -m spacy download en_core_web_sm

stop_service_if_running 5000 "spaCy service"
start_service "python3 main.py" "spaCy service"

# Frontend setup
change_dir "../../frontend"

if ! command_exists bun ; then
    echo "Bun is not installed. Please install Bun and try again."
    exit 1
fi

echo "Installing frontend packages..."
bun install

stop_service_if_running 3000 "Frontend server"
start_service "bun run start" "frontend development server"

# Wait for all background jobs to finish
wait
