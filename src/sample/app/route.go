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
	"net/http"
	"reflect"
	"runtime"
	"sample/controller"
	"strings"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type routeMap map[string]func(web.C, http.ResponseWriter, *http.Request)

var sampleRoute = routeMap{
	"request": controller.SampleTest,
}

var sampleRoute2 = routeMap{
	"request": controller.SampleTest2,
}

var webRoute = routeMap{
	"ping": controller.Ping,
}

/**************************************************************************************************/
/*!
 *  post routing
 */
/**************************************************************************************************/
func postRoute(prefix string, routeConf map[string]func(web.C, http.ResponseWriter, *http.Request)) *web.Mux {
	w := web.New()
	url := prefix + "/"

	// ROUTE TOP
	fmt.Println("[ROUTE]", url)
	for k, v := range routeConf {
		w.Post(url+k, v)
		// ROUTE PRINT
		space := strings.Repeat(" ", 30-len(k))
		vOf := reflect.ValueOf(v)
		fmt.Println("POST :", k, space, "->", runtime.FuncForPC(vOf.Pointer()).Name())
	}

	goji.Handle(url+"*", w)
	return w
}

/**************************************************************************************************/
/*!
 *  get routing
 */
/**************************************************************************************************/
func getRoute(prefix string, routeConf map[string]func(web.C, http.ResponseWriter, *http.Request)) *web.Mux {
	w := web.New()
	url := prefix + "/"

	// ROUTE TOP
	fmt.Println("[ROUTE]", url)
	for k, v := range routeConf {
		w.Get(url+k, v)
		// ROUTE PRINT
		space := strings.Repeat(" ", 30-len(k))
		vOf := reflect.ValueOf(v)
		fmt.Println("GET :", k, space, "->", runtime.FuncForPC(vOf.Pointer()).Name())
	}

	goji.Handle(url+"*", w)
	return w
}
