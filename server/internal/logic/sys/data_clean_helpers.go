package sys

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"hotgo/internal/consts"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

const dataCleanSecretMask = "******"

type dataFieldProfile struct {
	FieldPath    string
	FieldType    string
	TypeStats    map[string]int64
	SampleValues []string
	Count        int64
	NullCount    int64
}

func validateDataCleanConfig(config *gjson.Json) error {
	if config == nil {
		return nil
	}
	if err := validateDataCleanFilterConfig(config.GetJson("filter")); err != nil {
		return err
	}
	for _, item := range config.Get("steps").Array() {
		step := gjson.New(item)
		stepType := strings.ToLower(strings.TrimSpace(step.Get("type").String()))
		switch stepType {
		case "", "noop", "rename", "remove", "default", "set_field", "set_event_type",
			"to_string", "to_int", "to_float", "to_bool", "trim", "regex_replace",
			"parse_json", "mask", "drop_if":
		default:
			return gerror.Newf("清洗步骤类型无效: %s", stepType)
		}
	}
	return nil
}

func validateDataCleanFilterConfig(filter *gjson.Json) error {
	if filter == nil || !filter.Get("enabled").Bool() {
		return nil
	}
	for _, group := range filter.Get("groups").Array() {
		groupJson := gjson.New(group)
		for _, item := range groupJson.Get("conditions").Array() {
			operator := normalizeDataCleanFilterOperator(gjson.New(item).Get("operator").String())
			switch operator {
			case "eq", "ne", "contains", "not_contains", "matches", "empty", "not_empty",
				"gt", "gte", "lt", "lte", "between", "in", "not_in":
			default:
				return gerror.Newf("过滤操作符无效: %s", operator)
			}
		}
	}
	return nil
}

func normalizeDataCleanFilterOperator(operator string) string {
	operator = strings.TrimSpace(operator)
	operator = strings.ReplaceAll(operator, "-", "_")
	switch operator {
	case "notContains":
		return "not_contains"
	case "notEmpty":
		return "not_empty"
	case "notIn":
		return "not_in"
	default:
		return strings.ToLower(operator)
	}
}

func normalizeDataSinkConfig(config *gjson.Json) *gjson.Json {
	data := jsonToMap(config)
	if len(gjson.New(data).Get("targets").Array()) == 0 {
		data["targets"] = []map[string]interface{}{
			{"type": consts.DataSinkTargetTypeStrategyEngine, "enabled": true},
		}
	}
	if strings.TrimSpace(gjson.New(data).Get("failurePolicy").String()) == "" {
		data["failurePolicy"] = "continue"
	}
	if strings.TrimSpace(gjson.New(data).Get("mode").String()) == "" {
		data["mode"] = "parallel"
	}
	return gjson.New(data)
}

func validateDataSinkConfig(config *gjson.Json) error {
	config = normalizeDataSinkConfig(config)
	for _, item := range config.Get("targets").Array() {
		target := gjson.New(item)
		targetType := strings.ToLower(strings.TrimSpace(target.Get("type").String()))
		if targetType == "" {
			targetType = consts.DataSinkTargetTypeStrategyEngine
		}
		switch targetType {
		case consts.DataSinkTargetTypeStrategyEngine:
		case consts.DataSinkTargetTypeSink:
			if strings.TrimSpace(target.Get("sinkCode").String()) == "" {
				return gerror.New("输出目标分发需要选择数据源")
			}
		default:
			return gerror.Newf("输出目标类型无效: %s", targetType)
		}
	}
	return nil
}

func profileDataFields(payloads []map[string]interface{}, maxDepth int) []dataFieldProfile {
	profiles := make(map[string]*dataFieldProfile)
	for _, payload := range payloads {
		collectDataFields("", payload, 0, maxDepth, profiles)
	}
	list := make([]dataFieldProfile, 0, len(profiles))
	for _, profile := range profiles {
		profile.FieldType = mergeDataFieldType(profile.TypeStats)
		list = append(list, *profile)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].FieldPath < list[j].FieldPath
	})
	return list
}

func collectDataFields(prefix string, value interface{}, depth int, maxDepth int, profiles map[string]*dataFieldProfile) {
	if prefix != "" && maxDepth > 0 && depth >= maxDepth {
		appendDataFieldProfile(prefix, value, profiles)
		return
	}
	switch typed := value.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			if isDynamicDataFieldKey(key) {
				continue
			}
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			collectDataFields(path, typed[key], depth+1, maxDepth, profiles)
		}
	case []interface{}:
		if len(typed) == 0 {
			appendDataFieldProfile(prefix, value, profiles)
			return
		}
		for _, item := range typed {
			path := prefix
			if path != "" {
				path += "[]"
			}
			collectDataFields(path, item, depth+1, maxDepth, profiles)
		}
	default:
		appendDataFieldProfile(prefix, value, profiles)
	}
}

