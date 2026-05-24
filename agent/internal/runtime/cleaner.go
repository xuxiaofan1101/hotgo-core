package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"vogo-agent/internal/protocol"
)

type CleanResult struct {
	Payload map[string]any
	Dropped bool
}

func ApplyCleanPipeline(payload map[string]any, cleanConfig protocol.CleanConfig) (CleanResult, error) {
	active := cleanConfig.Enabled || len(cleanConfig.Steps) > 0 || configBool(cleanConfig.Filter, "enabled")
	if !active {
		return CleanResult{Payload: clonePayload(payload)}, nil
	}
	matched, err := filterMatches(payload, cleanConfig.Filter)
	if err != nil {
		return CleanResult{}, err
	}
	if !matched {
		return CleanResult{Payload: clonePayload(payload), Dropped: true}, nil
	}
	cleaned := clonePayload(payload)
	for _, step := range cleanConfig.Steps {
		stepType := strings.ToLower(strings.TrimSpace(configString(step, "type")))
		switch stepType {
		case "", "noop":
			continue
		case "rename":
			if value, ok := getPath(cleaned, configString(step, "from")); ok {
				setPath(cleaned, configString(step, "to"), value)
				deletePath(cleaned, configString(step, "from"))
			}
		case "remove", "drop_field":
			deletePath(cleaned, cleanField(step))
		case "default":
			field := cleanField(step)
			if _, ok := getPath(cleaned, field); !ok {
				setPath(cleaned, field, step["value"])
			}
		case "set_field":
			setPath(cleaned, cleanField(step), step["value"])
		case "to_string":
			transformPath(cleaned, cleanField(step), func(value any) any { return fmt.Sprint(value) })
		case "to_int":
			transformPath(cleaned, cleanField(step), func(value any) any { return toInt64(value) })
		case "to_float":
			transformPath(cleaned, cleanField(step), func(value any) any { return toFloat64(value) })
		case "to_bool":
			transformPath(cleaned, cleanField(step), func(value any) any { return toBool(value) })
		case "trim":
			transformPath(cleaned, cleanField(step), func(value any) any { return strings.TrimSpace(fmt.Sprint(value)) })
		case "regex_replace":
			if err := regexReplace(cleaned, step); err != nil {
				return CleanResult{}, err
			}
		case "parse_json":
			if err := parseJSONField(cleaned, cleanField(step)); err != nil {
				return CleanResult{}, err
			}
		case "mask":
			maskField(cleaned, step)
		case "drop_if":
			matched, err := filterConditionMatches(cleaned, step)
			if err != nil {
				return CleanResult{}, err
			}
			if matched {
				return CleanResult{Payload: cleaned, Dropped: true}, nil
			}
		default:
			return CleanResult{}, fmt.Errorf("不支持的数据清洗步骤: %s", stepType)
		}
	}
	return CleanResult{Payload: cleaned}, nil
}

func cleanField(step map[string]any) string {
	field := configString(step, "field")
	if field == "" {
		field = configString(step, "path")
	}
	return field
}

func filterMatches(payload map[string]any, filter map[string]any) (bool, error) {
	if len(filter) == 0 || !configBool(filter, "enabled") {
		return true, nil
	}
	results := make([]bool, 0)
	for _, condition := range mapSlice(filter, "conditions") {
		matched, err := filterConditionMatches(payload, condition)
		if err != nil {
			return false, err
		}
		results = append(results, matched)
	}
	for _, group := range mapSlice(filter, "groups") {
		matched, err := filterMatches(payload, group)
		if err != nil {
			return false, err
		}
		results = append(results, matched)
	}
	if len(results) == 0 {
		return true, nil
	}
	return combineResults(configString(filter, "logic"), results), nil
}

func filterConditionMatches(payload map[string]any, condition map[string]any) (bool, error) {
	operator := normalizeOperator(configString(condition, "operator"))
	field := configString(condition, "field", "path")
	value, exists := getPath(payload, field)
	expected := condition["value"]
	switch operator {
	case "eq":
		return fmt.Sprint(value) == fmt.Sprint(expected), nil
	case "ne":
		return fmt.Sprint(value) != fmt.Sprint(expected), nil
	case "contains":
		return strings.Contains(fmt.Sprint(value), fmt.Sprint(expected)), nil
	case "not_contains":
		return !strings.Contains(fmt.Sprint(value), fmt.Sprint(expected)), nil
	case "matches":
		re, err := regexp.Compile(fmt.Sprint(expected))
		if err != nil {
			return false, err
		}
		return re.MatchString(fmt.Sprint(value)), nil
	case "empty":
		return !exists || strings.TrimSpace(fmt.Sprint(value)) == "", nil
	case "not_empty":
		return exists && strings.TrimSpace(fmt.Sprint(value)) != "", nil
	case "gt", "gte", "lt", "lte":
		left := toFloat64(value)
		right := toFloat64(expected)
		switch operator {
		case "gt":
			return left > right, nil
		case "gte":
			return left >= right, nil
		case "lt":
			return left < right, nil
		default:
			return left <= right, nil
		}
	case "between":
		rangeValues := anySlice(expected)
		if len(rangeValues) < 2 {
			return false, nil
		}
		current := toFloat64(value)
		return current >= toFloat64(rangeValues[0]) && current <= toFloat64(rangeValues[1]), nil
	case "in", "not_in":
		found := false
		for _, item := range anySlice(expected) {
			if fmt.Sprint(item) == fmt.Sprint(value) {
				found = true
				break
			}
		}
		if operator == "not_in" {
			return !found, nil
		}
		return found, nil
	default:
		return false, fmt.Errorf("不支持的数据过滤操作符: %s", operator)
	}
}

