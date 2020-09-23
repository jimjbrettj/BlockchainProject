package MerkleTree

import (
	"crypto/sha256"
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
		Insert(leafs[i], trie)
	}
}

func findPrefix(s1 string, s2 string) string {
	return ""
}

func Insert(leaf *LeafNode, trie *Trie) {
	// If tree empty, insert at root
	if trie.Root == nil {
		node := TreeNode{}
		node.Hash = leaf.Hash
		node.RightEdge = ""
		node.LeftEdge = leaf.Key
		node.Left = leaf
		node.Right = nil
		trie.Root = &node
	} else {
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
