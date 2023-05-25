package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Task struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Down      bool      `json:"down"`
	Up        bool      `json:"up"`
}

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

func (h *Habitica) createRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-user", h.authorId)
	req.Header.Add("x-api-key", h.apiToken)
	req.Header.Add("x-client", h.ClientId)

	return req, nil
}

func (h *Habitica) GetTask(taskIdAlias string) (Task, error) {
	req, err := h.createRequest("GET",
		fmt.Sprintf(
			"%s/tasks/%s",
			h.Host,
			taskIdAlias,
		))
	if err != nil {
		return Task{}, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return Task{}, err
	}

	if response.StatusCode != http.StatusOK {
		return Task{}, fmt.Errorf(
			"habitica error: [%d]: %s",
			response.StatusCode,
			response.Status,
		)
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return Task{}, err
	}

	var task Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (h *Habitica) ScoreTask(taskIdAlias string) error {
	req, err := h.createRequest("POST",
		fmt.Sprintf(
			"%s/tasks/%s/score/up",
			h.Host,
			taskIdAlias,
		))
	if err != nil {
		return err
	}

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
