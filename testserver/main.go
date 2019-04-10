package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Static("/cc.xml", "testserver/cc.xml")
	e.Logger.Fatal(e.Start(":8080"))
}
