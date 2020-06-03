package providers

import (
	"Asgard/services"
)

var (
	UserService         *services.UserService
	AgentService        *services.AgentService
	GroupService        *services.GroupService
	AppService          *services.AppService
	JobService          *services.JobService
	TimingService       *services.TimingService
	MoniterService      *services.MonitorService
	ArchiveService      *services.ArchiveService
	OperationLogService *services.OperationLogService
	NoticeService       *services.NoticeService
)

func init() {
	UserService = services.NewUserService()
	AgentService = services.NewAgentService()
	GroupService = services.NewGroupService()
	AppService = services.NewAppService()
	JobService = services.NewJobService()
	TimingService = services.NewTimingService()
	MoniterService = services.NewMonitorService()
	ArchiveService = services.NewArchiveService()
	OperationLogService = services.NewOperationLogService()
	NoticeService = services.NewNoticeService()
}
