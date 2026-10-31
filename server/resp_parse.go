package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	INTEGER = ':'
	BULK    = '$'
	ERROR   = '-'
	ARRAY   = '*'
)

// Note:
/*
1. typ is used to determine the data type of Value holds
2. str used to hold the value of stringreceived from simple string
3. num to hold the value of integer
4. bulk to hold the bulk string
5. arr hold the array of values
*/
type Value struct {
	typ  string
	str  string
	num  int
	bulk string
	arr  []Value
}

type RESP struct {
	reader *bufio.Reader
}

func NewRESP(rd io.Reader) *RESP {
	return &RESP{reader: bufio.NewReader(rd)}
}

func (r *RESP) readLine() (line []byte, n int, err error) {
	for {
		a, err := r.reader.ReadByte()

		if err != nil {
			return nil, n, err
		}

		n += 1
		line = append(line, a)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, err
}

func (r *RESP) readInteger() (len int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	len_64_bit, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(len_64_bit), n, err
}

func (r *RESP) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}

}

func (r *RESP) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.arr = make([]Value, len)

	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.arr[i] = val
	}
	return v, nil

}
func (r *RESP) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readInteger()

	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)
	v.bulk = string(bulk)

	r.readLine()
	return v, nil
}
