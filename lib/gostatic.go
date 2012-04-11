package gostatic

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strings"
)

const (
  itemError itemType = iota
  itemEOF
  itemHeader
  itemLabel
  itemValue
  itemContent
)

var itemName = map[itemType]string{
  itemError:   "error",
  itemEOF:     "EOF",
  itemHeader:  "header",
  itemLabel:   "label",
  itemValue:   "value",
  itemContent: "content",
}

func (i item) String() string {
  switch i.typ {
  case itemEOF:
    return "EOF"
  case itemError:
    return i.val
  }
  return fmt.Sprintf("%q: %q", itemName[i.typ], i.val)
}

const headerDelim = "---"

func lexDocument(l *lexer) stateFn {
  if strings.HasPrefix(l.input[l.pos:], headerDelim) {
    l.pos += len(headerDelim)
    l.emit(itemHeader)
    return lexHeader
  } else {
    return l.errorf("File must begin with a header (%v).", headerDelim)
  }
  l.emit(itemEOF)
  return nil
}

func lexContent(l *lexer) stateFn {
  for {
    switch r := l.next(); {
    case r == eof:
      l.emit(itemContent)
      return l.errorf("done")
    }
  }
  return nil
}

func lexHeader(l *lexer) stateFn {
  for {
    if strings.HasPrefix(headerDelim, l.input[l.pos:]) {
      return lexContent
    }

    switch r := l.next(); {
    case r == ':':
      l.emit(itemLabel)
      return lexValue
    case isAlphaNumeric(r):
      // eat it up
    case r == eof:
      return l.errorf("Never found end of header (%v).", headerDelim)
    }
  }
  return nil
}

func lexValue(l *lexer) stateFn {
  for {
    switch r := l.next(); {
    case r == ' ':
      fmt.Println("ignoring space")
      l.ignore()
    default:
      fmt.Println("is not space")
      switch r {
      case '\n':
        fmt.Println("newline")
        l.backup()
        l.emit(itemValue)
        l.next()
        l.ignore()
        return lexHeader
      }
    case r == eof:
      return l.errorf("Reached eof while trying to parse a header value")
    }
  }
  return nil
}

func ParseArticle(filename string) []item {
  f, err := os.Open(filename)
  if err != nil {
    log.Fatalf("Unable to open file %v: %v\n", f, err)
  }
  in, err := ioutil.ReadAll(f)
  if err != nil {
    log.Fatalf("Error while reading file %v: %v\n", f, err)
  }
  l := lex("test", string(in))
  items := make([]item, 0)
  i := l.NextItem()
  for i.typ != itemEOF {
    items = append(items, i)
    i = l.NextItem()
  }
  return items
}
