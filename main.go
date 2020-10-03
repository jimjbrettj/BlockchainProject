package main

import (
	"./MerkleTree"
	"bufio"
	"fmt"
	"log"
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
	chain := MerkleTree.CreateChain()
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

		// Create array of constructed lead nodes
		leafs := make([]*MerkleTree.LeafNode, len(lines))
		for i, line := range lines {
			leafs[i] = MerkleTree.CreateLeafNode(line)
		}
		trie := MerkleTree.CreateTrie()
		block := MerkleTree.CreateBlock()
		//block.Difficulty == ??
		//block.Nonce == ??
		MerkleTree.Construct(leafs, trie)
		block.Tree = trie
		block.TreeHeadHash = trie.Root.Hash
		block.TimeStamp = uint64(time.Now().Unix())
		if lastBlock == nil {
			block.PreviousHash = "0"
		} else {
			block.PreviousHash = lastBlock.TreeHeadHash
		}
		lastBlock = block
		chain.Next = MerkleTree.CreateChain()
		chain.Next.Previous = chain
		chain.Block = block
		chain = chain.Next
	}
	if first == nil {
		return
	}
	chain = chain.Previous //Rewind due to pre creating and linking to the next node above

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
		currentBlock := chain.Block
		if currentBlock == nil {
			break
		}

		printBlock(file, currentBlock)
		if chain.Previous != nil {
			chain = chain.Previous
			continue
		} else {
			break
		}

	}
	defer file.Close()

}

func write_tree(trie *Trie) {
	var depth = 1 // two possible children
	var null_flag = false // break the while loop at depth of 0 children
	
	f, err := os.Create("merkle_tree.txt") // init txt
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// this first 'if' should be a for in x < depth, could we pass depth?
	if trie.Root != nil {
		_, err2 := fmt.fprintf(f, trie.Root)
		if err2 != nil {
			log.Fatal(err2)
		}
		if trie.Root.Left {
			_, err3 := fmt.fprintf(f, trie.Left)

			if err3 != nil {
				log.Fatal(err3)
			}
		}
		if trie.Root.Right {
			_, err4 := fmt.fprintf(f, trie.Right)

			if err4 != nil {
				log.Fatal(err4)
			}
		}


		depth++

		// if no child nodes null_flag = true
	}
}

func printBlock(file *os.File, block *MerkleTree.Block) {
	file.WriteString("BEGIN BLOCK\n")
	file.WriteString("BEGIN HEADER\n")
	file.WriteString(block.PreviousHash + "\n")
	file.WriteString(block.TreeHeadHash + "\n")
	file.WriteString(strconv.FormatUint(block.TimeStamp, 10) + "\n")
	file.WriteString(strconv.FormatUint(block.Difficulty, 10) + "\n")
	file.WriteString(strconv.FormatUint(uint64(block.Nonce), 10) + "\n")
	file.WriteString("END HEADER\n")
	// write_tree
	file.WriteString("END BLOCK\n\n")
}
