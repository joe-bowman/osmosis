package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibchooks "github.com/osmosis-labs/osmosis/v12/x/ibc-hooks"
	ibcratelimit "github.com/osmosis-labs/osmosis/v12/x/ibc-rate-limit"
	ibcratelimittypes "github.com/osmosis-labs/osmosis/v12/x/ibc-rate-limit/types"

	icahost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	// IBC Transfer: Defines the "transfer" IBC port
	transfer "github.com/cosmos/ibc-go/v3/modules/apps/transfer"

	_ "github.com/osmosis-labs/osmosis/v12/client/docs/statik"
	owasm "github.com/osmosis-labs/osmosis/v12/wasmbinding"
	epochskeeper "github.com/osmosis-labs/osmosis/v12/x/epochs/keeper"
	epochstypes "github.com/osmosis-labs/osmosis/v12/x/epochs/types"
	gammkeeper "github.com/osmosis-labs/osmosis/v12/x/gamm/keeper"
	gammtypes "github.com/osmosis-labs/osmosis/v12/x/gamm/types"
	incentiveskeeper "github.com/osmosis-labs/osmosis/v12/x/incentives/keeper"
	incentivestypes "github.com/osmosis-labs/osmosis/v12/x/incentives/types"
	lockupkeeper "github.com/osmosis-labs/osmosis/v12/x/lockup/keeper"
	lockuptypes "github.com/osmosis-labs/osmosis/v12/x/lockup/types"
	mintkeeper "github.com/osmosis-labs/osmosis/v12/x/mint/keeper"
	minttypes "github.com/osmosis-labs/osmosis/v12/x/mint/types"
	poolincentives "github.com/osmosis-labs/osmosis/v12/x/pool-incentives"
	poolincentiveskeeper "github.com/osmosis-labs/osmosis/v12/x/pool-incentives/keeper"
	poolincentivestypes "github.com/osmosis-labs/osmosis/v12/x/pool-incentives/types"
	"github.com/osmosis-labs/osmosis/v12/x/superfluid"
	superfluidkeeper "github.com/osmosis-labs/osmosis/v12/x/superfluid/keeper"
	superfluidtypes "github.com/osmosis-labs/osmosis/v12/x/superfluid/types"
	tokenfactorykeeper "github.com/osmosis-labs/osmosis/v12/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/osmosis-labs/osmosis/v12/x/tokenfactory/types"
	"github.com/osmosis-labs/osmosis/v12/x/twap"
	twaptypes "github.com/osmosis-labs/osmosis/v12/x/twap/types"
	"github.com/osmosis-labs/osmosis/v12/x/txfees"
	txfeeskeeper "github.com/osmosis-labs/osmosis/v12/x/txfees/keeper"
	txfeestypes "github.com/osmosis-labs/osmosis/v12/x/txfees/types"
	valsetpref "github.com/osmosis-labs/osmosis/v12/x/valset-pref"
	valsetpreftypes "github.com/osmosis-labs/osmosis/v12/x/valset-pref/types"
)

