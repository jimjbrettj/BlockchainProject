package MerkleTree

import (
	"crypto/sha256"
)

type Trie struct {
	Root *TreeNode
}

type TreeNode struct {
	Hash      string
	Left      interface{}
	Right     interface{}
	LeftEdge  string
	RightEdge string
	PrintID   int
}

type LeafNode struct {
	Key  string
	Hash string
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

func CreateTestTrie() *Trie {
	Trie := Trie{}

	// Define Leaves
	leaf1 := LeafNode{}
	leaf1.Key ="Jimbo"
	leaf1.Hash = Hash("Jimbo", "")

	leaf2 := LeafNode{}
	leaf2.Key ="Ollie"
	leaf2.Hash = Hash("Ollie", "")

	leaf3 := LeafNode{}
	leaf3.Key ="Mitch"
	leaf3.Hash = Hash("Mitch", "")

	leaf4 := LeafNode{}
	leaf4.Key ="Kess"
	leaf4.Hash = Hash("Kess", "")

	// Define tree nodes
	tree1 := TreeNode{}
	tree1.Hash = Hash(leaf1.Hash, leaf2.Hash)
	tree1.Left = leaf1
	tree1.Right = leaf2

	tree2 := TreeNode{}
	tree2.Hash = Hash(leaf3.Hash, leaf4.Hash)
	tree2.Left = leaf3
	tree2.Right = leaf4

	// Root
	root := TreeNode{}
	root.Hash = Hash(tree1.Hash, tree2.Hash)
	root.Left = tree1
	root.Right = tree2
	Trie.Root = &root

	return &Trie
}

// Init empty trie
func CreateTrie() *Trie {
	Trie := Trie{}
	node := TreeNode{}
	node.Hash = ""
	node.RightEdge = ""
	node.LeftEdge = ""
	node.Left = nil
	node.Right = nil
	node.PrintID = 1
	Trie.Root = &node
	return &Trie
}

// Constructs the trie. Ideally will return the root or the trie itself
func Construct(values []string, trie *Trie) {
	for i := range values {
		insert(trie.Root, values[i], len(values[i]), trie)
		//println(values[i])
	}
}

// Returns the min of 2 integers
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Returns the index of which to split the edge label
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

func GetNodesHash(node interface{}) string {
	switch node.(type) {
	case *LeafNode:
		return node.(*LeafNode).Key
	case *TreeNode:
		return node.(*TreeNode).Hash
	}
	return ""
}

func Belongs(key string, branch string) bool {
	return key[0] == branch[0]
}

func TreeNodeFromLeaves(leaf1 LeafNode, leaf2 LeafNode, index int) TreeNode {
	treeNode := TreeNode{}
	if leaf1.Key[index] < leaf2.Key[index] {
		treeNode.Left = leaf1
		treeNode.Right = leaf2
		treeNode.LeftEdge = leaf1.Key[index:]
		treeNode.RightEdge = leaf2.Key[index:]
		treeNode.Hash = Hash(leaf1.Key, leaf2.Key)
	} else {
		treeNode.Right = leaf1
		treeNode.Left = leaf2
		treeNode.RightEdge = leaf1.Key[index:]
		treeNode.LeftEdge = leaf2.Key[index:]
		treeNode.Hash = Hash(leaf2.Key, leaf1.Key)
	}
	return treeNode
}

func insertHelper(node interface{}, edge string, key string, index int, prefix int) interface{} {
	// Create node with inserted node as its child
	newLeafNode := LeafNode{}
	newLeafNode.Key = key
	switch node.(type) {
	case *LeafNode:
		return TreeNodeFromLeaves(node.(LeafNode), newLeafNode, index)
	case *TreeNode:
		oldNodeEdge := edge[index:]
		newLeafEdge := key[index+prefix:]
		hash := ""
		return &TreeNode{hash, node, newLeafNode, oldNodeEdge, newLeafEdge, 0}
	}
	return nil
}

// Currently not hashing as I go. Either need to add that logic or hash at the end. Honestly might be best to do at end.
func insert(node *TreeNode, key string, prefix int, trie *Trie) {
	index := findIndex(key[prefix:], node.LeftEdge)
	if index > 0 { // if index matches on left
		switch m := node.Left.(type) {
		case *LeafNode:
			if m.Key == key {
				return
			}
			// TODO make insert helper
			node.Left = insertHelper(m, node.LeftEdge, key, index, prefix)
			node.LeftEdge = node.LeftEdge[:index]
			return
		case *TreeNode:
			if index == len(node.LeftEdge) {
				insert(m, key, prefix+index, trie)
				return
			}
			// TODO make insert helper
			node.Left = insertHelper(m, node.LeftEdge, key, index, prefix)
			node.LeftEdge = node.LeftEdge[:index]
			return
		}
	}

	index = findIndex(key[prefix:], node.RightEdge)
	if index > 0 { // if index matches on right
		switch m := node.Right.(type) {
		case *LeafNode:
			if m.Key == key {
				return
			}
			// TODO make insert helper
			node.Right = insertHelper(m, node.RightEdge, key, index, prefix)
			node.RightEdge = node.RightEdge[:index]
			return
		case *TreeNode:
			if index == len(node.RightEdge) {
				insert(m, key, prefix+index, trie)
				return
			}
			// TODO make insert helper
			node.Right = insertHelper(m, node.RightEdge, key, index, prefix)
			node.RightEdge = node.RightEdge[:index]
			return
		}
	} else {
		// No prefix is shared on either side, insert extension code
		switch m := node.Right.(type) {
		case *LeafNode:
			if m.Key == "" {
				node.Right = CreateLeafNode(key)
				node.RightEdge = key
				return
			}
			goto ExtensionNode
		case *TreeNode:
			if node.RightEdge == "" {
				insert(m, key, prefix, trie)
				return
			}
			goto ExtensionNode
		}
	ExtensionNode:
		// TODO make insert helper
		node.Right = insertHelper(node.Right, node.RightEdge, key, index, prefix)
		node.RightEdge = ""
		return
	}
}

// Check the validity of entire Blockchain
func ValidateChain(block *Block) bool {
	// TODO For every block, validate all members

	/*
<<<<<<< HEAD
	Members are:
		Previous     *Block
		PreviousHash string
		TreeHeadHash string
		TimeStamp    uint64
		Difficulty   byte
		Nonce        int
		Tree         *Trie
	 */
	return false
}

// Checks the validity of a trie
func ValidateTrie(trie *Trie) bool {
	// TODO Recursive iteration of trie to validate hashes, if any dont match return false, else true
	return false
}
