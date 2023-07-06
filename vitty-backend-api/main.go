package main

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/cmd"

func main() {
	vittyApp := cmd.NewVittyCliApp()
	vittyApp.Run()
}
