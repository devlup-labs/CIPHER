package main

type ChunkMeta struct {
    FileName   string `json:"file_name"`
    ChunkIndex      int    `json:"index"`
    ChunkSize  int    `json:"chunk_size"`
    RandomKey      string
    CommitmentHash string
    MerkleLeafHash string
   
}

type Commitment struct {
    ChunkIndex int
    CommitmentHash string
}

//use later-for zero-knowledge proof
type Reveal struct {
    ChunkIndex int
    Key string
    ChunkData []byte
}

type DefaultFileChunker struct {
 chunkSize int
}

