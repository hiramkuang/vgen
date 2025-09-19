// internal/generator/generate.go
package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"text/template"

	vgenparser "github.com/hiramkuang/vgen/internal/parser"
)

// FieldInfo 保存从结构体字段中提取的信息
type FieldInfo struct {
	Name       string
	Rules      []vgenparser.Rule
	Validators []string
}

// StructInfo 保存结构体名称和其字段信息
type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

// GenerateValidator 为指定的 Go 文件生成 Validate() 方法
func GenerateValidator(filePath string) error {
	fmt.Printf("Debug: Parsing file %s\n", filePath)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}
	fmt.Printf("Debug: Parsed package %s\n", node.Name.Name)

	// 收集所有结构体信息
	var structInfos []StructInfo

	// 遍历文件中的所有声明
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// 遍历类型声明中的所有规格 (Spec)
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 检查是否为结构体类型
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 创建结构体信息
			structInfo := StructInfo{Name: typeSpec.Name.Name}

			// 遍历结构体的字段
			for _, field := range structType.Fields.List {
				// 忽略没有名字的字段（如嵌入结构体）
				if len(field.Names) == 0 {
					continue
				}

				// 获取字段名
				fieldName := field.Names[0].Name

				// 获取字段类型（用于生成更精确的校验代码）
				var fieldType string
				switch t := field.Type.(type) {
				case *ast.Ident:
					fieldType = t.Name
				case *ast.SelectorExpr:
					if ident, ok := t.X.(*ast.Ident); ok {
						fieldType = ident.Name + "." + t.Sel.Name
					}
				default:
					fieldType = "unsupported"
				}

				// 获取 vgen tag
				var tagValue string
				if field.Tag != nil {
					tagStr := strings.Trim(field.Tag.Value, "`")
					tagValue = reflect.StructTag(tagStr).Get("vgen")
				}

				// 如果没有 vgen tag，则跳过
				if tagValue == "" {
					continue
				}

				// 使用我们的 parser 解析 tag
				rules, err := vgenparser.ParseTag(tagValue)
				if err != nil {
					return fmt.Errorf("error parsing tag for field %s.%s: %w", structInfo.Name, fieldName, err)
				}

				// --- 核心：为每个规则生成校验代码片段 (优化后) ---
				var validators []string
				for _, rule := range rules {
					var code string
					switch rule.Name {
					case "required":
						switch fieldType {
						case "string":
							code = fmt.Sprintf("if s.%s == \"\" { errs = append(errs, fmt.Errorf(\"field %%s is required\", \"%s\")) }", fieldName, fieldName)
						case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
							code = fmt.Sprintf("if s.%s == 0 { errs = append(errs, fmt.Errorf(\"field %%s is required\", \"%s\")) }", fieldName, fieldName)
						default:
							code = fmt.Sprintf("// TODO: Implement 'required' check for type %s", fieldType)
						}
					case "min":
						switch fieldType {
						case "string":
							if v, err := rule.GetIntValue(); err == nil {
								code = fmt.Sprintf("if len(s.%s) < %d { errs = append(errs, fmt.Errorf(\"field %%s length must be at least %%d, got %%d\", \"%s\", %d, len(s.%s))) }", fieldName, v, fieldName, v, fieldName)
							} else {
								return fmt.Errorf("invalid 'min' value for string field %s.%s: %w", structInfo.Name, fieldName, err)
							}
						case "int":
							if v, err := rule.GetIntValue(); err == nil {
								code = fmt.Sprintf("if s.%s < %d { errs = append(errs, fmt.Errorf(\"field %%s must be at least %%d, got %%d\", \"%s\", %d, s.%s)) }", fieldName, v, fieldName, v, fieldName)
							} else {
								return fmt.Errorf("invalid 'min' value for int field %s.%s: %w", structInfo.Name, fieldName, err)
							}
						default:
							code = fmt.Sprintf("// TODO: Implement 'min' check for type %s", fieldType)
						}
					case "max":
						switch fieldType {
						case "string":
							if v, err := rule.GetIntValue(); err == nil {
								code = fmt.Sprintf("if len(s.%s) > %d { errs = append(errs, fmt.Errorf(\"field %%s length must be at most %%d, got %%d\", \"%s\", %d, len(s.%s))) }", fieldName, v, fieldName, v, fieldName)
							} else {
								return fmt.Errorf("invalid 'max' value for string field %s.%s: %w", structInfo.Name, fieldName, err)
							}
						case "int":
							if v, err := rule.GetIntValue(); err == nil {
								code = fmt.Sprintf("if s.%s > %d { errs = append(errs, fmt.Errorf(\"field %%s must be at most %%d, got %%d\", \"%s\", %d, s.%s)) }", fieldName, v, fieldName, v, fieldName)
							} else {
								return fmt.Errorf("invalid 'max' value for int field %s.%s: %w", structInfo.Name, fieldName, err)
							}
						default:
							code = fmt.Sprintf("// TODO: Implement 'max' check for type %s", fieldType)
						}
					case "email":
						if fieldType == "string" {
							code = fmt.Sprintf("if !isEmailValid(s.%s) { errs = append(errs, fmt.Errorf(\"field %%s is not a valid email\", \"%s\")) }", fieldName, fieldName)
						} else {
							return fmt.Errorf("rule 'email' is not applicable to field %s.%s of type %s", structInfo.Name, fieldName, fieldType)
						}
					// --- 新增规则开始 ---
					case "len":
						// len 规则适用于 string 和 slice
						if strings.HasPrefix(fieldType, "[]") || fieldType == "string" {
							if v, err := rule.GetIntValue(); err == nil {
								// 对于 string 和 slice，都使用 len() 函数
								code = fmt.Sprintf("if len(s.%s) != %d { errs = append(errs, fmt.Errorf(\"field %%s length must be %%d, got %%d\", \"%s\", %d, len(s.%s))) }", fieldName, v, fieldName, v, fieldName)
							} else {
								return fmt.Errorf("invalid 'len' value for field %s.%s: %w", structInfo.Name, fieldName, err)
							}
						} else {
							return fmt.Errorf("rule 'len' is not applicable to field %s.%s of type %s", structInfo.Name, fieldName, fieldType)
						}
					case "in":
						// in 规则目前主要适用于 string (可以扩展)
						if fieldType == "string" {
							// 获取值列表
							inValues := rule.GetInValues()
							if len(inValues) > 0 {
								// 生成一个 map 来进行 O(1) 查找，提高效率
								// map[值]bool
								allowedMap := make(map[string]bool, len(inValues))
								for _, val := range inValues {
									allowedMap[val] = true
								}
								// 将 map 转换为 Go 代码中的 map 字面量
								mapLiteral := "map[string]bool{"
								for _, val := range inValues { // 使用 range inValues 保留所有值和顺序
									// 对键进行转义，以防包含引号等特殊字符
									mapLiteral += fmt.Sprintf("%q: true,", val)
								}
								mapLiteral += "}"

								// 生成校验代码
								// 注意：我们在生成的代码中定义 map，以避免在包级别定义过多全局变量
								code = fmt.Sprintf(`
{
	allowedValues := %s
	if !allowedValues[s.%s] {
		errs = append(errs, fmt.Errorf("field %%s value '%%s' is not in the allowed list [%%s]", "%s", s.%s, %q))
	}
}`, mapLiteral, fieldName, fieldName, fieldName, strings.Join(inValues, ", "))
							} else {
								return fmt.Errorf("invalid 'in' value for field %s.%s: %w", structInfo.Name, fieldName, err)
							}
						} else {
							return fmt.Errorf("rule 'in' is not applicable to field %s.%s of type %s", structInfo.Name, fieldName, fieldType)
						}
					// --- 新增规则结束 ---
					default:
						return fmt.Errorf("unknown rule %s for field %s.%s", rule.Name, structInfo.Name, fieldName)
					}
					validators = append(validators, code)
				}

				// 保存字段信息
				structInfo.Fields = append(structInfo.Fields, FieldInfo{
					Name:       fieldName,
					Rules:      rules,
					Validators: validators,
				})
			}

			// 保存结构体信息
			structInfos = append(structInfos, structInfo)
		}
	}

	// 为每个结构体生成 Validate() 方法
	for _, structInfo := range structInfos {
		// 生成 Validate() 方法的代码
		var validateMethod strings.Builder
		validateMethod.WriteString(fmt.Sprintf("func (s *%s) Validate() error {\n", structInfo.Name))
		validateMethod.WriteString("var errs []error\n")

		// 为每个字段生成校验代码
		for _, fieldInfo := range structInfo.Fields {
			for _, validator := range fieldInfo.Validators {
				validateMethod.WriteString(validator + "\n")
			}
		}

		validateMethod.WriteString("if len(errs) > 0 {\n")
		validateMethod.WriteString("return fmt.Errorf(strings.Join(errs, \"\\n\"))\n")
		validateMethod.WriteString("}\n")
		validateMethod.WriteString("return nil\n")

		// 将生成的代码写入文件
		err := os.WriteFile(structInfo.Name+"_gen.go", []byte(validateMethod.String()), 0644)
		if err != nil {
			return fmt.Errorf("failed to write generated code to file: %w", err)
		}
	}

	return nil
}
