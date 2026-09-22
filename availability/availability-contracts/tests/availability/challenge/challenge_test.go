package challenge_test

import (
	"testing"

	challenge "cipher/availability/availability-contracts/challenge"
)

func TestCreateAndGenerateNextChallenge(t *testing.T) {
	created, err := challenge.CreateChallenge("contract-challenge", "provider-challenge", "file-challenge", 3, 0)
	if err != nil {
		t.Fatalf("CreateChallenge returned error: %v", err)
	}
	if created.ContractID != "contract-challenge" || created.ProviderID != "provider-challenge" || created.FileID != "file-challenge" || created.EpochID == "" || created.ChallengeID == "" {
		t.Fatalf("created challenge has incorrect context: %+v", created)
	}
	if created.ChunkID < 0 || created.ChunkID >= 3 || len(created.Nonce) != 32 {
		t.Fatalf("created challenge has invalid chunk or nonce: %+v", created)
	}
	state, err := challenge.GetEpochState(created.EpochID)
	if err != nil {
		t.Fatalf("GetEpochState returned error: %v", err)
	}
	next, err := challenge.GenerateNextChallenge("contract-challenge", "provider-challenge", "file-challenge", created.EpochID, state.ChallengeCounter, 1)
	if err != nil {
		t.Fatalf("GenerateNextChallenge returned error: %v", err)
	}
	if next.ChunkID == created.ChunkID || next.ChallengeCounter != 1 || next.RoundNumber != 1 {
		t.Fatalf("next challenge did not progress correctly: %+v", next)
	}
	roundHash := challenge.GenerateRoundHash(next.ProviderID, next.FileID, next.EpochID, next.ChallengeCounter, next.Nonce, next.RoundNumber)
	if next.TriggerStatus != challenge.CheckRoundTrigger(roundHash) {
		t.Fatal("challenge trigger status does not match its round hash")
	}
	if _, err := challenge.GenerateNextChallenge("contract-challenge", "provider-challenge", "file-challenge", created.EpochID, 0, 2); err == nil {
		t.Fatal("GenerateNextChallenge accepted a stale challenge counter")
	}
}

func TestCreateChallengeRejectsInvalidInput(t *testing.T) {
	if _, err := challenge.CreateChallenge("", "provider", "file", 1, 0); err == nil {
		t.Fatal("CreateChallenge accepted an empty contract ID")
	}
}
