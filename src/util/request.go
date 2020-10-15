package util

import (
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"strings"
)
//
//func Request_FormGet(ctx echo.Context, key string) string {
//	tmp :=ctx.QueryParam(key)
//	if tmp==""{
//		tmp=ctx.FormValue(key)
//	}
//	return tmp
//}
//
//func Request_FormGetInt(ctx echo.Context, key string) int64 {
//	tmpint:=ctx.QueryParam(key)
//	if tmpint==""{
//		tmpint=ctx.FormValue(key)
//	}
//	return String_GetInt64(tmpint)
//}


func Request_GetIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("Remote_addr"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

//取得客户端IP
func Request_GetIP_Fasthttp(ctx *fasthttp.RequestCtx) string {
	RemoteAddr := ctx.RemoteIP().String()
	X_Forwarded_For := string(ctx.Request.Header.Peek("X-Forwarded-For"))
	X_Real_Ip := string(ctx.Request.Header.Peek("X-Real-Ip"))
	return ClientPublicIP(X_Forwarded_For, X_Real_Ip, RemoteAddr)
}

////取得客户端IP
//func Request_GetIP_ECHO(ctx echo.Context) string {
//	RemoteAddr := ctx.RealIP()
//	X_Forwarded_For := string(ctx.Request().Header.Get("X-Forwarded-For"))
//	X_Real_Ip := string(ctx.Request().Header.Get("X-Real-Ip"))
//	return ClientPublicIP(X_Forwarded_For, X_Real_Ip, RemoteAddr)
//}

//取得客户端IP
func ClientPublicIP(X_Forwarded_For, X_Real_Ip, RemoteAddr string) string {
	var ip string
	//var tmp string
	for _, ip2 := range strings.Split(X_Forwarded_For, ",") {
		ip2 = strings.TrimSpace(ip2)
		if ip2 != "" && !HasLocalIPddr(ip2) {
			ip = ip2
			//if tmp == "" {
			//	tmp = ip
			//	continue
			//}
			////else {
			////	return ip
			////}
			//return ip
		}
	}

	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(X_Real_Ip)
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}

	if ip2, _, err := net.SplitHostPort(strings.TrimSpace(RemoteAddr)); err == nil {
		ip = ip2
		if !HasLocalIPddr(ip2) {
			return ip2
		}
	}

	return ip
}

// HasLocalIPddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPddr(ip string) bool {
	return !IsPublicIP(net.ParseIP(ip))
}

// IsPublicIP 检测 IP 地址是否是内网地址
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
