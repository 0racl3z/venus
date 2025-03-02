package v0api

import (
	"context"
	"io"
	"time"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	acrypto "github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/dline"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"
	"github.com/filecoin-project/venus/app/submodule/apitypes"
	"github.com/filecoin-project/venus/pkg/chain"
	syncTypes "github.com/filecoin-project/venus/pkg/chainsync/types"
	"github.com/filecoin-project/venus/pkg/crypto"
	"github.com/filecoin-project/venus/pkg/messagepool"
	"github.com/filecoin-project/venus/pkg/net"
	"github.com/filecoin-project/venus/pkg/specactors/builtin/miner"
	pstate "github.com/filecoin-project/venus/pkg/state"
	"github.com/filecoin-project/venus/pkg/types"
	"github.com/filecoin-project/venus/pkg/wallet"
)

type FullNodeStruct struct {
	IBlockServiceStruct
	IBlockStoreStruct
	IChainStruct
	IConfigStruct
	IDiscoveryStruct
	IMarketStruct
	IMiningStruct
	IMessagePoolStruct
	IMultiSigStruct
	INetworkStruct
	IPaychanStruct
	ISyncerStruct
	IWalletStruct
	IJwtAuthAPIStruct
}

type IAccountStruct struct {
	StateAccountKey func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) `perm:"read"`
}

type IActorStruct struct {
	ListActor     func(p0 context.Context) (map[address.Address]*types.Actor, error)                     `perm:"read"`
	StateGetActor func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*types.Actor, error) `perm:"read"`
}

type IBeaconStruct struct {
	BeaconGetEntry func(p0 context.Context, p1 abi.ChainEpoch) (*types.BeaconEntry, error) `perm:"read"`
}

type IBlockServiceStruct struct {
	DAGCat         func(p0 context.Context, p1 cid.Cid) (io.Reader, error)   `perm:"read"`
	DAGGetFileSize func(p0 context.Context, p1 cid.Cid) (uint64, error)      `perm:"read"`
	DAGGetNode     func(p0 context.Context, p1 string) (interface{}, error)  `perm:"read"`
	DAGImportData  func(p0 context.Context, p1 io.Reader) (ipld.Node, error) `perm:"read"`
}

type IBlockStoreStruct struct {
	ChainDeleteObj func(p0 context.Context, p1 cid.Cid) error                                 `perm:"read"`
	ChainHasObj    func(p0 context.Context, p1 cid.Cid) (bool, error)                         `perm:"read"`
	ChainReadObj   func(p0 context.Context, p1 cid.Cid) ([]byte, error)                       `perm:"read"`
	ChainStatObj   func(p0 context.Context, p1 cid.Cid, p2 cid.Cid) (apitypes.ObjStat, error) `perm:"read"`
}

type IChainStruct struct {
	IAccountStruct
	IActorStruct
	IBeaconStruct
	IMinerStateStruct
	IChainInfoStruct
}

