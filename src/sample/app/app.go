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
	"flag"
	"math/rand"
	"net/http"
	. "sample/controller"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

/**************************************************************************************************/
/*!
 *  稼働開始
 */
/**************************************************************************************************/
func Run() {

	rand.Seed(time.Now().UnixNano())

	// make route
	t := postRoute("/sample1", sampleRoute)
	t.Use(UseDefault)

	t2 := postRoute("/sample2", sampleRoute2)
	t2.Use(UserUserId)

	getRoute("/web", webRoute)

	// PORTを設定
	flag.Set("bind", ":9999")
	goji.Serve()
}

/**************************************************************************************************/
/*!
 *  AnalyzeType : NewUserの事前処理
 */
/**************************************************************************************************/
func UseDefault(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		BodyAnalyze(c, r, AnalyzeType(OnDefault))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/**************************************************************************************************/
/*!
 *  AnalyzeType : Defaultの事前処理
 */
/**************************************************************************************************/
func UserUserId(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		BodyAnalyze(c, r, AnalyzeType(OnUserId))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
