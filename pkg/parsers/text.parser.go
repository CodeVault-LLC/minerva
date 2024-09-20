package parsers

import (
	"bufio"
	"io"
)

type TextParser struct {
	ParseFunc func(string) (Item, bool)
}

func (p *TextParser) Parse(reader io.Reader, itemChan chan<- Item) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if item, ok := p.ParseFunc(line); ok {
			itemChan <- item
		}
	}
	return scanner.Err()
}
