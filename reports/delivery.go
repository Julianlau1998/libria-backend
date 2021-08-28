package reports

import (
	"libria/auth"
	"libria/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Delivery struct {
	reportService Service
}

func NewDelivery(reportService Service) Delivery {
	return Delivery{
		reportService: reportService,
	}
}

func (d *Delivery) Report(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	var topic models.Topic
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	contentType := c.QueryParam("contentType")
	err = d.reportService.Report(id, contentType, c)
	if err != nil {
		log.Warnf("reportDelivery.Report: %s", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, topic)
}

func (d *Delivery) Unreport(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) && headers.Get("userId") != "auth0|61140120b66da800691207c2" {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	contentType := c.QueryParam("contentType")
	err = d.reportService.Unreport(id, contentType, c)
	var topic models.Topic
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, topic)
}
