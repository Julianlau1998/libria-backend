package votes

import (
	"fmt"
	"libria/auth"
	"libria/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	voteService Service
}

func NewDelivery(voteService Service) Delivery {
	return Delivery{
		voteService: voteService,
	}
}

func (d *Delivery) GetAllByAnswer(c echo.Context) error {
	answerID := c.Param("answer_id")
	votes, err := d.voteService.GetAllByAnswer(answerID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, votes)
}

func (d *Delivery) Post(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	requestBody := new(models.Vote)
	if err := c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	vote, err := d.voteService.Post(requestBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, vote.ID)
}

func (d *Delivery) Update(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	requestBody := new(models.Vote)
	if err = c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	answer, err := d.voteService.Update(id, requestBody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, answer)
}
