package utils

import (
	"crypto/md5"
	"encoding/base64"
	"net"
	"net/http"
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

func GetRemoteIp(r *http.Request) string {
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

func Md5(val string) string {
	m := md5.New()
	m.Write([]byte(val))
	b := m.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(b)
}
