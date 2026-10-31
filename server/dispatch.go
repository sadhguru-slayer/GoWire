package main

import (
	"strings"

	keyvalues "server/keyValues"
)

func dispatch(parts []string, kv *keyvalues.KeyValues) (Value, string) {
	switch strings.ToUpper(parts[0]) {
	case "SET":
		kv.Set(parts[1], parts[2])
		return Value{
			typ: "string",
			str: "OK",
		}, "SET"

	case "GET":
		value, exists := kv.Get(parts[1])
		if !exists {
			return Value{
				typ: "null",
			}, ""
		}

		return Value{
			typ:  "bulk",
			bulk: value,
		}, "GET"

	case "DEL":
		if kv.Del(parts[1]) {
			return Value{
				typ: "integer",
				num: 1,
			}, "DEL"
		}

		return Value{
			typ: "integer",
			num: 0,
		}, ""

	case "EXISTS":
		if kv.Exists(parts[1]) {
			return Value{
				typ: "integer",
				num: 1,
			}, ""
		}

		return Value{
			typ: "integer",
			num: 0,
		}, "EXISTS"

	case "HSET":
		created := kv.HSet(parts[1], parts[2], parts[3])

		if created {
			return Value{
				typ: "integer",
				num: 1,
			}, "HSET"
		}

		return Value{
			typ: "integer",
			num: 0,
		}, ""

	case "HGET":
		value, exists := kv.HGet(parts[1], parts[2])

		if !exists {
			return Value{
				typ: "null",
			}, ""
		}

		return Value{
			typ:  "bulk",
			bulk: value,
		}, "HGET"

	case "HDEL":
		deleted := kv.HDel(parts[1], parts[2])

		if deleted {
			return Value{
				typ: "integer",
				num: 1,
			}, "HDEL"
		}

		return Value{
			typ: "integer",
			num: 0,
		}, ""

	case "HEXISTS":
		exists := kv.HExists(parts[1], parts[2])

		if exists {
			return Value{
				typ: "integer",
				num: 1,
			}, ""
		}

		return Value{
			typ: "integer",
			num: 0,
		}, "HEXISTS"
	}

	return Value{
		typ: "error",
		str: "unknown command",
	}, ""
}
