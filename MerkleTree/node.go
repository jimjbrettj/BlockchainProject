package MerkleTree

import (
	"crypto/sha256"
	"fmt"
)

type Trie struct {
	Root      *TreeNode
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

// Concats 2 strings and takes the hash of it
func Hash(l string, r string) string {
	s := l + r
	data := []byte(s)
	hash := sha256.Sum256(data)
	return string(hash[:])
}

// Init empty trie
func CreateTrie() *Trie {
	Trie := Trie{}
	return &Trie
}

func Construct(leafs []*LeafNode, trie *Trie) {
	for i := range leafs {
		Insert(leafs[i], leafs[i].Key, "", trie)
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func findIndex(s1 string, s2 string) int {
	i := 0
	min := min(len(s1), len(s2))
	for i < min {
		if s1[i] == s2[i] {
			i++
		} else {
			return i
		}
	}
	return i
}

func Insert(leaf *LeafNode, key string, prefix string, trie *Trie) {
	// If tree empty, insert at root
	if trie.Root == nil {
		node := TreeNode{}
		node.Hash = leaf.Hash
		node.RightEdge = ""
		node.LeftEdge = key
		node.Left = leaf
		node.Right = nil
		trie.Root = &node
	} else {
		// Check gcs on left
		fmt.Println("s1: ", key)
		fmt.Println("s2: ", trie.Root.LeftEdge)
		leftIndex := findIndex(key, trie.Root.LeftEdge)
		fmt.Println("Index: ", leftIndex)
		// CHeck again on right
		// If both are 0, create extension node. Insert left

		// Check if right is null
		if trie.Root.Right == nil && trie.Root.RightEdge < leaf.Key {
			trie.Root.Right = leaf
			trie.Root.RightEdge = leaf.Key
			switch trie.Root.Left.(type) {
				case *LeafNode:
					trie.Root.Hash = Hash(trie.Root.Left.(*LeafNode).Key, leaf.Key)
				case *TreeNode:
					trie.Root.Hash = Hash(trie.Root.Left.(*TreeNode).Hash, leaf.Key)
			}
		} else {
			// Add logic here
		}
	}
}
