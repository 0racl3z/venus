package init

import (
	"github.com/filecoin-project/venus/pkg/specactors"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/venus/pkg/specactors/adt"
	"github.com/filecoin-project/venus/pkg/specactors/builtin"
	"github.com/filecoin-project/venus/pkg/types"

	builtin0 "github.com/filecoin-project/specs-actors/actors/builtin"

	builtin2 "github.com/filecoin-project/specs-actors/v2/actors/builtin"

	builtin3 "github.com/filecoin-project/specs-actors/v3/actors/builtin"

	builtin4 "github.com/filecoin-project/specs-actors/v4/actors/builtin"

	builtin5 "github.com/filecoin-project/specs-actors/v5/actors/builtin"
)

func init() {

	builtin.RegisterActorState(builtin0.InitActorCodeID, func(store adt.Store, root cid.Cid) (cbor.Marshaler, error) {
		return load0(store, root)
	})

	builtin.RegisterActorState(builtin2.InitActorCodeID, func(store adt.Store, root cid.Cid) (cbor.Marshaler, error) {
		return load2(store, root)
	})

	builtin.RegisterActorState(builtin3.InitActorCodeID, func(store adt.Store, root cid.Cid) (cbor.Marshaler, error) {
		return load3(store, root)
	})

	builtin.RegisterActorState(builtin4.InitActorCodeID, func(store adt.Store, root cid.Cid) (cbor.Marshaler, error) {
		return load4(store, root)
	})

	builtin.RegisterActorState(builtin5.InitActorCodeID, func(store adt.Store, root cid.Cid) (cbor.Marshaler, error) {
		return load5(store, root)
	})
}

var (
	Address = builtin5.InitActorAddr
	Methods = builtin5.MethodsInit
)

func Load(store adt.Store, act *types.Actor) (State, error) {
	switch act.Code {

	case builtin0.InitActorCodeID:
		return load0(store, act.Head)

	case builtin2.InitActorCodeID:
		return load2(store, act.Head)

	case builtin3.InitActorCodeID:
		return load3(store, act.Head)

	case builtin4.InitActorCodeID:
		return load4(store, act.Head)

	case builtin5.InitActorCodeID:
		return load5(store, act.Head)

	}
	return nil, xerrors.Errorf("unknown actor code %s", act.Code)
}

func MakeState(store adt.Store, av specactors.Version, networkName string) (State, error) {
	switch av {

	case specactors.Version0:
		return make0(store, networkName)

	case specactors.Version2:
		return make2(store, networkName)

	case specactors.Version3:
		return make3(store, networkName)

	case specactors.Version4:
		return make4(store, networkName)

	case specactors.Version5:
		return make5(store, networkName)

	}
	return nil, xerrors.Errorf("unknown actor version %d", av)
}

func GetActorCodeID(av specactors.Version) (cid.Cid, error) {
	switch av {

	case specactors.Version0:
		return builtin0.InitActorCodeID, nil

	case specactors.Version2:
		return builtin2.InitActorCodeID, nil

	case specactors.Version3:
		return builtin3.InitActorCodeID, nil

	case specactors.Version4:
		return builtin4.InitActorCodeID, nil

	case specactors.Version5:
		return builtin5.InitActorCodeID, nil

	}

	return cid.Undef, xerrors.Errorf("unknown actor version %d", av)
}

type State interface {
	cbor.Marshaler

	ResolveAddress(address address.Address) (address.Address, bool, error)
	MapAddressToNewID(address address.Address) (address.Address, error)
	NetworkName() (string, error)

	ForEachActor(func(id abi.ActorID, address address.Address) error) error

	// Remove exists to support tooling that manipulates state for testing.
	// It should not be used in production code, as init actor entries are
	// immutable.
	Remove(addrs ...address.Address) error

	// Sets the network's name. This should only be used on upgrade/fork.
	SetNetworkName(name string) error

	// Sets the next ID for the init actor. This should only be used for testing.
	SetNextID(id abi.ActorID) error

	// Sets the address map for the init actor. This should only be used for testing.
	SetAddressMap(mcid cid.Cid) error

	AddressMap() (adt.Map, error)
	GetState() interface{}
}
