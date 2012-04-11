package gostatic

import (
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

const headerDelim = "---"

func lexDocument(l *lexer) stateFn {
  if strings.HasPrefix(l.input[l.pos:], headerDelim) {
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
      break
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
    default:
      switch r {
      case ' ', '\r', '\n':
        l.emit(itemValue)
        return lexContent
      }
    case r == eof:
      return l.errorf("Reached eof while trying to parse a header value")
    }
  }
  return nil
}

func ParseArticle(filename string) (c chan item) {
  f, err := os.Open(filename)
  if err != nil {
    log.Fatalf("Unable to open file %v: %v\n", f, err)
  }
  in, err := ioutil.ReadAll(f)
  if err != nil {
    log.Fatalf("Error while reading file %v: %v\n", f, err)
  }
  l := lex("test", string(in))
  go func() {
    c <- l.nextItem()
  }()
  return c
}
