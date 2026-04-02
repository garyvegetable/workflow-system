package expression

import (
	"regexp"
	"strings"
)

var allowedFields []string

func SetAllowedFields(fields []string) {
	allowedFields = fields
}

func ValidateExpression(expressionStr string) (bool, error) {
	// 白名单验证：只允许表单字段名
	// 提取表达式中的所有变量名
	re := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
	variables := re.FindAllString(expressionStr, -1)

	// 运算符白名单
	allowedOperators := []string{
		"+", "-", "*", "/", "%",
		">", "<", ">=", "<=", "==", "!=",
		"&&", "||", "!",
		"contains", "in",
	}

	for _, v := range variables {
		// 检查是否是运算符
		isOperator := false
		for _, op := range allowedOperators {
			if v == op {
				isOperator = true
				break
			}
		}
		if isOperator {
			continue
		}

		// 检查是否是允许的字段
		if len(allowedFields) > 0 {
			isAllowed := false
			for _, field := range allowedFields {
				if v == field {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				return false, nil
			}
		}
	}

	// 检查是否有危险函数
	dangerousPatterns := []string{
		"env.", "file.", "system.", "exec.", "os.",
		"require", "import", "__",
	}
	lowerExpr := strings.ToLower(expressionStr)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerExpr, pattern) {
			return false, nil
		}
	}

	return true, nil
}