type AppKeepers struct {
	// keepers, by order of initialization
	// "Special" keepers
	ParamsKeeper     *paramskeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	// "Normal" keepers
	AccountKeeper                *authkeeper.AccountKeeper
	BankKeeper                   *bankkeeper.BaseKeeper
	AuthzKeeper                  *authzkeeper.Keeper
	StakingKeeper                *stakingkeeper.Keeper
	DistrKeeper                  *distrkeeper.Keeper
	SlashingKeeper               *slashingkeeper.Keeper
	IBCKeeper                    *ibckeeper.Keeper
	ICAHostKeeper                *icahostkeeper.Keeper
	TransferKeeper               *ibctransferkeeper.Keeper
	EvidenceKeeper               *evidencekeeper.Keeper
	GAMMKeeper                   *gammkeeper.Keeper
	TwapKeeper                   *twap.Keeper
	LockupKeeper                 *lockupkeeper.Keeper
	EpochsKeeper                 *epochskeeper.Keeper
	IncentivesKeeper             *incentiveskeeper.Keeper
	MintKeeper                   *mintkeeper.Keeper
	PoolIncentivesKeeper         *poolincentiveskeeper.Keeper
	TxFeesKeeper                 *txfeeskeeper.Keeper
	SuperfluidKeeper             *superfluidkeeper.Keeper
	GovKeeper                    *govkeeper.Keeper
	WasmKeeper                   *wasm.Keeper
	ContractKeeper               *wasmkeeper.PermissionedKeeper
	TokenFactoryKeeper           *tokenfactorykeeper.Keeper
	ValidatorSetPreferenceKeeper *valsetpref.Keeper

	// IBC modules
	// transfer module
	RawIcs20TransferAppModule transfer.AppModule
	RateLimitingICS4Wrapper   *ibcratelimit.ICS4Wrapper
	TransferStack             *ibchooks.IBCMiddleware
	Ics20WasmHooks            *ibchooks.WasmHooks
	HooksICS4Wrapper          ibchooks.ICS4Middleware

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey
}

