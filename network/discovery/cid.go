package discovery

import (
	"cipher/network/content/core"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func contentIDToCID(id core.ContentID) (cid.Cid, error) {
	mh, err := multihash.Encode(id[:], multihash.SHA2_256)

	if err != nil {
		return cid.Undef, err
	}

	return cid.NewCidV1(cid.Raw, mh), nil
}
