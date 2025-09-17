// internal/parser/tag.go
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule 代表一个从 tag 中解析出的单条校验规则
type Rule struct {
	Name  string            // 规则名称，例如 "required", "min"
	Value string            // 规则的值，例如 "2", "50"
	Args  map[string]string // 未来可能支持的键值对参数 (预留)
}

// ParseTag 解析 vgen tag 字符串，例如 `vgen:"required,min=2,max=50"`
func ParseTag(tag string) ([]Rule, error) {
	var rules []Rule

	// 去掉首尾空格并按逗号分割
	parts := strings.Split(strings.TrimSpace(tag), ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		rule := Rule{}
		// 检查是否有等号，例如 "min=2"
		if eqIndex := strings.Index(part, "="); eqIndex != -1 {
			rule.Name = part[:eqIndex]
			rule.Value = part[eqIndex+1:]
		} else {
			// 没有等号，例如 "required"
			rule.Name = part
			rule.Value = ""
		}

		// 基本验证：规则名不能为空
		if rule.Name == "" {
			return nil, fmt.Errorf("invalid tag part: %s", part)
		}

		// 预留：可以在这里对 rule.Value 进行类型检查
		// 例如，如果 rule.Name 是 "min"，则 rule.Value 应该是数字
		// 我们将在后续步骤中实现

		rules = append(rules, rule)
	}

	return rules, nil
}

// GetIntValue 是一个辅助函数，用于安全地从 Rule.Value 获取整数值
func (r *Rule) GetIntValue() (int, error) {
	if r.Value == "" {
		return 0, fmt.Errorf("rule %s has no value", r.Name)
	}
	v, err := strconv.Atoi(r.Value)
	if err != nil {
		return 0, fmt.Errorf("rule %s: invalid integer value '%s'", r.Name, r.Value)
	}
	return v, nil
}

// GetInValues 解析 'in' 规则的值，返回一个字符串切片。
// 例如，对于 Value = "a,b,c"，将返回 []string{"a", "b", "c"}。
// 如果规则不是 'in' 或 Value 为空，则返回 nil。
func (r Rule) GetInValues() []string {
	if r.Name != "in" || r.Value == "" {
		return nil
	}
	// 按逗号分割，并去除每个值前后的空格
	parts := strings.Split(r.Value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" { // 忽略空字符串
			values = append(values, trimmed)
		}
	}
	return values
}