// InitNormalKeepers initializes all 'normal' keepers (account, app, bank, auth, staking, distribution, slashing, transfer, gamm, IBC router, pool incentives, governance, mint, txfees keepers).
func (appKeepers *AppKeepers) InitNormalKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	maccPerms map[string][]string,
	wasmDir string,
	wasmConfig wasm.Config,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
	blockedAddress map[string]bool,
) {
	// Add 'normal' keepers
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		appKeepers.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
	appKeepers.AccountKeeper = &accountKeeper
	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(banktypes.ModuleName),
		blockedAddress,
	)
	appKeepers.BankKeeper = &bankKeeper

	authzKeeper := authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
	)
	appKeepers.AuthzKeeper = &authzKeeper

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	)
	appKeepers.StakingKeeper = &stakingKeeper

	distrKeeper := distrkeeper.NewKeeper(
		appCodec, appKeepers.keys[distrtypes.StoreKey],
		appKeepers.GetSubspace(distrtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		blockedAddress,
	)
	appKeepers.DistrKeeper = &distrKeeper

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.GetSubspace(slashingtypes.ModuleName),
	)
	appKeepers.SlashingKeeper = &slashingKeeper

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
	)

	appKeepers.WireICS20PreWasmKeeper(appCodec, bApp)

	icaHostKeeper := icahostkeeper.NewKeeper(
		appCodec, appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.ScopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)
	appKeepers.ICAHostKeeper = &icaHostKeeper

	icaHostIBCModule := icahost.NewIBCModule(*appKeepers.ICAHostKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		// The transferIBC module is replaced by rateLimitingTransferModule
		AddRoute(ibctransfertypes.ModuleName, appKeepers.TransferStack)
	// Note: the sealing is done after creating wasmd and wiring that up

	// create evidence keeper with router
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)

	gammKeeper := gammkeeper.NewKeeper(
		appCodec, appKeepers.keys[gammtypes.StoreKey],
		appKeepers.GetSubspace(gammtypes.ModuleName),
		appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.DistrKeeper)
	appKeepers.GAMMKeeper = &gammKeeper

	appKeepers.TwapKeeper = twap.NewKeeper(
		appKeepers.keys[twaptypes.StoreKey],
		appKeepers.tkeys[twaptypes.TransientStoreKey],
		appKeepers.GetSubspace(twaptypes.ModuleName),
		appKeepers.GAMMKeeper)

	appKeepers.LockupKeeper = lockupkeeper.NewKeeper(
		appKeepers.keys[lockuptypes.StoreKey],
		// TODO: Visit why this needs to be deref'd
		*appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper, appKeepers.GetSubspace(lockuptypes.ModuleName))

	appKeepers.EpochsKeeper = epochskeeper.NewKeeper(appKeepers.keys[epochstypes.StoreKey])

	txFeesKeeper := txfeeskeeper.NewKeeper(
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.keys[txfeestypes.StoreKey],
		appKeepers.GAMMKeeper,
		appKeepers.GAMMKeeper,
		txfeestypes.FeeCollectorName,
		txfeestypes.NonNativeFeeCollectorName,
	)
	appKeepers.TxFeesKeeper = &txFeesKeeper

	appKeepers.IncentivesKeeper = incentiveskeeper.NewKeeper(
		appKeepers.keys[incentivestypes.StoreKey],
		appKeepers.GetSubspace(incentivestypes.ModuleName),
		appKeepers.BankKeeper,
		appKeepers.LockupKeeper,
		appKeepers.EpochsKeeper,
		appKeepers.DistrKeeper,
		appKeepers.TxFeesKeeper,
	)

	appKeepers.SuperfluidKeeper = superfluidkeeper.NewKeeper(
		appKeepers.keys[superfluidtypes.StoreKey], appKeepers.GetSubspace(superfluidtypes.ModuleName),
		*appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.StakingKeeper, appKeepers.DistrKeeper, appKeepers.EpochsKeeper, appKeepers.LockupKeeper, appKeepers.GAMMKeeper, appKeepers.IncentivesKeeper,
		lockupkeeper.NewMsgServerImpl(appKeepers.LockupKeeper))

	mintKeeper := mintkeeper.NewKeeper(
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper,
		appKeepers.EpochsKeeper,
		authtypes.FeeCollectorName,
	)
	appKeepers.MintKeeper = &mintKeeper

	poolIncentivesKeeper := poolincentiveskeeper.NewKeeper(
		appKeepers.keys[poolincentivestypes.StoreKey],
		appKeepers.GetSubspace(poolincentivestypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.IncentivesKeeper,
		appKeepers.DistrKeeper,
		appKeepers.GAMMKeeper,
	)
	appKeepers.PoolIncentivesKeeper = &poolIncentivesKeeper
	appKeepers.GAMMKeeper.SetPoolIncentivesKeeper(appKeepers.PoolIncentivesKeeper)

	tokenFactoryKeeper := tokenfactorykeeper.NewKeeper(
		appKeepers.keys[tokenfactorytypes.StoreKey],
		appKeepers.GetSubspace(tokenfactorytypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper.WithMintCoinsRestriction(tokenfactorytypes.NewTokenFactoryDenomMintCoinsRestriction()),
		appKeepers.DistrKeeper,
	)
	appKeepers.TokenFactoryKeeper = &tokenFactoryKeeper

	validatorSetPreferenceKeeper := valsetpref.NewKeeper(
		appKeepers.keys[valsetpreftypes.StoreKey],
		appKeepers.GetSubspace(valsetpreftypes.ModuleName),
		appKeepers.StakingKeeper,
	)

	appKeepers.ValidatorSetPreferenceKeeper = &validatorSetPreferenceKeeper

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate,osmosis,cosmwasm_1_1"

	wasmOpts = append(owasm.RegisterCustomPlugins(appKeepers.GAMMKeeper, appKeepers.BankKeeper, appKeepers.TwapKeeper, appKeepers.TokenFactoryKeeper), wasmOpts...)
	wasmOpts = append(owasm.RegisterStargateQueries(*bApp.GRPCQueryRouter(), appCodec), wasmOpts...)

	wasmKeeper := wasm.NewKeeper(
		appCodec,
		appKeepers.keys[wasm.StoreKey],
		appKeepers.GetSubspace(wasm.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.DistrKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)
	appKeepers.WasmKeeper = &wasmKeeper

	// Pass the contract keeper to all the structs (generally ICS4Wrappers for ibc middlewares) that need it
	appKeepers.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(appKeepers.WasmKeeper)
	appKeepers.RateLimitingICS4Wrapper.ContractKeeper = appKeepers.ContractKeeper
	appKeepers.Ics20WasmHooks.ContractKeeper = appKeepers.ContractKeeper

	// wire up x/wasm to IBC
	ibcRouter.AddRoute(wasm.ModuleName, wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper))
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(*appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(*appKeepers.DistrKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(*appKeepers.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(poolincentivestypes.RouterKey, poolincentives.NewPoolIncentivesProposalHandler(*appKeepers.PoolIncentivesKeeper)).
		AddRoute(txfeestypes.RouterKey, txfees.NewUpdateFeeTokenProposalHandler(*appKeepers.TxFeesKeeper)).
		AddRoute(superfluidtypes.RouterKey, superfluid.NewSuperfluidProposalHandler(*appKeepers.SuperfluidKeeper, *appKeepers.EpochsKeeper))

	// The gov proposal types can be individually enabled
	if len(wasmEnabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(appKeepers.WasmKeeper, wasmEnabledProposals))
	}

	govKeeper := govkeeper.NewKeeper(
		appCodec, appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName), appKeepers.AccountKeeper, appKeepers.BankKeeper,
		appKeepers.SuperfluidKeeper, govRouter)
	appKeepers.GovKeeper = &govKeeper
}

// Create the IBC Transfer Stack from bottom to top:
//
// * SendPacket. Originates from the transferKeeper and and goes up the stack:
// transferKeeper.SendPacket -> ibc_rate_limit.SendPacket -> ibc_hooks.SendPacket -> channel.SendPacket
// * RecvPacket, message that originates from core IBC and goes down to app, the flow is the other way
// channel.RecvPacket -> ibc_hooks.OnRecvPacket -> ibc_rate_limit.OnRecvPacket -> transfer.OnRecvPacket
//
// After this, the wasm keeper is required to be set on both
// appkeepers.WasmHooks AND appKeepers.RateLimitingICS4Wrapper
func (appKeepers *AppKeepers) WireICS20PreWasmKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp) {
	// Setup the ICS4Wrapper used by the hooks middleware
	wasmHooks := ibchooks.NewWasmHooks(nil) // The contract keeper needs to be set later
	appKeepers.Ics20WasmHooks = &wasmHooks
	appKeepers.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.Ics20WasmHooks,
	)

	// ChannelKeeper wrapper for rate limiting SendPacket(). The wasmKeeper needs to be added after it's created
	rateLimitingParams := appKeepers.GetSubspace(ibcratelimittypes.ModuleName)
	rateLimitingParams = rateLimitingParams.WithKeyTable(ibcratelimittypes.ParamKeyTable())
	rateLimitingICS4Wrapper := ibcratelimit.NewICS4Middleware(
		appKeepers.HooksICS4Wrapper,
		appKeepers.AccountKeeper,
		// wasm keeper we set later.
		nil,
		appKeepers.BankKeeper,
		rateLimitingParams,
	)
	appKeepers.RateLimitingICS4Wrapper = &rateLimitingICS4Wrapper

	// Create Transfer Keepers
	transferKeeper := ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		// The ICS4Wrapper is replaced by the rateLimitingICS4Wrapper instead of the channel
		appKeepers.RateLimitingICS4Wrapper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)
	appKeepers.TransferKeeper = &transferKeeper
	appKeepers.RawIcs20TransferAppModule = transfer.NewAppModule(*appKeepers.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(*appKeepers.TransferKeeper)

	// RateLimiting IBC Middleware
	rateLimitingTransferModule := ibcratelimit.NewIBCModule(transferIBCModule, appKeepers.RateLimitingICS4Wrapper)
	// Hooks Middleware
	hooksTransferModule := ibchooks.NewIBCMiddleware(&rateLimitingTransferModule, &appKeepers.HooksICS4Wrapper)
	appKeepers.TransferStack = &hooksTransferModule
}

