package shared

import (
	"errors"
	"fmt"
	"strings"
)

type Position struct {
	Idx  int    // index of character
	Ln   int    // line
	Col  int    // column
	Fn   string // filename
	Ftxt string // text
}

func NewPosition(i, l, c int, fn, ftxt string) *Position {
	return &Position{
		Idx:  i,
		Ln:   l,
		Col:  c,
		Fn:   fn,
		Ftxt: ftxt,
	}
}

func (p *Position) Advance(currentChar string) *Position {
	p.Idx++
	p.Col++

	if currentChar == "\n" {
		p.Ln++
		p.Col = 0
	}

	return p
}

func (p *Position) Copy() *Position {
	return NewPosition(p.Idx, p.Ln, p.Col, p.Fn, p.Ftxt)
}

func stringWithArrows(text string, posStart, posEnd *Position) string {
	result := ""

	// Calculate indices
	idxStart := strings.LastIndex(text[:posStart.Idx], "\n")
	if idxStart == -1 {
		idxStart = 0
	}
	idxEnd := strings.Index(text[idxStart:], "\n")
	if idxEnd == -1 {
		idxEnd = len(text)
	} else {
		idxEnd += idxStart
	}

	// Generate each line
	lineCount := posEnd.Ln - posStart.Ln + 1
	for i := 0; i < lineCount; i++ {
		// Calculate line columns
		line := text[idxStart:idxEnd]
		colStart := 0
		if i == 0 {
			colStart = posStart.Col
		}
		colEnd := len(line) - 1
		if i == lineCount-1 {
			colEnd = posEnd.Col
		}

		// Append to result
		result += line + "\n"
		result += strings.Repeat(" ", colStart) + strings.Repeat("^", colEnd-colStart)

		// Re-calculate indices
		idxStart = idxEnd
		idxEnd = strings.Index(text[idxStart:], "\n")
		if idxEnd == -1 {
			idxEnd = len(text)
		} else {
			idxEnd += idxStart
		}
	}

	return strings.ReplaceAll(result, "\t", "")
}

func baseError(posStart *Position, posEnd *Position, errorName string, details string) error {
	arrows := stringWithArrows(posStart.Ftxt, posStart, posEnd)
	errorMsg := fmt.Sprintf("%s: %s\nFile %s, line %d\n\n%s", errorName, details, posStart.Fn, posStart.Ln+1, arrows)
	return errors.New(errorMsg)
}

func IllegalCharError(posStart *Position, posEnd *Position, details string) error {
	return baseError(posStart, posEnd, "Illegal Character", details)
}

func InvalidSyntaxError(posStart *Position, posEnd *Position, details string) error {
	return baseError(posStart, posEnd, "Invalid Syntax", details)
}
