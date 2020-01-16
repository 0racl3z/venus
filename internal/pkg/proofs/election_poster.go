package proofs

import (
	"context"

	"github.com/filecoin-project/go-filecoin/internal/pkg/types"
	"github.com/filecoin-project/go-filecoin/internal/pkg/util/hasher"
	"github.com/filecoin-project/go-filecoin/internal/pkg/vm/address"
	sector "github.com/filecoin-project/go-sectorbuilder"
)

// SectorChallengeRatioDiv is the number of sectors per candidate partial
// ticket
const SectorChallengeRatioDiv = 25

// EPoStCandidate wraps the input data needed to verify an election PoSt
type EPoStCandidate struct {
	SectorID             uint64
	PartialTicket        []byte
	SectorChallengeIndex uint64
}

// ZipEPoStCandidates creates a slice of EPoSt candidates given slices for the fields
// Precondition: all slices have the same length
func ZipEPoStCandidates(sectorIDs, sectorChallengeIndexes []types.Uint64, partialTickets [][]byte) []*EPoStCandidate {
	candidates := make([]*EPoStCandidate, len(sectorIDs))
	for i := 0; i < len(sectorIDs); i++ {
		candidates[i] = &EPoStCandidate{
			SectorID:             uint64(sectorIDs[i]),
			PartialTicket:        partialTickets[i],
			SectorChallengeIndex: uint64(sectorChallengeIndexes[i]),
		}
	}
	return candidates
}

// ElectionPoster generates and verifies electoin PoSts
// Dragons: once we have a proper eposter this type should either be
// replaced or it should be a thin wrapper around the proper eposter
type ElectionPoster struct{}

// VerifyElectionPost returns the validity of the input PoSt proof
func (ep *ElectionPoster) VerifyElectionPost(ctx context.Context, sectorSize uint64, sectorInfo sector.SortedSectorInfo, challengeSeed []byte, proof []byte, candidates []*EPoStCandidate, proverID address.Address) (bool, error) {
	return true, nil
}

// ComputeElectionPoSt returns an election post proving that the partial
// tickets are linked to the sector commitments.
func (ep *ElectionPoster) ComputeElectionPoSt(sectorInfo sector.SortedSectorInfo, challengeSeed []byte, winners []*EPoStCandidate) ([]byte, error) {
	fakePoSt := make([]byte, 1)
	fakePoSt[0] = 0xe
	return fakePoSt, nil
}

// GenerateEPostCandidates generates election post candidates
func (ep *ElectionPoster) GenerateEPostCandidates(sectorInfo sector.SortedSectorInfo, challengeSeed []byte, faults []uint64) ([]*EPoStCandidate, error) {
	// Current fake behavior: generate one partial ticket per sector,
	// each partial ticket is the hash of the challengeSeed and sectorID
	var candidates []*EPoStCandidate
	hasher := hasher.NewHasher()
	for i, si := range sectorInfo.Values() {
		hasher.Int(si.SectorID)
		hasher.Bytes(challengeSeed)
		nextCandidate := &EPoStCandidate{
			SectorID:             si.SectorID,
			SectorChallengeIndex: uint64(i), // fake value of sector idx
			PartialTicket:        hasher.Hash(),
		}
		candidates = append(candidates, nextCandidate)
	}
	return candidates, nil
}

// ElectionPostChallengeCount is the total number of partial tickets allowed by
// the system
func (ep *ElectionPoster) ElectionPostChallengeCount(sectors, faults uint64) uint64 {
	if sectors-faults == 0 {
		return 0
	}
	// ceil(sectors / SectorChallengeRatioDiv)
	return (sectors-faults-1)/SectorChallengeRatioDiv + 1
}
