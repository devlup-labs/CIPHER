package main

import (
	"crypto/rand"
	"encoding/hex"
	

	"golang.org/x/crypto/sha3"
)

func GenerateHash(data []byte)string {
  // Generate a unique hash for the chunk data
  hash :=sha3.NewLegacyKeccak256()
  hash.Write(data) 
  result:=hash.Sum(nil)
  return hex.EncodeToString(result)
}

//A commitment is a way to promise something now,without revealing it,yet without being able to change it later.
// Interactions in a commitment scheme take place in two phases:
// the commit phase during which a value is chosen and committed to
// the reveal phase during which the value is revealed by the sender, then the receiver verifies its authenticity
// In the above metaphor, the commit phase is the sender putting the message in the box, and locking it. The reveal phase is the sender giving the key to the receiver, who uses it to open the box and verify its contents. The locked box is the commitment, and the key is the proof.
//Commitment schemes allow the prover to specify all the information in advance, and only reveal what should be revealed later in the proof--Zero-Knowledge proof.
//Receiver gets:commitment hash Receiver learns NOTHING about chunk-Hiding
//binding-because sender cannot change chunk without changing hash.

func GenerateCommitment(data []byte)(string,string,error){

	//random key generation
	key:=make([]byte,32)
	if _,err:=rand.Read(key);err!=nil{
		return "", "", err
	}
	
    keyHash:=hex.EncodeToString(key)

    //concatenate key + data--commitment = keccak256(key || data) --byte concatenation
	
	combined := make([]byte, 0, len(key)+len(data)) //to avoid allocationg new memory,copying old data,resizing everytime its better to mention capacity--len(key)+len(data)
    combined = append(combined, key...)
    combined = append(combined, data...)

	// generate hash
	commitment := GenerateHash(combined)

	return keyHash,commitment ,nil 

}

func VerifyCommitment( keyHex string,data []byte, commitment string,) bool {

    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return false
    }

    combined := make([]byte, 0, len(key)+len(data))
    combined = append(combined, key...)
    combined = append(combined, data...)

	computed:=GenerateHash(combined)

    return computed == commitment
}