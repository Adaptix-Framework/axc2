package adaptix

import "fmt"

type TaskStat struct {
	Count int `json:"count"`
	Bytes int `json:"bytes"`
}

func (s TaskStat) Empty() bool { return s.Count == 0 }

func (s TaskStat) SizeString() string { return FormatByteSize(s.Bytes) }

func (s TaskStat) add(nbytes int) TaskStat {
	s.Count++
	s.Bytes += nbytes
	return s
}

func (s TaskStat) plus(o TaskStat) TaskStat {
	return TaskStat{Count: s.Count + o.Count, Bytes: s.Bytes + o.Bytes}
}

type StatTasks struct {
	Local   TaskStat `json:"local"`
	Task    TaskStat `json:"task"`
	Browser TaskStat `json:"browser"`
	Job     TaskStat `json:"job"`
	Tunnel  TaskStat `json:"tunnel"`
	Proxy   TaskStat `json:"proxy"`
	Other   TaskStat `json:"other"`
	Packed  int      `json:"packed"`
}

func (s StatTasks) Of(taskType int) TaskStat {
	switch taskType {
	case TASK_TYPE_LOCAL:
		return s.Local
	case TASK_TYPE_TASK:
		return s.Task
	case TASK_TYPE_BROWSER:
		return s.Browser
	case TASK_TYPE_JOB:
		return s.Job
	case TASK_TYPE_TUNNEL:
		return s.Tunnel
	case TASK_TYPE_PROXY_DATA:
		return s.Proxy
	default:
		return s.Other
	}
}

func (s StatTasks) Commands() TaskStat { return s.Task.plus(s.Job) }

func (s StatTasks) Total() TaskStat {
	return s.Local.plus(s.Task).plus(s.Browser).plus(s.Job).plus(s.Tunnel).plus(s.Proxy).plus(s.Other)
}

func (s StatTasks) Empty() bool { return s.Total().Empty() }

func (s StatTasks) Select(types ...int) TaskStat {
	if len(types) == 0 {
		return s.Commands()
	}
	var out TaskStat
	for _, t := range types {
		out = out.plus(s.Of(t))
	}
	return out
}

func MakeStatTasks(tasks []TaskData) StatTasks {
	var st StatTasks
	for i := range tasks {
		n := len(tasks[i].Data)
		switch tasks[i].Type {
		case TASK_TYPE_LOCAL:
			st.Local = st.Local.add(n)
		case TASK_TYPE_TASK:
			st.Task = st.Task.add(n)
		case TASK_TYPE_BROWSER:
			st.Browser = st.Browser.add(n)
		case TASK_TYPE_JOB:
			st.Job = st.Job.add(n)
		case TASK_TYPE_TUNNEL:
			st.Tunnel = st.Tunnel.add(n)
		case TASK_TYPE_PROXY_DATA:
			st.Proxy = st.Proxy.add(n)
		default:
			st.Other = st.Other.add(n)
		}
	}
	return st
}

func FormatByteSize(n int) string {
	const (
		kb = 1024.0
		mb = kb * 1024
		gb = mb * 1024
	)
	size := float64(n)
	switch {
	case size >= gb:
		return fmt.Sprintf("%.2f Gb", size/gb)
	case size >= mb:
		return fmt.Sprintf("%.2f Mb", size/mb)
	case size >= 1000:
		return fmt.Sprintf("%.2f Kb", size/kb)
	default:
		return fmt.Sprintf("%d bytes", n)
	}
}
