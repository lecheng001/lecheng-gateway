package util

import (
	"github.com/valyala/fasthttp"
	"time"
)

//func Cookie_Get(c echo.Context, key string) string {
//	cookie, err := c.Cookie(key)
//	if err != nil {
//		LogError(err)
//		return ""
//	}
//	return cookie.Value
//}
//
//func Cookie_Set(c echo.Context, key string, val string, expires time.Duration) {
//	cookie := new(http.Cookie)
//	cookie.Name = key
//	cookie.Value = val
//	if expires <= 0 {
//		cookie.Expires = time.Now().Add(time.Hour * 100000)
//	} else {
//		cookie.Expires = time.Now().Add(expires)
//	}
//	cookie.Path = "/"
//	c.SetCookie(cookie)
//}
func Cookie_SetFasthttp( ctx *fasthttp.RequestCtx,  key string, val string, expires time.Duration) {
 	t:=time.Now()
	if expires <= 0 {
		t = time.Now().Add(time.Hour * 100000)
	} else {
		t = time.Now().Add(expires)
	}

	cookie := fasthttp.Cookie{}
	cookie.SetKey(key)
	cookie.SetValue(val)
	cookie.SetPath("/")
	cookie.SetExpire(t)
	ctx.Response.Header.SetCookie(&cookie)
}
