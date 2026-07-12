package adaptix

import (
	"encoding/base64"
	"fmt"
)

func GetStringArg(args map[string]any, key string) (string, error) {
	v, ok := args[key].(string)
	if !ok {
		return "", fmt.Errorf("parameter '%s' must be set", key)
	}
	return v, nil
}

func GetFloatArg(args map[string]any, key string) (float64, error) {
	v, ok := args[key].(float64)
	if !ok {
		return 0, fmt.Errorf("parameter '%s' must be set", key)
	}
	return v, nil
}

func GetBoolArg(args map[string]any, key string) bool {
	v, _ := args[key].(bool)
	return v
}

func GetIntArg(args map[string]any, key string) (int, error) {
	v, ok := args[key].(float64)
	if !ok {
		return 0, fmt.Errorf("parameter '%s' must be set", key)
	}
	return int(v), nil
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
	v, ok := args[key].(float64)
	if !ok {
		return defaultValue
	}
	return v
}

func MakeProxyTask(packData []byte, priority uint) TaskData {
	return TaskData{Type: TASK_TYPE_PROXY_DATA, Data: packData, Priority: priority, Sync: false}
}
