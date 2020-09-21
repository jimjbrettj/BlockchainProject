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

func Construct(leafs []*LeafNode, size int) *MerkleRoot {
	odd := false
	if size % 2 != 0 {
		odd = true
	}
	for i, leaf := range leafs {
		if i == len(leafs) - 1 && odd {
			fmt.Println("This is the last element")
		}

		fmt.Println(i, leaf.Key)
	}


	root := MerkleRoot{}
	return &root
}