package adaptix

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func GetStringArg(args map[string]any, key string) (string, error) {
	v, ok := args[key].(string)
	if !ok {
		return "", fmt.Errorf("parameter '%s' must be set", key)
	}
	return v, nil
}

func CoerceFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func GetFloatArg(args map[string]any, key string) (float64, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, fmt.Errorf("parameter '%s' must be set", key)
	}
	f, ok := CoerceFloat64(v)
	if !ok {
		return 0, fmt.Errorf("parameter '%s' must be set", key)
	}
	return f, nil
}

func GetBoolArg(args map[string]any, key string) bool {
	v, _ := args[key].(bool)
	if v {
		return true
	}
	if f, ok := CoerceFloat64(args[key]); ok {
		return f != 0
	}
	if s, ok := args[key].(string); ok {
		s = strings.ToLower(strings.TrimSpace(s))
		return s == "true" || s == "1" || s == "yes"
	}
	return false
}

func GetIntArg(args map[string]any, key string) (int, error) {
	f, err := GetFloatArg(args, key)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

func GetFileArg(args map[string]any, key string) ([]byte, error) {
	v, ok := args[key].(string)
	if !ok {
		return nil, fmt.Errorf("parameter '%s' must be set", key)
	}
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, fmt.Errorf("parameter '%s': invalid base64: %w", key, err)
	}
	return data, nil
}

func GetStringArgDefault(args map[string]any, key string, defaultValue string) string {
	v, ok := args[key].(string)
	if !ok {
		return defaultValue
	}
	return v
}

func GetFloatArgDefault(args map[string]any, key string, defaultValue float64) float64 {
	if f, ok := CoerceFloat64(args[key]); ok {
		return f
	}
	return defaultValue
}

func MakeProxyTask(packData []byte, priority uint) TaskData {
	return TaskData{Type: TASK_TYPE_PROXY_DATA, Data: packData, Priority: priority, Sync: false}
}
