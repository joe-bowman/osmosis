package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/osmosis-labs/osmosis/v12/osmoutils"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
	incentivestypes "github.com/osmosis-labs/osmosis/v12/x/incentives/types"
	lockuptypes "github.com/osmosis-labs/osmosis/v12/x/lockup/types"
	"github.com/osmosis-labs/osmosis/v12/x/pool-incentives/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey sdk.StoreKey

	paramSpace paramtypes.Subspace

	accountKeeper    types.AccountKeeper
	bankKeeper       types.BankKeeper
	incentivesKeeper types.IncentivesKeeper
	distrKeeper      types.DistrKeeper
	gammKeeper       types.GAMMKeeper
}

func NewKeeper(storeKey sdk.StoreKey, paramSpace paramtypes.Subspace, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, incentivesKeeper types.IncentivesKeeper, distrKeeper types.DistrKeeper, gammKeeper types.GAMMKeeper) Keeper {
	// ensure pool-incentives module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey: storeKey,

		paramSpace: paramSpace,

		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		incentivesKeeper: incentivesKeeper,
		distrKeeper:      distrKeeper,
		gammKeeper:       gammKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) CreatePoolGauges(ctx sdk.Context, poolId uint64) error {
	// Create the same number of gaugeges as there are LockableDurations
	for _, lockableDuration := range k.GetLockableDurations(ctx) {
		gaugeId, err := k.incentivesKeeper.CreateGauge(
			ctx,
			true,
			k.accountKeeper.GetModuleAddress(types.ModuleName),
			sdk.Coins{},
			lockuptypes.QueryCondition{
				LockQueryType: lockuptypes.ByDuration,
				Denom:         gammtypes.GetPoolShareDenom(poolId),
				Duration:      lockableDuration,
				Timestamp:     time.Time{},
			},
			// QUESTION: Should we set the startTime as the epoch start time that the modules share or the current block time?
			ctx.BlockTime(),
			1,
		)
		if err != nil {
			return err
		}

		k.SetPoolGaugeId(ctx, poolId, lockableDuration, gaugeId)
	}

	return nil
}

func (k Keeper) SetPoolGaugeId(ctx sdk.Context, poolId uint64, lockableDuration time.Duration, gaugeId uint64) {
	key := types.GetPoolGaugeIdStoreKey(poolId, lockableDuration)
	store := ctx.KVStore(k.storeKey)
	store.Set(key, sdk.Uint64ToBigEndian(gaugeId))

	key = types.GetPoolIdFromGaugeIdStoreKey(gaugeId, lockableDuration)
	store.Set(key, sdk.Uint64ToBigEndian(poolId))
}

func (k Keeper) GetPoolGaugeId(ctx sdk.Context, poolId uint64, lockableDuration time.Duration) (uint64, error) {
	key := types.GetPoolGaugeIdStoreKey(poolId, lockableDuration)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if len(bz) == 0 {
		return 0, sdkerrors.Wrapf(types.ErrNoGaugeIdExist, "gauge id for pool (%d) with duration (%s) not exist", poolId, lockableDuration.String())
	}

	return sdk.BigEndianToUint64(bz), nil
}

func (k Keeper) GetPoolIdFromGaugeId(ctx sdk.Context, gaugeId uint64, lockableDuration time.Duration) (uint64, error) {
	key := types.GetPoolIdFromGaugeIdStoreKey(gaugeId, lockableDuration)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if len(bz) == 0 {
		return 0, sdkerrors.Wrapf(types.ErrNoGaugeIdExist, "pool for gauge id (%d) with duration (%s) not exist", gaugeId, lockableDuration.String())
	}

	return sdk.BigEndianToUint64(bz), nil
}

func (k Keeper) SetLockableDurations(ctx sdk.Context, lockableDurations []time.Duration) {
	store := ctx.KVStore(k.storeKey)
	info := types.LockableDurationsInfo{LockableDurations: lockableDurations}
	osmoutils.MustSet(store, types.LockableDurationsKey, &info)
}

func (k Keeper) GetLockableDurations(ctx sdk.Context) []time.Duration {
	store := ctx.KVStore(k.storeKey)
	info := types.LockableDurationsInfo{}
	osmoutils.MustGet(store, types.LockableDurationsKey, &info)
	return info.LockableDurations
}

func (k Keeper) GetAllGauges(ctx sdk.Context) []incentivestypes.Gauge {
	gauges := k.incentivesKeeper.GetGauges(ctx)
	return gauges
}

func (k Keeper) IsPoolIncentivized(ctx sdk.Context, poolId uint64) bool {
	lockableDurations := k.GetLockableDurations(ctx)
	distrInfo := k.GetDistrInfo(ctx)

	candidateGaugeIds := []uint64{}
	for _, lockableDuration := range lockableDurations {
		gaugeId, err := k.GetPoolGaugeId(ctx, poolId, lockableDuration)
		if err == nil {
			candidateGaugeIds = append(candidateGaugeIds, gaugeId)
		}
	}

	for _, record := range distrInfo.Records {
		for _, gaugeId := range candidateGaugeIds {
			if record.GaugeId == gaugeId {
				return true
			}
		}
	}
	return false
}
