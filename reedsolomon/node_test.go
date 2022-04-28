package reedsolomon

import (
	"testing"

	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeRedundant(t *testing.T) {
	nd, err := NewNode(merkledag.NodeWithData([]byte("1234567890")))
	require.NoError(t, err)

	r := merkledag.NewRawNode([]byte("12345"))
	nd.AddRedundantNode(r)
	assert.Len(t, nd.RecoveryLinks(), 1)
	nd.RemoveRedundantNode(r.Cid())
	assert.Len(t, nd.RecoveryLinks(), 0)
}

func TestNodeMarshalUnmarshal(t *testing.T) {
	pdata := []byte("1234567890")
	in, err := NewNode(merkledag.NodeWithData(pdata))
	require.NoError(t, err)
	assert.Equal(t, pdata, in.Data())

	red := merkledag.NewRawNode([]byte("12345"))
	in.AddRedundantNode(red)

	data, err := MarshalNode(in)
	require.NoError(t, err)

	out, err := UnmarshalNode(data)
	require.NoError(t, err)

	out.SetCidBuilder(in.CidBuilder())
	assert.True(t, in.Cid().Equals(out.Cid()))
	assert.Equal(t, in.RecoveryLinks(), out.RecoveryLinks())
	assert.Equal(t, pdata, out.Data())
}

func TestNodeDecode(t *testing.T) {
	in, err := NewNode(merkledag.NodeWithData([]byte("1234567890")))
	require.NoError(t, err)
	require.Equal(t, Codec, in.Cid().Type())

	out, err := format.Decode(in)
	require.NoError(t, err)
	assert.True(t, in.Cid().Equals(out.Cid()))
}

func TestNodeCache(t *testing.T) {
	nd, err := NewNode(merkledag.NodeWithData([]byte("1234567890")))
	require.NoError(t, err)

	nd.AddRedundantNode(merkledag.NewRawNode([]byte("12345")))
	assert.Zero(t, nd.cid)
	assert.Nil(t, nd.cache)

	nd.Cid()
	assert.NotZero(t, nd.cid)
	assert.NotNil(t, nd.cache)

	nd.AddNodeLink("", merkledag.NewRawNode([]byte("12345")))
	assert.Nil(t, nd.cache)
}

func TestNode_Copy(t *testing.T) {
	nd, err := NewNode(merkledag.NodeWithData([]byte("1234567890")))
	require.NoError(t, err)

	r := merkledag.NewRawNode([]byte("12345"))
	nd.AddRedundantNode(r)

	cp := nd.Copy()
	nd.SetData([]byte{})
	nd.RemoveRedundantNode(r.Cid())

	assert.NotNil(t, cp.(*Node).RecoveryLinks())
	assert.NotNil(t, cp.(*Node).Data())
}
