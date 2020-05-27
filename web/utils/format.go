package utils

import (
	"fmt"
	"regexp"
	"time"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"

	"github.com/gin-gonic/gin"
)

var (
	TimeLocation = "Asia/Shanghai"
	TimeLayout   = "2006-01-02 15:04"
)

func FormatFloat(info float64) string {
	return fmt.Sprintf("%.4f", info)
}

func FormatTime(info time.Time) string {
	return info.Format("2006-01-02 15:04:05")
}

func ParseTime(str string) (time.Time, error) {
	locationName := TimeLocation
	if l, err := time.LoadLocation(locationName); err != nil {
		return time.Time{}, err
	} else {
		lt, _ := time.ParseInLocation(TimeLayout, str, l)
		return lt, nil
	}
}

func EmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func MobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func AppFormat(info *models.App) gin.H {
	data := gin.H{
		"ID":          info.ID,
		"Name":        info.Name,
		"GroupID":     info.GroupID,
		"AgentID":     info.AgentID,
		"Dir":         info.Dir,
		"Program":     info.Program,
		"Args":        info.Args,
		"StdOut":      info.StdOut,
		"StdErr":      info.StdErr,
		"AutoRestart": info.AutoRestart,
		"IsMonitor":   info.IsMonitor,
		"Status":      info.Status,
	}
	group := providers.GroupService.GetGroupByID(info.GroupID)
	data["GroupName"] = GroupNameFormat(group)
	agent := providers.AgentService.GetAgentByID(info.AgentID)
	data["AgentName"] = AgentNameFormat(agent)
	return data
}

func JobFormat(info *models.Job) gin.H {
	data := gin.H{
		"ID":        info.ID,
		"Name":      info.Name,
		"GroupID":   info.GroupID,
		"AgentID":   info.AgentID,
		"Dir":       info.Dir,
		"Program":   info.Program,
		"Args":      info.Args,
		"StdOut":    info.StdOut,
		"StdErr":    info.StdErr,
		"Spec":      info.Spec,
		"Timeout":   info.Timeout,
		"IsMonitor": info.IsMonitor,
		"Status":    info.Status,
	}
	group := providers.GroupService.GetGroupByID(info.GroupID)
	data["GroupName"] = GroupNameFormat(group)
	agent := providers.AgentService.GetAgentByID(info.AgentID)
	data["AgentName"] = AgentNameFormat(agent)
	return data
}

func TimingFormat(info *models.Timing) gin.H {
	data := gin.H{
		"ID":        info.ID,
		"Name":      info.Name,
		"GroupID":   info.GroupID,
		"AgentID":   info.AgentID,
		"Dir":       info.Dir,
		"Program":   info.Program,
		"Args":      info.Args,
		"StdOut":    info.StdOut,
		"StdErr":    info.StdErr,
		"Time":      info.Time.Format(TimeLayout),
		"Timeout":   info.Timeout,
		"IsMonitor": info.IsMonitor,
		"Status":    info.Status,
	}
	group := providers.GroupService.GetGroupByID(info.GroupID)
	data["GroupName"] = GroupNameFormat(group)
	agent := providers.AgentService.GetAgentByID(info.AgentID)
	data["AgentName"] = AgentNameFormat(agent)
	return data
}

func GroupNameFormat(group *models.Group) string {
	if group != nil {
		return group.Name
	} else {
		return ""
	}
}

func AgentNameFormat(agent *models.Agent) string {
	if agent != nil {
		if agent.Alias != "" {
			return agent.Alias
		} else {
			return fmt.Sprintf("%s:%s", agent.IP, agent.Port)
		}
	} else {
		return ""
	}
}

func MonitorFormat(moniters []models.Monitor) (cpus []string, memorys []string, times []string) {
	for _, moniter := range moniters {
		cpus = append(cpus, FormatFloat(moniter.CPU))
		memorys = append(memorys, FormatFloat(moniter.Memory))
		times = append(times, FormatTime(moniter.CreatedAt))
	}
	return
}

func ArchiveFormat(info *models.Archive) gin.H {
	data := gin.H{
		"ID":        info.ID,
		"UUID":      info.UUID,
		"PID":       info.PID,
		"BeginTime": FormatTime(info.BeginTime),
		"EndTime":   FormatTime(info.EndTime),
		"Status":    info.Status,
		"Signal":    info.Signal,
	}
	return data
}

func GetErrorMessage(code int) (message string) {
	var ok bool
	if constants.WEB_LANG == "cn" {
		message, ok = constants.ERROR_CN[code]
	} else {
		message, ok = constants.ERROR_EN[code]
	}
	if !ok {
		if constants.WEB_LANG == "cn" {
			return constants.ERROR_CN_NOFUND
		} else {
			return constants.ERROR_EN_NOFUND
		}
	}
	return
}
