package gostatic

import (
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
      return lexDocument
    }
  }
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
}

func main() {

}
