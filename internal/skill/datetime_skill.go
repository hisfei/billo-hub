package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ============================================
// DateTimeSkill - 生产级日期时间技能
// ============================================

// DateTimeParameters 定义日期时间技能的参数结构
type DateTimeParameters struct {
	Type       string                  `json:"type"`
	Properties DateTimeParamProperties `json:"properties"`
	Required   []string                `json:"required"`
}

// DateTimeParamProperties 定义参数属性
type DateTimeParamProperties struct {
	Action   DateTimeParamProperty `json:"action"`
	Timezone DateTimeParamProperty `json:"timezone"`
	Format   DateTimeParamProperty `json:"format"`
	Value    DateTimeParamProperty `json:"value"`
	Value2   DateTimeParamProperty `json:"value2"`
	Unit     DateTimeParamProperty `json:"unit"`
	Amount   DateTimeParamProperty `json:"amount"`
}

// DateTimeParamProperty 定义单个参数属性
type DateTimeParamProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// DateTimeArgs 定义日期时间操作参数
type DateTimeArgs struct {
	Action   string `json:"action"`
	Timezone string `json:"timezone,omitempty"`
	Format   string `json:"format,omitempty"`
	Value    string `json:"value,omitempty"`
	Value2   string `json:"value2,omitempty"`
	Unit     string `json:"unit,omitempty"`
	Amount   int    `json:"amount,omitempty"`
}

// DateTimeResult 定义日期时间操作结果
type DateTimeResult struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Action  string                 `json:"action"`
}

// DateTimeSkill 实现了生产级日期时间技能
type DateTimeSkill struct {
	name        string
	descName    string
	description string
	parameters  DateTimeParameters
}

// NewDateTimeSkill 创建新的日期时间技能实例
func NewDateTimeSkill() *DateTimeSkill {
	return &DateTimeSkill{
		name:     "datetime",
		descName: "日期时间处理",
		description: `日期时间处理技能提供全面的时间和日期操作能力。

支持的操作：
- now: 获取当前时间（支持时区）
- format: 时间格式化
- parse: 解析时间字符串
- add: 时间加减运算
- diff: 计算时间差
- timezone: 时区转换
- timestamp: 时间戳转换
- countdown: 计算倒计时
- weekday: 获取星期几
- is_workday: 判断是否工作日

时区支持：
- 使用IANA时区名称，如 "Asia/Shanghai"、"America/New_York"
- 支持 "local" 表示本地时区
- 支持 "UTC" 表示协调世界时

格式化占位符（遵循Go格式）：
- 2006: 四位年份
- 01: 两位月份
- 02: 两位日期
- 15: 24小时制
- 04: 分钟
- 05: 秒
- Monday: 星期全名
- Jan: 月份缩写`,
		parameters: DateTimeParameters{
			Type: "object",
			Properties: DateTimeParamProperties{
				Action: DateTimeParamProperty{
					Type:        "string",
					Description: "操作类型",
					Enum: []string{
						"now", "format", "parse", "add", "diff",
						"timezone", "timestamp", "countdown",
						"weekday", "is_workday", "start_of", "end_of",
					},
				},
				Timezone: DateTimeParamProperty{
					Type:        "string",
					Description: "时区名称，如 Asia/Shanghai",
					Default:     "local",
				},
				Format: DateTimeParamProperty{
					Type:        "string",
					Description: "时间格式，如 2006-01-02 15:04:05",
					Default:     "2006-01-02 15:04:05",
				},
				Value: DateTimeParamProperty{
					Type:        "string",
					Description: "时间值字符串或时间戳",
				},
				Value2: DateTimeParamProperty{
					Type:        "string",
					Description: "第二个时间值（用于diff操作）",
				},
				Unit: DateTimeParamProperty{
					Type:        "string",
					Description: "时间单位",
					Enum:        []string{"year", "month", "day", "hour", "minute", "second", "week"},
				},
				Amount: DateTimeParamProperty{
					Type:        "integer",
					Description: "数量（用于add操作）",
				},
			},
			Required: []string{"action"},
		},
	}
}

