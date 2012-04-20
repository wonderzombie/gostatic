package parser

import (
  "bufio"
  "os"
)

// Parser is a thin wrapper around a buffer. It caches the most recently read line. This allows the
// user to read a line without consuming it.
type Parser struct {
  reader *bufio.Reader
  Line   string
}

// Creates a new Parser for the given file.
func NewParser(f *os.File) *Parser {
  p := new(Parser)
  p.init(f)
  return p
}

// Initializes the Parser with a file.
func (p *Parser) init(f *os.File) {
  p.reader = bufio.NewReader(f)
}

// Reads a line from the file and caches it.
func (p *Parser) ReadLine() (string, error) {
  line, err := p.reader.ReadString('\n')
  if err != nil {
    return "", err
  }
  p.Line = line
  return line, nil
}
