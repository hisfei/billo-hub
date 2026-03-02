package skill

import (
	"billohub/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// --- 结构化返回结果 ---
type CreateTaskResult struct {
	Code     string `json:"code"`
	TaskID   string `json:"task_id"`
	Msg      string `json:"msg"`
	TaskType string `json:"task_type"`
	Spec     string `json:"spec"`
	NextRun  string `json:"next_run"`
	AgentID  string `json:"agent_id"`
}

type QueryTaskResult struct {
	Code  string                `json:"code"`
	Msg   string                `json:"msg"`
	Tasks []model.ScheduledTask `json:"tasks"`
}

type DeleteTaskResult struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	TaskID string `json:"task_id"`
}

// --- ScheduleTaskSkill ---
type ScheduleTaskSkill struct {
	agentId  string
	cron     *cron.Cron
	parser   cron.Parser
	callback model.AgentChatCallback
	storage  model.AgentStorage
	mu       sync.RWMutex
	// 用于在内存中快速查找 cron 的 EntryID
	cronEntryMap map[string]cron.EntryID
}

func (s *ScheduleTaskSkill) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}

func (s *ScheduleTaskSkill) ToJSON() (string, error) {
	res, err := json.Marshal(s)
	return string(res), err
}

// NewScheduleTaskSkill 构造函数，注入了 storage 依赖
func NewScheduleTaskSkill(agentId string, callback model.AgentChatCallback, storage model.AgentStorage) *ScheduleTaskSkill {

	// ✅ 修复：使用兼容 5位(标准) + 6位(带秒) 的解析器
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Second)

	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithParser(parser),
	)
	c.Start()

	return &ScheduleTaskSkill{
		agentId:      agentId,
		cron:         c,
		parser:       parser,
		callback:     callback,
		storage:      storage,
		cronEntryMap: make(map[string]cron.EntryID),
	}
}

// --- Skill 接口实现 ---
func (s *ScheduleTaskSkill) GetDescName() string { return "Scheduled Task" }
func (s *ScheduleTaskSkill) GetName() string     { return "schedule_task" }
func (s *ScheduleTaskSkill) GetDescription() string {
	return "Supports creating/querying/deleting periodic (cron) or one-time (once) scheduled tasks, and returns a structured JSON result containing a TaskID for subsequent operations."
}
func (s *ScheduleTaskSkill) GetParameters() any {
	// ... (参数定义保持不变)
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"cmd":       map[string]interface{}{"type": "string", "enum": []string{"create", "query", "delete"}, "description": "Command type: create, query, delete"},
			"task_type": map[string]interface{}{"type": "string", "enum": []string{"cron", "once"}, "description": "Required for create: cron or once"},
			"spec":      map[string]interface{}{"type": "string", "description": "Required for create: cron spec (e.g., '0 8 * * *') or duration (e.g., '10s')"},
			"message":   map[string]interface{}{"type": "string", "description": "Required for create: callback message"},
			"task_id":   map[string]interface{}{"type": "string", "description": "Required for delete; optional for query"},
		},
		"required": []string{"cmd"},
	}
}

type taskParam struct {
	Cmd      string `json:"cmd"`
	TaskType string `json:"task_type"`
	Spec     string `json:"spec"`
	Message  string `json:"message"`
	TaskID   string `json:"task_id"`
}

func (s *ScheduleTaskSkill) Execute(ctx context.Context, args string) (string, error) {
	var params taskParam
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return s.errorToJSON(fmt.Sprintf("parameter parsing failed: %v", err))
	}

	switch params.Cmd {
	case "create":
		return s.createTask(ctx, params)
	case "query":
		return s.queryTask(ctx, params.TaskID)
	case "delete":
		return s.deleteTask(ctx, params.TaskID)
	default:
		return s.errorToJSON(fmt.Sprintf("unsupported command type: %s", params.Cmd))
	}
}

// --- 核心业务逻辑 ---

