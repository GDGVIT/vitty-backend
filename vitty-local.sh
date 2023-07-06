#!/bin/sh

# Vitty
# Handy set of commands to run to get a new server up and running
command=$1

if [ -z "$command" ]; then
    echo "Please enter a command"
    echo "Available commands: up, down, restart, cli"
    exit 1
fi

# Start server command
if [ "$command" = "up" ]; then
    echo "Starting server"
    docker-compose -f docker-compose-local.yaml up -d --build
    exit 1
fi

# Stop server command
if [ "$command" = "down" ]; then
    echo "Stopping server"
    docker-compose -f docker-compose-local.yaml down
    exit 1
fi

# Restart server command
if [ "$command" = "restart" ]; then
    echo "Restarting server"
    docker-compose -f docker-compose-local.yaml down
    docker-compose -f docker-compose-local.yaml up -d --build
    exit 1
fi

# Management commands
if [ "$command" = "cli" ]; then
    shift # Discard the first argument
    docker-compose -f docker-compose-local.yaml run --rm vitty-api ./bin/vitty "$@"
    exit 1
fi