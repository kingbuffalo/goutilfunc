package goutilfunc

import (
	"encoding/json"
	"github.com/kingbuffalo/seelog"
	"math/rand"
	"runtime"
)

const RAND_CHAR_ARR_LEN = 62

var charArr [RAND_CHAR_ARR_LEN]byte

func init() {
	idx := 0
	for i := 0; i < 10; i++ {
		charArr[idx] = byte(48 + i)
		idx++
	}
	for i := 0; i < 26; i++ {
		charArr[idx] = byte(65 + i)
		idx++
	}
	for i := 0; i < 26; i++ {
		charArr[idx] = byte(97 + i)
		idx++
	}
}

func PanicJsonUnmarshalError(b []byte, err error) {
	se, ok := err.(*json.SyntaxError)
	if ok {
		if se != nil {
			s := se.Offset - 30
			if s < 0 {
				s = 0
			}
			e := se.Offset + 30
			if int(e) > len(b) {
				e = int64(len(b))
			}
			errMsg := string(b[s:e])
			panic(errMsg)
		}
	} else {
		panic(err)
	}
}

func RecoverPanic() {
	if err := recover(); err != nil {
		_ = seelog.Error(err)
		buf := make([]byte, 1<<11)
		runtime.Stack(buf, true)
		_ = seelog.Error(string(buf))
	}
}

func GenRandStr(strLen int) string {
	ret := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		rIdx := rand.Intn(RAND_CHAR_ARR_LEN)
		ret[i] = charArr[rIdx]
	}
	s := string(ret)
	return s
}
