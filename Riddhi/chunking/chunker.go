package main

import (
 //"encoding/hex"
 "fmt"
 "io"
 "os"
 "sync"

// "golang.org/x/crypto/sha3"
)

// ChunkFile splits a file into smaller chunks and returns metadata for each chunk.
// It reads the file sequentially and chunks it based on the specified chunk size.
func (c *DefaultFileChunker) ChunkFile(filePath string) ([]ChunkMeta, error) {
 var chunks []ChunkMeta // Store metadata for each chunk

 // Open the file for reading
 file, err := os.Open(filePath)
 if err != nil {
  return nil, err
 }
 defer file.Close()

 //Creation of fileid
 fileID := GenerateHash([]byte(filePath))

 // Create a buffer to hold the chunk data
 buffer := make([]byte, c.chunkSize)
 index := 0 // Initialize chunk index

 // Loop until EOF is reached
 for {
  // Read chunkSize bytes from the file into the buffer
  bytesRead, err := file.Read(buffer)
  if err != nil && err != io.EOF {
   return nil, err
  }
  if bytesRead == 0 {
   break // If bytesRead is 0, it means EOF is reached
  }

//   // Generate a unique hash for the chunk data
//   hash :=sha3.NewLegacyKeccak256()
//   hash.Write(buffer[:bytesRead]) //don't use hash.Write(buffer) since buffer always have fullsize
//   result:=hash.Sum(nil)
//   hashString := hex.EncodeToString(result[:])

  // Construct the chunk file name
  chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

  // Create a new chunk file and write the buffer data to it
  chunkFile, err := os.Create(chunkFileName)
  if err != nil {
   return nil, err
  }

  chunkCopy := make([]byte, bytesRead)
 copy(chunkCopy, buffer[:bytesRead])

  _, err = chunkFile.Write(chunkCopy)
  if err != nil {
   return nil, err
  }

  // Close the chunk file
  chunkFile.Close()

  //commitment engine
  key, commitment, err := GenerateCommitment(chunkCopy)
  if err != nil {
    return nil, err
   }

  //generate Leaf

  leaf := GenerateLeafHash(
    fileID,
    index,
    bytesRead,
    chunkCopy,
)

//Append metadata
  chunks = append(chunks, ChunkMeta{
    FileName:       chunkFileName,
    ChunkIndex:          index,
    ChunkSize:      bytesRead,
    RandomKey:      key,
    CommitmentHash: commitment,
    MerkleLeafHash: leaf,
})

  // Move to the next chunk
  index++
 }

 return chunks, nil
}


//for future
// ChunklargeFile splits a large file into smaller chunks in parallel and returns metadata for each chunk.
// It divides the file into chunks and processes them concurrently using multiple goroutines.
func (c *DefaultFileChunker) ChunklargeFile(filePath string) ([]ChunkMeta, error) {
 var wg sync.WaitGroup
 var mu sync.Mutex
 var chunks []ChunkMeta // Store metadata for each chunk

 // Open the file for reading
 file, err := os.Open(filePath)
 if err != nil {
  return nil, err
 }
 defer file.Close()

 // Get file information to determine the number of chunks
 fileInfo, err := file.Stat()
 if err != nil {
  return nil, err
 }

 numChunks := int(fileInfo.Size() / int64(c.chunkSize))
 if fileInfo.Size()%int64(c.chunkSize) != 0 {
  numChunks++
 }

 // Create channels to communicate between goroutines
 errChan := make(chan error, numChunks)
 indexChan := make(chan int, numChunks)

 // Populate the index channel with chunk indices
 for i := 0; i < numChunks; i++ {
  indexChan <- i
 }
 close(indexChan)

 // Start multiple goroutines to process chunks in parallel
 for i := 0; i < 4; i++ { // Number of parallel workers
  wg.Add(1)
  go func() {
   defer wg.Done()
   for index := range indexChan {
    // Calculate the offset for the current chunk
    offset := int64(index) * int64(c.chunkSize)
    buffer := make([]byte, c.chunkSize) // Create a buffer for chunk data

    // Read chunkSize bytes from the file into the buffer
    bytesRead, err := file.ReadAt(buffer,offset)
    if err != nil && err != io.EOF {
     errChan <- err
     return
    }

    // If bytesRead is 0, it means EOF is reached
    if bytesRead > 0 {
    //  // Generate a unique hash for the chunk data
    //  hash := sha3.NewLegacyKeccak256()
    //  hash.Write(buffer[:bytesRead])
    //  result := hash.Sum(nil)
    //  hashString := hex.EncodeToString(result[:])

     // Construct the chunk file name
     chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

     // Create a new chunk file and write the buffer data to it
     chunkFile, err := os.Create(chunkFileName)
     if err != nil {
      errChan <- err
      return
     }

    chunkCopy := make([]byte, bytesRead)
    copy(chunkCopy, buffer[:bytesRead])

     _, err = chunkFile.Write(chunkCopy)
     if err != nil {
      errChan <- err
      return
     }

     // Append metadata for the chunk to the chunks slice
     chunk := ChunkMeta{
      FileName: chunkFileName,ChunkIndex: index,ChunkSize: bytesRead,
     }
     mu.Lock()
     chunks = append(chunks, chunk)
     mu.Unlock()

     // Close the chunk file
     defer chunkFile.Close()

    }
   }
  }()
 }

 // Wait for all goroutines to finish
 go func() {
  wg.Wait()
  
  close(errChan)
 }()

 // Check for errors from goroutines
 for err := range errChan {
  if err != nil {
   return nil, err
  }
 }

 return chunks, nil
}