package app

/**************************************************************************************************/
/*!
 *  route.go
 *
 *  ルート管理ファイル
 *
 */
/**************************************************************************************************/
import (
	"fmt"
	"reflect"
	"runtime"
	"sample/controller"
	"strings"

	"github.com/labstack/echo"
)

type routeMap map[string]func(echo.Context) error

var sampleRoute = routeMap{
	"request": controller.SampleTest,
	"error":   controller.SampleError,
}

var sampleRoute2 = routeMap{
	"request": controller.SampleTest2,
}

var webRoute = routeMap{
	"ping":       controller.Ping,
	"ping_error": controller.PingError,
}

/**************************************************************************************************/
/*!
 *  post routing
 */
/**************************************************************************************************/
func postRoute(root *echo.Echo, prefix string, routeConf routeMap, middleware ...echo.MiddlewareFunc) {
	// ミドルウェアをまとめる
	var m []echo.MiddlewareFunc
	m = append(m, middleware...)

	// サブルート生成
	sub := root.Group(prefix, m...)

	// ROUTE TOP
	fmt.Println("[ROUTE]", prefix)
	for k, v := range routeConf {
		sub.POST("/"+k, v)
		// ROUTE PRINT
		space := strings.Repeat(" ", 30-len(k))
		vOf := reflect.ValueOf(v)
		fmt.Println("POST :", k, space, "->", runtime.FuncForPC(vOf.Pointer()).Name())
	}
}

/**************************************************************************************************/
/*
  GETで受け付けるルートを設定する

  argument
    root       : ルート直下
    prefix     : 接頭辞URL
    routeConf  : URLハンドラマップ
    middleware : 事前処理群

  return
    ディスパッチャ
*/
/**************************************************************************************************/
func getRoute(root *echo.Echo, prefix string, routeConf routeMap, middleware ...echo.MiddlewareFunc) {
	// ミドルウェアをまとめる
	var m []echo.MiddlewareFunc
	m = append(m, middleware...)

	// サブルート生成
	sub := root.Group(prefix)

	// ROUTE TOP
	fmt.Println("[ROUTE]", prefix)
	for k, v := range routeConf {
		sub.GET("/"+k, v)
		// ROUTE PRINT
		space := strings.Repeat(" ", 30-len(k))
		vOf := reflect.ValueOf(v)
		fmt.Println("GET :", k, space, "->", runtime.FuncForPC(vOf.Pointer()).Name())
	}
}
