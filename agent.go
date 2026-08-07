package adaptix

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Adaptix-Framework/axsafe"
)

type DeliveryFunc func(agentId int64, taskData TaskData) error

var ErrAgentRemoved = fmt.Errorf("agent removed")

type Agent struct {
	mu   sync.RWMutex
	data AgentData
	Fn   AgentFunctions

	active  atomic.Bool
	removed atomic.Bool

	HostedQueue  *axsafe.PriorityQueue
	RunningTasks axsafe.Map[int64, TaskData]
	RunningJobs  axsafe.Map[int64, *axsafe.Slice]
	pivotMu      sync.RWMutex
	PivotParent  *PivotData
	PivotChilds  *axsafe.Slice

	cmdGroupMu      sync.RWMutex
	cmdGroupEnabled map[string]bool
}

func NewAgent(data AgentData, fn AgentFunctions) *Agent {
	if fn.CreateCommand == nil {
		fn.CreateCommand = func(_ AgentData, _ map[string]any) (TaskData, ConsoleMessageData, error) {
			return TaskData{}, ConsoleMessageData{}, fmt.Errorf("CreateCommand not implemented")
		}
	}
	if fn.ProcessData == nil {
		fn.ProcessData = func(_ AgentData, _ []byte) error {
			return fmt.Errorf("ProcessData not implemented")
		}
	}
	if fn.PackTasks == nil {
		fn.PackTasks = func(_ AgentData, _ []TaskData) ([]byte, error) {
			return nil, fmt.Errorf("PackTasks not implemented")
		}
	}
	if fn.Encrypt == nil {
		fn.Encrypt = func(data []byte, _ []byte) ([]byte, error) {
			return nil, fmt.Errorf("Encrypt not implemented")
		}
	}
	if fn.Decrypt == nil {
		fn.Decrypt = func(data []byte, _ []byte) ([]byte, error) {
			return nil, fmt.Errorf("Decrypt not implemented")
		}
	}
	if fn.PivotPackData == nil {
		fn.PivotPackData = func(_ string, _ []byte) (TaskData, error) {
			return TaskData{}, fmt.Errorf("PivotPackData not implemented")
		}
	}

	a := &Agent{
		data:         data,
		Fn:           fn,
		HostedQueue:  axsafe.NewPriorityQueue(0x1000),
		RunningTasks: axsafe.NewMap[int64, TaskData](),
		RunningJobs:  axsafe.NewMap[int64, *axsafe.Slice](),
		PivotChilds:  axsafe.NewSlice(),
	}
	a.active.Store(true)
	return a
}

func (s *Agent) IsRemoved() bool  { return s.removed.Load() }
func (s *Agent) MarkRemoved()     { s.removed.Store(true) }
func (s *Agent) IsActive() bool   { return s.active.Load() }
func (s *Agent) SetActive(v bool) { s.active.Store(v) }

func (s *Agent) GetPivotParent() *PivotData {
	s.pivotMu.RLock()
	defer s.pivotMu.RUnlock()
	return s.PivotParent
}

func (s *Agent) SetPivotParent(p *PivotData) {
	s.pivotMu.Lock()
	defer s.pivotMu.Unlock()
	s.PivotParent = p
}

func (s *Agent) GetData() AgentData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}

func (s *Agent) UpdateData(fn func(*AgentData)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(&s.data)
}

func (s *Agent) SetCommandGroupEnabled(groupId string, enabled bool) {
	if groupId == "" {
		return
	}
	s.cmdGroupMu.Lock()
	defer s.cmdGroupMu.Unlock()
	if s.cmdGroupEnabled == nil {
		s.cmdGroupEnabled = make(map[string]bool)
	}
	s.cmdGroupEnabled[groupId] = enabled
}

func (s *Agent) IsCommandGroupEnabled(groupId string, defaultEnabled bool) bool {
	s.cmdGroupMu.RLock()
	defer s.cmdGroupMu.RUnlock()
	if s.cmdGroupEnabled == nil {
		return defaultEnabled
	}
	if v, ok := s.cmdGroupEnabled[groupId]; ok {
		return v
	}
	return defaultEnabled
}

func (s *Agent) GetCommandGroupOverrides() map[string]bool {
	s.cmdGroupMu.RLock()
	defer s.cmdGroupMu.RUnlock()
	if len(s.cmdGroupEnabled) == 0 {
		return map[string]bool{}
	}
	out := make(map[string]bool, len(s.cmdGroupEnabled))
	for k, v := range s.cmdGroupEnabled {
		out[k] = v
	}
	return out
}

func (s *Agent) ApplyCommandGroupOverrides(overrides map[string]bool) {
	s.cmdGroupMu.Lock()
	defer s.cmdGroupMu.Unlock()
	if len(overrides) == 0 {
		s.cmdGroupEnabled = nil
		return
	}
	s.cmdGroupEnabled = make(map[string]bool, len(overrides))
	for k, v := range overrides {
		s.cmdGroupEnabled[k] = v
	}
}

func (s *Agent) CommandGroupOverridesJSON() string {
	m := s.GetCommandGroupOverrides()
	if len(m) == 0 {
		return ""
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s *Agent) LoadCommandGroupOverridesJSON(raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" || raw == "null" {
		s.ApplyCommandGroupOverrides(nil)
		return
	}
	var m map[string]bool
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return
	}
	s.ApplyCommandGroupOverrides(m)
}

func (s *Agent) ProcessData(packed []byte) error {
	if s.removed.Load() {
		return ErrAgentRemoved
	}
	data := s.GetData()
	decrypted, err := s.Fn.Decrypt(packed, data.SessionKey)
	if err != nil {
		return err
	}
	return s.Fn.ProcessData(data, decrypted)
}

func (s *Agent) PackTasks(tasks []TaskData) ([]byte, error) {
	if s.removed.Load() {
		return nil, ErrAgentRemoved
	}
	data := s.GetData()
	packed, err := s.Fn.PackTasks(data, tasks)
	if err != nil {
		return nil, err
	}
	return s.Fn.Encrypt(packed, data.SessionKey)
}

func (s *Agent) EncryptData(data []byte) ([]byte, error) {
	if s.removed.Load() {
		return nil, ErrAgentRemoved
	}
	return s.Fn.Encrypt(data, s.GetData().SessionKey)
}

func (s *Agent) DecryptData(data []byte) ([]byte, error) {
	if s.removed.Load() {
		return nil, ErrAgentRemoved
	}
	return s.Fn.Decrypt(data, s.GetData().SessionKey)
}

// --- HookJob ---

type HookJob struct {
	Sent      bool
	Processed bool
	Job       TaskData
	Mu        sync.RWMutex
}
