package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v12/x/gamm/pool-models/balancer"
	balancertypes "github.com/osmosis-labs/osmosis/v12/x/gamm/pool-models/balancer"
	"github.com/osmosis-labs/osmosis/v12/x/gamm/pool-models/stableswap"
	"github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

// import (
// 	"math/rand"
// 	"time"

// 	"github.com/cosmos/cosmos-sdk/simapp"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/osmosis-labs/osmosis/v12/x/gamm/types"
// )

// func (suite *KeeperTestSuite) TestCleanupPool() {
// 	// Mint some assets to the accounts.
// 	for _, acc := range suite.TestAccs {
// 		suite.FundAcc(
// 			suite.App.BankKeeper,
// 			suite.Ctx,
// 			acc,
// 			sdk.NewCoins(
// 				sdk.NewCoin("uosmo", sdk.NewInt(1000000000)),
// 				sdk.NewCoin("foo", sdk.NewInt(1000)),
// 				sdk.NewCoin("bar", sdk.NewInt(1000)),
// 				sdk.NewCoin("baz", sdk.NewInt(1000)),
// 			),
// 		)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	poolId, err := suite.App.GAMMKeeper.CreateBalancerPool(suite.Ctx, acc1, defaultPoolParams, []types.PoolAsset{
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("foo", sdk.NewInt(1000)),
// 		},
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("bar", sdk.NewInt(1000)),
// 		},
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("baz", sdk.NewInt(1000)),
// 		},
// 	}, "")
// 	suite.NoError(err)

// 	for _, acc := range []sdk.AccAddress{acc2, acc3} {
// 		err = suite.App.GAMMKeeper.JoinPool(suite.Ctx, acc, poolId, types.OneShare.MulRaw(100), sdk.NewCoins(
// 			sdk.NewCoin("foo", sdk.NewInt(1000)),
// 			sdk.NewCoin("bar", sdk.NewInt(1000)),
// 			sdk.NewCoin("baz", sdk.NewInt(1000)),
// 		))
// 		suite.NoError(err)
// 	}

// 	pool, err := suite.App.GAMMKeeper.GetPool(suite.Ctx, poolId)
// 	suite.NoError(err)
// 	denom := pool.GetTotalShares().Denom
// 	totalAmount := sdk.ZeroInt()
// 	for _, acc := range suite.TestAccs {
// 		coin := suite.App.BankKeeper.GetBalance(suite.Ctx, acc, denom)
// 		suite.True(coin.Amount.Equal(types.OneShare.MulRaw(100)))
// 		totalAmount = totalAmount.Add(coin.Amount)
// 	}
// 	suite.True(totalAmount.Equal(types.OneShare.MulRaw(300)))

// 	err = suite.App.GAMMKeeper.CleanupBalancerPool(suite.Ctx, []uint64{poolId}, []string{})
// 	suite.NoError(err)
// 	for _, acc := range suite.TestAccs {
// 		for _, denom := range []string{"foo", "bar", "baz"} {
// 			amt := suite.App.BankKeeper.GetBalance(suite.Ctx, acc, denom)
// 			suite.True(amt.Amount.Equal(sdk.NewInt(1000)),
// 				"Expected equal %s: %d, %d", amt.Denom, amt.Amount.Int64(), 1000)
// 		}
// 	}
// }

// func (suite *KeeperTestSuite) TestCleanupPoolRandomized() {
// 	// address => deposited coins
// 	coinOf := make(map[string]sdk.Coins)
// 	denoms := []string{"foo", "bar", "baz"}

// 	// Mint some assets to the accounts.
// 	for _, acc := range suite.TestAccs {
// 		coins := make(sdk.Coins, 3)
// 		for i := range coins {
// 			amount := sdk.NewInt(rand.Int63n(1000))
// 			// give large amount of coins to the pool creator
// 			if i == 0 {
// 				amount = amount.MulRaw(10000)
// 			}
// 			coins[i] = sdk.Coin{denoms[i], amount}
// 		}
// 		coinOf[acc.String()] = coins
// 		coins = append(coins, sdk.NewCoin("uosmo", sdk.NewInt(1000000000)))

