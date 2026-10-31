package main

import (
	"fmt"
	"io"
	"net"
	keyvalues "server/keyValues"
	"strings"
)

func main() {
	kv := keyvalues.CreateKeyValue()
	listen, err := net.Listen("tcp", ":6140")

	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}

	defer listen.Close()
	fmt.Println("Listening to port :6140")

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println("Failed to load data:", err)
		return
	}
	defer aof.CloseAof()

	aof.ReadAof(func(value Value) {
		parts, err := valueToParts(value)
		if err != nil {
			fmt.Println("Invalid AOF command:", err)
			return
		}

		response, _ := dispatch(parts, kv)

		fmt.Println("Restored:", parts, response)
	})

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Client disconnected. ERROR:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr(), conn.LocalAddr())
		go func() {
			handleConnection(conn, kv, aof)
		}()
	}
}

func valueToParts(v Value) ([]string, error) {
	if v.typ != "array" {
		return nil, fmt.Errorf("command must be an array")
	}

	if len(v.arr) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	parts := make([]string, len(v.arr))

	for i, val := range v.arr {
		if val.typ != "bulk" {
			return nil, fmt.Errorf(
				"argument %d must be a bulk string", i,
			)
		}
		parts[i] = val.bulk
	}

	return parts, nil
}

func handleConnection(conn net.Conn, kv *keyvalues.KeyValues, aof *Aof) {
	defer conn.Close()

	resp := NewRESP(conn)
	writer := NewWriter(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			return
		}

		parts, err := valueToParts(value)
		if err != nil {
			errResp := Value{
				typ: "error",
				str: "invalid command format",
			}

			if err := writeResponse(writer, errResp); err != nil {
				return
			}
			continue
		}

		switch strings.ToUpper(parts[0]) {
		case "PING":
			if err := writeResponse(writer, Value{
				typ: "string",
				str: "PONG",
			}); err != nil {
				return
			}
			continue

		case "QUIT":
			if err := writeResponse(writer, Value{
				typ: "string",
				str: "OK",
			}); err != nil {
				return
			}
			return
		}

		if err := validateCommand(parts); err != nil {
			errResp := Value{
				typ: "error",
				str: err.Error(),
			}

			if err := writeResponse(writer, errResp); err != nil {
				return
			}
			continue
		}

		fmt.Printf("[%s] %s\n",
			conn.RemoteAddr(), strings.Join(parts, " "))

		response, write_methods := dispatch(parts, kv)
		fmt.Println(write_methods)
		allowed_methods := map[string]bool{
			"SET":  true,
			"DEL":  true,
			"HSET": true,
			"HDEL": true,
		}

		if allowed_methods[write_methods] {
			if err := aof.WriteAof(value); err != nil {
				fmt.Println("AOF write error:", err)
			}
		}

		if err := writeResponse(writer, response); err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}

func writeResponse(w *WRITE, response Value) error {
	return w.Write(response)
}
