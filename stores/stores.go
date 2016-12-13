package stores

import (
	"github.com/Everlane/evan/common"
)

// Compiler verification that the implementations conform to the interface.
func _verify() []common.Store {
	stores := make([]common.Store, 0)
	stores = append(stores, &ProcessLocalStore{})
	return stores
}
