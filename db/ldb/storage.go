package ldb

type Storage interface {
	Close()
	Set([]byte, []byte)
	Get([]byte) []byte
	Has([]byte) bool
	Del([]byte)
	Total(prefix []byte) int64
	List(prefix []byte, limit int, offset int, RemovePrefix bool) [][]byte
	ForEach(prefix []byte, RemovePrefix bool, fn FOR_EACH_FUNC)
	ForEachKey(prefix []byte, limit int, offset int, RemovePrefix bool, fn FOR_EACH_KEY_FUNC)
}
