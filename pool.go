package sb

import (
	"bytes"
	"sync"
)

const bufferSize = 255 // enough for most SQLs

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, bufferSize))
	},
}
