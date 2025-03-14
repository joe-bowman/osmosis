package apptesting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	gammkeeper "github.com/osmosis-labs/osmosis/v12/x/gamm/keeper"
	"github.com/osmosis-labs/osmosis/v12/x/gamm/pool-models/balancer"
	"github.com/osmosis-labs/osmosis/v12/x/gamm/pool-models/stableswap"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
)

var DefaultAcctFunds sdk.Coins = sdk.NewCoins(
	sdk.NewCoin("uosmo", sdk.NewInt(10000000000)),
	sdk.NewCoin("foo", sdk.NewInt(10000000)),
	sdk.NewCoin("bar", sdk.NewInt(10000000)),
	sdk.NewCoin("baz", sdk.NewInt(10000000)),
)

var DefaultPoolAssets = []balancer.PoolAsset{
	{
		Weight: sdk.NewInt(100),
		Token:  sdk.NewCoin("foo", sdk.NewInt(5000000)),
	},
	{
		Weight: sdk.NewInt(200),
		Token:  sdk.NewCoin("bar", sdk.NewInt(5000000)),
	},
	{
		Weight: sdk.NewInt(300),
		Token:  sdk.NewCoin("baz", sdk.NewInt(5000000)),
	},
	{
		Weight: sdk.NewInt(400),
		Token:  sdk.NewCoin("uosmo", sdk.NewInt(5000000)),
	},
}

var DefaultStableswapLiquidity = sdk.NewCoins(
	sdk.NewCoin("foo", sdk.NewInt(10000000)),
	sdk.NewCoin("bar", sdk.NewInt(10000000)),
	sdk.NewCoin("baz", sdk.NewInt(10000000)),
)

var ImbalancedStableswapLiquidity = sdk.NewCoins(
	sdk.NewCoin("foo", sdk.NewInt(10_000_000_000)),
	sdk.NewCoin("bar", sdk.NewInt(20_000_000_000)),
	sdk.NewCoin("baz", sdk.NewInt(30_000_000_000)),
)

// PrepareBalancerPoolWithCoins returns a balancer pool
// consisted of given coins with equal weight.
func (s *KeeperTestHelper) PrepareBalancerPoolWithCoins(coins ...sdk.Coin) uint64 {
	weights := make([]int64, len(coins))
	for i := 0; i < len(coins); i++ {
		weights[i] = 1
	}
	return s.PrepareBalancerPoolWithCoinsAndWeights(coins, weights)
}

// PrepareBalancerPoolWithCoins returns a balancer pool
// PrepareBalancerPoolWithCoinsAndWeights returns a balancer pool
// consisted of given coins with the specified weights.
func (s *KeeperTestHelper) PrepareBalancerPoolWithCoinsAndWeights(coins sdk.Coins, weights []int64) uint64 {
	var poolAssets []balancer.PoolAsset
	for i, coin := range coins {
		poolAsset := balancer.PoolAsset{
			Weight: sdk.NewInt(weights[i]),
			Token:  coin,
		}
		poolAssets = append(poolAssets, poolAsset)
	}

	return s.PrepareBalancerPoolWithPoolAsset(poolAssets)
}

// PrepareBalancerPool returns a Balancer pool's pool-ID with pool params set in PrepareBalancerPoolWithPoolParams.
func (s *KeeperTestHelper) PrepareBalancerPool() uint64 {
	poolId := s.PrepareBalancerPoolWithPoolParams(balancer.PoolParams{
		SwapFee: sdk.NewDec(0),
		ExitFee: sdk.NewDec(0),
	})

	spotPrice, err := s.App.GAMMKeeper.CalculateSpotPrice(s.Ctx, poolId, "foo", "bar")
	s.NoError(err)
	s.Equal(sdk.NewDec(2).String(), spotPrice.String())
	spotPrice, err = s.App.GAMMKeeper.CalculateSpotPrice(s.Ctx, poolId, "bar", "baz")
	s.NoError(err)
	s.Equal(sdk.NewDecWithPrec(15, 1).String(), spotPrice.String())
	spotPrice, err = s.App.GAMMKeeper.CalculateSpotPrice(s.Ctx, poolId, "baz", "foo")
	s.NoError(err)
	oneThird := sdk.NewDec(1).Quo(sdk.NewDec(3))
	sp := oneThird.MulInt(gammtypes.SpotPriceSigFigs).RoundInt().ToDec().QuoInt(gammtypes.SpotPriceSigFigs)
	s.Equal(sp.String(), spotPrice.String())

	return poolId
}

func (s *KeeperTestHelper) PrepareBasicStableswapPool() uint64 {
	// Mint some assets to the account.
	s.FundAcc(s.TestAccs[0], DefaultAcctFunds)

	params := stableswap.PoolParams{
		SwapFee: sdk.NewDec(0),
		ExitFee: sdk.NewDec(0),
	}

	msg := stableswap.NewMsgCreateStableswapPool(s.TestAccs[0], params, DefaultStableswapLiquidity, []uint64{}, "")
	poolId, err := s.App.GAMMKeeper.CreatePool(s.Ctx, msg)
	s.NoError(err)
	return poolId
}

