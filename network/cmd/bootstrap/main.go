package main

import (
	"cipher/network/identity"
	"cipher/network/transport"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func main() {
	// flag is used to pass in flags from the command line. Here, we define a flag for the port number.
	port := flag.Int(
		"p",
		4003,
		"Port for the bootstrap node (TCP)",
	)
	wsPort := flag.Int(
		"ws-port",
		0,
		"Port for the bootstrap node (WebSocket, 0 to disable)",
	)

	identityPath := flag.String("identity", "", "Custom path to identity key file (optional)")

	// parse the cmd line args passed
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var priv libp2pcrypto.PrivKey
	var err error
	if *identityPath != "" {
		priv, err = identity.LoadOrCreateFromPath(*identityPath)
	} else {
		priv, err = identity.LoadOrCreate()
	}
	if err != nil {
		log.Fatalf("Failed to load or create identity: %v", err)
	}

	h, kdht, err := transport.NewNode(ctx, *port, *wsPort, priv, "", false)
	if err != nil {
		log.Fatalf(
			"Failed to create bootstrap node: %v",
			err,
		)
	}

	defer h.Close()
	defer kdht.Close()

	log.Printf(
		"Bootstrap node started with Peer ID: %s",
		h.ID(),
	)

	log.Println("Bootstrap node listening on:")

	for _, addr := range h.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr, h.ID())
	}

	log.Println()
	log.Println("Bootstrap node is ready.")
	log.Println("Waiting for peers...")

	// this is just listening for shutdown signals from the OS -- SIGINT is Ctrl+C and SIGTERM is termination req from another program
	ch := make(chan os.Signal, 1)

	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-ch

	log.Println("Shutting down bootstrap node...")
}
