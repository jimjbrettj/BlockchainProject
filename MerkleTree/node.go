package MerkleTree

import (
	"crypto/sha256"
	"fmt"
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

// Init empty trie
func CreateTrie() *Trie {
	Trie := Trie{}
	return &Trie
}

func PrintLeft(node interface{}) {
	switch n := node.(type) {
	case *LeafNode:
		fmt.Println(n.Key)
		return
	case *TreeNode:
		fmt.Println(n.Hash)
		if n.Left != nil {
			fmt.Println("Left null")
			PrintLeft(n.Left)
		}
	}
	return
}

// Constructs the trie. Ideally will return the root or the trie itself
func (trie *Trie) Construct(values []string) {
	for i := range values {
		println(values[i])
		trie.Insert(values[i])
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
	switch n := node.(type) {
	case *LeafNode:
		return n.Hash
	case *TreeNode:
		return n.Hash
	default:
		println("DEFAULT CASE getNodeHash")
	}
	return ""
}

func Belongs(key string, branch string) bool {
	return key[0] == branch[0]
}

//func TreeNodeFromLeaves(leaf1 *LeafNode, leaf2 *LeafNode, index int) *TreeNode {
//	treeNode := TreeNode{}
//	if leaf1.Key[index] < leaf2.Key[index] {
//		treeNode.Left = leaf1
//		treeNode.Right = leaf2
//		treeNode.LeftEdge = leaf1.Key[index:]
//		treeNode.RightEdge = leaf2.Key[index:]
//		treeNode.Hash = Hash(leaf1.Key, leaf2.Key)
//	} else {
//		treeNode.Right = leaf1
//		treeNode.Left = leaf2
//		treeNode.RightEdge = leaf1.Key[index:]
//		treeNode.LeftEdge = leaf2.Key[index:]
//		treeNode.Hash = Hash(leaf2.Key, leaf1.Key)
//	}
//	return &treeNode
//}

func insertHelper(node interface{}, edge string, key string, index int, prefix int) *TreeNode {
	// Create node with inserted node as its child
	//newLeafNode := LeafNode{}
	//newLeafNode.Key = key
	//switch node.(type) {
	//case *LeafNode:
	//	return TreeNodeFromLeaves(node.(*LeafNode), &newLeafNode, index)
	//case *TreeNode:
	//	oldNodeEdge := edge[index:]
	//	newLeafEdge := key[index+prefix:]
	//	hash := ""
	//	return &TreeNode{hash, node, newLeafNode, oldNodeEdge, newLeafEdge, 0}
	//default:
	//	println("DEFAULT CASE insertHelper")
	//}
	oldNodeEdge := edge[index:]
	newLeafEdge := key[prefix+index:]
	newLeafNode := CreateLeafNode(key)
	return &TreeNode{"", node, newLeafNode, oldNodeEdge, newLeafEdge, 0}
	//return nil
}

func (trie *Trie) Insert(key string) {
	if trie.Root == nil {
		leaf := CreateLeafNode(key)
		dumb := CreateLeafNode("")
		trie.Root = &TreeNode{
			"",
			leaf,
			dumb,
			key,
			"",
			0,
		}
	} else {
		insert(trie.Root, key, 0)
	}
}

// Currently not hashing as I go. Either need to add that logic or hash at the end. Honestly might be best to do at end.
func insert(root *TreeNode, key string, prefix int) {
	index := findIndex(key[prefix:], root.LeftEdge)
	if index > 0 { // if index matches on left
		switch m := root.Left.(type) {
		case *LeafNode:
			if m.Key == key {
				return
			}
			root.Left = insertHelper(m, root.LeftEdge, key, index, prefix)
			root.LeftEdge = root.LeftEdge[:index]
			return
		case *TreeNode:
			if index == len(root.LeftEdge) {
				insert(m, key, prefix+index)
				return
			}
			root.Left = insertHelper(m, root.LeftEdge, key, index, prefix)
			root.LeftEdge = root.LeftEdge[:index]
			return
		default:
			println("DEFAULT CASE")
		}
	}

	index = findIndex(key[prefix:], root.RightEdge)
	if index > 0 { // if index matches on right
		switch m := root.Right.(type) {
		case *LeafNode:
			if m.Key == key {
				return
			}
			root.Right = insertHelper(m, root.RightEdge, key, index, prefix)
			root.RightEdge = root.RightEdge[:index]
			return
		case *TreeNode:
			if index == len(root.RightEdge) {
				insert(m, key, prefix+index)
				return
			}
			root.Right = insertHelper(m, root.RightEdge, key, index, prefix)
			root.RightEdge = root.RightEdge[:index]
			return
		default:
			println("DEFAULT CASE")
		}
	} else {
		// No prefix is shared on either side, insert extension code
		switch m := root.Right.(type) {
		case *LeafNode:
			if m.Key == "" {
				root.Right = CreateLeafNode(key)
				root.RightEdge = key
				return
			}
			goto ExtensionNode
		case *TreeNode:
			if root.RightEdge == "" {
				insert(m, key, prefix)
				return
			}
			goto ExtensionNode
		default:
			println("DEFAULT CASE")
		}
	ExtensionNode:
		root.Right = insertHelper(root.Right, root.RightEdge, key, index, prefix)
		root.RightEdge = ""
		return
	}
}

// Check the validity of entire Blockchain
func ValidateChain(block *Block) bool {
	return ValidateBlock(block)
}

// Checks that the current block is valid
func ValidateBlock(block *Block) bool {
	if block == nil {
		return true
	}
	validTrie := ValidateTrie(block.Tree)
	validHash := true
	return validTrie && validHash && ValidateBlock(block.Previous)
}

// Checks the validity of a trie
func ValidateTrie(trie *Trie) bool {
	return ValidateNode(trie.Root)
}

// Checks that the current node is valid
func ValidateNode(node interface{}) bool {
	switch n := node.(type) {
	case *LeafNode:
		savedHash := n.Hash
		actualHash := Hash(n.Key, "")
		return savedHash == actualHash;
	case *TreeNode:
		savedHash := n.Hash
		actualHash := Hash(GetNodesHash(n.Left), GetNodesHash(n.Right))
		if savedHash == actualHash {
			return ValidateNode(n.Left) && ValidateNode(n.Right)
		} else {
			return false
		}
	}
	return true
}