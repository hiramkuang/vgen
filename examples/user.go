// examples/user.go
package main

// User represents a user with validation rules.
type User struct {
	// ID     int    `vgen:"required"` // 如果需要 ID，取消注释并添加规则
	Name   string `vgen:"required,min=2,max=50"`
	Email  string `vgen:"required,email"`
	Age    int    `vgen:"required,min=0,max=150"`
	City   string `vgen:"len=5"`                      // 城市名必须是5个字符
	Status string `vgen:"in=active,pending,disabled"` // 状态只能是这三个值之一
	// 可以添加更多字段和规则进行测试
}
