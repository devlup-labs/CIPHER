package retrieval

import (
	"cipher/network/content/core"
	"cipher/network/content/engine"
	"cipher/network/content/manifest"
	"cipher/network/protocol/chunk"
	"cipher/network/transport"
	"context"
	"fmt"
	"log"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
)

// ResolveManifest takes a content ID, queries the DHT for providers who have that content, and attempts to retrieve the manifest from those providers, it connects to each provider, creates a chunk client, and requests the manifest. If successful, it returns the manifest; otherwise, it returns an error after trying all providers.
func ResolveManifest(
	ctx context.Context,
	id core.ContentID,
	kdht *dht.IpfsDHT,
	t *transport.Transport,
	eng *engine.ContentEngine,
	providers []peer.ID,
) (*manifest.Manifest, error) {

	var lastErr error

	for _, provider := range providers {

		client, err := chunk.NewClient(ctx, t, provider, eng)
		if err != nil {
			log.Printf(
				"[DHT] Failed to create chunk client for %s: %v",
				provider,
				err,
			)
			lastErr = err
			continue
		}

		manifestData, err := client.Resolve(ctx, id)
		client.Close()

		if err != nil {
			log.Printf(
				"[DHT] Failed to resolve manifest from provider %s: %v",
				provider,
				err,
			)
			lastErr = err
			continue
		}

		m, err := manifest.Deserialize(manifestData)
		if err != nil {
			log.Printf(
				"[DHT] Provider %s returned invalid manifest: %v",
				provider,
				err,
			)
			lastErr = err
			continue
		}

		log.Printf(
			"[DHT] Successfully resolved manifest from provider %s",
			provider,
		)

		return m, nil
	}

	return nil, fmt.Errorf(
		"failed to resolve manifest from all providers: %w",
		lastErr,
	)
}
