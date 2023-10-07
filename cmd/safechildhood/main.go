package main

import app "safechildhood/internal/pkg"

func main() {
	app := app.New([]string{"./configs/main.yaml"})

	app.Run()
}
