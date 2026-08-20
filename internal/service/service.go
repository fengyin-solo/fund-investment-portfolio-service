package service

import (
	"fundinvest/internal/config"
	"fundinvest/internal/store"
	"fundinvest/pkg/logger"
)

type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
