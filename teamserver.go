package adaptix

import (
	"context"
	"io"
	"net/http"
)

type WebSocketConn interface {
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type Teamserver interface {
	TsAgentGenID() int64
	TsAgentList() (string, error)
	TsAgentGetById(agentId int64) (AgentData, bool)
	TsAgentIsExists(agentId int64) bool
	TsAgentIdByUID(uid []byte) (int64, bool)
	TsAgentCreate(agentCrc string, agentUid []byte, beat []byte, listenerName string, ExternalIP string, Async bool) (AgentData, error)
	TsAgentUpdateData(newAgentData AgentData) error
	TsAgentTerminate(agentId int64, terminateTaskId int64) error
	TsAgentRemove(agentId int64) error
	TsAgentTickUpdate(ctx context.Context)

	TsAgentCommand(agentName string, agentId int64, clientName string, hookId string, handlerId string, cmdline string, ui bool, args map[string]any) error
	TsAgentConsoleOutput(agentId int64, client string, messageType int, message string, clearText string, store bool)
	TsAgentConsoleOutputClient(agentId int64, client string, messageType int, message string, clearText string)
	TsAgentConsoleErrorCommand(agentId int64, client string, cmdline string, message string, HookId string, HandlerId string)
	TsAgentConsoleLocalCommand(agentId int64, client string, cmdline string, message string, text string)

	TsAgentProcessData(agentId int64, bodyData []byte) error
	TsAgentBuildEmptyTasks(agentId int64) ([]byte, error)
	TsAgentGetHostedAll(agentId int64, maxDataSize int) ([]byte, error)
	TsAgentGetHostedTasks(agentId int64, maxDataSize int) ([]byte, error)
	TsAgentGetHostedTasksCount(agentId int64, count int, maxDataSize int) ([]byte, error)
	TsAgentUpdateDataPartial(agentId int64, updateData interface{}) error
	TsAgentSetTick(agentId int64, listenerName string) error
	TsAgentEncryptData(agentId int64, data []byte) ([]byte, error)
	TsAgentDecryptData(agentId int64, data []byte) ([]byte, error)
	TsAgentCommandGroupSet(agentId int64, groupId string, enabled bool) error
	TsAgentCommandGroupList(agentId int64) ([]map[string]interface{}, error)

	TsAgentBuildSyncOnce(agentName string, config string, listenersName []string, creator string, saveToStore bool, description string) ([]byte, string, error)
	TsAgentBuildCreateChannel(buildData string, wsconn WebSocketConn, creator string) error
	TsAgentBuildExecute(builderId string, workingDir string, env []string, program string, args ...string) error
	TsAgentBuildLog(builderId string, status int, message string) error
	TsAgentBuildSendFile(builderId string, filename string, content []byte) error
	TsAgentBuildClose(builderId string)

	TsPayloadRegister(agentType, filename string, content []byte, listeners []string, configJson, creator, buildId, watermark, description string) (PayloadData, error)
	TsPayloadList(showHidden bool) ([]byte, error)
	TsPayloadGetPage(offset, limit int, showHidden bool, filterExpr, sortCol, sortOrder string) ([]byte, error)
	TsPayloadGet(id int64) (PayloadData, error)
	TsPayloadDownload(id int64) (string, []byte, error)
	TsPayloadHide(ids []int64, hidden bool) error
	TsPayloadUpdateMeta(id int64, name, notes, artifact, arch string, hidden bool) (PayloadData, error)
	TsPayloadSetColor(ids []int64, background, foreground string, reset bool) error
	TsPayloadSetTag(ids []int64, tag string) error
	TsPayloadRemove(ids []int64, hard bool) error
	TsPayloadSync() ([]byte, error)
	TsPayloadImport(name, agentType, artifact, arch, creator string, listeners []string, content []byte, configJson string) (PayloadData, error)

	TsTaskGenID() int64
	TsTaskCreate(agentId int64, cmdline string, client string, taskData TaskData)
	TsTaskUpdate(agentId int64, data TaskData)
	TsTaskGetAvailableAll(agentId int64, availableSize int) ([]TaskData, error)
	TsTaskGetAvailableTasks(agentId int64, maxCount int, availableSize int) ([]TaskData, int, error)
	TsTaskRunningExists(agentId int64, taskId int64) bool
	TsTaskPostHook(hookData TaskData, jobIndex int) error
	TsTaskSave(hookData TaskData) error
	TsTasksPivotExists(agentId int64, first bool) bool

	TsTaskCancel(agentId int64, taskId int64) error
	TsTaskDelete(agentId int64, taskId int64) error
	TsTasksGetPage(agentId int64, offset, limit int, filterExpr, sortCol, sortOrder string, completedFilter *bool) ([]byte, error)
	TsConsoleGetPage(agentId int64, afterId int64, limit int, username string) ([]byte, error)
	TsConsoleSearch(agentId int64, query string, limit, offset int, username string) ([]byte, error)
	TsConsoleGetAround(agentId int64, centerId int64, limit int, username string) ([]byte, error)

	TsTunnelList() (string, error)
	TsTunnelClientStart(AgentId int64, Listen bool, Type int, Info string, Lhost string, Lport int, Client string, Thost string, Tport int, AuthUser string, AuthPass string) (int64, error)
	TsTunnelClientNewChannel(TunnelData string, wsconn WebSocketConn) error
	TsTunnelStart(TunnelId int64) (int64, error)
	TsTunnelDeactivate(TunnelId int64, clientName string) error
	TsTunnelClientStop(TunnelId int64, Client string) error
	TsTunnelStop(TunnelId int64) error
	TsTunnelClientCanControl(TunnelId int64, clientName string) error
	TsTunnelClientSetInfo(TunnelId int64, Info string, clientName string) error
	TsTunnelCreateSocks4(AgentId int64, Info string, Lhost string, Lport int) (int64, error)
	TsTunnelCreateSocks5(AgentId int64, Info string, Lhost string, Lport int, UseAuth bool, Username string, Password string) (int64, error)
	TsTunnelCreateLportfwd(AgentId int64, Info string, Lhost string, Lport int, Thost string, Tport int) (int64, error)
	TsTunnelCreateRportfwd(AgentId int64, Info string, Lport int, Thost string, Tport int) (int64, error)
	TsTunnelUpdateRportfwd(tunnelId int64, result bool) (int64, string, error)
	TsTunnelStopSocks(AgentId int64, Port int)
	TsTunnelStopLportfwd(AgentId int64, Port int)
	TsTunnelStopRportfwd(AgentId int64, Port int)
	TsTunnelGetPipe(AgentId int64, channelId int64) (*io.PipeReader, *io.PipeWriter, error)
	TsTunnelConnectionClose(channelId int64, writeOnly bool)
	TsTunnelConnectionHalt(channelId int64, errorCode byte)
	TsTunnelConnectionResume(AgentId int64, channelId int64, ioDirect bool)
	TsTunnelConnectionData(channelId int64, data []byte)
	TsTunnelConnectionAccept(tunnelId int64, channelId int64)
	TsTunnelPause(channelId int64)
	TsTunnelResume(channelId int64)
	TsTunnelChannelExists(channelId int64) bool

	TsTerminalGetPipe(AgentId int64, terminalId int64) (*io.PipeReader, *io.PipeWriter, error)
	TsTerminalConnResume(AgentId int64, terminalId int64, ioDirect bool)
	TsTerminalConnExists(terminalId int64) bool
	TsTerminalConnData(terminalId int64, data []byte)
	TsTerminalConnClose(terminalId int64, status string) error
	TsAgentTerminalCreateChannel(terminalData string, wsconn WebSocketConn) error
	TsAgentTerminalCloseChannel(terminalId int64, status string) error

	TsFileGenID() int64
	TsDownloadGet(fileId int64) (TransferData, error)
	TsDownloadAdd(AgentId int64, fileId int64, fileName string, fileSize int64) error
	TsDownloadUpdate(fileId int64, state int, data []byte) error
	TsDownloadClose(fileId int64, reason int) error
	TsDownloadSave(AgentId int64, fileId int64, filename string, content []byte) error
	TsDownloadList() (string, error)
	TsDownloadsGetPage(agentId int64, offset, limit int, filterExpr, sortCol, sortOrder string) ([]byte, error)
	TsDownloadSync(fileId int64) (string, []byte, error)
	TsDownloadDelete(fileId []int64) error
	TsDownloadSetTag(fileIds []int64, tag string) error
	TsDownloadGetFilepath(fileId int64) (string, error)

	TsUploadGet(fileId int64) (TransferData, error)
	TsUploadAdd(agentId int64, fileId int64, localPath string, remotePath string) error
	TsUploadAddContent(agentId int64, fileId int64, remotePath string, content []byte, canceled bool, kind int, artname string, arttype string) error
	TsUploadGetChunk(fileId int64, chunkSize int, needApprove bool) ([]byte, error)
	TsUploadApprove(fileId int64, approvedBytes int) error
	TsUploadClose(fileId int64, reason int) error
	TsUploadsGetPage(agentId int64, offset, limit int, filterExpr, sortCol, sortOrder string) ([]byte, error)
	TsUploadDelete(fileIds []int64) error
	TsUploadGetFilepath(fileId int64) (string, error)
	TsUploadGetFileContent(fileId int64) ([]byte, error)

	TsScreenGenID() int64
	TsScreenshotAdd(AgentId int64, Note string, Content []byte) error
	TsScreenshotList() (string, error)
	TsScreenshotsGetPage(offset, limit int, filterExpr, sortCol, sortOrder string) ([]byte, error)
	TsScreenshotGetImage(screenId int64) ([]byte, error)
	TsScreenshotDelete(screenId int64) error
	TsScreenshotNote(screenId int64, note string) error

	TsCredGenID() int64
	TsCredentilsList() (string, error)
	TsCredentialsGetPage(agentId int64, offset, limit int, filterExpr, sortCol, sortOrder string) ([]byte, error)
	TsCredentilsAdd(creds []map[string]interface{}) error
	TsCredentilsEdit(credId int64, username string, password string, realm string, credType string, tag string, storage string, host string) error
	TsCredentilsDelete(credsId []int64) error
	TsCredentialsSetTag(credsId []int64, tag string) error

	TsTargetGenID() int64
	TsTargetsList() (string, error)
	TsTargetsGetPage(offset, limit int, filterExpr, sortCol, sortOrder string) ([]byte, error)
	TsTargetsAdd(targets []map[string]interface{}) error
	TsTargetsCreateAlive(agentData AgentData) (int64, error)
	TsTargetsEdit(targetId int64, computer string, domain string, address string, os int, osDesk string, tag string, info string, alive bool) error
	TsTargetDelete(targetsId []int64) error
	TsTargetSetTag(targetsId []int64, tag string) error
	TsTargetRemoveSessions(agentsId []int64) error

	TsClientGuiDisksWindows(taskData TaskData, drives []ListingDrivesDataWin)
	TsClientGuiFilesStatus(taskData TaskData)
	TsClientGuiFilesWindows(taskData TaskData, path string, files []ListingFileDataWin)
	TsClientGuiFilesUnix(taskData TaskData, path string, files []ListingFileDataUnix)
	TsClientGuiProcessWindows(taskData TaskData, process []ListingProcessDataWin)
	TsClientGuiProcessUnix(taskData TaskData, process []ListingProcessDataUnix)

	TsSetAgentDeliveryFunc(AgentId int64, fn DeliveryFunc)
	TsRemoveAgentDeliveryFunc(AgentId int64)
	TsGetAgentDeliveryFunc(AgentId int64) DeliveryFunc

	TsPivotCreate(pivotId string, pAgentId int64, chAgentId int64, pivotName string, isRestore bool) error
	TsGetPivotInfoByName(pivotName string) (string, int64, int64)
	TsGetPivotInfoById(pivotId string) (string, int64, int64)
	TsPivotDelete(pivotId string) error

	TsListenerList() (string, error)
	TsListenerStart(listenerName string, configType string, config string, createTime int64, watermark string, customData []byte, tags string) error
	TsListenerEdit(listenerName string, configType string, config string, tags string) error
	TsListenerStop(listenerName string, configType string) error
	TsListenerPause(listenerName string, configType string) error
	TsListenerResume(listenerName string, configType string) error
	TsListenerGetProfile(listenerName string) (string, []byte, error)
	TsListenerInteralHandler(watermark string, data []byte) (int64, error)
	TsListenerConnector(listenerName string, data []byte) (int64, error)
	TsListenerSetTags(listenerName string, tags string) error

	TsServiceLoad(configPath string) error
	TsServiceUnload(serviceName string) error
	TsPluginServiceCall(serviceName string, operator string, function string, args string)
	TsPluginServiceCallWait(serviceName string, operator string, function string, args string, timeoutMs int) (resultJSON string, err error)
	TsServiceList() (string, error)

	TsAxScriptLoadUser(name string, script string) error
	TsAxScriptUnloadUser(name string) error
	TsAxScriptList() (string, error)
	TsAxScriptCommands() (string, error)
	TsAxScriptResolveHooks(agentName string, agentId int64, listenerRegName string, os int, cmdline string, args map[string]interface{}, client string) (string, string, bool, error)
	TsAxScriptIsServerHook(id string) bool
	TsAxScriptParseAndExecute(agentId int64, username string, cmdline string) error
	AxGetAgentContext(agentId int64) (agentName string, listenerRegName string, osType int, err error)

	TsEventHandlersList() (string, error) // JSON []EventHandlerInfo (full)
	TsEventHandlersGetPage(offset, limit int, q, event, source, group, enabled string) ([]byte, error)
	TsEventHandlerRegister(requestJSON string, operator string) (string, error)
	TsEventHandlerGet(id string) (string, error)
	TsEventHandlerEnable(id string) error
	TsEventHandlerDisable(id string) error
	TsEventHandlerRemove(id string) error
	TsEventMute(eventType string) error
	TsEventUnmute(eventType string) error
	TsEventMutesList() (string, error) // JSON []string
	TsEventHookSetEnabled(hookID string, enabled bool) error
	TsEventEmit(eventType string, text string) error
	TsEventEmitFrom(eventType string, source string, text string) error
	TsEventTypesList() (string, error)

	TsGroupList(scope string) []map[string]interface{}
	TsGroupCreate(parentId int64, name string, scope string) error
	TsGroupRename(groupId int64, name string) error
	TsGroupDelete(groupId int64) error
	TsGroupMembers(groupId int64, add []int64, remove []int64) error
	TsGroupMoveMember(agentId, fromGroupId, toGroupId int64) error
	TsGroupReparent(groupId int64, newParentId int64) error

	TsLogAdd(status LogStatus, level int, source, category string, format string, args ...any)
	TsLogWriter(status LogStatus, source, category string) io.Writer
	TsLogsGetPage(offset, limit int) ([]byte, error)
	TsLogsGetPageFiltered(offset, limit int, sourceFilter, categoryFilter, contains string) ([]byte, error)
	TsLogsGetPageBeforeId(beforeId int64, limit int) ([]byte, error)
	TsLogsGetPageBeforeIdFiltered(beforeId int64, limit int, sourceFilter, categoryFilter, contains string) ([]byte, error)

	TsChatSendMessage(username string, message string, replyToId int64, replyToName string)
	TsChatEditMessage(username string, id int64, newMessage string) error
	TsChatDeleteMessage(username string, id int64) error
	TsChatToggleReaction(username string, id int64, emoji string) error
	TsChatGetTodo() (string, string, int64)
	TsChatUpdateTodo(username string, content string) error
	TsChatHistory(limit int, beforeId int64) []byte
	TsChatSearch(query string, limit int, beforeId int64) []byte
	TsChatClear() error
	TsChatCount() int

	TsClientExists(username string) bool
	TsClientConnected(username string) bool
	TsClientDisconnect(username string)
	TsClientConnect(username string, socket WebSocketConn, clientType uint8, consoleTeamMode bool, subscriptions []string)
	TsClientSync(username string)
	TsClientSubscribe(username string, categories []string, consoleTeamMode *bool)

	CreateOTP(otpType string, data interface{}) (string, error)
	ValidateOTP(token string) (string, interface{}, bool)

	TsConvertCpToUTF8(input string, codePage int) string
	TsConvertUTF8toCp(input string, codePage int) string
	TsWin32Error(errorCode uint) string

	TsFrameHasPending(sessionId int64) bool
	TsFramePut(sessionId int64, index uint32, data []byte, totalSize uint32, chunkCount uint16) (bool, uint32, uint32, uint32, []byte)
	TsFramePutStream(sessionId int64, seqNum uint32, data []byte, isLast bool) (bool, []byte)
	TsFrameGetChunk(sessionId int64, reqOffset uint32, maxChunkSize int, encode func([]byte) []byte) (uint32, uint32, []byte, uint32, bool)
	TsFrameGetChunkSticky(sessionId int64, reqOffset uint32, maxChunkSize int, encode func([]byte) []byte) (uint32, uint32, []byte, uint32, bool)
	TsFrameAckDelivery(sessionId int64, ackOffset uint32, ackNonce uint32)
	TsFrameResetUpstream(sessionId int64)
	TsFrameResetDownstream(sessionId int64)

	TsExtenderDataSave(extenderName string, key string, data []byte) error
	TsExtenderDataLoad(extenderName string, key string) ([]byte, error)
	TsExtenderDataDelete(extenderName string, key string) error
	TsExtenderDataKeys(extenderName string) ([]string, error)
	TsExtenderDataDeleteAll(extenderName string) error

	TsEventHookRegister(eventType string, name string, phase int, priority int, handler func(event any) error) string
	TsEventHookOnPre(eventType string, name string, handler func(event any) error) string
	TsEventHookOnPost(eventType string, name string, handler func(event any) error) string
	TsEventHookUnregister(hookID string) bool
	TsEventHookUnregisterByName(name string) int

	TsPluginServiceSendDataAll(service string, data string)
	TsPluginServiceSendDataClient(operator string, service string, data string)
	TsPluginAgentCall(agentId int64, operator string, function string, args string)
	TsPluginAgentSendDataAll(agentId int64, data string)
	TsPluginAgentSendDataClient(operator string, agentId int64, data string)
	TsPluginListenerCall(listenerName string, operator string, function string, args string)
	TsPluginListenerSendDataAll(listenerName string, data string)
	TsPluginListenerSendDataClient(operator string, listenerName string, data string)

	TsEndpointRegister(method string, path string, handler func(username string, body []byte) (int, []byte)) error
	TsEndpointRegisterRaw(method string, path string, handler func(w http.ResponseWriter, r *http.Request, username string)) error
	TsEndpointUnregister(method string, path string) error
	TsEndpointExists(method string, path string) bool

	TsEndpointRegisterPublic(method string, path string, handler func(body []byte) (int, []byte)) error
	TsEndpointRegisterPublicRaw(method string, path string, handler func(w http.ResponseWriter, r *http.Request)) error
	TsEndpointUnregisterPublic(method string, path string) error
	TsEndpointExistsPublic(method string, path string) bool
}
