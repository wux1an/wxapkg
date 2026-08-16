package jsbeautifier

import (
	"fmt"
	"strings"
)

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

type output struct {
	indent_string      string
	baseIndentString   string
	indent_cache       []string
	baseIndentLength   int
	indent_length      int
	lines              []*outputline
	current_line       *outputline
	space_before_token bool
}

func (self *output) get_line_number() int {
	return len(self.lines)
}

func (self *output) add_new_line(force_new_line bool) bool {

	if len(self.lines) == 1 && self.just_added_newline() {
		return false
	}

	if force_new_line || !self.just_added_newline() {
		self.current_line = NewLine(self)
		self.lines = append(self.lines, self.current_line)

		return true
	}
	return false
}

func (self *output) get_code() string {
	var sweetCode strings.Builder

	for _, line := range self.lines {
		sweetCode.WriteString(line.String())
		sweetCode.WriteByte('\n')
	}

	return strings.TrimRight(sweetCode.String(), "\r\n\t ")
}

func (self *output) PrintLines() {
	for _, val := range self.lines {
		fmt.Println("Line: ", val, " len: ", val.get_character_count())
	}
}

func (self *output) set_indent(level int) bool {

	if len(self.lines) > 1 {
		for level >= len(self.indent_cache) {
			self.indent_cache = append(self.indent_cache, self.indent_cache[len(self.indent_cache)-1]+self.indent_string)
		}

		self.current_line.set_indent(level)
		return true
	}
	self.current_line.set_indent(0)
	return false
}

func (self *output) add_token(printable_token string) {
	self.add_space_before_token()
	self.current_line.push(printable_token)
}

func (self *output) add_space_before_token() {
	if self.space_before_token && !self.just_added_newline() {
		self.current_line.push(" ")
	}
	self.space_before_token = false
}

func (self *output) remove_redundant_indentation(frame Flags) {
	if frame.multiline_frame || frame.mode == ForInitializer || frame.mode == Conditional {
		return
	}

	index := frame.start_line_index
	for index < len(self.lines) {
		self.lines[index].remove_indent()
		index += 1
	}
}

func (self *output) trim(eat_newlines bool) {
	self.current_line.trim()

	for eat_newlines && len(self.lines) > 1 && self.current_line.is_empty() {
		self.lines = self.lines[:len(self.lines)-1]
		self.current_line = self.lines[len(self.lines)-1]
		self.current_line.trim()
	}
}
func (self *output) just_added_newline() bool {
	return self.current_line.is_empty()
}

func (self *output) just_added_blankline() bool {
	if self.just_added_newline() {
		if len(self.lines) == 1 {
			return true
		}

		line := self.lines[len(self.lines)-2]
		return line.is_empty()
	}

	return false
}

func NewOutput(indent_string string, baseIndentString string) *output {
	output := new(output)
	output.indent_string = indent_string
	output.baseIndentString = baseIndentString
	output.indent_length = len(indent_string)
	output.baseIndentLength = len(baseIndentString)
	output.current_line = nil
	output.indent_cache = append(output.indent_cache, baseIndentString)
	output.add_new_line(true)
	return output

}
