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
	"net/http"
	. "sample/controller"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"fmt"
)

/**************************************************************************************************/
/*!
 *  稼働開始
 */
/**************************************************************************************************/
func Run() {
	// make route
	t := postRoute("/sample1", sampleRoute)
	t.Use(UseNewUser)
	t.Use(UseCommon)

	t2 := postRoute("/sample2", sampleRoute2)
	t2.Use(UseNewUser)
	t2.Use(UseCommon)

	w := getRoute("/web", webRoute)
	w.Use(UseCommon)

	// PORTを設定
	flag.Set("bind", ":9999")
	goji.Serve()
}

/**************************************************************************************************/
/*!
 *  リクエスト毎の共通事前処理
 */
/**************************************************************************************************/
func UseCommon(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// アクセス表示
		fmt.Println(r.RequestURI)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/**************************************************************************************************/
/*!
 *  AnalyzeType : NewUserの事前処理
 */
/**************************************************************************************************/
func UseNewUser(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		BodyAnalyze(c, r, AnalyzeType(OnNewUser))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/**************************************************************************************************/
/*!
 *  AnalyzeType : Loginの事前処理
 */
/**************************************************************************************************/
func UseLogin(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		BodyAnalyze(c, r, AnalyzeType(OnLogin))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/**************************************************************************************************/
/*!
 *  AnalyzeType : Defaultの事前処理
 */
/**************************************************************************************************/
func UseDefault(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		BodyAnalyze(c, r, AnalyzeType(OnDefault))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}