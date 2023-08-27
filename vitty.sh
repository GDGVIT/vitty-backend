#!/bin/sh

# Vitty
# Handy set of commands to run to get a new server up and running
if [ "$1" = "local" ]; then
    shift # Discard the first argument
    environment="local"
    file="docker-compose-local.yaml"
else
    environment="production"
    file="docker-compose-prod.yaml"
fi
command=$1

if [ -z "$command" ]; then
    echo
    echo "██╗   ██╗██╗████████╗████████╗██╗   ██╗"
    echo "██║   ██║██║╚══██╔══╝╚══██╔══╝╚██╗ ██╔╝"
    echo "██║   ██║██║   ██║      ██║    ╚████╔╝ "
    echo "╚██╗ ██╔╝██║   ██║      ██║     ╚██╔╝  "
    echo " ╚████╔╝ ██║   ██║      ██║      ██║   "
    echo "  ╚═══╝  ╚═╝   ╚═╝      ╚═╝      ╚═╝   "  
    echo        
    echo "Environment: $environment"
    echo
    echo "Usage: vitty [command]"
    echo
    echo "Available commands:"
    echo "  up: Start the server"
    echo "  down: Stop the server"
    echo "  restart: Restart the server"
    echo "  cli: Run a command inside the container"
    exit 1
fi

# Start server command
if [ "$command" = "up" ]; then
    echo "Starting server"
    docker compose -f "$file" up -d --build
    exit 1
fi

# Stop server command
if [ "$command" = "down" ]; then
    echo "Stopping server"
    docker compose -f "$file" down
    exit 1
fi

# Restart server command
if [ "$command" = "restart" ]; then
    echo "Restarting server"
    docker compose -f "$file" down
    docker compose -f "$file" up -d --build
    exit 1
fi

# Management commands
if [ "$command" = "cli" ]; then
    shift # Discard the first argument
    docker compose -f $"file" run --rm vitty-api ./bin/vitty "$@"
    exit 1
fi