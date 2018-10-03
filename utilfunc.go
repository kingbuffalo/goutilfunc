package goutilfunc

import (
	"encoding/json"
)

func PanicJsonUnmarshalError(b []byte, err error) {
	se := err.(*json.SyntaxError)
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