// InitSpecialKeepers initiates special keepers (crisis appkeeper, upgradekeeper, params keeper)
func (appKeepers *AppKeepers) InitSpecialKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	wasmDir string,
	cdc *codec.LegacyAmino,
	invCheckPeriod uint,
	skipUpgradeHeights map[int64]bool,
	homePath string,
) {
	appKeepers.GenerateKeys()
	paramsKeeper := appKeepers.initParamsKeeper(appCodec, cdc, appKeepers.keys[paramstypes.StoreKey], appKeepers.tkeys[paramstypes.TStoreKey])
	appKeepers.ParamsKeeper = &paramsKeeper

	// set the BaseApp's parameter store
	bApp.SetParamStore(appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	appKeepers.ScopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.ScopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	appKeepers.CapabilityKeeper.Seal()

	// TODO: Make a SetInvCheckPeriod fn on CrisisKeeper.
	// IMO, its bad design atm that it requires this in state machine initialization
	crisisKeeper := crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName), invCheckPeriod, appKeepers.BankKeeper, authtypes.FeeCollectorName,
	)
	appKeepers.CrisisKeeper = &crisisKeeper

	upgradeKeeper := upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
	)
	appKeepers.UpgradeKeeper = &upgradeKeeper
}

// initParamsKeeper init params keeper and its subspaces.
func (appKeepers *AppKeepers) initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(incentivestypes.ModuleName)
	paramsKeeper.Subspace(lockuptypes.ModuleName)
	paramsKeeper.Subspace(poolincentivestypes.ModuleName)
	paramsKeeper.Subspace(superfluidtypes.ModuleName)
	paramsKeeper.Subspace(gammtypes.ModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)
	paramsKeeper.Subspace(tokenfactorytypes.ModuleName)
	paramsKeeper.Subspace(twaptypes.ModuleName)
	paramsKeeper.Subspace(ibcratelimittypes.ModuleName)

	return paramsKeeper
}

