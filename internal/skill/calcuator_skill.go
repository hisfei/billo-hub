package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// ============================================
// CalculatorSkill - 生产级计算器技能
// ============================================

// CalculatorParameters 定义计算器技能的参数结构
type CalculatorParameters struct {
	Type       string                    `json:"type"`
	Properties CalculatorParamProperties `json:"properties"`
	Required   []string                  `json:"required"`
}

// CalculatorParamProperties 定义参数属性
type CalculatorParamProperties struct {
	Expression CalculatorParamProperty `json:"expression"`
	Precision  CalculatorParamProperty `json:"precision"`
	AngleMode  CalculatorParamProperty `json:"angle_mode"`
	Action     CalculatorParamProperty `json:"action"`
	Value      CalculatorParamProperty `json:"value"`
	FromUnit   CalculatorParamProperty `json:"from_unit"`
	ToUnit     CalculatorParamProperty `json:"to_unit"`
	Value1     CalculatorParamProperty `json:"value1"`
	Value2     CalculatorParamProperty `json:"value2"`
	FromBase   CalculatorParamProperty `json:"from_base"`
	ToBase     CalculatorParamProperty `json:"to_base"`
}

// CalculatorParamProperty 定义单个参数属性
type CalculatorParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// CalculatorArgs 定义计算参数
type CalculatorArgs struct {
	Expression string `json:"expression"`
	Precision  int    `json:"precision,omitempty"`
	AngleMode  string `json:"angle_mode,omitempty"`
	Action     string `json:"action,omitempty"`
	Value      any    `json:"value,omitempty"`
	FromUnit   string `json:"from_unit,omitempty"`
	ToUnit     string `json:"to_unit,omitempty"`
	Value1     any    `json:"value1,omitempty"`
	Value2     any    `json:"value2,omitempty"`
	FromBase   int    `json:"from_base,omitempty"`
	ToBase     int    `json:"to_base,omitempty"`
}