type IChainInfoStruct struct {
	BlockTime                     func(p0 context.Context) time.Duration                                                                                             `perm:"read"`
	ChainGetBlock                 func(p0 context.Context, p1 cid.Cid) (*types.BlockHeader, error)                                                                   `perm:"read"`
	ChainGetBlockMessages         func(p0 context.Context, p1 cid.Cid) (*apitypes.BlockMessages, error)                                                              `perm:"read"`
	ChainGetMessage               func(p0 context.Context, p1 cid.Cid) (*types.UnsignedMessage, error)                                                               `perm:"read"`
	ChainGetMessagesInTipset      func(p0 context.Context, p1 types.TipSetKey) ([]apitypes.Message, error)                                                           `perm:"read"`
	ChainGetParentMessages        func(p0 context.Context, p1 cid.Cid) ([]apitypes.Message, error)                                                                   `perm:"read"`
	ChainGetParentReceipts        func(p0 context.Context, p1 cid.Cid) ([]*types.MessageReceipt, error)                                                              `perm:"read"`
	ChainGetRandomnessFromBeacon  func(p0 context.Context, p1 types.TipSetKey, p2 acrypto.DomainSeparationTag, p3 abi.ChainEpoch, p4 []byte) (abi.Randomness, error) `perm:"read"`
	ChainGetRandomnessFromTickets func(p0 context.Context, p1 types.TipSetKey, p2 acrypto.DomainSeparationTag, p3 abi.ChainEpoch, p4 []byte) (abi.Randomness, error) `perm:"read"`
	ChainGetReceipts              func(p0 context.Context, p1 cid.Cid) ([]types.MessageReceipt, error)                                                               `perm:"read"`
	ChainGetTipSet                func(p0 context.Context, p1 types.TipSetKey) (*types.TipSet, error)                                                                `perm:"read"`
	ChainGetTipSetByHeight        func(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) (*types.TipSet, error)                                             `perm:"read"`
	ChainHead                     func(p0 context.Context) (*types.TipSet, error)                                                                                    `perm:"read"`
	ChainList                     func(p0 context.Context, p1 types.TipSetKey, p2 int) ([]types.TipSetKey, error)                                                    `perm:"read"`
	ChainNotify                   func(p0 context.Context) <-chan []*chain.HeadChange                                                                                `perm:"read"`
	ChainSetHead                  func(p0 context.Context, p1 types.TipSetKey) error                                                                                 `perm:"read"`
	GetActor                      func(p0 context.Context, p1 address.Address) (*types.Actor, error)                                                                 `perm:"read"`
	GetEntry                      func(p0 context.Context, p1 abi.ChainEpoch, p2 uint64) (*types.BeaconEntry, error)                                                 `perm:"read"`
	GetFullBlock                  func(p0 context.Context, p1 cid.Cid) (*types.FullBlock, error)                                                                     `perm:"read"`
	GetParentStateRootActor       func(p0 context.Context, p1 *types.TipSet, p2 address.Address) (*types.Actor, error)                                               `perm:"read"`
	MessageWait                   func(p0 context.Context, p1 cid.Cid, p2 abi.ChainEpoch, p3 abi.ChainEpoch) (*chain.ChainMessage, error)                            `perm:"read"`
	ProtocolParameters            func(p0 context.Context) (*apitypes.ProtocolParams, error)                                                                         `perm:"read"`
	ResolveToKeyAddr              func(p0 context.Context, p1 address.Address, p2 *types.TipSet) (address.Address, error)                                            `perm:"read"`
	StateNetworkName              func(p0 context.Context) (apitypes.NetworkName, error)                                                                             `perm:"read"`
	StateNetworkVersion           func(p0 context.Context, p1 types.TipSetKey) (network.Version, error)                                                              `perm:"read"`
	StateSearchMsg                func(p0 context.Context, p1 cid.Cid) (*apitypes.MsgLookup, error)                                                                  `perm:"read"`
	StateSearchMsgLimited         func(p0 context.Context, p1 cid.Cid, p2 abi.ChainEpoch) (*apitypes.MsgLookup, error)                                               `perm:"read"`
	StateWaitMsg                  func(p0 context.Context, p1 cid.Cid, p2 uint64) (*apitypes.MsgLookup, error)                                                       `perm:"read"`
	StateWaitMsgLimited           func(p0 context.Context, p1 cid.Cid, p2 uint64, p3 abi.ChainEpoch) (*apitypes.MsgLookup, error)                                    `perm:"read"`
	StateGetReceipt               func(p0 context.Context, p1 cid.Cid, p2 types.TipSetKey) (*types.MessageReceipt, error)                                            `perm:"read"`
	VerifyEntry                   func(p0 *types.BeaconEntry, p1 *types.BeaconEntry, p2 abi.ChainEpoch) bool                                                         `perm:"read"`
}

type IConfigStruct struct {
	ConfigGet func(p0 context.Context, p1 string) (interface{}, error) `perm:"read"`
	ConfigSet func(p0 context.Context, p1 string, p2 string) error     `perm:"read"`
}

