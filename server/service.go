package server

import (
    "m5s/domain"
)

type ServerRepo interface {
    Update(metric *domain.Metric) error
}

type ServerService struct {
    repo ServerRepo
}

func NewServerService(repo ServerRepo) *ServerService {
    return &ServerService{
        repo: repo,
    }
}

func (ss *ServerService) Update(
    metricType string,
    name string,
    value string,
) error {
    metric, err := domain.NewMetric(metricType, name, value)
    if err != nil {
        return err
    }

    return ss.repo.Update(metric)
}
