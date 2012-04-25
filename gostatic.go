package main

import (
  "fmt"
  blackfriday "github.com/russross/blackfriday"
  parser "github.com/wonderzombie/gostatic/lib"
  "io"
  "log"
  "os"
  "path/filepath"
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

// TODO: this seems like a grab bag of data.
type MkdFileInfo struct {
  Content  string
  MetaInfo *Metadata
  OsInfo   os.FileInfo
  Path     string
}

func (m *MkdFileInfo) String() string {
  out := fmt.Sprintf("%v %v %v", m.MetaInfo.Title, m.OsInfo.Name(), m.Path)
  return out
}

type ContentFileInfo string

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
  for l, err = p.ReadLine(); err == nil; l, err = p.ReadLine() {
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
  }

  if err != nil {
    log.Println("Error while reading file:", err)
  }

  return
}

func ParseContent(p *parser.Parser) (content string) {
  c := make([]string, 0)

  for l, err := p.ReadLine(); err == nil; l, err = p.ReadLine() {
    c = append(c, l)
  }

  content = strings.Join(c, "\n")
  return
}

// TODO: better error handling.
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

func ListFiles(dir string) (infos []MkdFileInfo, misc []ContentFileInfo, err error) {
  wf := func(path string, info os.FileInfo, err error) error {
    // TODO: do more to separate out files by extension or type. Return those separately.
    if info.IsDir() {
      return nil
    }
    if !strings.HasSuffix(path, ".md") {
      misc = append(misc, ContentFileInfo(path))
      return nil
    }

    i := MkdFileInfo{}
    i.OsInfo = info
    i.Path = path
    infos = append(infos, i)
    return nil
  }

  err = filepath.Walk(dir, wf)
  if err != nil {
    log.Fatalf("Error while reading %v: %v\n", dir, err)
  }

  return
}

func CopyFile(srcPath, dstPath string) error {
  // TODO: just return error, OK? And maybe wrap it with our own message.
  f, err := os.Open(srcPath)
  if err != nil {
    log.Printf("Unable to open file %v for reading: %v\n", srcPath, err)
    return err
  }
  defer f.Close()

  out, err := os.Create(dstPath)
  if err != nil {
    log.Printf("Unable to create file %v for copying: %v\n", dstPath, err)
    return err
  }
  defer out.Close()

  _, err = io.Copy(out, f)
  if err != nil {
    log.Printf("Unable to copy file %v to destination %v: %v\n", srcPath, dstPath, err)
    return err
  }

  f.Close()
  out.Close()
  return nil
}

func ReplaceExt(file, ext string) string {
  tokens := strings.Split(file, ".")
  if len(tokens) == 1 {
    return file
  }

  i := len(tokens) - 1
  tokens[i] = ext

  return strings.Join(tokens, ".")
}

func main() {
  // TODO: parameterize _pages, _site, et al?
  infos, content, err := ListFiles("_pages")
  if err != nil {
    log.Fatalf("Error while listing files: %v\n", err)
    return
  }

  if len(infos) == 0 {
    log.Fatalln("Read zero files in _pages.")
    return
  }

  for _, info := range infos {
    f, err := os.Open(info.Path)
    if err != nil {
      log.Printf("Skipping file %v because of an error: %v", info.Path, err)
      continue
    }
    defer f.Close()

    m, rawContent := ReadFile(f)
    // Maybe ReadFile() should do this?
    content := string(blackfriday.MarkdownCommon([]byte(rawContent)))
    info.MetaInfo = m
    info.Content = content
    f.Close()
  }

  // Make _site dir.
  mode := os.FileMode(0755)
  err = os.Mkdir("_site", mode)
  if err != nil && !os.IsExist(err) {
    log.Fatalf("Error creating _site directory:", err)
  }

  // Process all of the md files and make them into html. Create dir structure.
  for _, info := range infos {
    // TODO: refactor this chunk into its own func.
    newPath := strings.Replace(info.Path, "_pages", "_site", 1)
    dir, file := filepath.Split(newPath)

    file = ReplaceExt(file, "html")
    newPath = filepath.Join(dir, file)
    log.Printf("%v -> %v", info.Path, newPath)

    // TODO: track which dirs we've already created?
    os.MkdirAll(dir, os.FileMode(0755))

    f, err := os.Create(newPath)
    if err != nil {
      log.Printf("Unable to create file %v: %v\n", newPath, err)
      continue
    }
    defer f.Close()
    _, err = f.WriteString(info.Content)
    if err != nil {
      log.Printf("Error while writing to file %v: %v\n", newPath, err)
      continue
    }
    f.Close()
  }

  // Finally, copy all the other files.
  for _, c := range content {
    miscFile := string(c)
    newPath := strings.Replace(miscFile, "_pages", "_site", 1)

    err = CopyFile(miscFile, newPath)
    if err != nil {
      // TODO: probably redundant.
      log.Println("Error while trying to copy file %v to %v: %v", miscFile, newPath, err)
    }
  }

}
