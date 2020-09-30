package MerkleTree

type Chain struct {
	Block    *Block
	Next     *Chain
	Previous *Chain
}

func CreateChain() *Chain {
	Chain := Chain{}
	return &Chain
}
