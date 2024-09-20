package parsers

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type ListType string

const (
	Domain ListType = "Domain"
	URL    ListType = "URL"
	IPv4   ListType = "IPv4"
	IPv6   ListType = "IPv6"
	Email  ListType = "Email"
)

type Item struct {
	Type     ListType
	Value    string
	Metadata map[string]interface{}
}

type Parser interface {
	Parse(reader io.Reader, itemChan chan<- Item) error
}

func ParseFile(filename string, parser Parser) (<-chan Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	itemChan := make(chan Item, 1000) // Buffered channel for better performance

	go func() {
		defer close(itemChan)
		defer file.Close()
		if err := parser.Parse(file, itemChan); err != nil {

			fmt.Printf("Error parsing file %s: %v\n", filename, err)
			return
		}
	}()

	return itemChan, nil
}

func ParseBytes(data []byte, parser Parser) (<-chan Item, error) {
	if parser == nil {
		return nil, fmt.Errorf("parser is nil")
	}

	itemChan := make(chan Item, 1000) // Buffered channel for better performance

	go func() {
		defer close(itemChan)
		if err := parser.Parse(bytes.NewReader(data), itemChan); err != nil {
			fmt.Printf("Error parsing bytes: %v\n", err)
			return
		}
	}()

	return itemChan, nil
}
