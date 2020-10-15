/**
开发人员：陈朝能
 */
package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Render struct {
}

func (_ Render) RendErr(ctx *fasthttp.RequestCtx, msg string) {
	tmp, err := json.Marshal(map[string]interface{}{"errcode": -1, "msg": msg, "content": nil})
	if err != nil {
		fmt.Fprint(ctx, err)

	}
 	ctx.Error(string(tmp),200)
}

func (_ Render) RendJson(ctx *fasthttp.RequestCtx, msg string) {
	tmp, err := json.Marshal(map[string]interface{}{"errcode": -1, "msg": msg, "content": nil})
	if err != nil {
		fmt.Fprint(ctx, err)

	}
	ctx.Response.SetBodyString(string(tmp))
}


func (_ Render) RendByte(ctx *fasthttp.RequestCtx, content []byte) {

	ctx.Response.SetBodyString(string(content))
}

func (_ Render) RendMap(ctx *fasthttp.RequestCtx, content  interface{}) {
	result,_:= json.Marshal(content)
	ctx.Response.SetBody(result)
}