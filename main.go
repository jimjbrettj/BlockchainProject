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
  }

  return lines, scanner.Err()
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

  leafs := make([]*MerkleTree.Node, len(lines))
  for i, line := range lines {
    leafs[i] = MerkleTree.CreateMerkleNode(line, nil, nil)
    fmt.Println(i, line)
  }

}
