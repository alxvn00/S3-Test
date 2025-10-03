package cmd

import (
	"context"
	"fmt"
	"v1/internal/infrastructure/config"
	"v1/internal/infrastructure/logger"
)

func Start(ctx context.Context) error {
	logger.Init()
	l := logger.GetLogger()
	defer func() {
		logger.Close()
	}()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("error create config: %w", err)
	}

	fmt.Println(l)   //
	fmt.Println(cfg) //
	return nil
}
