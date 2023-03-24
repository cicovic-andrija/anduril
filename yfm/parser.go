package yfm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type parser struct {
	reader    *bufio.Reader
	collected *bytes.Buffer
	read      int
	start     int
	end       int
}

func newParser(r io.Reader) *parser {
	return &parser{
		reader:    bufio.NewReader(r),
		collected: bytes.NewBuffer(nil),
		read:      0,
		start:     0,
		end:       0,
	}
}

func (p *parser) parse(v interface{}) error {
	for {
		line, eof, err := p.readLine()
		if err != nil || eof {
			return ErrNotFound
		}
		if line == "" {
			continue
		}
		if line != "---" {
			return ErrNotFound
		}
		p.start = p.read
		break
	}

	for {
		read_tmp := p.read
		line, eof, err := p.readLine()
		if err != nil {
			return err
		}

		if line != "---" {
			if eof {
				return ErrInvalidInput
			}
			continue
		}
		p.read = read_tmp
		p.end = p.read

		if err := yaml.Unmarshal(p.collected.Bytes()[p.start:p.end], v); err != nil {
			return fmt.Errorf("failed to decode input: %v", err)
		}

		return nil
	}
}

func (p *parser) readLine() (string, bool, error) {
	line, err := p.reader.ReadBytes('\n')
	eof := err == io.EOF
	if err != nil && !eof {
		return "", false, fmt.Errorf("failed to read input: %v", err)
	}
	n, _ := p.collected.Write(line)
	p.read += n
	return string(bytes.TrimSpace(line)), eof, nil
}
