package internal

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"habitica_functions/internal/services"
)

const PomodoroSessionSize = 4

var (
	redis        *services.UpstashDB
	habiticaApi  *services.Habitica
	singleTaskId string
	setTaskId    string
)

func initialize(ctx context.Context) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	if singleTaskId == "" {
		singleTaskId = strings.TrimSpace(os.Getenv("HABITICA_POMODORO_TASK_ID"))
	}

	if setTaskId == "" {
		setTaskId = strings.TrimSpace(os.Getenv("HABITICA_POMODORO_SET_TASK_ID"))
	}

	if redis == nil {
		redis = services.NewRedisClient(
			ctx,
			strings.TrimSpace(os.Getenv("REDIS_URL")),
		)

		log.Println("redis client created")
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

		log.Println("habitica client created")
	}
}

func FinishFocusSession() error {
	ctx := context.Background()

	initialize(ctx)
	var err error
	var wg sync.WaitGroup
	var count int

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("ticking counter...")
		count, err = redis.TickCounter(ctx)
		log.Println("counter current value:", count)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("scoring the single task...")
		err = habiticaApi.ScoreTask(singleTaskId)
	}()

	wg.Wait()
	if err != nil {
		return err
	}

	if count == PomodoroSessionSize {
		log.Println("pomodoro session finished")

		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Println("scoring the set task...")
			err = habiticaApi.ScoreTask(setTaskId)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Println("resetting counter...")
			err = redis.ResetCounter(ctx)
		}()

		wg.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
