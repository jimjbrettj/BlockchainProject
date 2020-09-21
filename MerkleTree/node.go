package MerkleTree

import "crypto/sha256"

type MerkleTree struct {
	Root *TreeNode
}

// Struct for tree nodes
//type Node struct {
//	Key       string
//	Hash      string
//	Left      *Node
//	Right     *Node
//	LeftEdge  string
//	RightEdge string
//}

type TreeNode struct {
	Hash      string
	//Left      *Node
	//Right     *Node
	LeftEdge  string
	RightEdge string
}

type LeafNode struct {
	Key       string
	Hash      string
}

func CreateLeafNode(key string) *LeafNode {
	node := LeafNode{}
	data := []byte(key)
	hash := sha256.Sum256(data)
	node.Key = key
	node.Hash = string(hash[:])
	return &node
}