// GetName 返回技能名称
func (s *DateTimeSkill) GetName() string {
	return s.name
}

// GetDescName 返回技能的描述性名称
func (s *DateTimeSkill) GetDescName() string {
	return s.descName
}

// GetDescription 返回技能的详细描述
func (s *DateTimeSkill) GetDescription() string {
	return s.description
}

// GetParameters 返回技能的参数定义
func (s *DateTimeSkill) GetParameters() any {
	return s.parameters
}

// Execute 执行日期时间操作
func (s *DateTimeSkill) Execute(ctx context.Context, args string) (string, error) {
	var dtArgs DateTimeArgs
	if err := json.Unmarshal([]byte(args), &dtArgs); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	result := &DateTimeResult{
		Action: dtArgs.Action,
		Data:   make(map[string]interface{}),
	}

	// 默认时区
	if dtArgs.Timezone == "" {
		dtArgs.Timezone = "local"
	}

	var err error
	switch dtArgs.Action {
	case "now":
		err = s.actionNow(dtArgs, result)
	case "format":
		err = s.actionFormat(dtArgs, result)
	case "parse":
		err = s.actionParse(dtArgs, result)
	case "add":
		err = s.actionAdd(dtArgs, result)
	case "diff":
		err = s.actionDiff(dtArgs, result)
	case "timezone":
		err = s.actionTimezone(dtArgs, result)
	case "timestamp":
		err = s.actionTimestamp(dtArgs, result)
	case "countdown":
		err = s.actionCountdown(dtArgs, result)
	case "weekday":
		err = s.actionWeekday(dtArgs, result)
	case "is_workday":
		err = s.actionIsWorkday(dtArgs, result)
	case "start_of":
		err = s.actionStartOf(dtArgs, result)
	case "end_of":
		err = s.actionEndOf(dtArgs, result)
	default:
		err = fmt.Errorf("不支持的操作: %s", dtArgs.Action)
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

// getLocation 获取时区
func (s *DateTimeSkill) getLocation(timezone string) (*time.Location, error) {
	switch strings.ToLower(timezone) {
	case "local", "":
		return time.Local, nil
	case "utc":
		return time.UTC, nil
	default:
		return time.LoadLocation(timezone)
	}
}

// parseTime 解析时间字符串
func (s *DateTimeSkill) parseTime(value, timezone string) (time.Time, error) {
	loc, err := s.getLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	// 尝试解析时间戳
	if ts, err := strconv.ParseInt(value, 10, 64); err == nil {
		// 判断是秒还是毫秒
		if ts > 1e12 {
			return time.UnixMilli(ts).In(loc), nil
		}
		return time.Unix(ts, 0).In(loc), nil
	}

	// 尝试多种格式解析
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"2006/01/02",
		"2006年01月02日",
		"2006年01月02日 15:04:05",
		"01/02/2006",
		"02-01-2006",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, value, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间: %s", value)
}

// actionNow 获取当前时间
func (s *DateTimeSkill) actionNow(args DateTimeArgs, result *DateTimeResult) error {
	loc, err := s.getLocation(args.Timezone)
	if err != nil {
		return err
	}

	now := time.Now().In(loc)
	format := args.Format
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	result.Data = map[string]interface{}{
		"datetime":     now.Format(format),
		"timestamp":    now.Unix(),
		"timestamp_ms": now.UnixMilli(),
		"date":         now.Format("2006-01-02"),
		"time":         now.Format("15:04:05"),
		"year":         now.Year(),
		"month":        int(now.Month()),
		"day":          now.Day(),
		"hour":         now.Hour(),
		"minute":       now.Minute(),
		"second":       now.Second(),
		"weekday":      now.Weekday().String(),
		"weekday_cn":   s.weekdayCN(now.Weekday()),
		"timezone":     now.Location().String(),
		"unix":         now.Unix(),
	}
	return nil
}

// actionFormat 格式化时间
func (s *DateTimeSkill) actionFormat(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" {
		return fmt.Errorf("需要提供value参数")
	}

	t, err := s.parseTime(args.Value, args.Timezone)
	if err != nil {
		return err
	}

	format := args.Format
	if format == "" {
		format = "2006-01-02 15:04:05"
	}

	result.Data = map[string]interface{}{
		"formatted": t.Format(format),
		"original":  args.Value,
	}
	return nil
}

// actionParse 解析时间字符串
func (s *DateTimeSkill) actionParse(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" {
		return fmt.Errorf("需要提供value参数")
	}

	t, err := s.parseTime(args.Value, args.Timezone)
	if err != nil {
		return err
	}

	result.Data = map[string]interface{}{
		"datetime":     t.Format("2006-01-02 15:04:05"),
		"timestamp":    t.Unix(),
		"timestamp_ms": t.UnixMilli(),
		"year":         t.Year(),
		"month":        int(t.Month()),
		"day":          t.Day(),
		"hour":         t.Hour(),
		"minute":       t.Minute(),
		"second":       t.Second(),
		"weekday":      t.Weekday().String(),
		"timezone":     t.Location().String(),
	}
	return nil
}

