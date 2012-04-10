package gostatic

import (
  "bytes"
  "fmt"
  "io"
  "launchpad.net/goyaml"
  "log"
  "os"
  "regexp"
  "strings"
)

var headerRe = regexp.MustCompile("-+")
var linkRe = regexp.MustCompile("[[(\\w+)}}")

type Metadata struct {
  Title    string
  Category string
  Tags     []string
}

func (m *Metadata) String() string {
  tt := strings.Join(m.Tags, "|")
  return fmt.Sprintf("Title: %v, Category: %v, Tags: %v", m.Title, m.Category, tt)
}

func metadataForBuffer(buf *bytes.Buffer) (m *Metadata) {
  header, err := buf.ReadBytes('\n')
  if err != nil {
    log.Fatalln(err)
  } else if !headerRe.Match(header) {
    log.Fatalln("Missing header")
  }

  var (
    l       []byte
    yaml    bytes.Buffer
    content bytes.Buffer
  )

  for !headerRe.Match(l) {
    if l, err = buf.ReadBytes('\n'); err != nil {
      log.Fatalln("Error while reading file:", err)
    }
    yaml.Write(l)
  }
  fmt.Println("Done parsing header.")

  fmt.Println("Parsing content.")
  for err != io.EOF {
    // Read word by word.
    l, err = buf.ReadBytes(' ')
    begin := bytes.IndexRune(l, '[')
    end := bytes.IndexRune(l, ']')
    if begin >= 0 && end >= 0 {
      link := l[begin:end]
      fmt.Println("Link discovered:", link)
    }
    content.Write(l)
  }

  goyaml.Unmarshal(yaml.Bytes(), &m)
  return
}

func main__() {
  filename := "test.md"
  f, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  buf := &bytes.Buffer{}
  _, err = buf.ReadFrom(f)
  if err != nil {
    log.Fatalln(err)
  }

  m := metadataForBuffer(buf)
  fmt.Println("Metadata:", m)

}
