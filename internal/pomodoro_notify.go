package internal

import (
	"context"
	"log"
	"os"
	"strings"

	"habitica_functions/internal/services"
)

const PomodoroSessionSize = 4

var (
	habitica     *services.Habitica
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

	if habitica == nil {
		habitica = services.NewHabiticaClient(
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

	log.Println("scoring a single task...")
	if err := habitica.ScoreTask(singleTaskId); err != nil {
		return err
	}
	task, err := habitica.GetTask(singleTaskId)
	if err != nil {
		return err
	}

	if int(task.Value)%PomodoroSessionSize == 0 {
		log.Println("pomodoro session finished, scoring the set task...")
		if err := habitica.ScoreTask(setTaskId); err != nil {
			return err
		}
	}

	return nil
}
