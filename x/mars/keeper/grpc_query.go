package keeper

import (
	"github.com/cosmonaut/mars/x/mars/types"
)

var _ types.QueryServer = Keeper{}
