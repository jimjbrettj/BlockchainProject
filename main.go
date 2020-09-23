package main

import (
  "CSE297/BlockchainProject/MerkleTree"
  "bufio"
  "fmt"
  "log"
  "os"
  "sort"
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
      fmt.Println("Tree: ", root.(*MerkleTree.TreeNode).Hash)
      if root.(*MerkleTree.TreeNode).Hash == root.(*MerkleTree.TreeNode).Left.(*MerkleTree.LeafNode).Hash {
        fmt.Println("Hashes are equal")
      }
      print(root.(*MerkleTree.TreeNode).Left)
      print(root.(*MerkleTree.TreeNode).Right)
  }

}

func main() {
  fmt.Print("Enter the filename: ")
  var filename string
  fmt.Scanln(&filename)
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
  MerkleTree.Construct(leafs, trie)
  print(trie.Root)

  //// split the file name to adhere to output format
  //splitFile := strings.Split(filename, ".")
  //outFile := splitFile[0] + ".out.txt"
  //// Create the file to write to
  //file, err := os.Create(outFile)
  //// Make sure creation completed
  //if err != nil {
  //  log.Fatal(err)
  //}
  //// Set the file to close when finished
  //defer file.Close()
  //
  //array := []string{"first", "second", "third", "fourth", "fifth"}
  //printed := writeArray(array, file)
  //
  //fmt.Println("Success? ", printed)
}