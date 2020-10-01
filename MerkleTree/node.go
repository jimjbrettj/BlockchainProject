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

// Constructs the trie. Ideally will return the root or the trie itself
func Construct(leafs []*LeafNode, trie *Trie) {
	for i := range leafs {
		RecursiveInsert(leafs[i], trie)
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
	return key[0] == branch[0];
}

func RecursiveInsert(leaf *LeafNode, trie *Trie) {
	// If root is nil, insert leaf as left child
	if trie.Root == nil {
		// Create TreeNode for root
		node := TreeNode{}
		node.Hash = leaf.Hash
		// Insert leaf node as left child
		node.Left = leaf
		node.LeftEdge = leaf.Key
		// Make nil right child
		node.Right = nil
		node.RightEdge = ""
		// Set root node
		trie.Root = &node
	} else {
		// Otherwise: decide if node belongs to left or right subtree
		if (Belongs(leaf.Key, trie.Root.LeftEdge)) {
			// LeafNode belongs to the left subtree
			index := findIndex(leaf.Key, trie.Root.LeftEdge)
			newNode := InsertHelper(leaf, trie, trie.Root.Left, index, trie.Root.LeftEdge[index:])
			trie.Root.Left = newNode
			trie.Root.LeftEdge = leaf.Key[:index]
			// Set the hash value of the root node
			if trie.Root.Right == nil {
				trie.Root.Hash = newNode.Hash
			} else {
				trie.Root.Hash = Hash(newNode.Hash, GetNodesHash(trie.Root.Right))
			}
		} else {
			// LeafNode belongs to the right subtree
			// If no right branch, insert leaf
			if trie.Root.Right == nil {
				trie.Root.Right = leaf
				trie.Root.RightEdge = leaf.Key
				// Set the hash value of the root node
				trie.Root.Hash = Hash(GetNodesHash(trie.Root.Left), leaf.Key)
			} else {
				index := findIndex(leaf.Key, trie.Root.RightEdge)
				newNode := InsertHelper(leaf, trie, trie.Root.Right, index, trie.Root.RightEdge[index:])
				trie.Root.Right = newNode
				trie.Root.RightEdge = leaf.Key[:index]
				// Set the hash value of the root node
				trie.Root.Hash = Hash(GetNodesHash(trie.Root.Left), newNode.Hash)
			}
		}
	}
}

func InsertHelper(leaf *LeafNode, trie *Trie, currentNode interface{}, index int, postfix string) TreeNode {
	// TODO: Postfix does not work properly.  My plan was for it to be used when a null string has to be used to decide which branch gets extended.  len(postfix) > 0 means left, else right
	fmt.Println(leaf.Key + " " + postfix) // Temp print statement
	// Check if the current node is a leaf or a tree node
	switch currentNode.(type) {
		case *LeafNode:
			return TreeNodeFromLeaves(leaf, currentNode.(*LeafNode), index)
		case *TreeNode:
			subKey := postfix + leaf.Key[index:]
			currTreeNode := currentNode.(TreeNode)
			// Decide if the node belongs to the left branch
			if Belongs(subKey, currTreeNode.LeftEdge) {
				nextIndex := findIndex(subKey, currTreeNode.LeftEdge)
				childNode := InsertHelper(leaf, trie, currTreeNode.Left, nextIndex, currTreeNode.LeftEdge[nextIndex:])
				currTreeNode.Left = childNode
				currTreeNode.LeftEdge = subKey[:nextIndex]
				// Set the hash value of the current node
				if currTreeNode.Right == nil {
					currTreeNode.Hash = childNode.Hash
				} else {
					currTreeNode.Hash = Hash(childNode.Hash, GetNodesHash(currTreeNode.Right))
				}
				return currTreeNode
			} else if currTreeNode.Right == nil {
				// If the right child of the current node is nil, we can insert it
				currTreeNode.Right = leaf
				currTreeNode.RightEdge = subKey
				// Set the hash value of the current node
				currTreeNode.Hash = Hash(GetNodesHash(currTreeNode.Left), leaf.Key)
				return currTreeNode
			} else if len(postfix) > 0 {
				// This doesn't work
				// need to extend the branch
				newCurrentNode := TreeNode{}
				newCurrentNode.Left = currTreeNode
				newCurrentNode.LeftEdge = postfix
				newCurrentNode.Right = leaf
				newCurrentNode.RightEdge = leaf.Key[index:]
				newCurrentNode.Hash = Hash(currTreeNode.Hash, leaf.Key)
				return newCurrentNode
			} else {
				childNode := InsertHelper(leaf, trie, currTreeNode.Right, index, "")
				currTreeNode.Right = childNode
				currTreeNode.RightEdge = ""
				currTreeNode.Hash = Hash(GetNodesHash(currTreeNode.Left), GetNodesHash(currTreeNode.Right))
				return currTreeNode
			}
	}
	return currentNode.(TreeNode)
}

func TreeNodeFromLeaves(leaf1 *LeafNode, leaf2 *LeafNode, index int) TreeNode {
	treeNode := TreeNode{}
	if (leaf1.Key[index] < leaf2.Key[index]) {
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
