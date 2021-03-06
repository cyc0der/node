package nats_dialog

import (
	"fmt"
	"github.com/mgutz/logxi/v1"
	"github.com/mysterium/node/communication"
	"github.com/mysterium/node/communication/nats"
	"github.com/mysterium/node/communication/nats_discovery"
	dto_discovery "github.com/mysterium/node/service_discovery/dto"
)

func NewDialogEstablisher(identity dto_discovery.Identity) *dialogEstablisher {
	return &dialogEstablisher{
		myIdentity: identity,
	}
}

const establisherLogPrefix = "[NATS.DialogEstablisher] "

type dialogEstablisher struct {
	myIdentity dto_discovery.Identity
}

func (establisher *dialogEstablisher) CreateDialog(contact dto_discovery.Contact) (communication.Dialog, error) {
	contactAddress, err := nats_discovery.NewAddressForContact(contact)
	if err != nil {
		return nil, err
	}

	log.Info(establisherLogPrefix, fmt.Sprintf("Connecting to: %#v", contactAddress))
	err = contactAddress.Connect()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to: %#v. %s", contact, err)
	}

	response, err := nats.NewSender(contactAddress).Request(&dialogCreateProducer{
		&dialogCreateRequest{
			IdentityId: establisher.myIdentity,
		},
	})
	if err != nil || response.(*dialogCreateResponse).Reason != 200 {
		return nil, fmt.Errorf("Dialog creation rejected: %s", err)
	}

	dialogAddress := nats_discovery.NewAddressNested(contactAddress, string(establisher.myIdentity))
	dialog := &dialog{nats.NewSender(dialogAddress), nats.NewReceiver(dialogAddress)}

	log.Info(establisherLogPrefix, fmt.Sprintf("Dialog established with: %#v", contact))
	return dialog, err
}
