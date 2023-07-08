@echo off

REM Vitty
REM Handy set of commands to run to get a new server up and running
set command=%1

if "%command%" == "local" (
    shift
    set file=docker-compose-local.yaml
    set environment=local
) else (
    set file=docker-compose-prod.yaml
    set environment=production
)

if "%command%" == "" (
    echo.
    echo ██╗   ██╗██╗████████╗████████╗██╗   ██╗
    echo ██║   ██║██║╚══██╔══╝╚══██╔══╝╚██╗ ██╔╝
    echo ██║   ██║██║   ██║      ██║    ╚████╔╝ 
    echo ╚██╗ ██╔╝██║   ██║      ██║     ╚██╔╝  
    echo  ╚████╔╝ ██║   ██║      ██║      ██║   
    echo   ╚═══╝  ╚═╝   ╚═╝      ╚═╝      ╚═╝   
    echo.                                   
    echo.
    echo Environment: %environment%
    echo.
    echo Usage: vitty [command]
    echo.
    echo Available commands:
    echo   up: Start the server
    echo   down: Stop the server
    echo   restart: Restart the server
    echo   cli: Run a command inside the container
    exit /b 1
)

REM Start server command
if "%command%" == "up" (
    echo Starting server
    docker-compose -f "%file%" up -d --build
    exit /b 1
)

REM Stop server command
if "%command%" == "down" (
    echo Stopping server
    docker-compose -f "%file%" down
    exit /b 1
)

REM Restart server command
if "%command%" == "restart" (
    echo Restarting server
    docker-compose -f "%file%" down
    docker-compose -f "%file%" up -d --build
    exit /b 1
)

REM Management commands
if "%command%" == "cli" (
    shift
    docker-compose -f "%file%" run --rm vitty-api ./bin/vitty %*
    exit /b 1
)