// CalculatorResult 定义计算结果
type CalculatorResult struct {
	Success       bool                   `json:"success"`
	Result        interface{}            `json:"result,omitempty"`
	Expression    string                 `json:"expression,omitempty"`
	Precision     int                    `json:"precision,omitempty"`
	Action        string                 `json:"action,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Warnings      []string               `json:"warnings,omitempty"`
	ExecutionTime float64                `json:"execution_time_ms,omitempty"`
}

// CalculatorSkill 实现了生产级计算器技能
type CalculatorSkill struct {
	name        string
	descName    string
	description string
	parameters  CalculatorParameters
}

// NewCalculatorSkill 创建新的计算器技能实例
func NewCalculatorSkill() *CalculatorSkill {
	return &CalculatorSkill{
		name:     "calculator",
		descName: "数学计算器",
		description: `生产级数学计算器，支持多种数学运算。

【基础运算】
- 加减乘除: +, -, *, /
- 幂运算: ^ 或 **
- 取模: %
- 取整: // (整数除法)

【科学函数】
- 三角函数: sin, cos, tan, asin, acos, atan
- 双曲函数: sinh, cosh, tanh
- 对数函数: log(以10为底), ln(自然对数), log2
- 指数函数: exp
- 其他: sqrt, cbrt, abs, floor, ceil, round

【常量】
- pi: 圆周率 π ≈ 3.14159
- e: 自然常数 e ≈ 2.71828
- phi: 黄金比例 φ ≈ 1.61803

【角度模式】
- deg: 角度制 (默认)
- rad: 弧度制

【高级功能】(通过action参数)
- unit_convert: 单位转换
- base_convert: 进制转换
- percentage: 百分比计算
- statistics: 统计计算
- compare: 数值比较

【示例】
- "2 + 3 * 4" = 14
- "sqrt(16)" = 4
- "sin(90)" = 1 (角度制)
- "log(100)" = 2
- "2^10" = 1024`,
		parameters: CalculatorParameters{
			Type: "object",
			Properties: CalculatorParamProperties{
				Expression: CalculatorParamProperty{
					Type:        "string",
					Description: "数学表达式",
				},
				Precision: CalculatorParamProperty{
					Type:        "integer",
					Description: "结果精度（小数位数）",
					Default:     10,
				},
				AngleMode: CalculatorParamProperty{
					Type:        "string",
					Description: "角度模式",
					Default:     "deg",
					Enum:        []string{"deg", "rad"},
				},
				Action: CalculatorParamProperty{
					Type:        "string",
					Description: "高级操作类型",
					Enum:        []string{"unit_convert", "base_convert", "percentage", "statistics", "compare", "factorial", "fibonacci", "is_prime"},
				},
				Value: CalculatorParamProperty{
					Type:        "number",
					Description: "操作值",
				},
				FromUnit: CalculatorParamProperty{
					Type:        "string",
					Description: "原始单位",
				},
				ToUnit: CalculatorParamProperty{
					Type:        "string",
					Description: "目标单位",
				},
				Value1: CalculatorParamProperty{
					Type:        "number",
					Description: "第一个值",
				},
				Value2: CalculatorParamProperty{
					Type:        "number",
					Description: "第二个值",
				},
				FromBase: CalculatorParamProperty{
					Type:        "integer",
					Description: "原进制 (2-36)",
				},
				ToBase: CalculatorParamProperty{
					Type:        "integer",
					Description: "目标进制 (2-36)",
				},
			},
			Required: []string{},
		},
	}
}

// GetName 返回技能名称
func (s *CalculatorSkill) GetName() string {
	return s.name
}

// GetDescName 返回技能的描述性名称
func (s *CalculatorSkill) GetDescName() string {
	return s.descName
}

// GetDescription 返回技能的详细描述
func (s *CalculatorSkill) GetDescription() string {
	return s.description
}

// GetParameters 返回技能的参数定义
func (s *CalculatorSkill) GetParameters() any {
	return s.parameters
}

// Execute 执行计算
func (s *CalculatorSkill) Execute(ctx context.Context, args string) (string, error) {
	var calcArgs CalculatorArgs
	if err := json.Unmarshal([]byte(args), &calcArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	result := &CalculatorResult{
		Precision: calcArgs.Precision,
	}

	if result.Precision <= 0 {
		result.Precision = 10
	}

	var err error

	// 判断是表达式计算还是高级操作
	if calcArgs.Action != "" {
		err = s.executeAction(calcArgs, result)
	} else if calcArgs.Expression != "" {
		err = s.evaluateExpression(calcArgs, result)
	} else {
		err = fmt.Errorf("需要提供expression或action参数")
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		result.Success = true
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return string(resultJSON), nil
}

// evaluateExpression 计算表达式
func (s *CalculatorSkill) evaluateExpression(args CalculatorArgs, result *CalculatorResult) error {
	expr := strings.TrimSpace(args.Expression)
	if expr == "" {
		return fmt.Errorf("表达式不能为空")
	}

	result.Expression = expr

	// 预处理表达式
	expr = s.preprocessExpression(expr, args.AngleMode)

	// 创建解析器并计算
	parser := newMathParser(expr, args.AngleMode == "rad")
	val, err := parser.parse()
	if err != nil {
		return err
	}

	// 格式化结果
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return fmt.Errorf("计算结果无效")
	}

	// 如果是整数，显示为整数
	if val == math.Trunc(val) && math.Abs(val) < 1e15 {
		result.Result = int64(val)
	} else {
		result.Result = s.round(val, result.Precision)
	}

	result.Data = map[string]interface{}{
		"angle_mode": args.AngleMode,
	}

	return nil
}

// preprocessExpression 预处理表达式
func (s *CalculatorSkill) preprocessExpression(expr, angleMode string) string {
	expr = strings.ToLower(expr)

	// 替换常量
	expr = strings.ReplaceAll(expr, "pi", fmt.Sprintf("(%.20f)", math.Pi))
	expr = strings.ReplaceAll(expr, "e", fmt.Sprintf("(%.20f)", math.E))
	expr = strings.ReplaceAll(expr, "phi", fmt.Sprintf("(%.20f)", math.Phi))

	// 替换 ** 为 ^
	expr = strings.ReplaceAll(expr, "**", "^")

	// 处理科学函数（角度转换）
	if angleMode == "deg" {
		// 在角度模式下，三角函数需要转换
		expr = regexp.MustCompile(`sin\(([^)]+)\)`).ReplaceAllStringFunc(expr, func(m string) string {
			submatch := regexp.MustCompile(`sin\(([^)]+)\)`).FindStringSubmatch(m)
			return fmt.Sprintf("sin_deg(%s)", submatch[1])
		})
		expr = regexp.MustCompile(`cos\(([^)]+)\)`).ReplaceAllStringFunc(expr, func(m string) string {
			submatch := regexp.MustCompile(`cos\(([^)]+)\)`).FindStringSubmatch(m)
			return fmt.Sprintf("cos_deg(%s)", submatch[1])
		})
		expr = regexp.MustCompile(`tan\(([^)]+)\)`).ReplaceAllStringFunc(expr, func(m string) string {
			submatch := regexp.MustCompile(`tan\(([^)]+)\)`).FindStringSubmatch(m)
			return fmt.Sprintf("tan_deg(%s)", submatch[1])
		})
	}

	return expr
}

// executeAction 执行高级操作
func (s *CalculatorSkill) executeAction(args CalculatorArgs, result *CalculatorResult) error {
	result.Action = args.Action

	switch args.Action {
	case "unit_convert":
		return s.unitConvert(args, result)
	case "base_convert":
		return s.baseConvert(args, result)
	case "percentage":
		return s.percentage(args, result)
	case "statistics":
		return s.statistics(args, result)
	case "compare":
		return s.compare(args, result)
	case "factorial":
		return s.factorial(args, result)
	case "fibonacci":
		return s.fibonacci(args, result)
	case "is_prime":
		return s.isPrime(args, result)
	default:
		return fmt.Errorf("不支持的操作: %s", args.Action)
	}
}

// ============================================
// 数学表达式解析器
// ============================================

type mathParser struct {
	tokens       []string
	pos          int
	angleModeRad bool
}

func newMathParser(expr string, angleModeRad bool) *mathParser {
	return &mathParser{
		tokens:       tokenizeMath(expr),
		angleModeRad: angleModeRad,
	}
}

func tokenizeMath(expr string) []string {
	var tokens []string
	var current strings.Builder

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		switch ch {
		case ' ', '\t', '\n':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case '+', '-', '*', '/', '%', '^', '(', ')', ',':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			// 处理负数
			if ch == '-' && (len(tokens) == 0 || tokens[len(tokens)-1] == "(" ||
				tokens[len(tokens)-1] == "," || tokens[len(tokens)-1] == "+" ||
				tokens[len(tokens)-1] == "-" || tokens[len(tokens)-1] == "*" ||
				tokens[len(tokens)-1] == "/" || tokens[len(tokens)-1] == "^") {
				current.WriteByte(ch)
			} else {
				tokens = append(tokens, string(ch))
			}
		case '.':
			current.WriteByte(ch)
		default:
			if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
				current.WriteByte(ch)
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (p *mathParser) parse() (float64, error) {
	return p.parseExpression()
}

func (p *mathParser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos]
		if op != "+" && op != "-" {
			break
		}
		p.pos++

		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}

		if op == "+" {
			left = left + right
		} else {
			left = left - right
		}
	}

	return left, nil
}

func (p *mathParser) parseTerm() (float64, error) {
	left, err := p.parsePower()
	if err != nil {
		return 0, err
	}

	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos]
		if op != "*" && op != "/" && op != "%" {
			break
		}
		p.pos++

		right, err := p.parsePower()
		if err != nil {
			return 0, err
		}

		switch op {
		case "*":
			left = left * right
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("除数不能为零")
			}
			left = left / right
		case "%":
			left = math.Mod(left, right)
		}
	}

	return left, nil
}

func (p *mathParser) parsePower() (float64, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}

	if p.pos < len(p.tokens) && p.tokens[p.pos] == "^" {
		p.pos++
		right, err := p.parsePower() // 右结合
		if err != nil {
			return 0, err
		}
		left = math.Pow(left, right)
	}

	return left, nil
}

func (p *mathParser) parseUnary() (float64, error) {
	if p.pos < len(p.tokens) {
		if p.tokens[p.pos] == "-" {
			p.pos++
			val, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			return -val, nil
		}
		if p.tokens[p.pos] == "+" {
			p.pos++
			return p.parseUnary()
		}
	}
	return p.parsePrimary()
}

func (p *mathParser) parsePrimary() (float64, error) {
	if p.pos >= len(p.tokens) {
		return 0, fmt.Errorf("表达式不完整")
	}

	token := p.tokens[p.pos]

	// 处理括号
	if token == "(" {
		p.pos++
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return 0, fmt.Errorf("缺少右括号")
		}
		p.pos++
		return val, nil
	}

	// 处理函数
	if isFunction(token) {
		return p.parseFunction(token)
	}

	// 处理数字
	p.pos++
	val, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的数字或标识符: %s", token)
	}
	return val, nil
}

func (p *mathParser) parseFunction(name string) (float64, error) {
	p.pos++ // 跳过函数名

	if p.pos >= len(p.tokens) || p.tokens[p.pos] != "(" {
		return 0, fmt.Errorf("函数 %s 缺少左括号", name)
	}
	p.pos++ // 跳过 (

	var args []float64
	for {
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		args = append(args, val)

		if p.pos >= len(p.tokens) {
			return 0, fmt.Errorf("函数 %s 缺少右括号", name)
		}

		if p.tokens[p.pos] == ")" {
			p.pos++
			break
		}
		if p.tokens[p.pos] == "," {
			p.pos++
			continue
		}
		return 0, fmt.Errorf("函数参数格式错误")
	}

	return p.callFunction(name, args)
}

func (p *mathParser) callFunction(name string, args []float64) (float64, error) {
	switch name {
	// 三角函数
	case "sin":
		if len(args) != 1 {
			return 0, fmt.Errorf("sin 需要1个参数")
		}
		return math.Sin(args[0]), nil
	case "cos":
		if len(args) != 1 {
			return 0, fmt.Errorf("cos 需要1个参数")
		}
		return math.Cos(args[0]), nil
	case "tan":
		if len(args) != 1 {
			return 0, fmt.Errorf("tan 需要1个参数")
		}
		return math.Tan(args[0]), nil
	case "asin":
		if len(args) != 1 {
			return 0, fmt.Errorf("asin 需要1个参数")
		}
		return math.Asin(args[0]), nil
	case "acos":
		if len(args) != 1 {
			return 0, fmt.Errorf("acos 需要1个参数")
		}
		return math.Acos(args[0]), nil
	case "atan":
		if len(args) != 1 {
			return 0, fmt.Errorf("atan 需要1个参数")
		}
		return math.Atan(args[0]), nil
	case "sin_deg":
		if len(args) != 1 {
			return 0, fmt.Errorf("sin 需要1个参数")
		}
		return math.Sin(args[0] * math.Pi / 180), nil
	case "cos_deg":
		if len(args) != 1 {
			return 0, fmt.Errorf("cos 需要1个参数")
		}
		return math.Cos(args[0] * math.Pi / 180), nil
	case "tan_deg":
		if len(args) != 1 {
			return 0, fmt.Errorf("tan 需要1个参数")
		}
		return math.Tan(args[0] * math.Pi / 180), nil

	// 双曲函数
	case "sinh":
		if len(args) != 1 {
			return 0, fmt.Errorf("sinh 需要1个参数")
		}
		return math.Sinh(args[0]), nil
	case "cosh":
		if len(args) != 1 {
			return 0, fmt.Errorf("cosh 需要1个参数")
		}
		return math.Cosh(args[0]), nil
	case "tanh":
		if len(args) != 1 {
			return 0, fmt.Errorf("tanh 需要1个参数")
		}
		return math.Tanh(args[0]), nil

	// 对数函数
	case "log":
		if len(args) != 1 {
			return 0, fmt.Errorf("log 需要1个参数")
		}
		return math.Log10(args[0]), nil
	case "ln":
		if len(args) != 1 {
			return 0, fmt.Errorf("ln 需要1个参数")
		}
		return math.Log(args[0]), nil
	case "log2":
		if len(args) != 1 {
			return 0, fmt.Errorf("log2 需要1个参数")
		}
		return math.Log2(args[0]), nil

	// 其他函数
	case "sqrt":
		if len(args) != 1 {
			return 0, fmt.Errorf("sqrt 需要1个参数")
		}
		if args[0] < 0 {
			return 0, fmt.Errorf("sqrt 参数不能为负数")
		}
		return math.Sqrt(args[0]), nil
	case "cbrt":
		if len(args) != 1 {
			return 0, fmt.Errorf("cbrt 需要1个参数")
		}
		return math.Cbrt(args[0]), nil
	case "abs":
		if len(args) != 1 {
			return 0, fmt.Errorf("abs 需要1个参数")
		}
		return math.Abs(args[0]), nil
	case "floor":
		if len(args) != 1 {
			return 0, fmt.Errorf("floor 需要1个参数")
		}
		return math.Floor(args[0]), nil
	case "ceil":
		if len(args) != 1 {
			return 0, fmt.Errorf("ceil 需要1个参数")
		}
		return math.Ceil(args[0]), nil
	case "round":
		if len(args) != 1 {
			return 0, fmt.Errorf("round 需要1个参数")
		}
		return math.Round(args[0]), nil
	case "exp":
		if len(args) != 1 {
			return 0, fmt.Errorf("exp 需要1个参数")
		}
		return math.Exp(args[0]), nil
	case "pow":
		if len(args) != 2 {
			return 0, fmt.Errorf("pow 需要2个参数")
		}
		return math.Pow(args[0], args[1]), nil
	case "max":
		if len(args) < 2 {
			return 0, fmt.Errorf("max 至少需要2个参数")
		}
		result := args[0]
		for _, v := range args[1:] {
			if v > result {
				result = v
			}
		}
		return result, nil
	case "min":
		if len(args) < 2 {
			return 0, fmt.Errorf("min 至少需要2个参数")
		}
		result := args[0]
		for _, v := range args[1:] {
			if v < result {
				result = v
			}
		}
		return result, nil

	default:
		return 0, fmt.Errorf("未知函数: %s", name)
	}
}

func isFunction(name string) bool {
	functions := []string{
		"sin", "cos", "tan", "asin", "acos", "atan",
		"sin_deg", "cos_deg", "tan_deg",
		"sinh", "cosh", "tanh",
		"log", "ln", "log2",
		"sqrt", "cbrt", "abs", "floor", "ceil", "round", "exp",
		"pow", "max", "min",
	}
	for _, f := range functions {
		if f == name {
			return true
		}
	}
	return false
}

// ============================================
// 高级操作实现
// ============================================

// unitConvert 单位转换
func (s *CalculatorSkill) unitConvert(args CalculatorArgs, result *CalculatorResult) error {
	value, ok := toFloat(args.Value)
	if !ok {
		return fmt.Errorf("无效的数值")
	}

	converter := getUnitConverter()
	converted, err := converter.Convert(value, args.FromUnit, args.ToUnit)
	if err != nil {
		return err
	}

	result.Result = s.round(converted, result.Precision)
	result.Data = map[string]interface{}{
		"original_value": value,
		"from_unit":      args.FromUnit,
		"to_unit":        args.ToUnit,
		"converted":      converted,
	}
	return nil
}

// baseConvert 进制转换
func (s *CalculatorSkill) baseConvert(args CalculatorArgs, result *CalculatorResult) error {
	value, ok := args.Value.(string)
	if !ok {
		// 尝试转数字
		if num, ok := toFloat(args.Value); ok {
			value = fmt.Sprintf("%.0f", num)
		} else {
			return fmt.Errorf("value 必须是字符串或数字")
		}
	}

	fromBase := args.FromBase
	toBase := args.ToBase

	if fromBase < 2 || fromBase > 36 {
		return fmt.Errorf("源进制必须在2-36之间")
	}
	if toBase < 2 || toBase > 36 {
		return fmt.Errorf("目标进制必须在2-36之间")
	}

	// 转换为十进制
	decimal, err := strconv.ParseInt(value, fromBase, 64)
	if err != nil {
		return fmt.Errorf("无法解析数值: %v", err)
	}

	// 转换为目标进制
	converted := strconv.FormatInt(decimal, toBase)

	result.Result = converted
	result.Data = map[string]interface{}{
		"original":  value,
		"from_base": fromBase,
		"to_base":   toBase,
		"decimal":   decimal,
		"converted": converted,
	}
	return nil
}

// percentage 百分比计算
func (s *CalculatorSkill) percentage(args CalculatorArgs, result *CalculatorResult) error {
	value1, ok1 := toFloat(args.Value1)
	value2, ok2 := toFloat(args.Value2)

	if !ok1 || !ok2 {
		return fmt.Errorf("需要提供有效的value1和value2")
	}

	if value2 == 0 {
		return fmt.Errorf("value2 不能为零")
	}

	percent := (value1 / value2) * 100
	change := value2 - value1
	changePercent := 0.0
	if value1 != 0 {
		changePercent = (change / value1) * 100
	}

	result.Data = map[string]interface{}{
		"value1":         value1,
		"value2":         value2,
		"percentage":     s.round(percent, result.Precision),
		"difference":     s.round(change, result.Precision),
		"change_percent": s.round(changePercent, result.Precision),
		"formatted":      fmt.Sprintf("%.2f%%", percent),
	}
	result.Result = s.round(percent, result.Precision)
	return nil
}

// statistics 统计计算
func (s *CalculatorSkill) statistics(args CalculatorArgs, result *CalculatorResult) error {
	var numbers []float64

	switch v := args.Value.(type) {
	case []interface{}:
		for _, item := range v {
			if num, ok := toFloat(item); ok {
				numbers = append(numbers, num)
			}
		}
	case string:
		// 解析逗号分隔的数字
		parts := strings.Split(v, ",")
		for _, part := range parts {
			if num, ok := toFloat(strings.TrimSpace(part)); ok {
				numbers = append(numbers, num)
			}
		}
	default:
		return fmt.Errorf("value 必须是数字数组或逗号分隔的字符串")
	}

	if len(numbers) == 0 {
		return fmt.Errorf("没有有效的数字")
	}

	// 计算统计值
	sum := 0.0
	min := numbers[0]
	max := numbers[0]

	for _, n := range numbers {
		sum += n
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	mean := sum / float64(len(numbers))

	// 计算方差和标准差
	variance := 0.0
	for _, n := range numbers {
		variance += math.Pow(n-mean, 2)
	}
	variance /= float64(len(numbers))
	stdDev := math.Sqrt(variance)

	// 计算中位数
	sorted := make([]float64, len(numbers))
	copy(sorted, numbers)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	median := sorted[len(sorted)/2]
	if len(sorted)%2 == 0 {
		median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	}

	result.Data = map[string]interface{}{
		"count":    len(numbers),
		"sum":      s.round(sum, result.Precision),
		"mean":     s.round(mean, result.Precision),
		"median":   s.round(median, result.Precision),
		"min":      s.round(min, result.Precision),
		"max":      s.round(max, result.Precision),
		"range":    s.round(max-min, result.Precision),
		"variance": s.round(variance, result.Precision),
		"std_dev":  s.round(stdDev, result.Precision),
	}
	result.Result = s.round(mean, result.Precision)
	return nil
}

// compare 数值比较
func (s *CalculatorSkill) compare(args CalculatorArgs, result *CalculatorResult) error {
	value1, ok1 := toFloat(args.Value1)
	value2, ok2 := toFloat(args.Value2)

	if !ok1 || !ok2 {
		return fmt.Errorf("需要提供有效的value1和value2")
	}

	diff := value1 - value2
	absDiff := math.Abs(diff)

	var relation string
	if value1 < value2 {
		relation = "less_than"
	} else if value1 > value2 {
		relation = "greater_than"
	} else {
		relation = "equal"
	}

	result.Data = map[string]interface{}{
		"value1":     value1,
		"value2":     value2,
		"difference": s.round(diff, result.Precision),
		"abs_diff":   s.round(absDiff, result.Precision),
		"relation":   relation,
		"ratio":      s.round(value1/value2, result.Precision),
	}
	result.Result = relation
	return nil
}

// factorial 阶乘
func (s *CalculatorSkill) factorial(args CalculatorArgs, result *CalculatorResult) error {
	n, ok := toFloat(args.Value)
	if !ok {
		return fmt.Errorf("无效的数值")
	}

	if n < 0 || n != math.Floor(n) {
		return fmt.Errorf("阶乘只能计算非负整数")
	}

	if n > 170 {
		return fmt.Errorf("数值过大，超出计算范围")
	}

	fact := uint64(1)
	for i := uint64(2); i <= uint64(n); i++ {
		fact *= i
	}

	result.Result = fact
	result.Data = map[string]interface{}{
		"n":         int(n),
		"factorial": fact,
	}
	return nil
}

// fibonacci 斐波那契
func (s *CalculatorSkill) fibonacci(args CalculatorArgs, result *CalculatorResult) error {
	n, ok := toFloat(args.Value)
	if !ok {
		return fmt.Errorf("无效的数值")
	}

	if n < 0 || n != math.Floor(n) {
		return fmt.Errorf("只能计算非负整数")
	}

	if n > 92 {
		return fmt.Errorf("数值过大，超出uint64范围")
	}

	fib := fibonacciN(int(n))

	result.Result = fib
	result.Data = map[string]interface{}{
		"n":         int(n),
		"fibonacci": fib,
	}
	return nil
}

func fibonacciN(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	a, b := uint64(0), uint64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// isPrime 判断素数
func (s *CalculatorSkill) isPrime(args CalculatorArgs, result *CalculatorResult) error {
	n, ok := toFloat(args.Value)
	if !ok {
		return fmt.Errorf("无效的数值")
	}

	if n < 2 || n != math.Floor(n) {
		result.Result = false
		result.Data = map[string]interface{}{
			"n":        int(n),
			"is_prime": false,
			"reason":   "素数必须是大于等于2的整数",
		}
		return nil
	}

	num := int(n)
	isPrime := true
	var divisor int

	if num == 2 {
		isPrime = true
	} else if num%2 == 0 {
		isPrime = false
		divisor = 2
	} else {
		for i := 3; i*i <= num; i += 2 {
			if num%i == 0 {
				isPrime = false
				divisor = i
				break
			}
		}
	}

	result.Result = isPrime
	result.Data = map[string]interface{}{
		"n":        num,
		"is_prime": isPrime,
	}
	if !isPrime && divisor > 0 {
		result.Data["divisor"] = divisor
		result.Data["reason"] = fmt.Sprintf("可被 %d 整除", divisor)
	}
	return nil
}

// ============================================
// 单位转换器
// ============================================

type unitConverter struct {
	categories map[string]map[string]float64
}

func getUnitConverter() *unitConverter {
	return &unitConverter{
		categories: map[string]map[string]float64{
			// 长度 (基准: 米)
			"length": {
				"m":  1,
				"km": 1000,
				"cm": 0.01,
				"mm": 0.001,
				"mi": 1609.344,
				"yd": 0.9144,
				"ft": 0.3048,
				"in": 0.0254,
				"里":  500,
				"丈":  3.333333333,
				"尺":  0.333333333,
				"寸":  0.033333333,
			},
			// 重量 (基准: 千克)
			"weight": {
				"kg": 1,
				"g":  0.001,
				"mg": 0.000001,
				"lb": 0.45359237,
				"oz": 0.028349523,
				"t":  1000,
				"斤":  0.5,
				"两":  0.05,
				"钱":  0.005,
			},
			// 温度 (特殊处理)
			"temperature": {},
			// 面积 (基准: 平方米)
			"area": {
				"m2":  1,
				"km2": 1000000,
				"cm2": 0.0001,
				"mm2": 0.000001,
				"ha":  10000,
				"ac":  4046.8564224,
				"ft2": 0.09290304,
				"in2": 0.00064516,
				"亩":   666.666666667,
				"顷":   66666.6666667,
			},
			// 体积 (基准: 升)
			"volume": {
				"l":     1,
				"ml":    0.001,
				"m3":    1000,
				"cm3":   0.001,
				"gal":   3.785411784,
				"qt":    0.946352946,
				"pt":    0.473176473,
				"cup":   0.2365882365,
				"oz_fl": 0.0295735296,
			},
			// 时间 (基准: 秒)
			"time": {
				"s":   1,
				"ms":  0.001,
				"min": 60,
				"h":   3600,
				"d":   86400,
				"w":   604800,
				"mo":  2592000, // 30天
				"y":   31536000,
			},
			// 速度 (基准: 米/秒)
			"speed": {
				"m/s":  1,
				"km/h": 0.277777778,
				"mph":  0.44704,
				"kn":   0.514444444,
				"ft/s": 0.3048,
				"马赫":   340.29,
			},
			// 数据存储 (基准: 字节)
			"data": {
				"B":  1,
				"KB": 1024,
				"MB": 1048576,
				"GB": 1073741824,
				"TB": 1099511627776,
				"PB": 1125899906842624,
			},
		},
	}
}

func (c *unitConverter) Convert(value float64, from, to string) (float64, error) {
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	if from == to {
		return value, nil
	}

	// 温度特殊处理
	if isTemperature(from) || isTemperature(to) {
		return c.convertTemperature(value, from, to)
	}

	// 查找单位所属类别
	fromCategory, fromFactor := c.findUnit(from)
	toCategory, toFactor := c.findUnit(to)

	if fromCategory == "" {
		return 0, fmt.Errorf("未知单位: %s", from)
	}
	if toCategory == "" {
		return 0, fmt.Errorf("未知单位: %s", to)
	}
	if fromCategory != toCategory {
		return 0, fmt.Errorf("无法转换不同类别的单位: %s -> %s", from, to)
	}

	// 转换: value * fromFactor / toFactor
	return value * fromFactor / toFactor, nil
}

func (c *unitConverter) findUnit(unit string) (string, float64) {
	for category, units := range c.categories {
		if factor, ok := units[unit]; ok {
			return category, factor
		}
	}
	return "", 0
}

func isTemperature(unit string) bool {
	unit = strings.ToLower(unit)
	return unit == "c" || unit == "f" || unit == "k" || unit == "℃" || unit == "℉"
}

func (c *unitConverter) convertTemperature(value float64, from, to string) (float64, error) {
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	// 先转换为摄氏度
	var celsius float64
	switch from {
	case "c", "℃":
		celsius = value
	case "f", "℉":
		celsius = (value - 32) * 5 / 9
	case "k":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("未知温度单位: %s", from)
	}

	// 再转换为目标单位
	switch to {
	case "c", "℃":
		return celsius, nil
	case "f", "℉":
		return celsius*9/5 + 32, nil
	case "k":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("未知温度单位: %s", to)
	}
}

// ============================================
// 辅助函数
// ============================================

func (s *CalculatorSkill) round(val float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(val*multiplier) / multiplier
}

func toFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint64:
		return float64(v), true
	case uint32:
		return float64(v), true
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		return f, err == nil
	default:
		return 0, false
	}
}

// ToJSON 序列化
func (s *CalculatorSkill) ToJSON() (string, error) {
	data := map[string]interface{}{
		"name":        s.name,
		"descName":    s.descName,
		"description": s.description,
		"parameters":  s.parameters,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonBytes), nil
}

// FromJSON 反序列化
func (s *CalculatorSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string               `json:"name"`
		DescName    string               `json:"descName"`
		Description string               `json:"description"`
		Parameters  CalculatorParameters `json:"parameters"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}
	s.name = data.Name
	s.descName = data.DescName
	s.description = data.Description
	s.parameters = data.Parameters
	return nil
}
