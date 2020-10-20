package main

import (
	"CSE297/BlockchainProject/MerkleTree"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

package main

import (
"./MerkleTree"
//"CSE297/BlockchainProject/MerkleTree"
"bufio"
"encoding/hex"

//"bytes"
"crypto/sha256"
"fmt"
"log"
//"math/big"
"math/rand"
"os"
"sort"
"strconv"
"strings"
"time"
)

func readFile(filename string) ([]string, error) {
	// Open and read file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(scanner.Text()) > 100 {
			log.Fatal("String input to large")
		}
	}

	return lines, scanner.Err()
}

func writeArray(tree []string, file *os.File) bool {
	// Write to the file
	_, err := file.WriteString("Starting to print the tree\n")
	// Make sure write was successful
	if err != nil {
		log.Fatal(err)
	}
	// Begin printing out the contents of the array
	for i := 0; i < len(tree); i++ {
		_, err2 := file.WriteString(tree[i] + "\n")
		// Make sure write was successful
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	return false
}

func print(root interface{}) {
	switch root.(type) {
	case *MerkleTree.LeafNode:
		fmt.Println("Leaf: ", root.(*MerkleTree.LeafNode).Key)
	case *MerkleTree.TreeNode:
		fmt.Println("TreeHash: ", root.(*MerkleTree.TreeNode).Hash)
		fmt.Println("TreeLeft: ", root.(*MerkleTree.TreeNode).LeftEdge)
		fmt.Println("TreeRight: ", root.(*MerkleTree.TreeNode).RightEdge)
		fmt.Println()
		print(root.(*MerkleTree.TreeNode).Left)
		print(root.(*MerkleTree.TreeNode).Right)
	}
}

func testTrie() {
	trie := MerkleTree.CreateTestTrie()
	file, err := os.Create("TrieTest.txt")
	// Make sure creation completed
	if err != nil {
		log.Fatal(err)
	}
	writeTree(trie, file)
	file.Close()
}

func main() {
	var first *string
	var lastBlock *MerkleTree.Block
	for {
		fmt.Print("Enter the filename or 'done' to exit: ")
		var filename string
		_, _ = fmt.Scanln(&filename)
		if filename == "done" {
			break
		}
		if first == nil {
			first = &filename
		}
		fmt.Println("Filename is: " + filename)

		lines, err := readFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		sort.Strings(lines)

		trie := MerkleTree.CreateTrie()
		block := MerkleTree.CreateBlock()

		////// Set block difficulty //////
		/// Commented out not necessary for now. Can do to do better but more
		// important things to do first
		//bigInt := big.NewInt(int64(1))
		//bigInt = bigInt.Lsh(bigInt, 106)
		//bigString := []byte(bigInt.String())
		//
		//for i := range bigString {
		//	bigString[i] = ^bigString[i]
		//}

		block.Difficulty = byte(128)
		MerkleTree.Construct(lines, trie)
		block.Tree = trie
		block.TreeHeadHash = trie.Root.Hash
		for {
			guess := rand.Intn(256)
			guessString := block.TreeHeadHash + string(guess)
			hash := sha256.Sum256([]byte(guessString))
			guessHash := hash[0]
			res := guessHash <= block.Difficulty
			if res {
				fmt.Println("Nonce guess valid")
				block.Nonce = guess
				break
			}
			fmt.Println("Nonce guess not valid")
		}
		block.TimeStamp = uint64(time.Now().Unix())
		if lastBlock == nil {
			block.PreviousHash = "0"
			block.Previous = nil
		} else {
			block.PreviousHash = lastBlock.TreeHeadHash
			block.Previous = lastBlock
		}
		lastBlock = block
	}
	if first == nil {
		return
	}

	testTrie()

	// split the file name to adhere to output format
	splitFile := strings.Split(*first, ".")
	outFile := splitFile[0] + ".block.out"
	// Create the file to write to
	file, err := os.Create(outFile)
	// Make sure creation completed
	if err != nil {
		log.Fatal(err)
	}
	// Set the file to close when finished
	for {
		currentBlock := lastBlock
		if lastBlock == nil {
			break
		}

		printBlock(file, currentBlock)
		if lastBlock.Previous != nil {
			lastBlock = lastBlock.Previous
			continue
		} else {
			break
		}

	}
	defer file.Close()

}

func printBlock(file *os.File, block *MerkleTree.Block) {
	file.WriteString("BEGIN BLOCK\n")
	file.WriteString("BEGIN HEADER\n")
	file.WriteString("PrevHash: " + block.PreviousHash + "\n")
	file.WriteString("RootHash: " + block.TreeHeadHash + "\n")
	file.WriteString("Time: " + strconv.FormatUint(block.TimeStamp, 10) + "\n")
	file.WriteString("Target: " + strconv.Itoa(int(block.Difficulty)) + "\n")
	file.WriteString("Nonce: " + strconv.FormatUint(uint64(block.Nonce), 10) + "\n")
	file.WriteString("END HEADER\n")
	writeTree(block.Tree, file)
	file.WriteString("END BLOCK\n\n")

}

/**
Returns 0 for null, 1 for TreeNode, 2 for LeafNode
*/
func getType(tree MerkleTree.TreeNode, left bool) int {
	if left {
		switch tree.Left.(type) {
		case MerkleTree.LeafNode:
			return 2
		case MerkleTree.TreeNode:
			return 1
		}
		return 0
	} else {
		switch tree.Right.(type) {
		case MerkleTree.LeafNode:
			return 2
		case MerkleTree.TreeNode:
			return 1
		}
		return 0
	}
}

func height(tree MerkleTree.TreeNode) int {

	var leftHeight = 1
	if getType(tree, true) == 1 {
		leftHeight = height(tree.Left.(MerkleTree.TreeNode))
	}

	var rightHeight = 1
	if getType(tree, false) == 1 {
		rightHeight = height(tree.Right.(MerkleTree.TreeNode))
	}

	if leftHeight > rightHeight {
		return leftHeight + 1
	} else {
		return rightHeight + 1
	}

}

func printNode(node MerkleTree.TreeNode, file *os.File) {
	file.WriteString("Node Id: " + strconv.Itoa(node.PrintID) + "\n")
	file.WriteString("Left ID " + strconv.Itoa(2*node.PrintID) + "\n")
	file.WriteString("Left edge " + node.LeftEdge + "\n")
	file.WriteString("Hash: " + hex.EncodeToString([]byte(node.Hash)) + "\n")
	file.WriteString("Right edge " + node.RightEdge + "\n")
	file.WriteString("Right ID" + strconv.Itoa(2*node.PrintID+1) + "\n")
	file.WriteString("\n\n")
}

func printLeaf(node MerkleTree.LeafNode, file *os.File) {
	file.WriteString("Leaf key " + node.Key + "\n")
	file.WriteString("Leaf hash " + hex.EncodeToString([]byte(node.Hash)) + "\n")
}


func generatePrintIDandHash(nodeIn *MerkleTree.TreeNode) {
	var node = *nodeIn
	var leftHash *string = nil
	var rightHash *string = nil
	if getType(node, true) == 1 {
		var left = node.Left.(MerkleTree.TreeNode)
		left.PrintID = 2 * node.PrintID
		generatePrintIDandHash(&left)
		node.Left = left
		leftHash = &left.Hash
	} else if getType(node, true) == 2 {
		var left = node.Left.(MerkleTree.LeafNode)
		leftHash = &left.Hash
	}

	if getType(node, false) == 1 {
		var right = node.Right.(MerkleTree.TreeNode)
		right.PrintID = 2*node.PrintID + 1
		generatePrintIDandHash(&right)
		node.Right = right
		rightHash = &right.Hash
	} else if getType(node, false) == 2 {
		var right = node.Right.(MerkleTree.LeafNode)
		rightHash = &right.Hash
	}
	if leftHash != nil && rightHash != nil {
		node.Hash = MerkleTree.Hash(*leftHash, *rightHash)
	} else if leftHash != nil {
		node.Hash = *leftHash
	} else if rightHash != nil {
		node.Hash = *rightHash
	}
	//Hopefully never doesn't meet one of these

	*nodeIn = node
}
func writeTree(tree *MerkleTree.Trie, file *os.File) {
	var root = *tree.Root
	root.PrintID = 1
	generatePrintIDandHash(&root)
	*tree.Root = root

	//var left = root.Left.(MerkleTree.TreeNode)
	//println(left.PrintID)
	var h = height(root)
	var i = 1
	for i <= h {
		printGivenLevel(*tree.Root, i, file)
		i++
	}

}
func printGivenLevel(tree MerkleTree.TreeNode, level int, file *os.File) {
	if level == 1 {
		printNode(tree, file)
	} else if level > 1 {
		if getType(tree, true) == 1 {
			printGivenLevel(tree.Left.(MerkleTree.TreeNode), level-1, file)
		} else if tree.Left != nil {
			fmt.Println("Got to left leaf")
			printLeaf(tree.Left.(MerkleTree.LeafNode), file)
		}

		if getType(tree, false) == 1 {
			printGivenLevel(tree.Right.(MerkleTree.TreeNode), level-1, file)
		} else if tree.Right != nil {
			fmt.Println("Got to right leaf")
			printLeaf(tree.Right.(MerkleTree.LeafNode), file)
		}
	}
}

package MerkleTree

type Block struct {
	Previous     *Block
	PreviousHash string
	TreeHeadHash string
	TimeStamp    uint64
	Difficulty   byte
	Nonce        int
	Tree         *Trie
}

func CreateBlock() *Block {
	Block := Block{}
	return &Block
}

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

func CreateLeafNode(key string) LeafNode {
	node := LeafNode{}
	data := []byte(key)
	hash := sha256.Sum256(data)
	node.Key = key
	node.Hash = string(hash[:])
	return node
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
	//leaf1.Hash = Hash("Jimbo", "")

	leaf2 := LeafNode{}
	leaf2.Key ="Ollie"
	//leaf2.Hash = Hash("Ollie", "")

	leaf3 := LeafNode{}
	leaf3.Key ="Mitch"
	//leaf3.Hash = Hash("Mitch", "")

	leaf4 := LeafNode{}
	leaf4.Key ="Kess"
	//leaf4.Hash = Hash("Kess", "")

	// Define tree nodes
	tree1 := TreeNode{}
	//tree1.Hash = Hash(leaf1.Hash, leaf2.Hash)
	tree1.Left = leaf1
	tree1.Right = leaf2

	tree2 := TreeNode{}
	//tree2.Hash = Hash(leaf3.Hash, leaf4.Hash)
	tree2.Left = leaf3
	tree2.Right = leaf4

	// Root
	root := TreeNode{}
	//root.Hash = Hash(tree1.Hash, tree2.Hash)
	root.Left = tree1
	root.Right = tree2
	//root.PrintID = 1
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
	//node.PrintID = 1
	Trie.Root = &node
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
			fmt.Println("LEft null")
			PrintLeft(n.Left)
		}

	}
	return
}

