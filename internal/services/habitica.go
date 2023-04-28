package services

import (
	"fmt"
	"io"
	"net/http"
)

type HabiticaService interface {
	ScoreTask(taskIdAlias string) error
}

type HabiticaConfig struct {
	Host     string
	AuthorId string
	ApiToken string
	ClientId string
}

type Habitica struct {
	Host     string
	authorId string
	apiToken string
	ClientId string
}

func NewHabiticaClient(config *HabiticaConfig) *Habitica {
	return &Habitica{
		Host:     config.Host,
		authorId: config.AuthorId,
		apiToken: config.ApiToken,
		ClientId: config.ClientId,
	}
}

func (h *Habitica) ScoreTask(taskIdAlias string) error {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/tasks/%s/score/up",
			h.Host,
			taskIdAlias,
		),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("x-api-user", h.authorId)
	req.Header.Add("x-api-key", h.apiToken)
	req.Header.Add("x-client", h.ClientId)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"habitica error: [%d]: %s",
			response.StatusCode,
			response.Status,
		)
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)

	return nil
}
