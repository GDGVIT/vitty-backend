<p align="center">
<a href="https://dscvit.com">
	<img src="https://user-images.githubusercontent.com/30529572/92081025-fabe6f00-edb1-11ea-9169-4a8a61a5dd45.png" alt="DSC VIT"/>
</a>
	<h2 align="center"> Vitty Backend</h2>
	<h4 align="center"> An application that extracts and reads data from a timetable and sends push notifications to user to attend scheduled classes.</h4>
</p>

---
[![Join Us](https://img.shields.io/badge/Join%20Us-Developer%20Student%20Clubs-red)](https://dsc.community.dev/vellore-institute-of-technology/)
[![Discord Chat](https://img.shields.io/discord/760928671698649098.svg)](https://discord.gg/498KVdSKWR)

[![UI ](https://img.shields.io/badge/User%20Interface-Link%20to%20UI-orange?style=flat-square&logo=appveyor)](https://www.figma.com/file/3ILW1qy1qIjiJ5S78zyIqh/VITTY?node-id=1%3A4)


## Table of Contents
- [Key Features](#key-features)
- [Dependencies](#dependencies)
	- [Deployment](#deployment)
	- [Oauth2](#oauth2)
	- [Web API](#web-api)
	- [Database](#database)
- [Setting-up and Installation](#setting-up-and-installation)
	- [Prerequisites](#prerequisites)
	- [Configure the environment variables](#configure-the-environment-variables)
- [Usage](#usage)
	- [Running the application](#running-the-application)
	- [Stopping the application](#stopping-the-application)
	- [Using management commands of the CLI application](#using-management-commands-of-the-cli-application)
- [Developer](#developer)


<br>

## Key Features
- [x] User Authenitication using Google Oauth2
- [x] Timetable detection from VTOP timetable page text
- [x] Add and remove friends using friend requests
- [x] User search to friend other users
- [x] Display friends' timetable
- [x] Show mutual friends
- [x] Friend suggestions based on mutual friends
- [x] CLI application to manage the web api
<!-- - [ ]  < feature >
- [ ]  < feature > -->

<br>
<br>

## Tech Stack and Dependencies
### Deployment
- [Docker](https://www.docker.com/)

### Oauth2
- [Google Oauth Credentials](https://developers.google.com/identity/protocols/oauth2/web-server)

### Web API
- **Language** - [Go](https://go.dev/)
- **Framework** - [Fiber](https://gofiber.io/)
- **ORM** - [Gorm](https://gorm.io/)
- **CLI** framework - [urfave/cli](https://cli.urfave.org/)
- **JWT** - [golang-jwt](https://golang-jwt.github.io/jwt/)

### Database
- [PostgreSQL](https://www.postgresql.org/)  

<br>
<br>

## Setting-up and Installation  
### Prerequisites
- Download and install [Docker](https://docs.docker.com/get-docker/) and [Docker compose](https://docs.docker.com/compose/install/)
- Get [Google Oauth2](https://developers.google.com/identity/protocols/oauth2/web-server) Credentials JSON file for a web application 

### Configure the environment variables
Configure the env files present in `/vitty-backend-api/.env/`   
For local environment use `.local` and for production use `.production`    

Following environment variables need to be configured -
- `FIBER_PORT` - Value of port used by the web api in the form `:<PORT>`, default value is `:3000`
- `DEBUG` - Set to `true` for local environment, `false` for production environment
- `POSTGRES_URL` - Set to `postgres://<POSTGRES_USER>:<POSTGRES_PASSWORD>@postgres:<POSTGRES_PORT>/<POSTGRES_DB>`
- `POSTGRES_USER` - Username for postgres database
- `POSTGRES_PASSWORD` - Password for postgres database
- `POSTGRES_DB` - Name of the postgres database
- `POSTGRES_HOST` - Hostname for postgres database, default value is `postgres`
- `POSTGRES_PORT` - Port for postgres database, default value is `5432`
- `OAUTH_CALLBACK_URL` - Callback URL for Google Oauth2, `http://<backend-url>/api/auth/google/callback`
- `JWT_SECRET` - JWT secret key that will be used to sign the tokens

<br>
<br>


## Usage
### Running the application
Use the following command to run the application -  

#### MacOS and Linux
```zsh
./vitty.sh up
```

#### Windows
```cmd
vitty.bat up
```

<br>

### Stopping the application
Use the following command to stop the application -

#### MacOS and Linux
```zsh
./vitty.sh down
```

#### Windows
```cmd
vitty.bat down
```

<br>

### Using management commands of the CLI application
Use the following command to run the CLI application -

#### MacOS and Linux
```zsh
./vitty.sh cli <command>
```

#### Windows
```cmd
vitty.bat cli <command>
```

> **NOTE:** Replace script file names with `.vitty-local.sh` and `.vitty-local.bat` for local environment

<br>
<br>

## Developer

<table>
	<tr align="center">
		<td>
		Dhruv Shah
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/88224695" width="150" height="150" alt="Dhruv Shah">
		</p>
			<p align="center">
				<a href = "https://github.com/Dhruv9449">
					<img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/>
				</a>
				<a href = "https://www.linkedin.com/in/Dhruv9449" target="_blank">
					<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
				</a>
			</p>
		</td>
	</tr>
</table>

<p align="center">
	Made with :heart: by <a href="https://dscvit.com">DSC VIT</a>
</p>
