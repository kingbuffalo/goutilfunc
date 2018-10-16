package lru

import (
	"github.com/kingbuffalo/goutilfunc"
	"math/rand"
	"testing"
	"time"
)

func Test_normalUse(t *testing.T) {
	lru := NewStrLRU(10)
	lru.GetValue("abc")
	lru.AddValue("abc", []byte{1, 2})
	b := lru.GetValue("abc").([]byte)
	if b[0] != 1 || b[1] != 2 {
		t.Error("error")
	}
	lru.RmValue("abc")
	if lru.GetValue("abc") != nil {
		t.Error("error")
	}
	lru.AddValue("cba", []byte{2, 2})
	lru.AddValue("cba", []byte{3, 3})

	b = lru.GetValue("cba").([]byte)
	if b[0] != 3 || b[1] != 3 {
		t.Log(b)
		t.Error("error")
	}
	lru.RmValue("cba")

	strArr := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	for i, v := range strArr {
		lru.AddValue(v, []byte{byte(i)})
	}
	var iter *StrLRUNode
	lrunode := lru.RangeReverse(iter)
	var idx byte = 0
	for lrunode != nil {
		iter = lrunode
		d := iter.data.([]byte)

		if d[0] != idx {
			t.Error("error")
		}

		lrunode = lru.RangeReverse(iter)
		idx++
	}
	lru.AddValue("a", 1000)
	v, ok := lru.GetValue("a").(int)
	if !ok {
		t.Error("set error")
	}
	if v != 1000 {
		t.Error("set error")
	}
}

func genRandByte() []byte {
	n := rand.Int31()
	b1 := byte(n & 0xff)
	b2 := byte(n >> 8 & 0xff)
	b3 := byte(n >> 16 & 0xff)
	b4 := byte(n >> 24 & 0xff)
	return []byte{b1, b2, b3, b4}
}

func Benchmark_getset(b *testing.B) {
	lru := NewStrLRU(10)
	rand.Seed(time.Now().Unix())
	for i := 0; i < b.N; i++ {
		key := goutilfunc.GenRandStr(10)
		value := genRandByte()
		lru.AddValue(key, value)

		gkey := goutilfunc.GenRandStr(10)
		lru.GetValue(gkey)
	}
}
