package cli_test

import (
	"fmt"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ignite/mars/testutil/network"
	"github.com/ignite/mars/testutil/nullify"
	"github.com/ignite/mars/x/mars/client/cli"
	"github.com/ignite/mars/x/mars/types"
)

func (s *CLITestSuite) networkWithUserObjects(n int) (*network.Network, []types.User) {
	s.T().Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		user := types.User{
			Id: uint64(i),
		}
		nullify.Fill(&user)
		state.UserList = append(state.UserList, user)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	s.Require().NoError(err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(s.T(), cfg), state.UserList
}

func (s *CLITestSuite) TestShowUser() {
	net, objs := s.networkWithUserObjects(2)
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc string
		id   string
		args []string
		err  error
		obj  types.User
	}{
		{
			desc: "found",
			id:   fmt.Sprintf("%d", objs[0].Id),
			args: common,
			obj:  objs[0],
		},
		{
			desc: "not found",
			id:   "not_found",
			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		s.Run(tc.desc, func() {
			clientCtx := s.clientCtx

			args := []string{tc.id}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdShowUser(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				s.Require().True(ok)
				s.Require().ErrorIs(stat.Err(), tc.err)
			} else {
				s.Require().NoError(err)
				var resp types.QueryGetUserResponse
				s.Require().NoError(net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().NotNil(resp.User)
				s.Require().Equal(
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.User),
				)
			}
		})
	}
}

func (s *CLITestSuite) TestListUser() {
	net, objs := s.networkWithUserObjects(5)
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			clientCtx := s.clientCtx

			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdListUser(), args)
			s.Require().NoError(err)
			var resp types.QueryAllUserResponse
			s.Require().NoError(net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.User), step)
			s.Require().Subset(
				nullify.Fill(objs),
				nullify.Fill(resp.User),
			)
		}
	})
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			clientCtx := s.clientCtx

			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdListUser(), args)
			s.Require().NoError(err)
			var resp types.QueryAllUserResponse
			s.Require().NoError(net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.User), step)
			s.Require().Subset(
				nullify.Fill(objs),
				nullify.Fill(resp.User),
			)
			next = resp.Pagination.NextKey
		}
	})
	s.Run("Total", func() {
		clientCtx := s.clientCtx

		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdListUser(), args)
		s.Require().NoError(err)
		var resp types.QueryAllUserResponse
		s.Require().NoError(net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().NoError(err)
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(
			nullify.Fill(objs),
			nullify.Fill(resp.User),
		)
	})
}
