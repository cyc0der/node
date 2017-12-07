package identity

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/mysterium/node/service_discovery/dto"
	"github.com/stretchr/testify/assert"
)

func newManager(accountValue string) *IdentityManager {
	return &IdentityManager{
		KeystoreManager: &KeyStoreFake{
			AccountsMock: []accounts.Account{
				IdentityToAccount(accountValue),
			},
		},
	}
}

func newManagerWithError(errorMock error) *IdentityManager {
	return &IdentityManager{
		KeystoreManager: &KeyStoreFake{
			ErrorMock: errorMock,
		},
	}
}

func Test_CreateNewIdentity(t *testing.T) {
	manager := newManager("0x000000000000000000000000000000000000000A")
	identity, err := manager.CreateNewIdentity("")

	assert.NoError(t, err)
	assert.Equal(t, *identity, dto.Identity("0x000000000000000000000000000000000000bEEF"))
	assert.Len(t, manager.KeystoreManager.Accounts(), 2)
}

func Test_CreateNewIdentityError(t *testing.T) {
	im := newManagerWithError(errors.New("Identity create failed"))
	identity, err := im.CreateNewIdentity("")

	assert.EqualError(t, err, "Identity create failed")
	assert.Nil(t, identity)
}

func Test_GetIdentities(t *testing.T) {
	manager := newManager("0x000000000000000000000000000000000000000A")

	assert.Equal(
		t,
		[]dto.Identity{
			dto.Identity("0x000000000000000000000000000000000000000A"),
		},
		manager.GetIdentities(),
	)
}

func Test_GetIdentity(t *testing.T) {
	manager := newManager("0x000000000000000000000000000000000000000A")

	assert.Equal(
		t,
		dto.Identity("0x000000000000000000000000000000000000000A"),
		*manager.GetIdentity("0x000000000000000000000000000000000000000A"),
	)
	assert.Equal(
		t,
		dto.Identity("0x000000000000000000000000000000000000000A"),
		*manager.GetIdentity("0x000000000000000000000000000000000000000a"),
	)
	assert.Nil(
		t,
		manager.GetIdentity("0x000000000000000000000000000000000000000B"),
	)
}

func Test_HasIdentity(t *testing.T) {
	manager := newManager("0x000000000000000000000000000000000000000A")

	assert.True(t, manager.HasIdentity("0x000000000000000000000000000000000000000A"))
	assert.True(t, manager.HasIdentity("0x000000000000000000000000000000000000000a"))
	assert.False(t, manager.HasIdentity("0x000000000000000000000000000000000000000B"))
}
