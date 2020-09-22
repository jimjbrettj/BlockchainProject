package main

import (
  "CSE297/BlockchainProject/MerkleTree"
  "bufio"
  "fmt"
  "log"
  "os"
  "sort"
  "strings"
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
<<<<<<< HEAD
    _, err2 := file.WriteString(tree[i] + "\n")
=======
		_, err2 := file.WriteString(tree[i] + "\n")
>>>>>>> origin/rkWriteFile
    // Make sure write was successful
    if err2 != nil {
      log.Fatal(err2)
    }
  }
  return false;
}

func writeTree(root *MerkleTree.TreeNode, file *os.File) bool {
  // Write to the file
  _, err2 := file.WriteString("Hello GoLang")
  // Make sure write was successful
  if err2 != nil {
<<<<<<< HEAD
    log.Fatal(err2)
=======
      log.Fatal(err2)
>>>>>>> origin/rkWriteFile
  }

  // Return
  return false
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

  leafs := make([]*MerkleTree.LeafNode, len(lines))
  for i, line := range lines {
    leafs[i] = MerkleTree.CreateLeafNode(line)
  }
  root := MerkleTree.Construct(leafs, len(lines))
  fmt.Println("Root is: ", root)

  // split the file name to adhere to output format
  splitFile := strings.Split(filename, ".")
  outFile := splitFile[0] + ".out.txt"
  // Create the file to write to
  file, err := os.Create(outFile)
  // Make sure creation completed
  if err != nil {
    log.Fatal(err)
    // Set the file to close when finished
    defer file.Close()

    array := []string{"first", "second", "third", "fourth", "fifth"}
    printed := writeArray(array, file)

    fmt.Println("Success? ", printed)
  }
}