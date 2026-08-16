package jsbeautifier

import "strings"

// Copyright (c) 2014 Ditashi Sayomi

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

type outputline struct {
	parent          *output
	character_count int
	indent_count    int
	items           []string
	empty           bool
}

func (self *outputline) get_character_count() int {
	return self.character_count
}

func (self *outputline) is_empty() bool {
	return self.empty
}

func (self *outputline) set_indent(level int) {
	self.character_count = self.parent.baseIndentLength + level*self.parent.indent_length
	self.indent_count = level
}

func (self *outputline) last() string {
	if !self.is_empty() {
		return self.items[len(self.items)-1]
	} else {
		return ""
	}
}

func (self *outputline) push(input string) {
	self.items = append(self.items, input)
	self.character_count += len(input)
	self.empty = false
}

func (self *outputline) remove_indent() {
	if self.indent_count > 0 {
		self.indent_count -= 1
		self.character_count -= self.parent.indent_length
	}
}

func (self *outputline) trim() {
	for self.last() == " " {
		self.items = self.items[:len(self.items)-1]

		self.character_count -= 1
	}
	self.empty = len(self.items) == 0
}

func (self *outputline) String() string {
	if self.is_empty() {
		return ""
	}

	var result strings.Builder
	if self.character_count > 0 {
		result.Grow(self.character_count)
	}
	if self.indent_count >= 0 {
		result.WriteString(self.parent.indent_cache[self.indent_count])
	}
	for _, item := range self.items {
		result.WriteString(item)
	}
	return result.String()
}

func NewLine(parent *output) *outputline {
	return &outputline{parent, 0, -1, make([]string, 0), true}
}
