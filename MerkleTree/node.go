package MerkleTree

import (
	"crypto/sha256"
	"fmt"
)
type MerkleRoot struct {
	Root      TreeNode
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

func Construct(leafs []*LeafNode, size int) *TreeNode {
	//odd := false
	// Adds empty block if odd
	root := TreeNode{}
	if size % 2 != 0 {
		//odd = true
		node := LeafNode{}
		node.Key = ""
		node.Hash = ""
		leafs = append(leafs, &node)
	}

	for i := 0; i < len(leafs); i++{
		fmt.Println(i, leafs[i].Key)
	}
	return &root
}