// actionAdd 时间加减
func (s *DateTimeSkill) actionAdd(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" {
		return fmt.Errorf("需要提供value参数")
	}
	if args.Unit == "" {
		return fmt.Errorf("需要提供unit参数")
	}

	t, err := s.parseTime(args.Value, args.Timezone)
	if err != nil {
		return err
	}

	amount := args.Amount
	var newTime time.Time

	switch args.Unit {
	case "year":
		newTime = t.AddDate(amount, 0, 0)
	case "month":
		newTime = t.AddDate(0, amount, 0)
	case "day":
		newTime = t.AddDate(0, 0, amount)
	case "week":
		newTime = t.AddDate(0, 0, amount*7)
	case "hour":
		newTime = t.Add(time.Duration(amount) * time.Hour)
	case "minute":
		newTime = t.Add(time.Duration(amount) * time.Minute)
	case "second":
		newTime = t.Add(time.Duration(amount) * time.Second)
	default:
		return fmt.Errorf("不支持的时间单位: %s", args.Unit)
	}

	result.Data = map[string]interface{}{
		"original":  t.Format("2006-01-02 15:04:05"),
		"result":    newTime.Format("2006-01-02 15:04:05"),
		"operation": fmt.Sprintf("%+d %s", amount, args.Unit),
		"timestamp": newTime.Unix(),
	}
	return nil
}

// actionDiff 计算时间差
func (s *DateTimeSkill) actionDiff(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" || args.Value2 == "" {
		return fmt.Errorf("需要提供value和value2参数")
	}

	t1, err := s.parseTime(args.Value, args.Timezone)
	if err != nil {
		return err
	}

	t2, err := s.parseTime(args.Value2, args.Timezone)
	if err != nil {
		return err
	}

	diff := t2.Sub(t1)
	absDiff := diff
	if diff < 0 {
		absDiff = -diff
	}

	// 计算各单位的差值
	years := t2.Year() - t1.Year()
	months := int(t2.Month()) - int(t1.Month())
	days := t2.Day() - t1.Day()

	// 调整月份和天数
	if days < 0 {
		months--
	}
	if months < 0 {
		years--
		months += 12
	}

	result.Data = map[string]interface{}{
		"value1":         t1.Format("2006-01-02 15:04:05"),
		"value2":         t2.Format("2006-01-02 15:04:05"),
		"diff_ns":        diff.Nanoseconds(),
		"diff_ms":        diff.Milliseconds(),
		"diff_sec":       int(diff.Seconds()),
		"diff_min":       int(diff.Minutes()),
		"diff_hour":      int(diff.Hours()),
		"diff_day":       int(absDiff.Hours() / 24),
		"diff_week":      int(absDiff.Hours() / 24 / 7),
		"years":          years,
		"months":         months,
		"total_days":     int(absDiff.Hours() / 24),
		"is_positive":    diff >= 0,
		"human_readable": s.humanizeDuration(absDiff),
	}
	return nil
}

