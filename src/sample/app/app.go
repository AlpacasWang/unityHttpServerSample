package app

/**************************************************************************************************/
/*!
 *  app.go
 *
 *  アプリエントリポイント
 *
 */
/**************************************************************************************************/
import (
	"math/rand"
	. "sample/controller"
	"time"

	"github.com/labstack/echo"
)

/**************************************************************************************************/
/*!
 *  稼働開始
 */
/**************************************************************************************************/
func Run() {

	rand.Seed(time.Now().UnixNano())

	ec := echo.New()

	// make route
	postRoute(ec, "/sample1", sampleRoute, apiAnalyze(OnDefault))
	postRoute(ec, "/sample2", sampleRoute2, apiAnalyze(OnUserId))

	getRoute(ec, "/web", webRoute)

	// PORTを設定
	ec.Logger.Fatal(ec.Start(":9999"))
}

func apiAnalyze(analyzeType AnalyzeType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			BodyAnalyze(c, analyzeType)
			return next(c)
		}
	}
}
