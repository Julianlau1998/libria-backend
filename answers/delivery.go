package answers

import (
	"fmt"
	"libria/auth"
	"libria/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	answerService Service
}

func NewDelivery(answerService Service) Delivery {
	return Delivery{
		answerService: answerService,
	}
}

func (d *Delivery) GetAll(c echo.Context) error {
	answers, err := d.answerService.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, answers)
}

func (d *Delivery) GetAllByTopic(c echo.Context) error {
	req := c.Request()
	userId := req.Header.Get("UserId")
	topicID := c.Param("topic_id")
	answers, err := d.answerService.GetAllByTopic(topicID, userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, answers)
}

func (d *Delivery) GetById(c echo.Context) error {
	id := c.Param("id")
	answer, err := d.answerService.GetById(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, answer)
}

func (d *Delivery) Post(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	requestBody := new(models.Answer)
	if err := c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	answer, err := d.answerService.Post(requestBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, answer.ID)
}

func (d *Delivery) Update(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	requestBody := new(models.Answer)
	if err = c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if headers.Get("UserId") != requestBody.UserID {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	answer, err := d.answerService.Update(id, requestBody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, answer)
}

func (d *Delivery) Delete(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	requestBody := new(models.Answer)
	if err = c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Print(requestBody.UserID)

	if headers.Get("UserId") != requestBody.UserID {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	answer, err := d.answerService.Delete(id)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, answer)
}
