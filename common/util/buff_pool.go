package util

import (
	"bytes"
	"sync"
)

//---------------------------------------
var buffPool sync.Pool

func Borrow() *bytes.Buffer {
	var buffer *bytes.Buffer
	item := buffPool.Get()
	if item == nil {
		var byteSlice []byte
		byteSlice = make([]byte, 0, 10*1024)
		buffer = bytes.NewBuffer(byteSlice)

	} else {
		buffer = item.(*bytes.Buffer)
	}
	return buffer
}

func Return(buffer *bytes.Buffer) {
	buffer.Reset()
	buffPool.Put(buffer)
}