// 		suite.FundAcc(
// 			suite.App.BankKeeper,
// 			suite.Ctx,
// 			acc,
// 			coins.Sort(),
// 		)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	initialAssets := []types.PoolAsset{}
// 	for _, coin := range coinOf[acc1.String()] {
// 		initialAssets = append(initialAssets, types.PoolAsset{Weight: types.OneShare.MulRaw(100), Token: coin})
// 	}
// 	poolId, err := suite.App.GAMMKeeper.CreateBalancerPool(suite.Ctx, acc1, defaultPoolParams, initialAssets, "")
// 	suite.NoError(err)

// 	for _, acc := range []sdk.AccAddress{acc2, acc3} {
// 		err = suite.App.GAMMKeeper.JoinPool(suite.Ctx, acc, poolId, types.OneShare, coinOf[acc.String()])
// 		suite.NoError(err)
// 	}

// 	err = suite.App.GAMMKeeper.CleanupBalancerPool(suite.Ctx, []uint64{poolId}, []string{})
// 	suite.NoError(err)
// 	for _, acc := range suite.TestAccs {
// 		for _, coin := range coinOf[acc.String()] {
// 			amt := suite.App.BankKeeper.GetBalance(suite.Ctx, acc, coin.Denom)
// 			// the refund could have rounding error
// 			suite.True(amt.Amount.Sub(coin.Amount).Abs().LTE(sdk.NewInt(2)),
// 				"Expected equal %s: %d, %d", amt.Denom, amt.Amount.Int64(), coin.Amount.Int64())
// 		}
// 	}
// }

// func (suite *KeeperTestSuite) TestCleanupPoolErrorOnSwap() {
// 	suite.Ctx = suite.Ctx.WithBlockTime(time.Unix(1000, 1000))
// 	suite.FundAcc(
// 		suite.App.BankKeeper,
// 		suite.Ctx,
// 		acc1,
// 		sdk.NewCoins(
// 			sdk.NewCoin("uosmo", sdk.NewInt(1000000000)),
// 			sdk.NewCoin("foo", sdk.NewInt(1000)),
// 			sdk.NewCoin("bar", sdk.NewInt(1000)),
// 			sdk.NewCoin("baz", sdk.NewInt(1000)),
// 		),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	poolId, err := suite.App.GAMMKeeper.CreateBalancerPool(suite.Ctx, acc1, defaultPoolParams, []types.PoolAsset{
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("foo", sdk.NewInt(1000)),
// 		},
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("bar", sdk.NewInt(1000)),
// 		},
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("baz", sdk.NewInt(1000)),
// 		},
// 	}, "")
// 	suite.NoError(err)

// 	err = suite.App.GAMMKeeper.CleanupBalancerPool(suite.Ctx, []uint64{poolId}, []string{})
// 	suite.NoError(err)

// 	_, _, err = suite.App.GAMMKeeper.SwapExactAmountIn(suite.Ctx, acc1, poolId, sdk.NewCoin("foo", sdk.NewInt(1)), "bar", sdk.NewInt(1))
// 	suite.Error(err)
// }

// func (suite *KeeperTestSuite) TestCleanupPoolWithLockup() {
// 	suite.Ctx = suite.Ctx.WithBlockTime(time.Unix(1000, 1000))
// 	suite.FundAcc(
// 		suite.App.BankKeeper,
// 		suite.Ctx,
// 		acc1,
// 		sdk.NewCoins(
// 			sdk.NewCoin("uosmo", sdk.NewInt(1000000000)),
// 			sdk.NewCoin("foo", sdk.NewInt(1000)),
// 			sdk.NewCoin("bar", sdk.NewInt(1000)),
// 			sdk.NewCoin("baz", sdk.NewInt(1000)),
// 		),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	poolId, err := suite.App.GAMMKeeper.CreateBalancerPool(suite.Ctx, acc1, defaultPoolParams, []types.PoolAsset{
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("foo", sdk.NewInt(1000)),
// 		},
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("bar", sdk.NewInt(1000)),
// 		},
// 		{
// 			Weight: sdk.NewInt(100),
// 			Token:  sdk.NewCoin("baz", sdk.NewInt(1000)),
// 		},
// 	}, "")
// 	suite.NoError(err)

// 	_, err = suite.App.LockupKeeper.LockTokens(suite.Ctx, acc1, sdk.Coins{sdk.NewCoin(types.GetPoolShareDenom(poolId), types.InitPoolSharesSupply)}, time.Hour)
// 	suite.NoError(err)

