package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/benc-uk/go-rest-api/pkg/sse"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/frasmataz/go-scrum-poker/config"
	"github.com/frasmataz/go-scrum-poker/internal/util"
)

type Message struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

var messages = []Message{
	{ID: 1, Text: "Welcome to the GO + HTMX chat!", Time: time.Now()},
}

var htmlRenderer = util.NewHTMLRenderer("templates")

func main() {
	config, err := config.GetConfigFromFlags(os.Args[1:])
	if err != nil {
		panic(err)
	}

	broker := sse.NewBroker[Message]()
	broker.MessageAdapter = sseHandler

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = htmlRenderer

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.POST("/add", func(c echo.Context) error {
		text := c.FormValue("text")
		newMessage := Message{
			ID:   len(messages) + 1,
			Text: text,
			Time: time.Now(),
		}
		messages = append(messages, newMessage)
		broker.SendToGroup("*", newMessage)

		return c.Render(http.StatusOK, "message.html", newMessage)
	})

	e.GET("/messages", func(c echo.Context) error {
		return c.Render(http.StatusOK, "messages.html", map[string]any{"messages": messages})
	})

	e.GET("/events", func(c echo.Context) error {
		log.Printf("SSE client connected, ip: %v", c.RealIP())

		return broker.Stream(time.Now().String(), c.Response().Writer, *c.Request())
	})

	e.Static("/static", "./static")
	e.Start(fmt.Sprintf("%v:%v", config.Host, config.Port))
}

func sseHandler(msg Message, clientID string) sse.SSE {
	log.Printf("MESSAGE: %v", msg)
	sse := sse.SSE{
		Event: "message",
		Data:  "",
	}

	msgHTML, _ := htmlRenderer.RenderToString("message.html", msg)
	log.Printf("HTML: %v", msgHTML)

	// Write the HTML response, but we need to strip out newlines from the template for SSE
	sse.Data = strings.Replace(msgHTML, "\n", "", -1)

	return sse
}
