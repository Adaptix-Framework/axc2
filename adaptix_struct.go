package adaptix

import (
	"os"
)

const (
	OS_UNKNOWN = 0
	OS_WINDOWS = 1
	OS_LINUX   = 2
	OS_MAC     = 3

	TASK_TYPE_LOCAL      = 0
	TASK_TYPE_TASK       = 1
	TASK_TYPE_BROWSER    = 2
	TASK_TYPE_JOB        = 3
	TASK_TYPE_TUNNEL     = 4
	TASK_TYPE_PROXY_DATA = 5

	MESSAGE_INFO    = 5
	MESSAGE_ERROR   = 6
	MESSAGE_SUCCESS = 7

	BUILD_LOG_NONE    = 0
	BUILD_LOG_INFO    = 1
	BUILD_LOG_ERROR   = 2
	BUILD_LOG_SUCCESS = 3

	DOWNLOAD_STATE_RUNNING  = 1
	DOWNLOAD_STATE_STOPPED  = 2
	DOWNLOAD_STATE_FINISHED = 3
	DOWNLOAD_STATE_CANCELED = 4

	TUNNEL_TYPE_SOCKS4      = 1
	TUNNEL_TYPE_SOCKS5      = 2
	TUNNEL_TYPE_SOCKS5_AUTH = 3
	TUNNEL_TYPE_LOCAL_PORT  = 4
	TUNNEL_TYPE_REVERSE     = 5

	ADDRESS_TYPE_IPV4   = 1
	ADDRESS_TYPE_DOMAIN = 3
	ADDRESS_TYPE_IPV6   = 4

	SOCKS5_SUCCESS                 byte = 0
	SOCKS5_SERVER_FAILURE          byte = 1
	SOCKS5_NOT_ALLOWED_RULESET     byte = 2
	SOCKS5_NETWORK_UNREACHABLE     byte = 3
	SOCKS5_HOST_UNREACHABLE        byte = 4
	SOCKS5_CONNECTION_REFUSED      byte = 5
	SOCKS5_TTL_EXPIRED             byte = 6
	SOCKS5_COMMAND_NOT_SUPPORTED   byte = 7
	SOCKS5_ADDR_TYPE_NOT_SUPPORTED byte = 8
)

type PluginService interface {
	Call(operator string, function string, args string)
}

type PluginListener interface {
	Create(name, config string, customData []byte) (ExtenderListener, ListenerData, []byte, error)
}

type ExtenderListener interface {
	Start() error
	Edit(config string) (ListenerData, []byte, error)
	Stop() error
	GetProfile() ([]byte, error)
	InternalHandler(data []byte) (string, error)
}

type PluginAgent interface {
	GenerateProfiles(profile BuildProfile) ([][]byte, error)
	BuildPayload(profile BuildProfile, agentProfiles [][]byte) ([]byte, string, error)

	GetExtender() ExtenderAgent
	CreateAgent(beat []byte) (AgentData, ExtenderAgent, error)
}

type ExtenderAgent interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(data []byte, key []byte) ([]byte, error)

	PackTasks(agentData AgentData, tasks []TaskData) ([]byte, error)
	PivotPackData(pivotId string, data []byte) (TaskData, error)

	CreateCommand(agentData AgentData, args map[string]any) (TaskData, ConsoleMessageData, error)
	ProcessData(agentData AgentData, decryptedData []byte) error

	TunnelCallbacks() TunnelCallbacks
	TerminalCallbacks() TerminalCallbacks
}

type TunnelCallbacks struct {
	ConnectTCP func(channelId, tunnelType, addressType int, address string, port int) TaskData
	ConnectUDP func(channelId, tunnelType, addressType int, address string, port int) TaskData
	WriteTCP   func(channelId int, data []byte) TaskData
	WriteUDP   func(channelId int, data []byte) TaskData
	Close      func(channelId int) TaskData
	Reverse    func(tunnelId, port int) TaskData
}

type TerminalCallbacks struct {
	Start func(terminalId int, program string, sizeH, sizeW, oemCP int) TaskData
	Write func(terminalId, oemCP int, data []byte) TaskData
	Close func(terminalId int) TaskData
}

type ListenerData struct {
	Name       string `json:"l_name"`
	RegName    string `json:"l_reg_name"`
	Protocol   string `json:"l_protocol"`
	Type       string `json:"l_type"`
	BindHost   string `json:"l_bind_host"`
	BindPort   string `json:"l_bind_port"`
	AgentAddr  string `json:"l_agent_addr"`
	CreateTime int64  `json:"a_create_time"`
	Status     string `json:"l_status"`
	Data       string `json:"l_data"`
	Watermark  string `json:"l_watermark"`
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
	CustomData   []byte `json:"a_custom_data"`
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
	HandlerId   string `json:"t_handler_id"`
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

type TransportProfile struct {
	Watermark string `json:"watermark"`
	Profile   []byte `json:"profile"`
}

type BuildProfile struct {
	BuilderId        string             `json:"build_id"`
	AgentConfig      string             `json:"agent_params"`
	ListenerProfiles []TransportProfile `json:"listener_profiles"`
}
