package goutilfunc

import (
	"encoding/json"
	"github.com/kingbuffalo/seelog"
	"runtime"
)

func PanicJsonUnmarshalError(b []byte, err error) {
	se := err.(*json.SyntaxError)
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
}

func RecoverPanic() {
	if err := recover(); err != nil {
		_ = seelog.Error(err)
		buf := make([]byte, 1<<11)
		runtime.Stack(buf, true)
		_ = seelog.Error(string(buf))
	}
}
