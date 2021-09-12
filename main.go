package main

import (
	"database/sql"
	"fmt"
	"libria/answers"
	"libria/comments"
	"libria/reports"
	"libria/topics"
	"libria/utility"
	"libria/votes"
	"os"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		dynamicCORSConfig := middleware.CORSConfig{
			AllowOrigins: []string{"https://libria-app.com", "https://libria-app.netlify.app"},
			AllowHeaders: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		}
		CORSMiddleware := middleware.CORSWithConfig(dynamicCORSConfig)
		CORSHandler := CORSMiddleware(next)
		return CORSHandler(ctx)
	}
}

var dbClient *sql.DB

func startup() {
	dbClient = utility.NewDbClient()
	for utility.Migrate(dbClient) != nil {
		fmt.Printf("Verbindung Fehlgeschlagen, %s", utility.Migrate(dbClient))
		time.Sleep(20 * time.Second)
	}
}

func main() {
	startup()
	VoteRepository := votes.NewRepository(dbClient)
	VoteService := votes.NewService(VoteRepository)
	VoteDelivery := votes.NewDelivery(VoteService)

	AnswerRepository := answers.NewRepository(dbClient)
	AnswerService := answers.NewService(AnswerRepository, VoteService)
	AnswerDelivery := answers.NewDelivery(AnswerService)

	TopicRepository := topics.NewRepository(dbClient)
	TopicService := topics.NewService(TopicRepository, AnswerService)
	TopicDelivery := topics.NewDelivery(TopicService)

	CommentRepository := comments.NewRepository(dbClient)
	CommentService := comments.NewService(CommentRepository)
	CommentDelivery := comments.NewDelivery(CommentService)

	ReportRepository := reports.NewRepository(dbClient)
	ReportService := reports.NewService(ReportRepository)
	ReportDelivery := reports.NewDelivery(ReportService)

	e := echo.New()
	e.Use(CORSMiddlewareWrapper)

	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/api/topics", TopicDelivery.GetAll)
	e.GET("/api/reported/topics", TopicDelivery.GetReported)
	e.GET("/api/topic/:id", TopicDelivery.GetById)
	e.GET("/api/topicname/:topicName", TopicDelivery.GetByTopicName)
	e.GET("/api/randomTopic", TopicDelivery.GetRandom)
	e.POST("/api/topic", TopicDelivery.Post)
	e.PUT("/api/topic/:id", TopicDelivery.Update)
	e.PUT("/api/topicBestAnswer/:id", TopicDelivery.UpdateBestAnswer)
	e.DELETE("/api/topic/:id", TopicDelivery.Delete)

	e.GET("/api/answers", AnswerDelivery.GetAll)
	e.GET("/api/reported/answers", AnswerDelivery.GetReported)
	e.GET("/api/answers/:topic_id", AnswerDelivery.GetAllByTopic)
	e.GET("/api/answer/:id", AnswerDelivery.GetById)
	e.POST("/api/answer", AnswerDelivery.Post)
	e.PUT("/api/answer/:id", AnswerDelivery.Update)
	e.DELETE("/api/answer/:id", AnswerDelivery.Delete)

	e.GET("/api/comments", CommentDelivery.GetAll)
	e.GET("/api/reported/comments", CommentDelivery.GetReported)
	e.GET("/api/comments/:answer_id", CommentDelivery.GetAllByAnswer)
	e.GET("/api/comment/:id", CommentDelivery.GetById)
	e.POST("/api/comment", CommentDelivery.Post)
	e.PUT("/api/comment/:id", CommentDelivery.Update)
	e.DELETE("/api/comment/:id", CommentDelivery.Delete)

	e.PUT("/api/report/:id", ReportDelivery.Report)
	e.PUT("/api/unreport/:id", ReportDelivery.Unreport)

	e.GET("/api/votes/:answer_id", VoteDelivery.GetAllByAnswer)
	e.POST("/api/vote", VoteDelivery.Post)
	e.PUT("/api/vote/:id", VoteDelivery.Update)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
