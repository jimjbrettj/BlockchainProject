package main

import (
	"./MerkleTree"
	"reflect"

	//"CSE297/BlockchainProject/MerkleTree"
	"bufio"
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

func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func height(tree MerkleTree.TreeNode) int {

	var leftHeight = 1
	if IsInstanceOf(&tree.Left, (*MerkleTree.TreeNode)(nil)) {
		leftHeight = height(tree.Left.(MerkleTree.TreeNode))
	}

	var rightHeight = 1
	if IsInstanceOf(&tree.Right, (*MerkleTree.TreeNode)(nil)) {
		rightHeight = height(tree.Right.(MerkleTree.TreeNode))
	}

	if leftHeight > rightHeight {
		return leftHeight + 1
	} else {
		return rightHeight + 1
	}

}

func printNode(node MerkleTree.TreeNode, file *os.File) {
	file.WriteString(strconv.Itoa(node.PrintID) + "\n")
	file.WriteString(strconv.Itoa(2*node.PrintID) + "\n")
	file.WriteString(node.LeftEdge + "\n")
	file.WriteString(node.Hash + "\n")
	file.WriteString(node.RightEdge + "\n")
	file.WriteString(strconv.Itoa(2*node.PrintID+1) + "\n")
	file.WriteString("\n\n")
}

func printLeaf(node MerkleTree.LeafNode, file *os.File) {
	file.WriteString(node.Key + "\n")
	file.WriteString(node.Hash + "\n")
}

func generatePrintID(node MerkleTree.TreeNode) {
	if node.Left != nil {
		var left1 = node.Left.(MerkleTree.TreeNode)
		left1.PrintID = 2 * node.PrintID
	}

	if IsInstanceOf(&node.Left, (*MerkleTree.TreeNode)(nil)) {
		var left = node.Left.(MerkleTree.TreeNode)
		left.PrintID = 2 * node.PrintID
		generatePrintID(left)
	}

	if IsInstanceOf(&node.Right, (*MerkleTree.TreeNode)(nil)) {
		var right = node.Right.(MerkleTree.TreeNode)
		right.PrintID = 2*node.PrintID + 1
		generatePrintID(right)
	}

}
func writeTree(tree *MerkleTree.Trie, file *os.File) {

	var root = *tree.Root
	root.PrintID = 1
	generatePrintID(root)
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
		if IsInstanceOf(&tree.Left, (*MerkleTree.TreeNode)(nil)) {
			printGivenLevel(tree.Left.(MerkleTree.TreeNode), level-1, file)
		} else if tree.Left != nil {
			printLeaf(tree.Left.(MerkleTree.LeafNode), file)
		}

		if IsInstanceOf(&tree.Right, (*MerkleTree.TreeNode)(nil)) {
			printGivenLevel(tree.Right.(MerkleTree.TreeNode), level-1, file)
		} else if tree.Right != nil {
			printLeaf(tree.Right.(MerkleTree.LeafNode), file)
		}

	}
}
