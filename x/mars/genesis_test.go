package mars_test

import (
	"testing"

	keepertest "github.com/Pantani/mars/testutil/keeper"
	"github.com/Pantani/mars/testutil/nullify"
	"github.com/Pantani/mars/x/mars"
	"github.com/Pantani/mars/x/mars/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MarsKeeper(t)
	mars.InitGenesis(ctx, k, genesisState)
	got := mars.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
