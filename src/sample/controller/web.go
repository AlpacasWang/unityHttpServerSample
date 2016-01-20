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
	"net/http"
	"github.com/zenazn/goji/web"
)

func Ping(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte("pong"))
}
