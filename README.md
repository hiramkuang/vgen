# VGen - Go Struct Validation Generator

[![Go Report Card](https://goreportcard.com/badge/github.com/hiramkuang/vgen)](https://goreportcard.com/report/github.com/hiramkuang/vgen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`vgen` 是一个简单而强大的 Go 代码生成工具，它可以根据你在结构体字段中定义的标签（tags），自动生成类型安全的 `Validate()` 方法。从此告别繁琐、重复的手写验证逻辑。

## 特性

- **基于标签**: 通过简单的 `vgen:"..."` 标签来定义验证规则。
- **代码生成**: 自动生成 `_validator.go` 文件，包含所有验证逻辑。
- **类型安全**: 生成的代码是纯 Go 代码，易于调试和集成。
- **可扩展**: 轻松添加自定义验证规则。
- **命令行友好**: 提供了简单易用的命令行工具来处理文件或整个目录。

## 安装

### 方式一：使用 `go install` (推荐)

```bash
go install github.com/hiramkuang/vgen/cmd/vgen@latest
```

这会将 `vgen` 命令安装到你的 `$GOPATH/bin` 目录下。请确保该目录已添加到你的系统环境变量 `PATH` 中。

### 方式二：从源码构建

1.  克隆此仓库到你的本地机器。
2.  进入项目根目录。
3.  运行构建命令：

    ```bash
    # Windows
    go build -o vgen.exe ./cmd/vgen
    # Linux/macOS
    go build -o vgen ./cmd/vgen
    ```
4.  将生成的 `vgen` 或 `vgen.exe` 可执行文件移动到你的 `PATH` 环境变量包含的目录中，或者直接在项目根目录下使用。

## 快速开始

### 1. 定义你的结构体

在你的 `.go` 文件中定义一个结构体，并使用 `vgen` 标签添加验证规则。

```go
// examples/user.go
package main

// User 代表一个用户，包含验证规则。
type User struct {
    Name   string `vgen:"required,min=2,max=50"`
    Email  string `vgen:"required,email"`
    Age    int    `vgen:"required,min=0,max=150"`
    City   string `vgen:"len=5"`
    Status string `vgen:"in=active,pending,disabled"`
}
```

### 2. 生成验证代码

使用 `vgen` 命令处理你的 Go 文件或所在目录。

```bash
# 处理单个文件
vgen examples/user.go

# 或者处理整个目录
vgen examples/
```

这将在 `examples` 目录下生成一个名为 `user_validator.go` 的文件。

### 3. 使用生成的验证器

在你的代码中，可以直接调用生成的 `Validate()` 方法。

```go
// examples/main.go
package main

import (
    "fmt"
    "log"
)

func main() {
    user := &User{
        Name:  "A", // 名字太短
        Email: "not-an-email", // 邮箱格式错误
        Age:   -5, // 年龄小于0
        City:  "NY", // 长度不等于5
        Status: "unknown", // 无效状态
    }

    if err := user.Validate(); err != nil {
        // 验证失败，打印错误信息
        log.Printf("Validation failed: %v", err)
        // 示例输出:
        // Validation failed: field Name length must be at least 2, got 1; field Email is not a valid email address; field Age must be at least 0, got -5; field City length must be 5, got 2; field Status must be one of [active pending disabled], got unknown
    } else {
        fmt.Println("User is valid!")
    }
}
```

### 4. 运行你的程序

确保包含了生成的 `_validator.go` 文件一起编译。

```bash
go run ./examples
```

## 命令行工具 (`vgen`)

### 用法

```bash
vgen [flags] <FILE_OR_DIR_PATH>
```

### 参数

-   `<FILE_OR_DIR_PATH>`: (必需) 要处理的单个 Go 源文件路径或包含 Go 源文件的目录路径。

### 标志 (Flags)

-   `-h, --help`: 显示帮助信息。
-   `-o, --output string`: 指定生成文件的输出目录。默认与输入文件在同一目录。
-   `-r, --recursive`: 如果输入是目录，则递归处理所有子目录。
-   `-v, --verbose`: 启用详细输出模式，显示处理过程中的调试信息。

### 示例

```bash
# 处理单个文件
vgen path/to/your/file.go

# 处理目录下的所有 .go 文件 (不递归)
vgen path/to/your/directory

# 递归处理目录及其子目录下的所有 .go 文件
vgen -r path/to/your/directory

# 将生成的文件输出到 'generated' 目录
vgen -o generated path/to/your/file.go
```

## 支持的验证规则

| 规则 | 描述 | 适用类型 | 示例 |
| :--- | :--- | :--- | :--- |
| `required` | 字段不能为空 (字符串为 `""`, 指针为 `nil`, 切片/映射为 `len == 0`) | 所有类型 | `vgen:"required"` |
| `min` | 数值最小值 / 字符串或切片/映射的最小长度 | `string`, `int/*`, `uint/*`, `float*`, `[]T`, `map[K]V` | `vgen:"min=18"` |
| `max` | 数值最大值 / 字符串或切片/映射的最大长度 | `string`, `int/*`, `uint/*`, `float*`, `[]T`, `map[K]V` | `vgen:"max=100"` |
| `len` | 字符串或切片/映射的精确长度 | `string`, `[]T`, `map[K]V` | `vgen:"len=5"` |
| `email` | 验证字符串是否为有效的电子邮件地址 | `string` | `vgen:"email"` |
| `in` | 验证字符串值是否在给定的列表中 | `string` | `vgen:"in=active,pending,disabled"` |

## 开发与贡献

我们欢迎任何形式的贡献！

### 项目结构

```
vgen/
├── cmd/vgen/             # CLI 命令入口
│   └── main.go
├── examples/             # 示例代码
├── internal/
│   ├── generator/        # 代码生成核心逻辑
│   │   └── generate.go
│   └── parser/           # 标签解析逻辑
│       └── tag.go
└── go.mod                # Go 模块文件
```

### 如何贡献

1.  Fork 此仓库。
2.  创建你的特性分支 (`git checkout -b feature/AmazingFeature`)。
3.  提交你的更改 (`git commit -m 'Add some AmazingFeature'`)。
4.  推送到分支 (`git push origin feature/AmazingFeature`)。
5.  开启一个 Pull Request。

### 添加新规则

1.  在 `internal/parser/tag.go` 的 `ParseTag` 函数中添加新规则的解析逻辑。
2.  在 `internal/generator/generate.go` 中添加生成对应验证代码的逻辑。
3.  在 `examples` 目录中添加测试用例。
4.  更新 `README.md` 中的“支持的验证规则”表格。

## 许可证

本项目采用 MIT 许可证。