type IDiscoveryStruct struct {
}

type IJwtAuthAPIStruct struct {
	AuthNew func(p0 context.Context, p1 []auth.Permission) ([]byte, error)                                             `perm:"read"`
	Verify  func(p0 context.Context, p1 string, p2 string, p3 string, p4 string, p5 string) ([]auth.Permission, error) `perm:"read"`
}

type IMarketStruct struct {
	StateMarketParticipants func(p0 context.Context, p1 types.TipSetKey) (map[string]apitypes.MarketBalance, error) `perm:"read"`
}

type IMessagePoolStruct struct {
	DeleteByAdress             func(p0 context.Context, p1 address.Address) error                                                                                 `perm:"read"`
	GasBatchEstimateMessageGas func(p0 context.Context, p1 []*types.EstimateMessage, p2 uint64, p3 types.TipSetKey) ([]*types.EstimateResult, error)              `perm:"read"`
	GasEstimateFeeCap          func(p0 context.Context, p1 *types.UnsignedMessage, p2 int64, p3 types.TipSetKey) (big.Int, error)                                 `perm:"read"`
	GasEstimateGasLimit        func(p0 context.Context, p1 *types.UnsignedMessage, p2 types.TipSetKey) (int64, error)                                             `perm:"read"`
	GasEstimateGasPremium      func(p0 context.Context, p1 uint64, p2 address.Address, p3 int64, p4 types.TipSetKey) (big.Int, error)                             `perm:"read"`
	GasEstimateMessageGas      func(p0 context.Context, p1 *types.UnsignedMessage, p2 *types.MessageSendSpec, p3 types.TipSetKey) (*types.UnsignedMessage, error) `perm:"read"`
	MpoolBatchPush             func(p0 context.Context, p1 []*types.SignedMessage) ([]cid.Cid, error)                                                             `perm:"read"`
	MpoolBatchPushMessage      func(p0 context.Context, p1 []*types.UnsignedMessage, p2 *types.MessageSendSpec) ([]*types.SignedMessage, error)                   `perm:"read"`
	MpoolBatchPushUntrusted    func(p0 context.Context, p1 []*types.SignedMessage) ([]cid.Cid, error)                                                             `perm:"read"`
	MpoolDeleteByAdress        func(p0 context.Context, p1 address.Address) error                                                                                 `perm:"read"`
	MpoolClear                 func(p0 context.Context, p1 bool) error                                                                                            `perm:"read"`
	MpoolGetConfig             func(p0 context.Context) (*messagepool.MpoolConfig, error)                                                                         `perm:"read"`
	MpoolGetNonce              func(p0 context.Context, p1 address.Address) (uint64, error)                                                                       `perm:"read"`
	MpoolPending               func(p0 context.Context, p1 types.TipSetKey) ([]*types.SignedMessage, error)                                                       `perm:"read"`
	MpoolPublishByAddr         func(p0 context.Context, p1 address.Address) error                                                                                 `perm:"read"`
	MpoolPublishMessage        func(p0 context.Context, p1 *types.SignedMessage) error                                                                            `perm:"read"`
	MpoolPush                  func(p0 context.Context, p1 *types.SignedMessage) (cid.Cid, error)                                                                 `perm:"read"`
	MpoolPushMessage           func(p0 context.Context, p1 *types.UnsignedMessage, p2 *types.MessageSendSpec) (*types.SignedMessage, error)                       `perm:"read"`
	MpoolPushUntrusted         func(p0 context.Context, p1 *types.SignedMessage) (cid.Cid, error)                                                                 `perm:"read"`
	MpoolSelect                func(p0 context.Context, p1 types.TipSetKey, p2 float64) ([]*types.SignedMessage, error)                                           `perm:"read"`
	MpoolSelects               func(p0 context.Context, p1 types.TipSetKey, p2 []float64) ([][]*types.SignedMessage, error)                                       `perm:"read"`
	MpoolSetConfig             func(p0 context.Context, p1 *messagepool.MpoolConfig) error                                                                        `perm:"read"`
	MpoolSub                   func(p0 context.Context) (<-chan messagepool.MpoolUpdate, error)                                                                   `perm:"read"`
}

