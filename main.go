package main

import (
  "fmt"
  gostatic "github.com/wonderzombie/gostatic/lib"
)

func main() {
  fmt.Println("hello gostatic")
  items := gostatic.ParseArticle("test.md")

  fmt.Println(items)
}
