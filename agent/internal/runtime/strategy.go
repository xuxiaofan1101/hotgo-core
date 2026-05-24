package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	"vogo-agent/internal/protocol"
)

type strategySink struct {
	registry *ConnectorRegistry
	rules    []agentStrategyRule
	sinks    map[string]protocol.OutputConfig
}

type agentStrategyRule struct {
	Id                  int64                 `json:"id"`
	Name                string                `json:"name"`
	EventType           string                `json:"eventType"`
	Decision            string                `json:"decision"`
	ConditionExpression string                `json:"conditionExpression"`
	Priority            int                   `json:"priority"`
	StopOnMatch         bool                  `json:"stopOnMatch"`
	Actions             []agentStrategyAction `json:"actions"`
}

type agentStrategyAction struct {
	ActionType   string         `json:"actionType"`
	ActionParams map[string]any `json:"actionParams"`
}

func newStrategySink(config protocol.OutputConfig, registry *ConnectorRegistry) (Sink, error) {
	rules, err := decodeStrategyRules(config.Config["rules"])
	if err != nil {
		return nil, err
	}
	sinks, err := decodeStrategySinks(config.Config["sinks"])
	if err != nil {
		return nil, err
	}
	sort.SliceStable(rules, func(i, j int) bool {
		if rules[i].Priority == rules[j].Priority {
			return rules[i].Id < rules[j].Id
		}
		return rules[i].Priority > rules[j].Priority
	})
	return &strategySink{registry: registry, rules: rules, sinks: sinks}, nil
}

func (s *strategySink) Write(ctx context.Context, payload map[string]any) error {
	eventType := strings.TrimSpace(fmt.Sprint(payload["eventType"]))
	for _, rule := range s.rules {
		if eventType != "" && strings.TrimSpace(rule.EventType) != "" && rule.EventType != eventType {
			continue
		}
		matched, err := evaluateStrategyExpression(rule.ConditionExpression, payload)
		if err != nil {
			return err
		}
		if !matched {
			continue
		}
		if err := s.executeActions(ctx, rule, payload); err != nil {
			return err
		}
		if rule.StopOnMatch || rule.Decision != "" && rule.Decision != "allow" {
			break
		}
	}
	return nil
}

func (s *strategySink) executeActions(ctx context.Context, rule agentStrategyRule, payload map[string]any) error {
	for _, action := range rule.Actions {
		actionType := strings.TrimSpace(action.ActionType)
		params := mergeActionParams(action.ActionParams, payload)
		switch actionType {
		case "output.send", "security.alert":
			if err := s.executeOutputAction(ctx, params, payload); err != nil {
				return err
			}
		case "webhook.call":
			sink, err := newHTTPSink(protocol.OutputConfig{Type: "http", Config: params})
			if err != nil {
				return err
			}
			if err := sink.Write(ctx, payload); err != nil {
				return err
			}
		case "system.write_log", "risk.deny", "risk.review", "risk.challenge", "risk.mark_risk", "security.raise_risk":
			continue
		default:
			return fmt.Errorf("agent 暂不支持策略动作: %s", actionType)
		}
	}
	return nil
}

func (s *strategySink) executeOutputAction(ctx context.Context, params map[string]any, payload map[string]any) error {
	sinkCode := strings.TrimSpace(fmt.Sprint(params["sinkCode"]))
	if sinkCode == "" {
		return fmt.Errorf("策略输出动作缺少 sinkCode")
	}
	output, ok := s.sinks[sinkCode]
	if !ok {
		return fmt.Errorf("策略输出目标未下发: %s", sinkCode)
	}
	if override, ok := params["config"].(map[string]any); ok {
		output.Config = mergeConfig(output.Config, override)
	}
	outPayload := clonePayload(payload)
	if message := strings.TrimSpace(fmt.Sprint(params["message"])); message != "" && message != "<nil>" {
		outPayload["message"] = message
	}
	if extra, ok := params["payload"].(map[string]any); ok {
		for key, value := range extra {
			outPayload[key] = value
		}
	}
	sink, err := s.registry.BuildSink(output)
	if err != nil {
		return err
	}
	return sink.Write(ctx, outPayload)
}

func decodeStrategyRules(value any) ([]agentStrategyRule, error) {
	if value == nil {
		return nil, nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var rules []agentStrategyRule
	if err := json.Unmarshal(raw, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func decodeStrategySinks(value any) (map[string]protocol.OutputConfig, error) {
	result := map[string]protocol.OutputConfig{}
	if value == nil {
		return result, nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func evaluateStrategyExpression(expression string, payload map[string]any) (bool, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		expression = "true"
	}
	env, err := cel.NewEnv(strategyEnvOptions(payload)...)
	if err != nil {
		return false, err
	}
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return false, issues.Err()
	}
	if !ast.OutputType().IsExactType(cel.BoolType) {
		return false, fmt.Errorf("策略表达式必须返回 bool")
	}
	program, err := env.Program(ast)
	if err != nil {
		return false, err
	}
	out, _, err := program.Eval(strategyActivation(payload))
	if err != nil {
		return false, err
	}
	matched, ok := out.Value().(bool)
	if !ok {
		return false, fmt.Errorf("策略表达式执行结果不是 bool")
	}
	return matched, nil
}

func strategyEnvOptions(payload map[string]any) []cel.EnvOption {
	keys := map[string]struct{}{
		"payload": {},
		"data":    {},
	}
	for key := range payload {
		keys[key] = struct{}{}
	}
	options := make([]cel.EnvOption, 0, len(keys)+1)
	for key := range keys {
		options = append(options, cel.Variable(key, cel.DynType))
	}
	return append(options, strategyIpInCidrFunction())
}

func strategyActivation(payload map[string]any) map[string]any {
	activation := map[string]any{
		"payload": payload,
		"data":    payload,
	}
	for key, value := range payload {
		activation[key] = normalizeStrategyValue(value)
	}
	return activation
}

func normalizeStrategyValue(value any) any {
	switch typed := value.(type) {
	case json.Number:
		if intValue, err := typed.Int64(); err == nil {
			return intValue
		}
		if floatValue, err := typed.Float64(); err == nil {
			return floatValue
		}
	}
	return value
}

func strategyIpInCidrFunction() cel.EnvOption {
	return cel.Function("ipInCidr",
		cel.Overload("agent_ip_in_cidr_string_string",
			[]*cel.Type{cel.StringType, cel.StringType},
			cel.BoolType,
			cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				ip := net.ParseIP(strings.TrimSpace(fmt.Sprint(lhs.Value())))
				if ip == nil {
					return types.Bool(false)
				}
				_, network, err := net.ParseCIDR(strings.TrimSpace(fmt.Sprint(rhs.Value())))
				return types.Bool(err == nil && network.Contains(ip))
			}),
		),
	)
}

func mergeActionParams(params map[string]any, payload map[string]any) map[string]any {
	result := map[string]any{}
	for key, value := range params {
		if text, ok := value.(string); ok {
			result[key] = renderStrategyTemplate(text, payload)
			continue
		}
		result[key] = value
	}
	if _, ok := result["eventType"]; !ok {
		result["eventType"] = payload["eventType"]
	}
	return result
}

func renderStrategyTemplate(template string, payload map[string]any) string {
	result := template
	for key, value := range payload {
		result = strings.ReplaceAll(result, "{{"+key+"}}", fmt.Sprint(value))
	}
	result = strings.ReplaceAll(result, "{{now}}", time.Now().Format(time.RFC3339))
	return result
}
