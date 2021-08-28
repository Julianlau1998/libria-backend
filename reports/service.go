package reports

import (
	"github.com/labstack/echo/v4"
)

type Service struct {
	reportRepo Repository
}

func NewService(reportRepo Repository) Service {
	return Service{
		reportRepo: reportRepo,
	}
}

func (s *Service) Report(id string, contentType string, c echo.Context) error {
	return s.reportRepo.Report(id, contentType)
}

func (s *Service) Unreport(id string, contentType string, c echo.Context) error {
	return s.reportRepo.Unreport(id, contentType)
}
