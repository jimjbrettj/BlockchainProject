package MerkleTree

import (
	"crypto/sha256"
	"fmt"
)
type MerkleTree struct {
	Root      []interface{}
}

type TreeNode struct {
	Hash      string
	Left      interface{}
	Right     interface{}
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

func hash(l string, r string) string {
	s := l + r
	data := []byte(s)
	hash := sha256.Sum256(data)
	return string(hash[:])
}

func InitTree() *MerkleTree {
	MerkleTree := MerkleTree{}
	return &MerkleTree
}


func CreateMerkleTree(leafs []*LeafNode, tree *MerkleTree) *TreeNode{
	// 1) Add all leaf nodes to root[0]
	root := tree.Root[0].([]*LeafNode)
	//tree.Root[0] = leafs
	for i := 0; i < len(leafs); i++ {
		root[i] = leafs[i]
		fmt.Println(root[i])
	}

	// 2) While root[0].length > 1
		// Create temp array of Nodes

		// For i = 0; i < root[0].length; i += 2
			// if (i < root[0].length - 1 && i % 2 == 0)
				// Create new node with hash of i and i+1 nodes
			// else
				// Create and push node at i
			// Push temp array to front of root

	// return the root of the tree
	node := TreeNode{}
	return &node
}
