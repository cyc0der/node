package identity

import (
	"errors"
	"github.com/mysterium/node/service_discovery/dto"
)

type identityHandler struct {
	manager IdentityManagerInterface
	cache   *identityCache
}

func NewNodeIdentityHandler(keystore keystoreInterface, cacheDir string) *identityHandler {
	return &identityHandler{
		manager: NewIdentityManager(keystore),
		cache:   NewIdentityCache(cacheDir, "remember.json"),
	}
}

func (ih *identityHandler) Select(nodeKey string) (id dto.Identity, err error) {
	if len(nodeKey) > 0 {
		id, err = ih.manager.GetIdentity(nodeKey)
		if err != nil {
			return id, err
		}

		ih.cache.StoreIdentity(id)
		return
	}

	id, err = ih.getIdentityFromCache()
	if err != nil {
		return
	}

	ih.cache.StoreIdentity(id)

	return
}

func (ih *identityHandler) Create() (id dto.Identity, err error) {
	// if all fails, create a new one
	id, err = ih.manager.CreateNewIdentity("")
	if err != nil {
		return id, err
	}

	ih.cache.StoreIdentity(id)

	return
}

func (ih *identityHandler) getIdentityFromCache() (identity dto.Identity, err error) {
	identity, err = ih.cache.GetIdentity()

	if err != nil || !ih.manager.HasIdentity(string(identity)) {
		return identity, errors.New("identity not found in cache")
	}

	return identity, nil
}
