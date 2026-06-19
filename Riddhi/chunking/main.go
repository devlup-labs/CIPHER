package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
    chunker := DefaultFileChunker{
        chunkSize: 32768,
    }

    chunks, err := chunker.ChunkFile("test.pdf")
    if err != nil {
        panic(err)
    }

   for _, chunk := range chunks {
    fmt.Printf(
        "Chunk %d | Size: %d\n",
        chunk.ChunkIndex,
        chunk.ChunkSize,
       
    )
}

for _, chunk := range chunks {

    data, err := os.ReadFile(chunk.FileName)
    if err != nil {
        panic(err)
    }

    hash := GenerateHash(data)

    fmt.Printf(
        "Chunk %d | Hash: %s\n",
        chunk.ChunkIndex,
        hash,
    )
}

data := []byte("hello world")
//data1:=[]byte("hello worlds")

key, commitment, err := GenerateCommitment(data)
if err!=nil{
    log.Fatal(err)
}

valid := VerifyCommitment(
    key,
    data,
    commitment,
)

fmt.Println(valid)

leaf := GenerateLeafHash(
    "file123",
    0,
    len(data),
    data,
)

fmt.Println(leaf)



for _, chunk := range chunks {

    fmt.Println("Chunk:", chunk.ChunkIndex)

    fmt.Println("Commitment:", chunk.CommitmentHash)

    fmt.Println("Leaf:", chunk.MerkleLeafHash)

    fmt.Println()
}

}