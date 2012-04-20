package main

import (
  "fmt"
  blackfriday "github.com/russross/blackfriday"
  parser "github.com/wonderzombie/gostatic/lib"
  "log"
  "os"
  "strings"
)

const (
  header = "---"
  delim  = ":"
)

type Metadata struct {
  Title    string
  Category string
  Tags     []string
}

func (m *Metadata) String() string {
  tt := strings.Join(m.Tags, "|")
  return fmt.Sprintf("Title: %v, Category: %v, Tags: %v", m.Title, m.Category, tt)
}

func ParseMetadata(p *parser.Parser) (m *Metadata) {
  l, err := p.ReadLine()
  if err != nil {
    log.Fatalf("Error while trying to parse file: %v\n", err)
  }

  if !strings.HasPrefix(l, header) {
    log.Println("Warning: file had no header. Assuming no metadata.")
    return
  }
  m = new(Metadata)

  l, err = p.ReadLine()
  if err != nil {
    log.Fatalf("Error while reading file:", err)
  }

  for err == nil {
    if strings.HasPrefix(l, header) {
      // Done parsing metadata.
      return
    }

    l = strings.TrimSpace(l)
    l = strings.Trim(l, "\n")
    // Split on the first colon, returning two items total.
    ll := strings.SplitN(l, ":", 2)
    if len(ll) != 2 {
      log.Fatalf("Malformed metadata: %v", l)
    }
    label, data := strings.TrimSpace(ll[0]), strings.TrimSpace(ll[1])

    switch label {
    case "title":
      m.Title = data
    case "category":
      m.Category = data
    case "tags":
      tags := strings.Split(data, ",")
      for _, t := range tags {
        m.Tags = append(m.Tags, strings.TrimSpace(t))
      }
    default:
      log.Println("Ignoring line:", l)
    }

    l, err = p.ReadLine()
  }

  if err != nil {
    log.Println("Error while reading file:", err)
  }

  return
}

func ParseContent(p *parser.Parser) (content string) {
  c := make([]string, 0)

  l, err := p.ReadLine()
  for err == nil {
    c = append(c, l)
    l, err = p.ReadLine()
  }

  content = strings.Join(c, "\n")
  return
}

func ReadFile(f *os.File) (m *Metadata, content string) {
  p := parser.NewParser(f)
  m = ParseMetadata(p)

  if m == nil {
    log.Println("No metadata.")
  }

  // Done parsing metadata.
  content = ParseContent(p)
  return
}

func main() {
  filename := "test.md"
  f, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  m, rawContent := ReadFile(f)

  content := string(blackfriday.MarkdownCommon([]byte(rawContent)))
  log.Println(m)
  log.Println(content)
}
