package adaptix

import (
	"os"
)

type ListenerData struct {
	Name      string `json:"l_name"`
	RegName   string `json:"l_reg_name"`
	Protocol  string `json:"l_protocol"`
	Type      string `json:"l_type"`
	BindHost  string `json:"l_bind_host"`
	BindPort  string `json:"l_bind_port"`
	AgentAddr string `json:"l_agent_addr"`
	Status    string `json:"l_status"`
	Data      string `json:"l_data"`
	Watermark string `json:"l_watermark"`
}

type AgentData struct {
	Crc          string `json:"a_crc"`
	Id           string `json:"a_id"`
	Name         string `json:"a_name"`
	SessionKey   []byte `json:"a_session_key"`
	Listener     string `json:"a_listener"`
	Async        bool   `json:"a_async"`
	ExternalIP   string `json:"a_external_ip"`
	InternalIP   string `json:"a_internal_ip"`
	GmtOffset    int    `json:"a_gmt_offset"`
	Sleep        uint   `json:"a_sleep"`
	Jitter       uint   `json:"a_jitter"`
	Pid          string `json:"a_pid"`
	Tid          string `json:"a_tid"`
	Arch         string `json:"a_arch"`
	Elevated     bool   `json:"a_elevated"`
	Process      string `json:"a_process"`
	Os           int    `json:"a_os"`
	OsDesc       string `json:"a_os_desc"`
	Domain       string `json:"a_domain"`
	Computer     string `json:"a_computer"`
	Username     string `json:"a_username"`
	Impersonated string `json:"a_impersonated"`
	OemCP        int    `json:"a_oemcp"`
	ACP          int    `json:"a_acp"`
	CreateTime   int64  `json:"a_create_time"`
	LastTick     int    `json:"a_last_tick"`
	KillDate     int    `json:"a_killdate"`
	WorkingTime  int    `json:"a_workingtime"`
	Tags         string `json:"a_tags"`
	Mark         string `json:"a_mark"`
	Color        string `json:"a_color"`
	TargetId     string `json:"a_target"`
}

type TaskDataTunnel struct {
	ChannelId int
	Data	  TaskData
}

type TaskData struct {
	Type        int    `json:"t_type"`
	TaskId      string `json:"t_task_id"`
	AgentId     string `json:"t_agent_id"`
	Client      string `json:"t_client"`
	HookId      string `json:"t_hook_id"`
	User        string `json:"t_user"`
	Computer    string `json:"t_computer"`
	StartDate   int64  `json:"t_start_date"`
	FinishDate  int64  `json:"t_finish_date"`
	Data        []byte `json:"t_data"`
	CommandLine string `json:"t_command_line"`
	MessageType int    `json:"t_message_type"`
	Message     string `json:"t_message"`
	ClearText   string `json:"t_clear_text"`
	Completed   bool   `json:"t_completed"`
	Sync        bool   `json:"t_sync"`
}

type ConsoleMessageData struct {
	Message string `json:"m_message"`
	Status  int    `json:"m_status"`
	Text    string `json:"m_text"`
}

type ListingFileDataWin struct {
	IsDir    bool   `json:"b_is_dir"`
	Size     int64  `json:"b_size"`
	Date     int64  `json:"b_date"`
	Filename string `json:"b_filename"`
}

type ListingFileDataUnix struct {
	IsDir    bool   `json:"b_is_dir"`
	Mode     string `json:"b_mode"`
	User     string `json:"b_user"`
	Group    string `json:"b_group"`
	Size     int64  `json:"b_size"`
	Date     string `json:"b_date"`
	Filename string `json:"b_filename"`
}

type ListingProcessDataWin struct {
	Pid         uint   `json:"b_pid"`
	Ppid        uint   `json:"b_ppid"`
	SessionId   uint   `json:"b_session_id"`
	Arch        string `json:"b_arch"`
	Context     string `json:"b_context"`
	ProcessName string `json:"b_process_name"`
}

type ListingProcessDataUnix struct {
	Pid         uint   `json:"b_pid"`
	Ppid        uint   `json:"b_ppid"`
	TTY         string `json:"b_tty"`
	Context     string `json:"b_context"`
	ProcessName string `json:"b_process_name"`
}

type ListingDrivesDataWin struct {
	Name string `json:"b_name"`
	Type string `json:"b_type"`
}

type ChatData struct {
	Username string `json:"c_username"`
	Message  string `json:"c_message"`
	Date     int64  `json:"c_date"`
}

type DownloadData struct {
	FileId     string `json:"d_file_id"`
	AgentId    string `json:"d_agent_id"`
	AgentName  string `json:"d_agent_name"`
	User       string `json:"d_user"`
	Computer   string `json:"d_computer"`
	RemotePath string `json:"d_remote_path"`
	LocalPath  string `json:"d_local_path"`
	TotalSize  int    `json:"d_total_size"`
	RecvSize   int    `json:"d_recv_size"`
	Date       int64  `json:"d_date"`
	State      int    `json:"d_state"`
	File       *os.File
}

type ScreenData struct {
	ScreenId  string `json:"s_screen_id"`
	User      string `json:"s_user"`
	Computer  string `json:"s_computer"`
	LocalPath string `json:"s_local_path"`
	Note      string `json:"s_note"`
	Date      int64  `json:"s_date"`
	Content   []byte `json:"s_content"`
}

type TunnelData struct {
	TunnelId  string `json:"p_tunnel_id"`
	AgentId   string `json:"p_agent_id"`
	Computer  string `json:"p_computer"`
	Username  string `json:"p_username"`
	Process   string `json:"p_process"`
	Type      string `json:"p_type"`
	Info      string `json:"p_info"`
	Interface string `json:"p_interface"`
	Port      string `json:"p_port"`
	Client    string `json:"p_client"`
	Fhost     string `json:"p_fhost"`
	Fport     string `json:"p_fport"`
	AuthUser  string `json:"p_auth_user"`
	AuthPass  string `json:"p_auth_pass"`
}

type PivotData struct {
	PivotId       string `json:"p_pivot_id"`
	PivotName     string `json:"p_pivot_name"`
	ParentAgentId string `json:"p_parent_agent_id"`
	ChildAgentId  string `json:"p_child_agent_id"`
}

type CredsData struct {
	CredId   string `json:"c_creds_id"`
	Username string `json:"c_username"`
	Password string `json:"c_password"`
	Realm    string `json:"c_realm"`
	Type     string `json:"c_type"`
	Tag      string `json:"c_tag"`
	Date     int64  `json:"c_date"`
	Storage  string `json:"c_storage"`
	AgentId  string `json:"c_agent_id"`
	Host     string `json:"c_host"`
}

type TargetData struct {
	TargetId string   `json:"t_target_id"`
	Computer string   `json:"t_computer"`
	Domain   string   `json:"t_domain"`
	Address  string   `json:"t_address"`
	Os       int      `json:"t_os"`
	OsDesk   string   `json:"t_os_desk"`
	Tag      string   `json:"t_tag"`
	Info     string   `json:"t_info"`
	Date     int64    `json:"t_date"`
	Alive    bool     `json:"t_alive"`
	Agents   []string `json:"t_agents"`
}
