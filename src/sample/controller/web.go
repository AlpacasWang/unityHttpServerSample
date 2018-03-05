package controller

/**************************************************************************************************/
/*!
 *  web.go
 *
 *  webに返答が行くサンプル
 *
 */
/**************************************************************************************************/

import (
	"github.com/labstack/echo"
	"net/http"
)

func Ping(c echo.Context) error{
	c.String(http.StatusOK,"pong");
	return nil;
}

func PingError(c echo.Context) error{
	c.String(http.StatusBadRequest,"pong error");
	return nil;
}
