package clientid

import (
	"crypto/md5"
	"encoding/base64"
	"net"
	"net/http"

	"dxkite.cn/log"
)

var remoteIpFrom = []string{
	"http-x-real-ip",
	"http-client-ip",
	"http-x-forwarded-for",
	"http-x-forwarded",
	"http-x-cluster-client-ip",
	"http-forwarded-for",
	"http-forwarded",
	"remote-addr",
}

const ClientIdKey = "client-id"

func getRemoteAddr(r *http.Request) string {
	addr := r.RemoteAddr
	for _, v := range remoteIpFrom {
		vv := r.Header.Get(v)
		if vv == "" {
			continue
		}
		if v := net.ParseIP(vv); v != nil {
			return vv
		}
	}

	h, _, _ := net.SplitHostPort(addr)
	return h
}

func getMd5(val string) string {
	m := md5.New()
	m.Write([]byte(val))
	b := m.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(b)
}

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := ""
		if v, err := r.Cookie(ClientIdKey); err == nil {
			clientId = v.Value
		}

		if clientId == "" {
			addr := getRemoteAddr(r)
			ua := r.UserAgent()
			log.Info("client", addr, ua)
			clientId := getMd5(addr + "-" + ua)
			cookie := &http.Cookie{
				Name:  ClientIdKey,
				Value: clientId,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}

		r.Header.Set(ClientIdKey, clientId)
		handler.ServeHTTP(w, r)
	})
}
