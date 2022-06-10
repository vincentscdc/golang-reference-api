package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"bnpl/internal/configuration"
	"bnpl/internal/port/grpc/protos/creditline/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// configuration
	currEnv := "local"
	if e := os.Getenv("APP_ENVIRONMENT"); e != "" {
		currEnv = e
	}

	cfg, err := configuration.GetConfig(currEnv)
	if err != nil {
		if errors.Is(err, configuration.MissingBaseConfigError{}) {
			log.Printf("getConfig: %v", err)

			return
		}

		log.Printf("getConfig: %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to server")

		return
	}

	defer conn.Close()

	c := creditline.NewPayLaterServiceClient(conn)

	ctx := context.Background()

	reply, err := c.GetCreditLine(ctx, &creditline.GetCreditLineRequest{})
	if err != nil {
		log.Error().Err(err).Msg("failed to get server credit line")

		return
	}

	log.Info().Msgf("reply: %v", reply)
}
