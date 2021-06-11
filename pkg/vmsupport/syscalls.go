package vmsupport

import (
	"context"
	"errors"
	"fmt"
	goruntime "runtime"
	"sync"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	rt5 "github.com/filecoin-project/specs-actors/v5/actors/runtime"
	proof5 "github.com/filecoin-project/specs-actors/v5/actors/runtime/proof"
	"github.com/filecoin-project/venus/pkg/util/ffiwrapper"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/venus/pkg/crypto"
	"github.com/filecoin-project/venus/pkg/slashing"
	"github.com/filecoin-project/venus/pkg/state"
	"github.com/filecoin-project/venus/pkg/vm"
)

var log = logging.Logger("vmsupport")

type faultChecker interface {
	VerifyConsensusFault(ctx context.Context, h1, h2, extra []byte, view slashing.FaultStateView) (*rt5.ConsensusFault, error)
}

// Syscalls contains the concrete implementation of VM system calls, including connection to
// proof verification and blockchain inspection.
// Errors returned by these methods are intended to be returned to the actor code to respond to: they must be
// entirely deterministic and repeatable by other implementations.
// Any non-deterministic error will instead trigger a panic.
// TODO: determine a more robust mechanism for distinguishing transient runtime failures from deterministic errors
// in VM and supporting code. https://github.com/filecoin-project/venus/issues/3844
type Syscalls struct {
	faultChecker faultChecker
	verifier     ffiwrapper.Verifier
}

func NewSyscalls(faultChecker faultChecker, verifier ffiwrapper.Verifier) *Syscalls {
	return &Syscalls{
		faultChecker: faultChecker,
		verifier:     verifier,
	}
}

func (s *Syscalls) VerifySignature(ctx context.Context, view vm.SyscallsStateView, signature crypto.Signature, signer address.Address, plaintext []byte) error {
	return state.NewSignatureValidator(view).ValidateSignature(ctx, plaintext, signer, signature)
}

func (s *Syscalls) HashBlake2b(data []byte) [32]byte {
	return blake2b.Sum256(data)
}

func (s *Syscalls) ComputeUnsealedSectorCID(_ context.Context, proof abi.RegisteredSealProof, pieces []abi.PieceInfo) (cid.Cid, error) {
	return ffiwrapper.GenerateUnsealedCID(proof, pieces)
}

func (s *Syscalls) VerifySeal(_ context.Context, info proof5.SealVerifyInfo) error {
	ok, err := s.verifier.VerifySeal(info)
	if err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("seal invalid")
	}
	return nil
}

var BatchSealVerifyParallelism = 2 * goruntime.NumCPU()

func (s *Syscalls) BatchVerifySeals(ctx context.Context, vis map[address.Address][]proof5.SealVerifyInfo) (map[address.Address][]bool, error) {
	out := make(map[address.Address][]bool)

	sema := make(chan struct{}, BatchSealVerifyParallelism)

	var wg sync.WaitGroup
	for addr, seals := range vis {
		results := make([]bool, len(seals))
		out[addr] = results

		for i, seal := range seals {
			wg.Add(1)
			go func(ma address.Address, ix int, svi proof5.SealVerifyInfo, res []bool) {
				defer wg.Done()
				sema <- struct{}{}

				if err := s.VerifySeal(ctx, svi); err != nil {
					log.Warnw("seal verify in batch failed", "miner", ma, "index", ix, "err", err)
					res[ix] = false
				} else {
					res[ix] = true
				}

				<-sema
			}(addr, i, seal, results)
		}
	}
	wg.Wait()

	return out, nil
}

func (s *Syscalls) VerifyAggregateSeals(aggregate proof5.AggregateSealVerifyProofAndInfos) error {
	ok, err := s.verifier.VerifyAggregateSeals(aggregate)
	if err != nil {
		return xerrors.Errorf("failed to verify aggregated PoRep: %w", err)
	}
	if !ok {
		return fmt.Errorf("invalid aggregate proof")
	}

	return nil
}

func (s *Syscalls) VerifyPoSt(ctx context.Context, info proof5.WindowPoStVerifyInfo) error {
	ok, err := s.verifier.VerifyWindowPoSt(ctx, info)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("window PoSt verification failed")
	}
	return nil
}

func (s *Syscalls) VerifyConsensusFault(ctx context.Context, h1, h2, extra []byte, view vm.SyscallsStateView) (*rt5.ConsensusFault, error) {
	return s.faultChecker.VerifyConsensusFault(ctx, h1, h2, extra, view)
}
