package main

import app "safechildhood/internal/pkg"

func main() {
	app := app.New("./configs/main.yaml")

	app.Run()
}
