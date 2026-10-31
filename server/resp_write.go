package main

import (
	"fmt"
	"io"
	"strconv"
)

// Marshal section

func (v Value) Marshall() []byte {
	switch v.typ {
	case "array":
		return v.marshallArray()
	case "bulk":
		return v.marshallBulk()
	case "string":
		return v.marshallString()
	case "integer":
		return v.marshallInteger()
	case "error":
		return v.marshallError()
	case "null":
		return v.marshallNull()
	default:
		return []byte{}
	}
}

func (v Value) marshallArray() []byte {
	len := len(v.arr)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.arr[i].Marshall()...)
	}
	return bytes
}

func (v Value) marshallBulk() []byte {
	bulk := v.bulk
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallString() []byte {
	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallInteger() []byte {
	return []byte(fmt.Sprintf(":%d\r\n", v.num))
}

func (v Value) marshallError() []byte {
	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

// Wrtie section

type WRITE struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *WRITE {
	return &WRITE{writer: w}
}

func (w *WRITE) Write(v Value) error {
	var bytes = v.Marshall()
	for len(bytes) > 0 {
		n, err := w.writer.Write(bytes)
		if err != nil {
			return err
		}

		if n == 0 {
			return io.ErrShortWrite
		}

		bytes = bytes[n:]
	}

	return nil
}
