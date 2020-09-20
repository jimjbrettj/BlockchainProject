package main

import (
  "fmt"
  "os"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Println("Error, no filename given")
    os.Exit(1)
  }
  arg := os.Args[1]
  fmt.Println("Filename is: " + arg)
}
