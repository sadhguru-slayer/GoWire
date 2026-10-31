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

func NewReader(rd io.Reader) *RESP {
	return &RESP{
		reader: bufio.NewReader(rd),
	}
}

// Reads a line terminated by CRLF.
func (r *RESP) readLine() ([]byte, error) {
	var line []byte

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, err
		}

		line = append(line, b)

		if len(line) >= 2 &&
			line[len(line)-2] == '\r' &&
			line[len(line)-1] == '\n' {
			return line[:len(line)-2], nil
		}
	}
}

func (r *RESP) readInteger() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %w", line, err)
	}

	return n, nil
}

// Reads one complete RESP value.
func (r *RESP) Read() (Value, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case STRING:
		line, err := r.readLine()
		if err != nil {
			return Value{}, err
		}

		return Value{
			typ: STRING_NAME,
			str: string(line),
		}, nil

	case ERROR:
		line, err := r.readLine()
		if err != nil {
			return Value{}, err
		}

		return Value{
			typ: ERROR_NAME,
			str: string(line),
		}, nil

	case INTEGER:
		num, err := r.readInteger()
		if err != nil {
			return Value{}, err
		}

		return Value{
			typ: INTEGER_NAME,
			num: num,
		}, nil

	case BULK:
		return r.readBulk()

	case ARRAY:
		return r.readArray()

	default:
		return Value{}, fmt.Errorf(
			"unsupported RESP type: %q", typ,
		)
	}
}

func (r *RESP) readBulk() (Value, error) {
	length, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	// RESP2 null bulk string: $-1\r\n
	if length == -1 {
		return Value{
			typ:  BULK_NAME,
			bulk: "",
		}, nil
	}

	if length < -1 {
		return Value{}, fmt.Errorf("invalid bulk length: %d", length)
	}

	bulk := make([]byte, length)

	if _, err := io.ReadFull(r.reader, bulk); err != nil {
		return Value{}, err
	}

	var terminator [2]byte
	if _, err := io.ReadFull(r.reader, terminator[:]); err != nil {
		return Value{}, err
	}

	if terminator != [2]byte{'\r', '\n'} {
		return Value{}, fmt.Errorf("invalid bulk string terminator")
	}

	return Value{
		typ:  BULK_NAME,
		bulk: string(bulk),
	}, nil
}

func (r *RESP) readArray() (Value, error) {
	length, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	if length < 0 {
		return Value{}, fmt.Errorf("invalid array length: %d", length)
	}

	v := Value{
		typ: ARRAY_NAME,
		arr: make([]Value, length),
	}

	for i := 0; i < length; i++ {
		item, err := r.Read()
		if err != nil {
			return Value{}, err
		}

		v.arr[i] = item
	}

	return v, nil
}

const (
	STRING_NAME  = "string"
	INTEGER_NAME = "integer"
	ERROR_NAME   = "error"
	BULK_NAME    = "bulk"
	ARRAY_NAME   = "array"
)
