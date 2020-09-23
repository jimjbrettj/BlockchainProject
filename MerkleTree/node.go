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
func hash(l string, r string) string {
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
		Insert(leafs[i], trie)
	}
}

func findPrefix(s1 string, s2 string) string {
	return ""
}

func Insert(leaf *LeafNode, trie *Trie) {
	if trie.Root == nil {
		node := TreeNode{}
		node.Hash = leaf.Hash
		node.RightEdge = ""
		node.LeftEdge = leaf.Key
		node.Left = leaf
		node.Right = nil
		trie.Root = &node
	} else {
		if trie.Root.Right == nil && trie.Root.RightEdge < leaf.Key {
			fmt.Println("Made it here")
			trie.Root.Right = leaf
			trie.Root.RightEdge = leaf.Key
			switch t := trie.Root.Left.(type) {
				default: fmt.Println("Type is: ", t)
				case *LeafNode:
					trie.Root.Hash = hash(trie.Root.Left.(*LeafNode).Hash, leaf.Hash)
				case *TreeNode:
					trie.Root.Hash = hash(trie.Root.Left.(*TreeNode).Hash, leaf.Hash)
			}
		}
	}
}