type IMinerStateStruct struct {
	StateCirculatingSupply             func(p0 context.Context, p1 types.TipSetKey) (abi.TokenAmount, error)                                                           `perm:"read"`
	StateListActors                    func(p0 context.Context, p1 types.TipSetKey) ([]address.Address, error)                                                         `perm:"read"`
	StateListMiners                    func(p0 context.Context, p1 types.TipSetKey) ([]address.Address, error)                                                         `perm:"read"`
	StateLookupID                      func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error)                                       `perm:"read"`
	StateMarketBalance                 func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (apitypes.MarketBalance, error)                                `perm:"read"`
	StateMarketDeals                   func(p0 context.Context, p1 types.TipSetKey) (map[string]pstate.MarketDeal, error)                                              `perm:"read"`
	StateMarketStorageDeal             func(p0 context.Context, p1 abi.DealID, p2 types.TipSetKey) (*apitypes.MarketDeal, error)                                       `perm:"read"`
	StateMinerActiveSectors            func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]*miner.SectorOnChainInfo, error)                            `perm:"read"`
	StateMinerAvailableBalance         func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (big.Int, error)                                               `perm:"read"`
	StateMinerDeadlines                func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]apitypes.Deadline, error)                                   `perm:"read"`
	StateMinerFaults                   func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (bitfield.BitField, error)                                     `perm:"read"`
	StateMinerInfo                     func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (miner.MinerInfo, error)                                       `perm:"read"`
	StateMinerInitialPledgeCollateral  func(p0 context.Context, p1 address.Address, p2 miner.SectorPreCommitInfo, p3 types.TipSetKey) (big.Int, error)                 `perm:"read"`
	StateMinerPartitions               func(p0 context.Context, p1 address.Address, p2 uint64, p3 types.TipSetKey) ([]apitypes.Partition, error)                       `perm:"read"`
	StateMinerPower                    func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*apitypes.MinerPower, error)                                  `perm:"read"`
	StateMinerPreCommitDepositForPower func(p0 context.Context, p1 address.Address, p2 miner.SectorPreCommitInfo, p3 types.TipSetKey) (big.Int, error)                 `perm:"read"`
	StateMinerProvingDeadline          func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*dline.Info, error)                                           `perm:"read"`
	StateMinerRecoveries               func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (bitfield.BitField, error)                                     `perm:"read"`
	StateMinerSectorAllocated          func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (bool, error)                             `perm:"read"`
	StateMinerSectorCount              func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (apitypes.MinerSectors, error)                                 `perm:"read"`
	StateMinerSectorSize               func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (abi.SectorSize, error)                                        `perm:"read"`
	StateMinerSectors                  func(p0 context.Context, p1 address.Address, p2 *bitfield.BitField, p3 types.TipSetKey) ([]*miner.SectorOnChainInfo, error)     `perm:"read"`
	StateMinerWorkerAddress            func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error)                                       `perm:"read"`
	StateSectorExpiration              func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorExpiration, error)          `perm:"read"`
	StateSectorGetInfo                 func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorOnChainInfo, error)         `perm:"read"`
	StateSectorPartition               func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorLocation, error)            `perm:"read"`
	StateSectorPreCommitInfo           func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (miner.SectorPreCommitOnChainInfo, error) `perm:"read"`
	StateVMCirculatingSupplyInternal   func(p0 context.Context, p1 types.TipSetKey) (chain.CirculatingSupply, error)                                                   `perm:"read"`
}

type IMiningStruct struct {
	MinerCreateBlock func(p0 context.Context, p1 *apitypes.BlockTemplate) (*types.BlockMsg, error)                                         `perm:"read"`
	MinerGetBaseInfo func(p0 context.Context, p1 address.Address, p2 abi.ChainEpoch, p3 types.TipSetKey) (*apitypes.MiningBaseInfo, error) `perm:"read"`
}

