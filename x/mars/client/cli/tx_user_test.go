package cli_test

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ignite/mars/x/mars/client/cli"
)

func (s *CLITestSuite) TestCreateUser() {
	fields := []string{"xyz"}
	tests := []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addr.String()),
				fmt.Sprintf("--%s=%s", flags.FlagChainID, s.clientCtx.ChainID),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync), // sync mode as there are no funds yet
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10))).String()),
			},
		},
	}
	for _, tc := range tests {
		s.Run(tc.desc, func() {
			var (
				cmd       = cli.CmdCreateUser()
				clientCtx = s.clientCtx
			)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append(fields, tc.args...))
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)

			var txResp sdk.TxResponse
			s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &txResp), out.String())
			s.Require().Equal(tc.code, txResp.Code)
		})
	}
}

func (s *CLITestSuite) TestUpdateUser() {
	fields := []string{"xyz"}
	common := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addr.String()),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, s.clientCtx.ChainID),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync), // sync mode as there are no funds yet
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10))).String()),
	}
	_, err := clitestutil.ExecTestCLICmd(s.clientCtx, cli.CmdCreateUser(), append(fields, common...))
	s.Require().NoError(err)

	tests := []struct {
		desc string
		id   string
		args []string
		code uint32
		err  error
	}{
		{
			desc: "valid",
			id:   "0",
			args: common,
		},
		{
			desc: "key not found",
			id:   "1",
			args: common,
			code: sdkerrors.ErrKeyNotFound.ABCICode(),
		},
		{
			desc: "invalid key",
			id:   "invalid",
			args: common,
			err:  strconv.ErrSyntax,
		},
	}
	for _, tc := range tests {
		s.Run(tc.desc, func() {
			var (
				cmd       = cli.CmdUpdateUser()
				clientCtx = s.clientCtx
			)

			args := []string{tc.id}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)

			var txResp sdk.TxResponse
			s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &txResp))
			s.Require().Equal(tc.code, txResp.Code)
		})
	}
}

func (s *CLITestSuite) TestDeleteUser() {
	fields := []string{"xyz"}
	common := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.addr.String()),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, s.clientCtx.ChainID),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync), // sync mode as there are no funds yet
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10))).String()),
	}
	_, err := clitestutil.ExecTestCLICmd(s.clientCtx, cli.CmdCreateUser(), append(fields, common...))
	s.Require().NoError(err)

	tests := []struct {
		desc string
		id   string
		args []string
		code uint32
		err  error
	}{
		{
			desc: "valid",
			id:   "0",
			args: common,
		},
		{
			desc: "key not found",
			id:   "1",
			args: common,
			code: sdkerrors.ErrKeyNotFound.ABCICode(),
		},
		{
			desc: "invalid key",
			id:   "invalid",
			args: common,
			err:  strconv.ErrSyntax,
		},
	}
	for _, tc := range tests {
		s.Run(tc.desc, func() {
			var (
				cmd       = cli.CmdDeleteUser()
				clientCtx = s.clientCtx
			)

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, append([]string{tc.id}, tc.args...))
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)

			var txResp sdk.TxResponse
			s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &txResp))
			s.Require().Equal(tc.code, txResp.Code)
		})
	}
}
