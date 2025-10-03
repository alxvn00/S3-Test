package app

import (
	"context"
	"log"
	"v1/internal/cmd"
)

func main() {
	ctx := context.Background()

	if err := cmd.Start(ctx); err != nil {
		log.Fatalf("error starting app: %v", err)
	}
}