// actionTimezone 时区转换
func (s *DateTimeSkill) actionTimezone(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" {
		return fmt.Errorf("需要提供value参数")
	}
	if args.Timezone == "" {
		return fmt.Errorf("需要提供目标时区")
	}

	t, err := s.parseTime(args.Value, "local")
	if err != nil {
		return err
	}

	targetLoc, err := s.getLocation(args.Timezone)
	if err != nil {
		return err
	}

	converted := t.In(targetLoc)

	result.Data = map[string]interface{}{
		"original":       t.Format("2006-01-02 15:04:05"),
		"original_tz":    t.Location().String(),
		"converted":      converted.Format("2006-01-02 15:04:05"),
		"converted_tz":   converted.Location().String(),
		"timezone":       args.Timezone,
		"offset_seconds": converted.Unix() - t.Unix(),
	}
	return nil
}

// actionTimestamp 时间戳转换
func (s *DateTimeSkill) actionTimestamp(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" {
		return fmt.Errorf("需要提供value参数")
	}

	// 判断是转换为时间戳还是从时间戳转换
	var t time.Time
	var err error

	if ts, e := strconv.ParseInt(args.Value, 10, 64); e == nil {
		// 输入是时间戳
		if ts > 1e12 {
			t = time.UnixMilli(ts)
		} else {
			t = time.Unix(ts, 0)
		}
		result.Data = map[string]interface{}{
			"timestamp":  args.Value,
			"datetime":   t.Format("2006-01-02 15:04:05"),
			"date":       t.Format("2006-01-02"),
			"time":       t.Format("15:04:05"),
			"iso8601":    t.Format(time.RFC3339),
			"unix":       t.Unix(),
			"unix_milli": t.UnixMilli(),
			"unix_nano":  t.UnixNano(),
		}
	} else {
		// 输入是时间字符串
		t, err = s.parseTime(args.Value, args.Timezone)
		if err != nil {
			return err
		}
		result.Data = map[string]interface{}{
			"datetime":   t.Format("2006-01-02 15:04:05"),
			"unix":       t.Unix(),
			"unix_milli": t.UnixMilli(),
			"unix_nano":  t.UnixNano(),
			"iso8601":    t.Format(time.RFC3339),
		}
	}
	return nil
}

// actionCountdown 计算倒计时
func (s *DateTimeSkill) actionCountdown(args DateTimeArgs, result *DateTimeResult) error {
	if args.Value == "" {
		return fmt.Errorf("需要提供value参数（目标时间）")
	}

	target, err := s.parseTime(args.Value, args.Timezone)
	if err != nil {
		return err
	}

	now := time.Now()
	loc, _ := s.getLocation(args.Timezone)
	now = now.In(loc)

	diff := target.Sub(now)

	result.Data = map[string]interface{}{
		"target_time":    target.Format("2006-01-02 15:04:05"),
		"current_time":   now.Format("2006-01-02 15:04:05"),
		"is_future":      diff > 0,
		"total_seconds":  int(diff.Seconds()),
		"total_minutes":  int(diff.Minutes()),
		"total_hours":    int(diff.Hours()),
		"total_days":     int(diff.Hours() / 24),
		"remaining":      s.formatCountdown(diff),
		"human_readable": s.humanizeCountdown(diff),
	}
	return nil
}

// actionWeekday 获取星期几
func (s *DateTimeSkill) actionWeekday(args DateTimeArgs, result *DateTimeResult) error {
	var t time.Time
	var err error

	if args.Value == "" {
		t = time.Now()
	} else {
		t, err = s.parseTime(args.Value, args.Timezone)
		if err != nil {
			return err
		}
	}

	result.Data = map[string]interface{}{
		"date":        t.Format("2006-01-02"),
		"weekday":     t.Weekday().String(),
		"weekday_cn":  s.weekdayCN(t.Weekday()),
		"weekday_num": int(t.Weekday()),
		"is_weekend":  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
	}
	return nil
}

