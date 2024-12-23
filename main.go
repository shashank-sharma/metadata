package main

import (
	"fmt"

	"github.com/shashank-sharma/metadata/internal/app"
)

func main() {
	application, err := app.New("dashboard-metadata")
	if err != nil {
		fmt.Println("Failed to start Application: ", err)
		return
	}
	application.Start()
}
