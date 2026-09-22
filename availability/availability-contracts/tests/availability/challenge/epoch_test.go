package challenge_test

import (
	"testing"

	challenge "cipher/availability/availability-contracts/challenge"
)

func TestEpochSelectionAndReset(t *testing.T) {
	epochID, err := challenge.CreateEpoch("provider-epoch", "file-epoch", 5)
	if err != nil {
		t.Fatalf("CreateEpoch returned error: %v", err)
	}
	state, err := challenge.GetEpochState(epochID)
	if err != nil {
		t.Fatalf("GetEpochState returned error: %v", err)
	}
	if state.UncoveredChunkCount != 5 || state.ChallengeCounter != 0 {
		t.Fatalf("initial state = %+v, want five uncovered chunks and counter zero", state)
	}

	selected := make(map[int]bool)
	for value := byte(0); value < 5; value++ {
		chunkID, err := challenge.SelectUncoveredChunk(epochID, []byte{value})
		if err != nil {
			t.Fatalf("SelectUncoveredChunk returned error: %v", err)
		}
		if selected[chunkID] {
			t.Fatalf("chunk %d was selected twice", chunkID)
		}
		selected[chunkID] = true
	}
	if _, err := challenge.SelectUncoveredChunk(epochID, []byte{0}); err == nil {
		t.Fatal("selection from an empty epoch succeeded")
	}

	nextEpochID, err := challenge.StartNextEpoch("provider-epoch", "file-epoch")
	if err != nil {
		t.Fatalf("StartNextEpoch returned error: %v", err)
	}
	nextState, err := challenge.GetEpochState(nextEpochID)
	if err != nil {
		t.Fatalf("GetEpochState for next epoch returned error: %v", err)
	}
	if nextState.UncoveredChunkCount != 5 || nextState.ChallengeCounter != 0 || nextState.RoundNumber != 1 {
		t.Fatalf("next epoch state = %+v, want reset chunks, counter zero, round one", nextState)
	}
}

func TestCheckRoundTrigger(t *testing.T) {
	if !challenge.CheckRoundTrigger([]byte{0}) {
		t.Fatal("hash value 0 should trigger a new round")
	}
	if challenge.CheckRoundTrigger([]byte{1}) || challenge.CheckRoundTrigger([]byte{2}) {
		t.Fatal("hash values 1 and 2 should not trigger a new round")
	}
}

func TestEpochRejectsInvalidOperations(t *testing.T) {
	if _, err := challenge.CreateEpoch("", "file", 1); err == nil {
		t.Fatal("CreateEpoch accepted an empty provider ID")
	}
	if _, err := challenge.CreateEpoch("provider", "file", 0); err == nil {
		t.Fatal("CreateEpoch accepted zero chunks")
	}
	if _, err := challenge.SelectUncoveredChunk("unknown-epoch", []byte{1}); err == nil {
		t.Fatal("SelectUncoveredChunk accepted an unknown epoch")
	}
	if _, err := challenge.SelectUncoveredChunk("unknown-epoch", nil); err == nil {
		t.Fatal("SelectUncoveredChunk accepted an empty hash")
	}
	if _, err := challenge.GetEpochState("unknown-epoch"); err == nil {
		t.Fatal("GetEpochState accepted an unknown epoch")
	}
	if _, err := challenge.StartNextEpoch("new-provider", "new-file"); err == nil {
		t.Fatal("StartNextEpoch succeeded without a prior epoch")
	}
}
