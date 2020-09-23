package MerkleTree

import (
	"crypto/sha256"
)
type Trie struct {
	Root      interface{}
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

func Insert(leaf *LeafNode, trie *Trie) {
	//fmt.Println(leaf.Key)
	if trie.Root == nil {
		trie.Root = leaf
		//fmt.Println("Root: ", trie.Root.(*LeafNode).Key)
	} else {
		switch trie.Root.(type) {
			case *LeafNode:
				//fmt.Println("Leaf")
				node := TreeNode{}
				node.Hash = hash(trie.Root.(*LeafNode).Hash, leaf.Hash)
				node.RightEdge = ""
				node.LeftEdge = ""
				if trie.Root.(*LeafNode).Key < leaf.Key {
					node.Left = trie.Root
					node.Right = leaf
				} else {
					node.Right = trie.Root
					node.Left = leaf
				}
				trie.Root = node
				//fmt.Println("Root left: ", trie.Root.(TreeNode).Left.(*LeafNode).Key)
				//fmt.Println("Root right: ", trie.Root.(TreeNode).Right.(*LeafNode).Key)
			case *TreeNode:
				//fmt.Println("Tree")
		}
	}
}
