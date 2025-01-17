package agent

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"m5s/domain"
	"m5s/internal/repository"
)

type Repo interface {
	Update(metric *domain.Metric) error
	GetMetricsList() []*domain.Metric
}

type Service struct {
	repo           Repo
	stat           *domain.Statistics
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewAgentService(
	pollInterval time.Duration,
	reportInterval time.Duration,
	serverAddr string,
) *Service {
	repo := repository.NewInMemStorage()

	return &Service{
		repo:           repo,
		stat:           domain.NewStatistics(),
		serverAddr:     serverAddr,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
}

func (as *Service) StartPoller() {
	log.Printf("Start poller with %d duration", as.pollInterval)

	ticker := time.NewTicker(as.pollInterval)

	for range ticker.C {
		as.stat.Refresh()

		for name, value := range as.stat.CurrentValues {
			metric := domain.NewGauge(name, value)
			if err := as.repo.Update(metric); err != nil {
				log.Fatalf("%v", err)
			}
		}

		pollCountMetric := domain.NewCounter("PollCount", 1)
		if err := as.repo.Update(pollCountMetric); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func (as *Service) StartReporter() {
	log.Printf("Start reporter with %d duration", as.reportInterval)

	ticker := time.NewTicker(as.reportInterval)

	for range ticker.C {
		for _, metric := range as.repo.GetMetricsList() {
			if err := as.makeRequest(metric); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}

func (as *Service) makeRequest(metric *domain.Metric) error {
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"http://%s/update/%s/%s/%s",
			as.serverAddr, metric.Type, metric.Name, metric,
		),
		nil,
	)
	if err != nil {
		return fmt.Errorf("execute http request: %v", err)
	}

	request.Header.Set("Content-Type", "text/plain")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
