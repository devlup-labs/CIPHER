package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cipher/internal/content/core"
	"cipher/internal/content/crypto"
	"cipher/internal/content/engine"
	"cipher/internal/content/manifest"
	"cipher/internal/content/storage"
	"cipher/internal/content/verifier"
)

func main() {
	ingestFile := flag.String("ingest", "", "File to ingest and chunk")
	reassembleOut := flag.String("out", "", "Output file for reassembled data (requires -manifest)")
	manifestFile := flag.String("manifest", "test_files/manifest.json", "Manifest JSON file (output for ingest, input for reassemble)")
	flag.Parse()

	if *ingestFile == "" && *reassembleOut == "" {
		log.Fatal("Must specify either -ingest <file> or -out <file>")
	}

	// Initialize Content Engine components
	config := core.EngineConfig{ChunkSize: 256 * 1024} // 256KB chunks for manual testing
	enc := crypto.NewChaCha20Encryptor()
	dig := verifier.NewSHA256Digest()
	keys := engine.NewLocalKeyProvider()

	storeDir := "./test_files/content_store"
	if err := storage.NewFSStorage(storeDir); err != nil {
		log.Fatalf("Failed to create store dir: %v", err)
	}
	store := storage.NewFSStore(storeDir)

	eng := engine.NewContentEngine(config, enc, dig, store, store, keys)
	ctx := context.Background()

	if *ingestFile != "" {
		log.Printf("Ingesting %s...", *ingestFile)
		f, err := os.Open(*ingestFile)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer f.Close()

		m, err := eng.Ingest(ctx, f, manifest.TypeFile)
		if err != nil {
			log.Fatalf("Failed to ingest file: %v", err)
		}

		// Save manifest to disk
		manifestData, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			log.Fatalf("Failed to serialize manifest: %v", err)
		}
		if err := os.WriteFile(*manifestFile, manifestData, 0644); err != nil {
			log.Fatalf("Failed to write manifest: %v", err)
		}

		log.Printf("Ingest complete!")
		log.Printf("Manifest saved to: %s", *manifestFile)
		log.Printf("Chunks saved to  : %s/", storeDir)
		log.Printf("Total chunks     : %d", len(m.ChunkIDs))

		// For manual testing, we persist the generated content key so it can be reassembled later
		// (Normally this would be handled securely or retrieved over network)
		key, _ := keys.Get(ctx, m.Descriptor.ID)
		keyPath := filepath.Join(storeDir, fmt.Sprintf("%x.key", m.Descriptor.ID))
		os.WriteFile(keyPath, key, 0600)
		log.Printf("Test content key saved to: %s", keyPath)
	}

	if *reassembleOut != "" {
		log.Printf("Reassembling from manifest %s...", *manifestFile)

		manifestData, err := os.ReadFile(*manifestFile)
		if err != nil {
			log.Fatalf("Failed to read manifest: %v", err)
		}

		var m manifest.Manifest
		if err := json.Unmarshal(manifestData, &m); err != nil {
			log.Fatalf("Failed to parse manifest: %v", err)
		}

		// Load the test content key back into the key provider
		keyPath := filepath.Join(storeDir, fmt.Sprintf("%x.key", m.Descriptor.ID))
		key, err := os.ReadFile(keyPath)
		if err == nil {
			keys.Put(ctx, m.Descriptor.ID, key)
		} else {
			log.Printf("Warning: Could not load test key from %s: %v", keyPath, err)
		}

		out, err := os.Create(*reassembleOut)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer out.Close()

		if err := eng.Reassemble(ctx, &m, out); err != nil {
			log.Fatalf("Failed to reassemble: %v", err)
		}

		log.Printf("Reassembly complete! File written to: %s", *reassembleOut)
	}
}
