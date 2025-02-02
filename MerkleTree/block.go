package MerkleTree

type Block struct {
	Previous     *Block
	PreviousHash string
	TreeHeadHash string
	TimeStamp    uint64
	Difficulty   byte
	Nonce        int
	Tree         *Trie
}

func CreateBlock() *Block {
	Block := Block{}
	return &Block
}
