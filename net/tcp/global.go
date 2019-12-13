package tcp

import (
	"time"
)

var (
	KEEP_ALIVE              = time.Second * 30
	RECONNECT_AFTER_SECONDS = int64(3)
)
