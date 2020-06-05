package utils

import (
	"Asgard/models"
	"Asgard/providers"
)

func OpetationLog(userID, _type, relatedID, action int64) {
	opetation := new(models.Operation)
	opetation.UserID = userID
	opetation.Type = _type
	opetation.RelatedID = relatedID
	opetation.Action = action
	providers.OperationService.CreateOperation(opetation)
}
