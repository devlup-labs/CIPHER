// Package types contains the protocol objects shared by Availability packages.
package types

import "time"

type ContractID string
type EpochID string
type ChallengeID string

type ContractState string

const (
	Proposed     ContractState = "PROPOSED"
	Agreed       ContractState = "AGREED"
	Transferring ContractState = "TRANSFERRING"
	Ready        ContractState = "READY"
	Active       ContractState = "ACTIVE"
	Completed    ContractState = "COMPLETED"
	Expired      ContractState = "EXPIRED"
	Terminated   ContractState = "TERMINATED"
	Disputed     ContractState = "DISPUTED"
)

func (s ContractState) IsTerminal() bool {
	return s == Completed || s == Expired || s == Terminated || s == Disputed
}

type AvailabilityContract struct {
	ID            ContractID
	PublisherID   string
	ProviderID    string
	FileID        string
	PaymentAmount int64
	FundedAmount  int64
	Duration      time.Duration
	CreatedAt     time.Time
	EndsAt        time.Time
	State         ContractState
	Challenges    []ChallengeID
	Results       []ChallengeResult
	Settlement    SettlementState
}

type EpochStatus string

const (
	EpochActive EpochStatus = "ACTIVE"
	EpochEnded  EpochStatus = "ENDED"
)

type Epoch struct {
	ID               EpochID
	ProviderID       string
	FileID           string
	TotalChunks      int
	ChallengeCounter uint64
	Uncovered        []int
	Status           EpochStatus
	RoundNumber      uint64
	CreatedAt        time.Time
}

type EpochState struct {
	EpochID             EpochID
	ChallengeCounter    uint64
	UncoveredChunkCount int
	EpochStatus         EpochStatus
	RoundNumber         uint64
}

type Challenge struct {
	ChallengeID      ChallengeID
	ContractID       ContractID
	ProviderID       string
	FileID           string
	EpochID          EpochID
	ChunkID          int
	Nonce            []byte
	ChallengeCounter uint64
	RoundNumber      uint64
	TriggerStatus    bool
	CreatedAt        time.Time
}

type ChallengeResult struct {
	ChallengeID ChallengeID
	Succeeded   bool
	RecordedAt  time.Time
}

type SettlementStatus string

const (
	SettlementPending    SettlementStatus = "PENDING"
	SettlementCompleted  SettlementStatus = "COMPLETED"
	SettlementWithheld   SettlementStatus = "WITHHELD"
	SettlementExpired    SettlementStatus = "EXPIRED"
	SettlementTerminated SettlementStatus = "TERMINATED"
)

type SettlementState struct {
	Status           SettlementStatus
	ReleasedAmount   int64
	WithheldAmount   int64
	SuccessfulProofs int
	FailedProofs     int
	FinalizedAt      time.Time
}
