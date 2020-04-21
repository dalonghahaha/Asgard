package controllers

import (
	"Asgard/models"
	"Asgard/web/utils"
)

func formatArchive(info *models.Archive) map[string]interface{} {
	data := map[string]interface{}{
		"ID":        info.ID,
		"UUID":      info.UUID,
		"PID":       info.PID,
		"BeginTime": utils.FormatTime(info.BeginTime),
		"EndTime":   utils.FormatTime(info.EndTime),
		"Status":    info.Status,
		"Signal":    info.Signal,
	}
	return data
}
