package adaptix

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Adaptix-Framework/axsafe"
)

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

func (s *Agent) SetData(d AgentData) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.data = d
}

func (s *Agent) UpdateData(fn func(*AgentData)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(&s.data)
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

// HookJob

type HookJob struct {
	Sent      bool
	Processed bool
	Job       TaskData
	Mu        sync.RWMutex
}
