package reedsolomon

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"

	"github.com/Yangzhengtang/FEC"
)

var log = logging.Logger("recovery")

// Custom codec for Reed-Solomon recovery Nodes.
const Codec uint64 = 0x700 // random number // TODO Register in IPFS codec table.

func init() {
	// register global decoder
	format.Register(Codec, DecodeNode)

	// register codec
	cid.Codecs["recovery-reedsolomon"] = Codec
	cid.CodecToStr[Codec] = "recovery-reedsolomon"
}

type reedSolomon struct {
	dag format.DAGService
}

// NewEncoder creates new Reed-Solomon Encoder.
func NewEncoder(dag format.DAGService) recovery.Encoder {
	return &reedSolomon{dag: dag}
}

func (rs *reedSolomon) Encode(ctx context.Context, nd format.Node, r recovery.Recoverability) (recovery.Node, error) {
	rd, ok := nd.(recovery.Node)
	if ok {
		return rd, nil
	}

	return Encode(ctx, rs.dag, nd, r)
}
