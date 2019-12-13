package errors

import (
	"fmt"
	"net/http"
)

var http_code_tab map[int]int = map[int]int{
	400: http.StatusBadRequest,
	401: http.StatusUnauthorized,
	403: http.StatusForbidden,
	404: http.StatusNotFound,
	405: http.StatusMethodNotAllowed,
	409: http.StatusConflict,
	429: http.StatusTooManyRequests,
	500: http.StatusInternalServerError,
	502: http.StatusBadGateway,
	504: http.StatusGatewayTimeout,
}

var http_msg_tab map[int]string = map[int]string{
	400: "400 Bad Request",
	401: "401 Unauthorized",
	403: "403 Forbidden",
	404: "404 Status Not Found",
	405: "405 Method Not Allowed",
	409: "409 Conflict",
	429: "429 Too Many Requests",
	500: "500 Internal Server Error",
	502: "502 Bad Gateway",
	504: "504 Gateway Timeout",
}

func Send(rw http.ResponseWriter, code int) {

	c, h := http_code_tab[code]

	if !h {
		code = 500
		c = http.StatusInternalServerError
	}

	m, _ := http_msg_tab[code]

	rw.WriteHeader(c)
	fmt.Fprint(rw, m)
}
