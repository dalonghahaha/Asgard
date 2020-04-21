package utils

import (
	"Asgard/models"
	"fmt"
	"regexp"
	"time"
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

func MonitorFormat(moniters []models.Monitor) (cpus []string, memorys []string, times []string) {
	for _, moniter := range moniters {
		cpus = append(cpus, FormatFloat(moniter.CPU))
		memorys = append(memorys, FormatFloat(moniter.Memory))
		times = append(times, FormatTime(moniter.CreatedAt))
	}
	return
}