// SetupHooks sets up hooks for modules.
func (appKeepers *AppKeepers) SetupHooks() {
	// For every module that has hooks set on it,
	// you must check InitNormalKeepers to ensure that its not passed by de-reference
	// e.g. *app.StakingKeeper doesn't appear

	// Recall that SetHooks is a mutative call.
	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
			appKeepers.SuperfluidKeeper.Hooks(),
		),
	)

	appKeepers.GAMMKeeper.SetHooks(
		gammtypes.NewMultiGammHooks(
			// insert gamm hooks receivers here
			appKeepers.PoolIncentivesKeeper.Hooks(),
			appKeepers.TwapKeeper.GammHooks(),
		),
	)

	appKeepers.LockupKeeper.SetHooks(
		lockuptypes.NewMultiLockupHooks(
			// insert lockup hooks receivers here
			appKeepers.SuperfluidKeeper.Hooks(),
		),
	)

	appKeepers.IncentivesKeeper.SetHooks(
		incentivestypes.NewMultiIncentiveHooks(
		// insert incentive hooks receivers here
		),
	)

	appKeepers.MintKeeper.SetHooks(
		minttypes.NewMultiMintHooks(
			// insert mint hooks receivers here
			appKeepers.PoolIncentivesKeeper.Hooks(),
		),
	)

	appKeepers.EpochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(
			// insert epoch hooks receivers here
			appKeepers.TxFeesKeeper.Hooks(),
			appKeepers.TwapKeeper.EpochHooks(),
			appKeepers.SuperfluidKeeper.Hooks(),
			appKeepers.IncentivesKeeper.Hooks(),
			appKeepers.MintKeeper.Hooks(),
		),
	)

	appKeepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// insert governance hooks receivers here
		),
	)
}

// TODO: We need to automate this, by bundling with a module struct...
func KVStoreKeys() []string {
	return []string{
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		ibchost.StoreKey,
		icahosttypes.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		gammtypes.StoreKey,
		twaptypes.StoreKey,
		lockuptypes.StoreKey,
		incentivestypes.StoreKey,
		epochstypes.StoreKey,
		poolincentivestypes.StoreKey,
		authzkeeper.StoreKey,
		txfeestypes.StoreKey,
		superfluidtypes.StoreKey,
		wasm.StoreKey,
		tokenfactorytypes.StoreKey,
		valsetpreftypes.StoreKey,
	}
}
