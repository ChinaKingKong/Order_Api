package main

import (
	"log"
	"order_api/app"
)

func main() {
	application := app.NewApp()
	
	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}