func appendDataFieldProfile(path string, value interface{}, profiles map[string]*dataFieldProfile) {
	path = strings.TrimSpace(path)
	if path == "" || len(path) > 512 {
		return
	}
	fieldType := dataFieldType(value)
	profile := profiles[path]
	if profile == nil {
		profile = &dataFieldProfile{
			FieldPath: path,
			FieldType: fieldType,
			TypeStats: map[string]int64{},
		}
		profiles[path] = profile
	}
	profile.TypeStats[fieldType]++
	profile.Count++
	if value == nil {
		profile.NullCount++
	}
	sample := maskDataFieldSample(path, fmt.Sprint(value))
	if sample != "" && !stringSliceContains(profile.SampleValues, sample) && len(profile.SampleValues) < 5 {
		profile.SampleValues = append(profile.SampleValues, sample)
	}
}

func dataFieldType(value interface{}) string {
	switch value.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "string"
	}
}

func mergeDataFieldType(stats map[string]int64) string {
	if len(stats) == 0 {
		return "unknown"
	}
	preferred := "unknown"
	var count int64
	for fieldType, current := range stats {
		if current > count || current == count && fieldType < preferred {
			preferred = fieldType
			count = current
		}
	}
	return preferred
}

func parseDataCleanPayloads(raw string) ([]map[string]interface{}, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	var object map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &object); err == nil && object != nil {
		return []map[string]interface{}{object}, nil
	}
	var objects []map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &objects); err == nil {
		return objects, nil
	}
	scanner := bufio.NewScanner(bytes.NewReader([]byte(raw)))
	scanner.Buffer(make([]byte, 1024), 1024*1024*10)
	var list []map[string]interface{}
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			return nil, gerror.Newf("第 %d 行不是合法 JSON 对象: %v", lineNo, err)
		}
		list = append(list, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func limitDataCleanSamplePayloads(payloads []map[string]interface{}, limit int) []map[string]interface{} {
	if limit <= 0 || len(payloads) <= limit {
		return payloads
	}
	return payloads[:limit]
}

func appendDataCleanSamplePayloads(target []map[string]interface{}, items []map[string]interface{}, limit int) []map[string]interface{} {
	for _, item := range items {
		if len(target) >= limit {
			break
		}
		target = append(target, item)
	}
	return target
}

func mergeDataCleanSourceConfig(connectorConfig *gjson.Json, taskConfig *gjson.Json) map[string]interface{} {
	merged := jsonToMap(connectorConfig)
	for key, value := range jsonToMap(taskConfig) {
		merged[key] = value
	}
	return merged
}

func dataCleanConfigString(config *gjson.Json, path string) string {
	if config == nil {
		return ""
	}
	return strings.TrimSpace(config.Get(path).String())
}

func jsonToMap(config *gjson.Json) map[string]interface{} {
	if config == nil || config.IsNil() {
		return map[string]interface{}{}
	}
	return gconv.Map(config.Interface())
}

func sampleValueToStrings(value *gjson.Json) []string {
	if value == nil || value.IsNil() {
		return nil
	}
	values := gconv.Strings(value.Interface())
	if len(values) > 0 {
		return values
	}
	text := strings.TrimSpace(value.String())
	if text == "" {
		return nil
	}
	return []string{text}
}

func maskDataFieldSample(path, value string) string {
	if len(value) > 256 {
		value = value[:256]
	}
	lower := strings.ToLower(path)
	if strings.Contains(lower, "password") ||
		strings.Contains(lower, "secret") ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "authorization") ||
		strings.Contains(lower, "accesskey") ||
		strings.Contains(lower, "access_key") {
		return dataCleanSecretMask
	}
	return value
}

func isDynamicDataFieldKey(key string) bool {
	key = strings.TrimSpace(key)
	if key == "" {
		return true
	}
	if regexp.MustCompile(`^\d+$`).MatchString(key) {
		return true
	}
	if regexp.MustCompile(`^[0-9a-fA-F]{16,}$`).MatchString(key) {
		return true
	}
	return false
}

func stringSliceContains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