type IMultiSigStruct struct {
	MsigAddApprove     func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address, p6 bool) (cid.Cid, error)                               `perm:"read"`
	MsigAddCancel      func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 bool) (cid.Cid, error)                                                   `perm:"read"`
	MsigAddPropose     func(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 bool) (cid.Cid, error)                                                              `perm:"read"`
	MsigApprove        func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address) (cid.Cid, error)                                                                                `perm:"read"`
	MsigApproveTxnHash func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address, p4 address.Address, p5 types.BigInt, p6 address.Address, p7 uint64, p8 []byte) (cid.Cid, error) `perm:"read"`
	MsigCancel         func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address, p4 types.BigInt, p5 address.Address, p6 uint64, p7 []byte) (cid.Cid, error)                     `perm:"read"`
	MsigCreate         func(p0 context.Context, p1 uint64, p2 []address.Address, p3 abi.ChainEpoch, p4 types.BigInt, p5 address.Address, p6 types.BigInt) (cid.Cid, error)                         `perm:"read"`
	MsigGetVested      func(p0 context.Context, p1 address.Address, p2 types.TipSetKey, p3 types.TipSetKey) (types.BigInt, error)                                                                  `perm:"read"`
	MsigPropose        func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt, p4 address.Address, p5 uint64, p6 []byte) (cid.Cid, error)                                `perm:"read"`
	MsigRemoveSigner   func(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 bool) (cid.Cid, error)                                                              `perm:"read"`
	MsigSwapApprove    func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address, p6 address.Address) (cid.Cid, error)                    `perm:"read"`
	MsigSwapCancel     func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address) (cid.Cid, error)                                        `perm:"read"`
	MsigSwapPropose    func(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 address.Address) (cid.Cid, error)                                                   `perm:"read"`
}

type INetworkStruct struct {
	NetAddrsListen            func(p0 context.Context) (peer.AddrInfo, error)                                  `perm:"read"`
	NetworkConnect            func(p0 context.Context, p1 []string) (<-chan net.ConnectionResult, error)       `perm:"read"`
	NetworkFindPeer           func(p0 context.Context, p1 peer.ID) (peer.AddrInfo, error)                      `perm:"read"`
	NetworkFindProvidersAsync func(p0 context.Context, p1 cid.Cid, p2 int) <-chan peer.AddrInfo                `perm:"read"`
	NetworkGetBandwidthStats  func(p0 context.Context) metrics.Stats                                           `perm:"admin"`
	NetworkGetClosestPeers    func(p0 context.Context, p1 string) (<-chan peer.ID, error)                      `perm:"read"`
	NetworkGetPeerAddresses   func(p0 context.Context) []ma.Multiaddr                                          `perm:"admin"`
	NetworkGetPeerID          func(p0 context.Context) peer.ID                                                 `perm:"admin"`
	NetworkPeers              func(p0 context.Context, p1 bool, p2 bool, p3 bool) (*net.SwarmConnInfos, error) `perm:"read"`
	Version                   func(p0 context.Context) (apitypes.Version, error)                               `perm:"read"`
}

