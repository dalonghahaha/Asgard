package controllers

import "Asgard/models"

func formatArchive(info *models.Archive) map[string]interface{} {
	data := map[string]interface{}{
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