func normalizeOperator(operator string) string {
	switch strings.ToLower(strings.TrimSpace(operator)) {
	case "=", "==", "equals":
		return "eq"
	case "!=", "<>", "not_equals":
		return "ne"
	case ">":
		return "gt"
	case ">=":
		return "gte"
	case "<":
		return "lt"
	case "<=":
		return "lte"
	default:
		return strings.ToLower(strings.TrimSpace(operator))
	}
}

func combineResults(logic string, results []bool) bool {
	if strings.EqualFold(strings.TrimSpace(logic), "or") {
		for _, item := range results {
			if item {
				return true
			}
		}
		return false
	}
	for _, item := range results {
		if !item {
			return false
		}
	}
	return true
}

func regexReplace(payload map[string]any, step map[string]any) error {
	re, err := regexp.Compile(configString(step, "pattern"))
	if err != nil {
		return err
	}
	replacement := configString(step, "replacement")
	transformPath(payload, cleanField(step), func(value any) any {
		return re.ReplaceAllString(fmt.Sprint(value), replacement)
	})
	return nil
}

func parseJSONField(payload map[string]any, field string) error {
	value, ok := getPath(payload, field)
	if !ok {
		return nil
	}
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "" {
		return nil
	}
	var parsed any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return err
	}
	setPath(payload, field, parsed)
	return nil
}

func maskField(payload map[string]any, step map[string]any) {
	mode := strings.ToLower(strings.TrimSpace(configString(step, "mode")))
	replacement := configString(step, "replacement")
	if replacement == "" {
		replacement = "***"
	}
	transformPath(payload, cleanField(step), func(value any) any {
		if mode == "hash" {
			sum := sha256.Sum256([]byte(fmt.Sprint(value)))
			return hex.EncodeToString(sum[:])
		}
		return replacement
	})
}

func getPath(payload map[string]any, path string) (any, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, false
	}
	current := any(payload)
	for _, part := range strings.Split(path, ".") {
		obj, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = obj[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func setPath(payload map[string]any, path string, value any) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	parts := strings.Split(path, ".")
	current := payload
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)
		if !ok {
			next = map[string]any{}
			current[part] = next
		}
		current = next
	}
	current[parts[len(parts)-1]] = value
}

func deletePath(payload map[string]any, path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	parts := strings.Split(path, ".")
	current := payload
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)
		if !ok {
			return
		}
		current = next
	}
	delete(current, parts[len(parts)-1])
}

func transformPath(payload map[string]any, path string, transform func(any) any) {
	value, ok := getPath(payload, path)
	if !ok {
		return
	}
	setPath(payload, path, transform(value))
}

func clonePayload(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	data, err := json.Marshal(input)
	if err != nil {
		output := make(map[string]any, len(input))
		for key, value := range input {
			output[key] = value
		}
		return output
	}
	var output map[string]any
	if err := json.Unmarshal(data, &output); err != nil {
		return map[string]any{}
	}
	return output
}

func mapSlice(config map[string]any, key string) []map[string]any {
	value, ok := config[key]
	if !ok || value == nil {
		return nil
	}
	switch typed := value.(type) {
	case []map[string]any:
		return typed
	case []any:
		result := make([]map[string]any, 0, len(typed))
		for _, item := range typed {
			if mapped, ok := item.(map[string]any); ok {
				result = append(result, mapped)
			}
		}
		return result
	default:
		return nil
	}
}

func anySlice(value any) []any {
	switch typed := value.(type) {
	case []any:
		return typed
	case []string:
		result := make([]any, 0, len(typed))
		for _, item := range typed {
			result = append(result, item)
		}
		return result
	default:
		return []any{typed}
	}
}

func toInt64(value any) int64 {
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int64:
		return typed
	case float64:
		return int64(typed)
	case json.Number:
		parsed, _ := typed.Int64()
		return parsed
	default:
		parsed, _ := strconv.ParseInt(strings.TrimSpace(fmt.Sprint(value)), 10, 64)
		return parsed
	}
}

func toFloat64(value any) float64 {
	switch typed := value.(type) {
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case float64:
		return typed
	case json.Number:
		parsed, _ := typed.Float64()
		return parsed
	default:
		parsed, _ := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(value)), 64)
		return parsed
	}
}

func toBool(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		parsed, _ := strconv.ParseBool(strings.TrimSpace(typed))
		return parsed
	default:
		return fmt.Sprint(value) == "1"
	}
}