func (s *KeeperTestHelper) PrepareImbalancedStableswapPool() uint64 {
	// Mint some assets to the account.
	s.FundAcc(s.TestAccs[0], ImbalancedStableswapLiquidity)

	params := stableswap.PoolParams{
		SwapFee: sdk.NewDec(0),
		ExitFee: sdk.NewDec(0),
	}

	msg := stableswap.NewMsgCreateStableswapPool(s.TestAccs[0], params, ImbalancedStableswapLiquidity, []uint64{1, 1, 1}, "")
	poolId, err := s.App.GAMMKeeper.CreatePool(s.Ctx, msg)
	s.NoError(err)
	return poolId
}

// PrepareBalancerPoolWithPoolParams sets up a Balancer pool with poolParams.
func (s *KeeperTestHelper) PrepareBalancerPoolWithPoolParams(poolParams balancer.PoolParams) uint64 {
	// Mint some assets to the account.
	s.FundAcc(s.TestAccs[0], DefaultAcctFunds)

	msg := balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams, DefaultPoolAssets, "")
	poolId, err := s.App.GAMMKeeper.CreatePool(s.Ctx, msg)
	s.NoError(err)
	return poolId
}

// PrepareBalancerPoolWithPoolAsset sets up a Balancer pool with an array of assets.
func (s *KeeperTestHelper) PrepareBalancerPoolWithPoolAsset(assets []balancer.PoolAsset) uint64 {
	// Add coins for pool creation fee + coins needed to mint balances
	fundCoins := sdk.Coins{sdk.NewCoin("uosmo", sdk.NewInt(10000000000))}
	for _, a := range assets {
		fundCoins = fundCoins.Add(a.Token)
	}
	s.FundAcc(s.TestAccs[0], fundCoins)

	msg := balancer.NewMsgCreateBalancerPool(s.TestAccs[0], balancer.PoolParams{
		SwapFee: sdk.ZeroDec(),
		ExitFee: sdk.ZeroDec(),
	}, assets, "")
	poolId, err := s.App.GAMMKeeper.CreatePool(s.Ctx, msg)
	s.NoError(err)
	return poolId
}

func (s *KeeperTestHelper) RunBasicSwap(poolId uint64) {
	denoms, err := s.App.GAMMKeeper.GetPoolDenoms(s.Ctx, poolId)
	s.Require().NoError(err)

	swapIn := sdk.NewCoins(sdk.NewCoin(denoms[0], sdk.NewInt(1000)))
	s.FundAcc(s.TestAccs[0], swapIn)

	msg := gammtypes.MsgSwapExactAmountIn{
		Sender:            s.TestAccs[0].String(),
		Routes:            []gammtypes.SwapAmountInRoute{{PoolId: poolId, TokenOutDenom: denoms[1]}},
		TokenIn:           swapIn[0],
		TokenOutMinAmount: sdk.ZeroInt(),
	}

	gammMsgServer := gammkeeper.NewMsgServerImpl(s.App.GAMMKeeper)
	_, err = gammMsgServer.SwapExactAmountIn(sdk.WrapSDKContext(s.Ctx), &msg)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) RunBasicExit(poolId uint64) {
	shareInAmount := sdk.NewInt(100)
	tokenOutMins := sdk.NewCoins()

	msg := gammtypes.MsgExitPool{
		Sender:        s.TestAccs[0].String(),
		PoolId:        poolId,
		ShareInAmount: shareInAmount,
		TokenOutMins:  tokenOutMins,
	}

	gammMsgServer := gammkeeper.NewMsgServerImpl(s.App.GAMMKeeper)
	_, err := gammMsgServer.ExitPool(sdk.WrapSDKContext(s.Ctx), &msg)
	s.Require().NoError(err)
}

func (s *KeeperTestHelper) RunBasicJoin(poolId uint64) {
	pool, _ := s.App.GAMMKeeper.GetPoolAndPoke(s.Ctx, poolId)
	denoms, err := s.App.GAMMKeeper.GetPoolDenoms(s.Ctx, poolId)
	s.Require().NoError(err)

	tokenIn := sdk.NewCoins()
	for _, denom := range denoms {
		tokenIn = tokenIn.Add(sdk.NewCoin(denom, sdk.NewInt(10000000)))
	}

	s.FundAcc(s.TestAccs[0], sdk.NewCoins(tokenIn...))

	totalPoolShare := pool.GetTotalShares()
	msg := gammtypes.MsgJoinPool{
		Sender:         s.TestAccs[0].String(),
		PoolId:         poolId,
		ShareOutAmount: totalPoolShare.Quo(sdk.NewInt(100000)),
		TokenInMaxs:    tokenIn,
	}

	gammMsgServer := gammkeeper.NewMsgServerImpl(s.App.GAMMKeeper)
	_, err = gammMsgServer.JoinPool(sdk.WrapSDKContext(s.Ctx), &msg)
	s.Require().NoError(err)
}
