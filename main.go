package main

import (
	"./MerkleTree"
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

func writeTree(root *MerkleTree.TreeNode, file *os.File) bool {
	// Write to the file
	_, err2 := file.WriteString("Hello GoLang")
	// Make sure write was successful
	if err2 != nil {
		log.Fatal(err2)
	}

	// Return
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
	//chain := MerkleTree.CreateChain()
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
			if res  {
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

		//chain.Root = block
		//chain = chain.Next
	}
	if first == nil {
		return
	}
	//chain = chain.Previous //Rewind due to pre creating and linking to the next node above

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
	file.WriteString("Target: " + string(int(block.Difficulty)) + "\n")
	file.WriteString("Nonce: " + strconv.FormatUint(uint64(block.Nonce), 10) + "\n")
	file.WriteString("END HEADER\n")
	file.WriteString("END BLOCK\n\n")

}
