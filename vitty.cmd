@echo off

REM Vitty
REM Handy set of commands to run to get a new server up and running
set command=%1

if "%command%"=="" (
    echo Please enter a command
    echo Available commands: up, down, restart, manage
    exit /b 1
)

REM Start server command
if "%command%"=="up" (
    echo Starting server
    docker-compose -f docker-compose-prod.yaml up -d --build
    exit /b 1
)

REM Stop server command
if "%command%"=="down" (
    echo Stopping server
    docker-compose -f docker-compose-prod.yaml down
    exit /b 1
)

REM Restart server command
if "%command%"=="restart" (
    echo Restarting server
    docker-compose -f docker-compose-prod.yaml down
    docker-compose -f docker-compose-prod.yaml up -d --build
    exit /b 1
)

REM Management commands
if "%command%"=="manage" (
    shift
    docker-compose -f docker-compose-prod.yaml run --rm vitty-api ./bin/vitty %*
    exit /b 1
)
