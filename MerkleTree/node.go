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

// Constructs the trie. Ideally will return the root or the trie itself
func Construct(leafs []*LeafNode, trie *Trie) {
	for i := range leafs {
		Insert(leafs[i], leafs[i].Key, "", trie)
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

// Currently not hashing as I go. Either need to add that logic or hash at the end. Honestly might be best to do at end.
func Insert(leaf *LeafNode, key string, prefix string, trie *Trie) {
	/*
		Insert takes a leaf node as a param and adds it to the trie.
		This method will use recursion by passing different nodes to iterate through the tree.
	    The key param hold the value we wish to insert.
		Prefix is the substring used for creating the edge label and trie is the trie object

		1) Check left edge of root and see if any substring matches. If so we add on left side
		2) Do the same for the right side
		3) If neither side has a matching substring, then create extension node to be added on the right side.
	 */

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
		/*
			This checks left side for common substring and if found, add it. This code snippet can be used as a
			reference to show how I deal with properly changing the edge labels based on the index and how I add a
			node.

			This works for my specific case, prob need to add logic (prob recursion tbh but maybe not) to keep searching
		 */
		index := findIndex(key, trie.Root.LeftEdge)
		if index != 0 {
			prefix = key[:index]
			postfix := key[index:]
			replacement := trie.Root.LeftEdge[:index]

			// Create new tree node
			node := TreeNode{}
			node.Hash = ""
			node.RightEdge = postfix
			node.LeftEdge = prefix
			node.Left = trie.Root.Left
			node.Right = leaf

			trie.Root.Left = node
			trie.Root.LeftEdge = replacement
		}

		/*
		    Check again on right side.
	 		Will need to use recursion to find where to add node. Then can do in a similar way as I did above I believe.
			Code I have hear is temporary. Will need to change to check the index on the right edge like this:
				index := findIndex(key, trie.Root.RightEdge)
			Then use recursion to find where to insert
		 */
		if trie.Root.Right == nil && index == 0 {
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

		// If neither edge matches, create an extension node on the right side.
	}
}
