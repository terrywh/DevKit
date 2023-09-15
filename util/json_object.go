package util

import (
	"fmt"
	"strconv"
	"strings"
)

type JSONObject map[string]interface{}

func (o JSONObject) Get(key string) interface{} {
	return o.get(strings.Split(key, "."), map[string]interface{}(o))
}

func (o JSONObject) GetString(key string) string {
	return fmt.Sprint(o.get(strings.Split(key, "."), map[string]interface{}(o)))
}

func (o JSONObject) get(key []string, src interface{}) interface{} {
	if len(key) > 0 {
		if ctr, ok := src.(map[string]interface{}); ok {
			return o.get(key[1:], ctr[key[0]])
		} else if ctr, ok := src.([]interface{}); ok {
			idx, _ := strconv.Atoi(key[0])
			return o.get(key[1:], ctr[idx])
		} else {
			return nil
		}
	} else {
		return src
	}
}