// RegisterTask 用于从数据库加载任务并注册到 cron
func (s *ScheduleTaskSkill) RegisterTask(ctx context.Context, task model.ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var cronSpec string
	isOnce := task.TaskType == "once"

	if isOnce {
		// 对于重启的一次性任务，如果执行时间已过，则直接删除
		nextRun, err := time.Parse("2006-01-02 15:04:05", task.Spec)
		if err != nil || time.Now().After(nextRun) {
			return s.storage.DeleteScheduledTask(task.ID)
		}
		cronSpec = nextRun.Format("05 04 15 02 01 *")
	} else {
		cronSpec = task.Spec
	}

	entryID, err := s.cron.AddFunc(cronSpec, func() {
		if s.callback != nil {
			// 创建一个新的 context，避免使用旧的
			taskCtx := context.WithValue(context.Background(), model.CtxMessageKey, model.CtxMessage{
				AgentID: task.AgentID,
				ChatId:  task.ChatID,
			})
			s.callback.Chat(taskCtx, fmt.Sprintf("[%s] Scheduled task triggered: %s", time.Now().Format("2006-01-02 15:04:05"), task.Message))
		}
		if isOnce {
			s.deleteTask(context.Background(), task.ID)
		}
	})

	if err != nil {
		return err
	}

	s.cronEntryMap[task.ID] = entryID
	return nil
}

func (s *ScheduleTaskSkill) createTask(ctx context.Context, params taskParam) (string, error) {
	if params.TaskType == "" || params.Spec == "" || params.Message == "" {
		return s.errorToJSON("task_type, spec, and message are required to create a task")
	}

	ctxMsg, _ := ctx.Value(model.CtxMessageKey).(model.CtxMessage)
	taskID := fmt.Sprintf("task_%s_%d", s.agentId, time.Now().UnixMilli())

	task := &model.ScheduledTask{
		ID:       taskID,
		AgentID:  s.agentId,
		ChatID:   ctxMsg.ChatId,
		TaskType: params.TaskType,
		Spec:     params.Spec,
		Message:  params.Message,
		IsActive: true,
	}

	// 注册到 cron
	if err := s.RegisterTask(ctx, *task); err != nil {
		return s.errorToJSON(fmt.Sprintf("failed to register task: %v", err))
	}

	// 持久化到数据库
	if err := s.storage.SaveScheduledTask(task); err != nil {
		s.cron.Remove(s.cronEntryMap[taskID]) // 如果数据库保存失败，回滚 cron 注册
		return s.errorToJSON(fmt.Sprintf("failed to save task to database: %v", err))
	}

	s.mu.RLock()
	entry := s.cron.Entry(s.cronEntryMap[taskID])
	s.mu.RUnlock()

	successResult := CreateTaskResult{
		Code:     "success",
		TaskID:   taskID,
		Msg:      "Task created successfully",
		TaskType: params.TaskType,
		Spec:     params.Spec,
		NextRun:  entry.Next.Format("2006-01-02 15:04:05"),
		AgentID:  s.agentId,
	}
	return s.resultToJSON(successResult)
}

func (s *ScheduleTaskSkill) queryTask(ctx context.Context, taskID string) (string, error) {
	// 在这个设计中，查询直接访问数据库，因为它是单一事实来源
	tasks, err := s.storage.LoadAllScheduledTasks()
	if err != nil {
		return s.errorToJSON(fmt.Sprintf("failed to query tasks: %v", err))
	}

	// 如果指定了 taskID，则过滤
	if taskID != "" {
		for _, t := range tasks {
			if t.ID == taskID {
				tasks = []model.ScheduledTask{t}
				break
			}
		}
	}

	successResult := QueryTaskResult{
		Code:  "success",
		Msg:   fmt.Sprintf("Found %d tasks", len(tasks)),
		Tasks: tasks,
	}
	return s.resultToJSON(successResult)
}

func (s *ScheduleTaskSkill) deleteTask(ctx context.Context, taskID string) (string, error) {
	if taskID == "" {
		return s.errorToJSON("task_id is required to delete a task")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 从 cron 中移除
	if entryID, ok := s.cronEntryMap[taskID]; ok {
		s.cron.Remove(entryID)
		delete(s.cronEntryMap, taskID)
	}

	// 从数据库中删除
	if err := s.storage.DeleteScheduledTask(taskID); err != nil {
		return s.errorToJSON(fmt.Sprintf("failed to delete task from database: %v", err))
	}

	successResult := DeleteTaskResult{
		Code:   "success",
		Msg:    "Task deleted successfully",
		TaskID: taskID,
	}
	return s.resultToJSON(successResult)
}

// --- 辅助函数 ---
func (s *ScheduleTaskSkill) errorToJSON(msg string) (string, error) {
	errResult := map[string]string{"code": "error", "msg": msg}
	errJSON, _ := json.Marshal(errResult)
	return string(errJSON), fmt.Errorf(msg)
}

func (s *ScheduleTaskSkill) resultToJSON(v interface{}) (string, error) {
	resultJSON, err := json.Marshal(v)
	if err != nil {
		return s.errorToJSON(fmt.Sprintf("failed to marshal result: %v", err))
	}
	return string(resultJSON), nil
}
