package internal

import (
	"context"
	"os"
	"strings"
	"sync"

	"habitica-functions/internal/services"
)

const PomodoroSessionSize = 4

type PomodoroConfig struct {
	SingleTaskID string
	SetTaskID    string
}

var (
	redis       *services.UpstashDB
	habiticaApi *services.Habitica
)

func initialize(ctx context.Context) {
	if redis == nil {
		redis = services.NewRedisClient(
			ctx,
			strings.TrimSpace(os.Getenv("REDIS_URL")),
		)
	}

	if habiticaApi == nil {
		habiticaApi = services.NewHabiticaClient(
			&services.HabiticaConfig{
				Host:     strings.TrimSpace(os.Getenv("HABITICA_HOST")),
				AuthorId: strings.TrimSpace(os.Getenv("HABITICA_AUTHOR_ID")),
				ApiToken: strings.TrimSpace(os.Getenv("HABITICA_API_KEY")),
				ClientId: strings.TrimSpace(os.Getenv("HABITICA_CLIENT_ID")),
			},
		)
	}
}

func FinishFocusSession(config *PomodoroConfig) error {
	ctx := context.Background()

	initialize(ctx)
	var err error
	var wg sync.WaitGroup
	var count int

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err = redis.TickCounter(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = habiticaApi.ScoreTask(config.SingleTaskID)
	}()

	wg.Wait()
	if err != nil {
		return err
	}

	if count == PomodoroSessionSize {
		wg.Add(1)
		go func() {
			err = habiticaApi.ScoreTask(config.SetTaskID)
		}()

		wg.Add(1)
		go func() {
			err = redis.ResetCounter(ctx)
		}()

		wg.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
