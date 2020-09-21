package main

import (
  "CSE297/BlockchainProject/MerkleTree"
  "bufio"
  "fmt"
  "log"
  "os"
)


func main() {
  fmt.Print("Enter the filename: ")
  var filename string
  fmt.Scanln(&filename)
  fmt.Println("Filename is: " + filename)

  // Open and read file
  file, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    var line string = scanner.Text()
    node := MerkleTree.CreateMerkleNode(line, nil, nil)
    println(node.Key)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}