// Constructs the trie. Ideally will return the root or the trie itself
func Construct(values []string, trie *Trie) {
	for i := range values {
		println(values[i])
		Insert(trie.Root, values[i], len(values[i]), trie)
		fmt.Println("Root left: ", trie.Root.Left)
	}
	//fmt.Println("Root left: ", trie.Root.Left.(LeafNode).Key)
	//fmt.Println("Root right: ", trie.Root.Right.(LeafNode).Key)
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
	case LeafNode:
		return n.Hash
	case TreeNode:
		return n.Hash
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
	case LeafNode:
		return TreeNodeFromLeaves(node.(LeafNode), newLeafNode, index)
	case TreeNode:
		oldNodeEdge := edge[index:]
		newLeafEdge := key[index+prefix:]
		hash := ""
		return &TreeNode{hash, node, newLeafNode, oldNodeEdge, newLeafEdge, 0}
	}

	return nil
}

// Currently not hashing as I go. Either need to add that logic or hash at the end. Honestly might be best to do at end.
func Insert(node *TreeNode, key string, prefix int, trie *Trie) {
	if trie.Root.Left == nil {
		fmt.Println("Left is nil insert")
		newNode := LeafNode{}
		newNode.Key = key
		newNode.Hash = Hash(key, "")
		trie.Root.Left = newNode
		trie.Root.LeftEdge = key
		return
	}
	index := findIndex(key[prefix:], node.LeftEdge)
	fmt.Println("Index: ", index)
	if index > 0 { // if index matches on left
		switch m := node.Left.(type) {
		case *LeafNode:
			fmt.Println("*LEAF")
			if m.Key == key {
				return
			}
			// TODO make insert helper
			node.Left = insertHelper(m, node.LeftEdge, key, index, prefix)
			node.LeftEdge = node.LeftEdge[:index]
			return
		case *TreeNode:
			fmt.Println("*TREE")
			if index == len(node.LeftEdge) {
				Insert(m, key, prefix+index, trie)
				return
			}
			// TODO make insert helper
			node.Left = insertHelper(m, node.LeftEdge, key, index, prefix)
			node.LeftEdge = node.LeftEdge[:index]
			return
		case LeafNode:
			fmt.Println("LEAF")
			if m.Key == key {
				return
			}
			// TODO make insert helper
			node.Left = insertHelper(m, node.LeftEdge, key, index, prefix)
			node.LeftEdge = node.LeftEdge[:index]
			return
		case TreeNode:
			fmt.Println("TREE")
			if index == len(node.LeftEdge) {
				Insert(&m, key, prefix+index, trie)
				return
			}
			// TODO make insert helper
			node.Left = insertHelper(m, node.LeftEdge, key, index, prefix)
			node.LeftEdge = node.LeftEdge[:index]
			return

		}
	}

	if trie.Root.Right == nil {
		fmt.Println("Right is nil insert", key)
		newNode := LeafNode{}
		newNode.Key = key
		newNode.Hash = Hash(key, "")
		trie.Root.Right = newNode
		trie.Root.RightEdge = key
		return
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
				Insert(m, key, prefix+index, trie)
				return
			}
			// TODO make insert helper
			node.Right = insertHelper(m, node.RightEdge, key, index, prefix)
			node.RightEdge = node.RightEdge[:index]
			return
		case LeafNode:
			if m.Key == key {
				return
			}
			// TODO make insert helper
			node.Right = insertHelper(m, node.RightEdge, key, index, prefix)
			node.RightEdge = node.RightEdge[:index]
			return
		case TreeNode:
			if index == len(node.RightEdge) {
				Insert(&m, key, prefix+index, trie)
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
				Insert(m, key, prefix, trie)
				return
			}
			goto ExtensionNode
		case LeafNode:
			if m.Key == "" {
				node.Right = CreateLeafNode(key)
				node.RightEdge = key
				return
			}
			goto ExtensionNode
		case TreeNode:
			if node.RightEdge == "" {
				Insert(&m, key, prefix, trie)
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
	return ValidateBlock(block)
}

// Checks that the current block is valid
func ValidateBlock(block *Block) bool {
	if block == nil {
		return true
	}
	validTrie := ValidateTrie(block.Tree)
	// TODO: Validate hash of previous block
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
