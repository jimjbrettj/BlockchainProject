package MerkleTree

import "crypto/sha256"

type MerkleTree struct {
	Root *Node
}

// Struct for tree nodes
type Node struct {
	Key       string
	Hash      string
	Left      *Node
	Right     *Node
	LeftEdge  string
	RightEdge string
}

func CreateMerkleNode(key string, left, right *Node) *Node {
	node := Node{}
	data := []byte(key)

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Key = key
		node.Hash = string(hash[:])
		node.LeftEdge = ""
		node.RightEdge = ""
	} else {
		hash := left.Hash + right.Hash
		node.Key = ""
		node.Hash = hash
		// TODO: handle edges
		node.LeftEdge = "left"
		node.RightEdge = "right"

	}
	node.Left = left
	node.Right = right
	return &node
}