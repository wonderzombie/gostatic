package gostatic

import (
  "fmt"
  "unicode"
  "unicode/utf8"
)

type itemType int

type item struct {
  typ itemType
  val string
}

const (
  eof = -1
)

type stateFn func(*lexer) stateFn

type lexer struct {
  name  string
  input string
  start int
  pos   int
  width int
  items chan item
  state stateFn
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
  l := &lexer{
    name:  name,
    input: input,
    state: lexDocument,
    items: make(chan item, 2), // Two items sufficient.
  }
  return l
}

// nextItem returns the next item from the input.
func (l *lexer) NextItem() item {
  for {
    select {
    case item := <-l.items:
      return item
    default:
      l.state = l.state(l)
    }
  }
  panic("not reached")
}

func (l *lexer) emit(t itemType) {
  i := item{t, l.input[l.start:l.pos]}
  fmt.Println("Emitting ", i)
  l.items <- i
  l.start = l.pos
}

func (l *lexer) next() (r rune) {
  if l.pos >= len(l.input) {
    l.width = 0
    return eof
  }
  r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
  l.pos += l.width
  return r
}

func (l *lexer) ignore() {
  l.start = l.pos
}

func (l *lexer) backup() {
  l.pos -= l.width
}

// func (l *lexer) peek() int {
//     rune := l.next()
//     l.backup()
//     return rune
// }

// // accept consumes the next rune
// // if it's from the valid set.
// func (l *lexer) accept(valid string) bool {
//     if strings.IndexRune(valid, l.next()) >= 0 {
//         return true
//     }
//     l.backup()
//     return false
// }

// // acceptRun consumes a run of runes from the valid set.
// func (l *lexer) acceptRun(valid string) {
//     for strings.IndexRune(valid, l.next()) >= 0 {
//     }
//     l.backup()
// }

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
  l.items <- item{
    itemError,
    fmt.Sprintf(format, args...),
  }
  return nil
}

func isSpace(r rune) bool {
  switch r {
  case ' ', '\t', '\n', '\r':
    return true
  }
  return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
  return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
