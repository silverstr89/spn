package simulation

import (
	"errors"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/spn/testutil/sample"
	"github.com/tendermint/spn/x/profile/keeper"
	"github.com/tendermint/spn/x/profile/types"
)

// FindCoordinatorAccount find a sim account for a coordinator that exists or not
func FindCoordinatorAccount(
	r *rand.Rand,
	ctx sdk.Context,
	k keeper.Keeper,
	accs []simtypes.Account,
	exist bool,
) (simtypes.Account, bool) {
	// Randomize the set for coordinator operation entropy
	r.Shuffle(len(accs), func(i, j int) {
		accs[i], accs[j] = accs[j], accs[i]
	})

	for _, acc := range accs {
		coordByAddress, err := k.GetCoordinatorByAddress(ctx, acc.Address.String())
		found := !errors.Is(err, types.ErrCoordAddressNotFound)
		if found == exist {
			coord, found := k.GetCoordinator(ctx, coordByAddress.CoordinatorID)
			if found && !coord.Active {
				continue
			}
			return acc, true
		}
	}
	return simtypes.Account{}, false
}

// SimulateMsgUpdateValidatorDescription simulates a MsgUpdateValidatorDescription message
func SimulateMsgUpdateValidatorDescription(ak types.AccountKeeper, bk types.BankKeeper, _ keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Select a random account
		simAccount, _ := simtypes.RandomAcc(r, accs)

		desc := sample.ValidatorDescription(sample.String(50))
		msg := types.NewMsgUpdateValidatorDescription(
			simAccount.Address.String(),
			desc.Identity,
			desc.Moniker,
			desc.Website,
			desc.SecurityContact,
			desc.Details,
		)
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgAddValidatorOperatorAddress simulates a MsgAddValidatorOperatorAddress message
func SimulateMsgAddValidatorOperatorAddress(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgAddValidatorOperatorAddress{}

		// TODO: Handling the AddValidatorOperatorAddress simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "AddValidatorOperatorAddress simulation not implemented"), nil, nil
	}
}

// SimulateMsgCreateCoordinator simulates a MsgCreateCoordinator message
func SimulateMsgCreateCoordinator(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Find an account with no coordinator
		simAccount, found := FindCoordinatorAccount(r, ctx, k, accs, false)
		if !found {
			// No message if all coordinator created
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateCoordinator, "skip coordinator creation"), nil, nil
		}

		msg := types.NewMsgCreateCoordinator(
			simAccount.Address.String(),
			sample.String(30),
			sample.String(30),
			sample.String(30),
		)
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgUpdateCoordinatorDescription simulates a MsgUpdateCoordinatorDescription message
func SimulateMsgUpdateCoordinatorDescription(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Find an account with coordinator associated
		simAccount, found := FindCoordinatorAccount(r, ctx, k, accs, true)
		if !found {
			// No message if no coordinator
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateCoordinatorDescription, "skip update coordinator description"), nil, nil
		}

		desc := sample.CoordinatorDescription()
		msg := types.NewMsgUpdateCoordinatorDescription(
			simAccount.Address.String(),
			desc.Identity,
			desc.Website,
			desc.Details,
		)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgUpdateCoordinatorAddress simulates a MsgUpdateCoordinatorAddress message
func SimulateMsgUpdateCoordinatorAddress(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Select a random account
		coord, found := FindCoordinatorAccount(r, ctx, k, accs, true)
		if !found {
			// No message if no coordinator
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateCoordinatorAddress, "skip update coordinator address"), nil, nil
		}
		simAccount, found := FindCoordinatorAccount(r, ctx, k, accs, false)
		if !found && coord.Address.String() != simAccount.Address.String() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateCoordinatorAddress, "skip update coordinator address"), nil, nil
		}
		msg := types.NewMsgUpdateCoordinatorAddress(coord.Address.String(), simAccount.Address.String())
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      coord,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgDisableCoordinator simulates a MsgDisableCoordinator message
func SimulateMsgDisableCoordinator(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Find an account with coordinator associated
		// avoid delete coordinator associated a chain (id 0,1,2)
		simAccount, found := FindCoordinatorAccount(r, ctx, k, accs[3:], true)
		if !found {
			// No message if no coordinator
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDisableCoordinator, "skip update coordinator delete"), nil, nil
		}

		msg := types.NewMsgDisableCoordinator(simAccount.Address.String())
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
