package runtime

import (
	"fmt"
	"sort"
	"strings"

	"vogo-agent/internal/protocol"
)

type FieldCollector struct {
	known         map[string]struct{}
	maxParseDepth int
}

func NewFieldCollector(snapshot protocol.FieldSnapshot) *FieldCollector {
	known := make(map[string]struct{}, len(snapshot.KnownFieldHashes))
	for _, hash := range snapshot.KnownFieldHashes {
		if hash = strings.TrimSpace(hash); hash != "" {
			known[hash] = struct{}{}
		}
	}
	return &FieldCollector{known: known, maxParseDepth: snapshot.MaxParseDepth}
}

func (c *FieldCollector) Collect(payload map[string]any) []protocol.FieldSample {
	if c == nil {
		return nil
	}
	samples := make([]protocol.FieldSample, 0)
	c.collectInto("", payload, 0, &samples)
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Path < samples[j].Path
	})
	return samples
}

func (c *FieldCollector) collectInto(prefix string, value any, depth int, samples *[]protocol.FieldSample) {
	if prefix != "" && c.maxParseDepth > 0 && depth >= c.maxParseDepth {
		c.appendSample(prefix, value, samples)
		return
	}
	if obj, ok := value.(map[string]any); ok {
		keys := make([]string, 0, len(obj))
		for key := range obj {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			c.collectInto(path, obj[key], depth+1, samples)
		}
		return
	}
	c.appendSample(prefix, value, samples)
}

func (c *FieldCollector) appendSample(prefix string, value any, samples *[]protocol.FieldSample) {
	fieldType := inferFieldType(value)
	hash := protocol.FieldHash(prefix, fieldType)
	if _, ok := c.known[hash]; ok {
		return
	}
	c.known[hash] = struct{}{}
	*samples = append(*samples, protocol.FieldSample{
		Path:      prefix,
		Type:      fieldType,
		Sample:    maskFieldSample(prefix, fmt.Sprint(value)),
		Hash:      hash,
		Count:     1,
		NullCount: boolToInt64(value == nil),
	})
}

func inferFieldType(value any) string {
	switch value.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return "string"
	}
}

func maskFieldSample(path, sample string) string {
	lower := strings.ToLower(path)
	for _, keyword := range []string{"password", "token", "secret", "authorization", "apikey", "api_key"} {
		if strings.Contains(lower, keyword) {
			return "***"
		}
	}
	return sample
}

func boolToInt64(value bool) int64 {
	if value {
		return 1
	}
	return 0
}
