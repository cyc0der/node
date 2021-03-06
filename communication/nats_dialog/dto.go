package nats_dialog

import (
	"github.com/mysterium/node/communication"
	dto_discovery "github.com/mysterium/node/service_discovery/dto"
)

// Consume is trying to establish new dialog with Provider
const ENDPOINT_DIALOG_CREATE = communication.RequestType("dialog-create")

var (
	responseOK              = dialogCreateResponse{200, "OK"}
	responseInvalidIdentity = dialogCreateResponse{400, "Invalid identity"}
)

type dialogCreateRequest struct {
	IdentityId dto_discovery.Identity `json:"identity_id"`
}

type dialogCreateResponse struct {
	Reason        uint   `json:"reason"`
	ReasonMessage string `json:"reasonMessage"`
}
