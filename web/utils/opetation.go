package utils

import (
	"Asgard/models"
	"Asgard/providers"
)

func OpetationLog(userID, _type, relatedID, action int64) {
	opetationLog := new(models.OperationLog)
	opetationLog.UserID = userID
	opetationLog.Type = _type
	opetationLog.RelatedID = relatedID
	opetationLog.Action = action
	providers.OperationLogService.CreateOperationLog(opetationLog)
}
