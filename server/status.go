package server

import (
	logging "fetch-me-if-you-read-me/logger"
	"fetch-me-if-you-read-me/model"

	"net/http"

	"go.uber.org/zap"
)

type status struct {
	logger *zap.SugaredLogger
	model  *model.Model
}

func (s *status) statusHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("Checking status")
	err := s.model.CheckStatus()

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Debugf("Status probe failed with %v", http.StatusInternalServerError)
	} else {

		w.WriteHeader(http.StatusOK)
	}
}

func newStatus(logger *logging.Logger, model *model.Model) *status {
	return &status{
		logger: logger.Log,
		model:  model,
	}
}
