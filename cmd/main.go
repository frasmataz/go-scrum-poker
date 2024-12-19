package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aohorodnyk/mimeheader"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/frasmataz/go-scrum-poker/config"
	scrum_poker "github.com/frasmataz/go-scrum-poker/internal"
	"github.com/frasmataz/go-scrum-poker/internal/util"
)

var htmlRenderer = util.NewHTMLRenderer("templates")

func main() {
	config, err := config.GetConfigFromFlags(os.Args[1:])
	if err != nil {
		panic(err)
	}

	// broker := sse.NewBroker[Message]()
	// broker.MessageAdapter = sseHandler

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = htmlRenderer
	e.Static("/static", "./static")

	gameController := scrum_poker.GameController{}

	e.GET("/", func(c echo.Context) error {
		acceptHeader := mimeheader.ParseAcceptHeader(c.Request().Header.Get("Accept"))

		if acceptHeader.Match("text/html") {
			return c.Render(http.StatusOK, "index.html", nil)
		} else if acceptHeader.Match("application/json") {
			return c.JSON(http.StatusOK, nil)
		}

		return echo.NewHTTPError(http.StatusNotAcceptable, "Expected 'text/html' or 'application/json'")
	})

	e.GET("/game", func(c echo.Context) error {
		acceptHeader := mimeheader.ParseAcceptHeader(c.Request().Header.Get("Accept"))

		if acceptHeader.Match("text/html") {
			return c.Render(http.StatusOK, "game.html", gameController.Games)
		} else if acceptHeader.Match("application/json") {
			return c.JSON(http.StatusOK, gameController.Games)
		}

		return nil
	})

	e.Start(fmt.Sprintf("%v:%v", config.Host, config.Port))
}

// func sseHandler(msg Message, clientID string) sse.SSE {
// 	log.Printf("MESSAGE: %v", msg)
// 	sse := sse.SSE{
// 		Event: "message",
// 		Data:  "",
// 	}

// 	msgHTML, _ := htmlRenderer.RenderToString("message.html", msg)
// 	log.Printf("HTML: %v", msgHTML)

// 	// Write the HTML response, but we need to strip out newlines from the template for SSE
// 	sse.Data = strings.Replace(msgHTML, "\n", "", -1)

// 	return sse
// }
