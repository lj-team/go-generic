package allow

import (
	"net/http"
)

func Default(rw http.ResponseWriter) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, x-custom-parameter, content-type, x-request-id, authorization, accept-language, accept-charset, pragma, user-agent")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE, HEAD")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
}

func Cookies(rw http.ResponseWriter, origin string) {
	rw.Header().Set("Access-Control-Allow-Origin", origin)
	rw.Header().Set("Access-Control-Allow-Headers", "accept, x-custom-parameter, content-type, x-request-id, authorization, accept-language, accept-charset, pragma, user-agent, cookie, set-cookie")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE, HEAD")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
}
