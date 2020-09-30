package MerkleTree

type Block struct {
	PreviousHash string
	TreeHeadHash string
	TimeStamp    uint64
	Difficulty   uint64
	Nonce        uint32
	Tree         *Trie
}

func CreateBlock() *Block {
	Block := Block{}
	return &Block
}
