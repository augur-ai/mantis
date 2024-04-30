/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Module file printer.

package modfile

import (
	"bytes"
	"fmt"
	"strings"
)

// Format returns a go.mod file as a byte slice, formatted in standard style.
func Format(f *FileSyntax) []byte {
	pr := &printer{}
	pr.file(f)
	return pr.Bytes()
}

// A printer collects the state during printing of a file or expression.
type printer struct {
	bytes.Buffer           // output buffer
	comment      []Comment // pending end-of-line comments
	margin       int       // left margin (indent), a number of tabs
}

// printf prints to the buffer.
func (p *printer) printf(format string, args ...interface{}) {
	fmt.Fprintf(p, format, args...)
}

// indent returns the position on the current line, in bytes, 0-indexed.
func (p *printer) indent() int {
	b := p.Bytes()
	n := 0
	for n < len(b) && b[len(b)-1-n] != '\n' {
		n++
	}
	return n
}

// newline ends the current line, flushing end-of-line comments.
func (p *printer) newline() {
	if len(p.comment) > 0 {
		p.printf(" ")
		for i, com := range p.comment {
			if i > 0 {
				p.trim()
				p.printf("\n")
				for i := 0; i < p.margin; i++ {
					p.printf("\t")
				}
			}
			p.printf("%s", strings.TrimSpace(com.Token))
		}
		p.comment = p.comment[:0]
	}

	p.trim()
	p.printf("\n")
	for i := 0; i < p.margin; i++ {
		p.printf("\t")
	}
}

// trim removes trailing spaces and tabs from the current line.
func (p *printer) trim() {
	// Remove trailing spaces and tabs from line we're about to end.
	b := p.Bytes()
	n := len(b)
	for n > 0 && (b[n-1] == '\t' || b[n-1] == ' ') {
		n--
	}
	p.Truncate(n)
}

// file formats the given file into the print buffer.
func (p *printer) file(f *FileSyntax) {
	for _, com := range f.Before {
		p.printf("%s", strings.TrimSpace(com.Token))
		p.newline()
	}

	for i, stmt := range f.Stmt {
		switch x := stmt.(type) {
		case *CommentBlock:
			// comments already handled
			p.expr(x)

		default:
			p.expr(x)
			p.newline()
		}

		for _, com := range stmt.Comment().After {
			p.printf("%s", strings.TrimSpace(com.Token))
			p.newline()
		}

		if i+1 < len(f.Stmt) {
			p.newline()
		}
	}
}

func (p *printer) expr(x Expr) {
	// Emit line-comments preceding this expression.
	if before := x.Comment().Before; len(before) > 0 {
		// Want to print a line comment.
		// Line comments must be at the current margin.
		p.trim()
		if p.indent() > 0 {
			// There's other text on the line. Start a new line.
			p.printf("\n")
		}
		// Re-indent to margin.
		for i := 0; i < p.margin; i++ {
			p.printf("\t")
		}
		for _, com := range before {
			p.printf("%s", strings.TrimSpace(com.Token))
			p.newline()
		}
	}

	switch x := x.(type) {
	default:
		panic(fmt.Errorf("printer: unexpected type %T", x))

	case *CommentBlock:
		// done

	case *LParen:
		p.printf("(")
	case *RParen:
		p.printf(")")

	case *Line:
		sep := ""
		for _, tok := range x.Token {
			p.printf("%s%s", sep, tok)
			sep = " "
		}

	case *LineBlock:
		for _, tok := range x.Token {
			p.printf("%s ", tok)
		}
		p.expr(&x.LParen)
		p.margin++
		for _, l := range x.Line {
			p.newline()
			p.expr(l)
		}
		p.margin--
		p.newline()
		p.expr(&x.RParen)
	}

	// Queue end-of-line comments for printing when we
	// reach the end of the line.
	p.comment = append(p.comment, x.Comment().Suffix...)
}
