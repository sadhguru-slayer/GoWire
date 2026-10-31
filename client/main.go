package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func EncodeCommand(command string) string {
	parts := strings.Fields(command)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))

	for _, part := range parts {
		builder.WriteString(fmt.Sprintf("$%d\r\n", len(part)))
		builder.WriteString(part)
		builder.WriteString("\r\n")
	}

	return builder.String()
}

func printValue(v Value) {
	switch v.typ {
	case "string":
		fmt.Println(v.str)

	case "integer":
		fmt.Println(v.num)

	case "error":
		fmt.Println("ERR:", v.str)

	case "bulk":
		fmt.Println(v.bulk)

	case "array":
		for _, item := range v.arr {
			printValue(item)
		}

	default:
		fmt.Printf("%+v\n", v)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:6140")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to localhost:6140")

	input := bufio.NewScanner(os.Stdin)

	// Reuse one decoder for the lifetime of this connection.
	resp := NewReader(conn)

	for {
		fmt.Print("> ")

		if !input.Scan() {
			break
		}

		command := input.Text()
		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		encoded := EncodeCommand(command)

		if _, err := io.WriteString(conn, encoded); err != nil {
			fmt.Println("Write error:", err)
			return
		}

		// Decode exactly one complete RESP response.
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Response error:", err)
			return
		}

		printValue(value)

		if strings.EqualFold(parts[0], "QUIT") {
			return
		}
	}

	if err := input.Err(); err != nil {
		fmt.Println("Input error:", err)
	}
}
