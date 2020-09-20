package main

import (
  "fmt"
)

func main() {
  // Println function is used to
  // display output in the next line
  fmt.Print("Enter the filename: ")

  // var then variable name then variable type
  var filename string

  // Taking input from user
  fmt.Scanln(&filename)
  fmt.Println("Filename is: " + filename)
}
