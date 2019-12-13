package ldb

type FOR_EACH_KEY_FUNC func([]byte) bool
type FOR_EACH_FUNC func([]byte, []byte) bool
