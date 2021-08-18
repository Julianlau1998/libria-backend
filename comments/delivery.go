package comments

import (
	"fmt"
	"libria/auth"
	"libria/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	commentService Service
}

func NewDelivery(commentService Service) Delivery {
	return Delivery{
		commentService: commentService,
	}
}

func (d *Delivery) GetAll(c echo.Context) error {
	comments, err := d.commentService.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, comments)
}

func (d *Delivery) GetAllByAnswer(c echo.Context) error {
	req := c.Request()
	userId := req.Header.Get("UserId")
	answerID := c.Param("answer_id")
	comments, err := d.commentService.GetAllByAnswer(answerID, userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, comments)
}

func (d *Delivery) GetById(c echo.Context) error {
	id := c.Param("id")
	comment, err := d.commentService.GetById(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, comment)
}

func (d *Delivery) Post(c echo.Context) error {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	requestBody := new(models.Comment)
	if err := c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if headers.Get("UserId") != requestBody.UserID {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	comment, err := d.commentService.Post(requestBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, comment.ID)
}

func (d *Delivery) Update(c echo.Context) (err error) {
	req := c.Request()
	headers := req.Header
	if !auth.IsAuthorized(headers["Cookie"]) {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	requestBody := new(models.Comment)
	if err = c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if headers.Get("UserId") != requestBody.UserID {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}
	comment, err := d.commentService.Update(id, requestBody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, comment)
}

func (d *Delivery) Delete(c echo.Context) (err error) {
	// req := c.Request()
	// headers := req.Header
	// if !auth.IsAuthorized(headers["Cookie"]) {
	// 	return c.String(http.StatusUnauthorized, "unauthorized")
	// }
	id := c.Param("id")
	requestBody := new(models.Comment)
	if err = c.Bind(requestBody); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// if headers.Get("UserId") != requestBody.UserID {
	// 	return c.String(http.StatusUnauthorized, "unauthorized")
	// }
	comment, err := d.commentService.Delete(id)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, comment)
}
