package clientid

import (
	"net/http"

	"dxkite.cn/explorer/src/core/utils"
	"dxkite.cn/log"
)

func Middleware(handler http.Handler, clientIdKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := ""
		if v, err := r.Cookie(clientIdKey); err == nil {
			clientId = v.Value
		}

		if clientId == "" {
			addr := utils.GetRemoteIp(r)
			ua := r.UserAgent()
			log.Info("client", addr, ua)
			clientId := utils.Md5(addr + "-" + ua)
			cookie := &http.Cookie{
				Name:  clientIdKey,
				Value: clientId,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}

		r.Header.Set(clientIdKey, clientId)
		handler.ServeHTTP(w, r)
	})
}