type IPaychanStruct struct {
	PaychAllocateLane           func(p0 context.Context, p1 address.Address) (uint64, error)                                                               `perm:"read"`
	PaychAvailableFunds         func(p0 context.Context, p1 address.Address) (*apitypes.ChannelAvailableFunds, error)                                      `perm:"read"`
	PaychAvailableFundsByFromTo func(p0 context.Context, p1 address.Address, p2 address.Address) (*apitypes.ChannelAvailableFunds, error)                  `perm:"read"`
	PaychCollect                func(p0 context.Context, p1 address.Address) (cid.Cid, error)                                                              `perm:"read"`
	PaychGet                    func(p0 context.Context, p1 address.Address, p2 address.Address, p3 big.Int) (*apitypes.ChannelInfo, error)                `perm:"read"`
	PaychGetWaitReady           func(p0 context.Context, p1 cid.Cid) (address.Address, error)                                                              `perm:"read"`
	PaychList                   func(p0 context.Context) ([]address.Address, error)                                                                        `perm:"read"`
	PaychNewPayment             func(p0 context.Context, p1 address.Address, p2 address.Address, p3 []apitypes.VoucherSpec) (*apitypes.PaymentInfo, error) `perm:"read"`
	PaychSettle                 func(p0 context.Context, p1 address.Address) (cid.Cid, error)                                                              `perm:"read"`
	PaychStatus                 func(p0 context.Context, p1 address.Address) (*types.PaychStatus, error)                                                   `perm:"read"`
	PaychVoucherAdd             func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 big.Int) (big.Int, error)              `perm:"read"`
	PaychVoucherCheckSpendable  func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 []byte) (bool, error)                  `perm:"read"`
	PaychVoucherCheckValid      func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher) error                                                `perm:"read"`
	PaychVoucherCreate          func(p0 context.Context, p1 address.Address, p2 big.Int, p3 uint64) (*apitypes.VoucherCreateResult, error)                 `perm:"read"`
	PaychVoucherList            func(p0 context.Context, p1 address.Address) ([]*paych.SignedVoucher, error)                                               `perm:"read"`
	PaychVoucherSubmit          func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 []byte) (cid.Cid, error)               `perm:"read"`
}

type ISyncerStruct struct {
	ChainSyncHandleNewTipSet func(p0 context.Context, p1 *types.ChainInfo) error                                                    `perm:"read"`
	ChainTipSetWeight        func(p0 context.Context, p1 types.TipSetKey) (big.Int, error)                                          `perm:"read"`
	Concurrent               func(p0 context.Context) int64                                                                         `perm:"read"`
	SetConcurrent            func(p0 context.Context, p1 int64) error                                                               `perm:"read"`
	StateCall                func(p0 context.Context, p1 *types.UnsignedMessage, p2 types.TipSetKey) (*apitypes.InvocResult, error) `perm:"read"`
	SyncState                func(p0 context.Context) (*apitypes.SyncState, error)                                                  `perm:"read"`
	SyncSubmitBlock          func(p0 context.Context, p1 *types.BlockMsg) error                                                     `perm:"read"`
	SyncerTracker            func(p0 context.Context) *syncTypes.TargetTracker                                                      `perm:"read"`
}

type IWalletStruct struct {
	HasPassword          func(p0 context.Context) bool                                                                         `perm:"admin"`
	LockWallet           func(p0 context.Context) error                                                                        `perm:"admin"`
	SetPassword          func(p0 context.Context, p1 []byte) error                                                             `perm:"admin"`
	UnLockWallet         func(p0 context.Context, p1 []byte) error                                                             `perm:"admin"`
	WalletAddresses      func(p0 context.Context) []address.Address                                                            `perm:"admin"`
	WalletBalance        func(p0 context.Context, p1 address.Address) (abi.TokenAmount, error)                                 `perm:"read"`
	WalletDefaultAddress func(p0 context.Context) (address.Address, error)                                                     `perm:"write"`
	WalletExport         func(p0 address.Address, p1 string) (*crypto.KeyInfo, error)                                          `perm:"admin"`
	WalletHas            func(p0 context.Context, p1 address.Address) (bool, error)                                            `perm:"write"`
	WalletImport         func(p0 *crypto.KeyInfo) (address.Address, error)                                                     `perm:"admin"`
	WalletNewAddress     func(p0 address.Protocol) (address.Address, error)                                                    `perm:"write"`
	WalletSetDefault     func(p0 context.Context, p1 address.Address) error                                                    `perm:"admin"`
	WalletSign           func(p0 context.Context, p1 address.Address, p2 []byte, p3 wallet.MsgMeta) (*crypto.Signature, error) `perm:"sign"`
	WalletSignMessage    func(p0 context.Context, p1 address.Address, p2 *types.UnsignedMessage) (*types.SignedMessage, error) `perm:"sign"`
	WalletState          func(p0 context.Context) int                                                                          `perm:"admin"`
}