// 	for _, lock := range suite.App.LockupKeeper.GetLocksDenom(suite.Ctx, types.GetPoolShareDenom(poolId)) {
// 		err = suite.App.LockupKeeper.ForceUnlock(suite.Ctx, lock)
// 		suite.NoError(err)
// 	}

// 	err = suite.App.GAMMKeeper.CleanupBalancerPool(suite.Ctx, []uint64{poolId}, []string{})
// 	suite.NoError(err)
// 	for _, coin := range []string{"foo", "bar", "baz"} {
// 		amt := suite.App.BankKeeper.GetBalance(suite.Ctx, acc1, coin)
// 		// the refund could have rounding error
// 		suite.True(amt.Amount.Equal(sdk.NewInt(1000)) || amt.Amount.Equal(sdk.NewInt(1000).SubRaw(1)),
// 			"Expected equal %s: %d, %d", amt.Denom, amt.Amount.Int64(), sdk.NewInt(1000).Int64())
// 	}
// }

// TestGetPoolAndPoke tests that the right pools is returned from GetPoolAndPoke.
// For the pools implementing the weighted extension, asserts that PokePool is called.
func (suite *KeeperTestSuite) TestGetPoolAndPoke() {
	const (
		startTime = 1000
		blockTime = startTime + 100
	)

	// N.B.: We make a copy because SmoothWeightChangeParams get mutated.
	// We would like to avoid mutating global pool assets that are used in other tests.
	defaultPoolAssetsCopy := make([]balancertypes.PoolAsset, 2)
	copy(defaultPoolAssetsCopy, defaultPoolAssets)

	startPoolWeightAssets := []balancertypes.PoolAsset{
		{
			Weight: defaultPoolAssets[0].Weight.Quo(sdk.NewInt(2)),
			Token:  defaultPoolAssets[0].Token,
		},
		{
			Weight: defaultPoolAssets[1].Weight.Mul(sdk.NewInt(3)),
			Token:  defaultPoolAssets[1].Token,
		},
	}

	tests := map[string]struct {
		isPokePool bool
		poolId     uint64
	}{
		"weighted pool - change weights": {
			isPokePool: true,
			poolId: suite.prepareCustomBalancerPool(defaultAcctFunds, startPoolWeightAssets, balancer.PoolParams{
				SwapFee: defaultSwapFee,
				ExitFee: defaultExitFee,
				SmoothWeightChangeParams: &balancer.SmoothWeightChangeParams{
					StartTime:          time.Unix(startTime, 0), // start time is before block time so the weights should change
					Duration:           time.Hour,
					InitialPoolWeights: startPoolWeightAssets,
					TargetPoolWeights:  defaultPoolAssetsCopy,
				},
			}),
		},
		"non weighted pool": {
			poolId: suite.prepareCustomStableswapPool(
				defaultAcctFunds,
				stableswap.PoolParams{
					SwapFee: defaultSwapFee,
					ExitFee: defaultExitFee,
				},
				sdk.NewCoins(sdk.NewCoin(defaultAcctFunds[0].Denom, defaultAcctFunds[0].Amount.QuoRaw(2)), sdk.NewCoin(defaultAcctFunds[1].Denom, defaultAcctFunds[1].Amount.QuoRaw(2))),
				[]uint64{1, 1},
			),
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			k := suite.App.GAMMKeeper
			ctx := suite.Ctx.WithBlockTime(time.Unix(blockTime, 0))

			pool, err := k.GetPoolAndPoke(ctx, tc.poolId)

			suite.Require().NoError(err)
			suite.Require().Equal(tc.poolId, pool.GetId())

			if tc.isPokePool {
				pokePool, ok := pool.(types.WeightedPoolExtension)
				suite.Require().True(ok)

				poolAssetWeight0, err := pokePool.GetTokenWeight(startPoolWeightAssets[0].Token.Denom)
				suite.Require().NoError(err)

				poolAssetWeight1, err := pokePool.GetTokenWeight(startPoolWeightAssets[1].Token.Denom)
				suite.Require().NoError(err)

				suite.Require().NotEqual(startPoolWeightAssets[0].Weight, poolAssetWeight0)
				suite.Require().NotEqual(startPoolWeightAssets[1].Weight, poolAssetWeight1)
				return
			}

			_, ok := pool.(types.WeightedPoolExtension)
			suite.Require().False(ok)
		})
	}
}
