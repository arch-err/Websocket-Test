package main

import (
	"fmt"
	"html/template"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

var PORT int = 8080

type Templates struct {
	templates *template.Template
}

type Page struct {
	Title string
}

func newPage() Page {
	return Page{Title: ""}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{templates: template.Must(template.ParseGlob("views/*.html"))}

}

func wsock(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			t := time.Now()
			currentTime := t.Format("2006-01-02 15:04:05")
			err := websocket.Message.Send(ws, currentTime)
			time.Sleep(1 * time.Second)
			if err != nil {
				c.Logger().Error(err)
			}
			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate()

	page := newPage()
	e.GET("/", func(c echo.Context) error {
		page.Title = "index"
		return c.Render(200, "index", page)
	})

	e.GET("/*", func(c echo.Context) error {
		return c.String(404, "404 | Not Found")
	})
	e.GET("/ws/file.txt", wsock)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
