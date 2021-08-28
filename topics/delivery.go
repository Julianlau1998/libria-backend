package topics

import (
	"libria/auth"
	"libria/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	topicService Service
}

func NewDelivery(topicService Service) Delivery {
	return Delivery{
		topicService: topicService,
	}
}

func (d *Delivery) GetAll(c echo.Context) error {
	topics, err := d.topicService.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, topics)
}

func (d *Delivery) GetReported(c echo.Context) error {
	topics, err := d.topicService.GetReported()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, topics)
}

func (d *Delivery) GetById(c echo.Context) error {
	id := c.Param("id")
	topic, err := d.topicService.GetById(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, topic)
}

func (d *Delivery) GetRandom(c echo.Context) error {
	topics, err := d.topicService.GetRandom()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, topics)
}

func (d *Delivery) Post(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	requestBody := new(models.Topic)
	if err := c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if headers.Get("UserId") != requestBody.UserID {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	topic, err := d.topicService.Post(requestBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, topic.ID)
}

func (d *Delivery) Update(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	requestBody := new(models.Topic)
	if err = c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if headers.Get("UserId") != requestBody.UserID {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	topic, err := d.topicService.Update(id, requestBody)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, topic)
}

func (d *Delivery) UpdateBestAnswer(c echo.Context) (err error) {
	id := c.Param("id")
	bestAnswer, err := d.topicService.UpdateBestAnswer(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, bestAnswer)
}

func (d *Delivery) Delete(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	var isAdmin bool
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	userId := headers.Get("userId")
	id := c.Param("id")
	for i := 0; i < len(auth.Admins()); i++ {
		if userId == auth.Admins()[i] {
			isAdmin = true
		}
	}
	if isAdmin == true {
		topic, err := d.topicService.DeleteAsAdmin(id, userId)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, topic)
	}

	topic, err := d.topicService.Delete(id, userId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, topic)
}