// actionIsWorkday 判断是否工作日
func (s *DateTimeSkill) actionIsWorkday(args DateTimeArgs, result *DateTimeResult) error {
	var t time.Time
	var err error

	if args.Value == "" {
		t = time.Now()
	} else {
		t, err = s.parseTime(args.Value, args.Timezone)
		if err != nil {
			return err
		}
	}

	weekday := t.Weekday()
	isWeekend := weekday == time.Saturday || weekday == time.Sunday

	// 简单判断：周一到周五为工作日
	// 实际应用中可以加入节假日判断
	result.Data = map[string]interface{}{
		"date":       t.Format("2006-01-02"),
		"weekday":    weekday.String(),
		"weekday_cn": s.weekdayCN(weekday),
		"is_workday": !isWeekend,
		"is_weekend": isWeekend,
	}
	return nil
}

// actionStartOf 获取时间段的开始
func (s *DateTimeSkill) actionStartOf(args DateTimeArgs, result *DateTimeResult) error {
	if args.Unit == "" {
		return fmt.Errorf("需要提供unit参数")
	}

	var t time.Time
	var err error

	if args.Value == "" {
		t = time.Now()
	} else {
		t, err = s.parseTime(args.Value, args.Timezone)
		if err != nil {
			return err
		}
	}

	var start time.Time
	switch args.Unit {
	case "day":
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "week":
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
	case "month":
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case "year":
		start = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	case "hour":
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case "minute":
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	default:
		return fmt.Errorf("不支持的单位: %s", args.Unit)
	}

	result.Data = map[string]interface{}{
		"original": t.Format("2006-01-02 15:04:05"),
		"start_of": start.Format("2006-01-02 15:04:05"),
		"unit":     args.Unit,
	}
	return nil
}

// actionEndOf 获取时间段的结束
func (s *DateTimeSkill) actionEndOf(args DateTimeArgs, result *DateTimeResult) error {
	if args.Unit == "" {
		return fmt.Errorf("需要提供unit参数")
	}

	var t time.Time
	var err error

	if args.Value == "" {
		t = time.Now()
	} else {
		t, err = s.parseTime(args.Value, args.Timezone)
		if err != nil {
			return err
		}
	}

	var end time.Time
	switch args.Unit {
	case "day":
		end = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	case "week":
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		end = time.Date(t.Year(), t.Month(), t.Day()-weekday+7, 23, 59, 59, 999999999, t.Location())
	case "month":
		end = time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, t.Location())
	case "year":
		end = time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
	case "hour":
		end = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 999999999, t.Location())
	case "minute":
		end = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, 999999999, t.Location())
	default:
		return fmt.Errorf("不支持的单位: %s", args.Unit)
	}

	result.Data = map[string]interface{}{
		"original": t.Format("2006-01-02 15:04:05"),
		"end_of":   end.Format("2006-01-02 15:04:05"),
		"unit":     args.Unit,
	}
	return nil
}

// 辅助函数
func (s *DateTimeSkill) weekdayCN(weekday time.Weekday) string {
	names := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return names[int(weekday)]
}

func (s *DateTimeSkill) formatCountdown(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d天 %02d:%02d:%02d", days, hours, minutes, seconds)
}

func (s *DateTimeSkill) humanizeCountdown(d time.Duration) string {
	isNegative := d < 0
	if isNegative {
		d = -d
		return fmt.Sprintf("已过去 %s", s.humanizeDuration(d))
	}
	return fmt.Sprintf("还有 %s", s.humanizeDuration(d))
}

func (s *DateTimeSkill) humanizeDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 && days == 0 {
		parts = append(parts, fmt.Sprintf("%d分钟", minutes))
	}
	if len(parts) == 0 {
		return "不到1分钟"
	}
	return strings.Join(parts, " ")
}

// ToJSON 序列化
func (s *DateTimeSkill) ToJSON() (string, error) {
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
func (s *DateTimeSkill) FromJSON(jsonStr string) error {
	var data struct {
		Name        string             `json:"name"`
		DescName    string             `json:"descName"`
		Description string             `json:"description"`
		Parameters  DateTimeParameters `json:"parameters"`
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
