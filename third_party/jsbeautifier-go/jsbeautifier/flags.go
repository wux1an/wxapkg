package jsbeautifier

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

type Flags struct {
	mode      Mode
	parent    *Flags
	last_text string

	multiline_frame        bool
	start_line_index       int
	indentation_level      int
	line_indent_level      int
	if_block               bool
	do_block               bool
	in_case                bool
	declaration_statement  bool
	declaration_assignment bool
	do_while               bool
	ternary_depth          int
	last_word              string
	else_block             bool
	in_case_statement      bool
	case_body              bool
}

func (self *Flags) apply_base(flags_base *Flags, added_newline bool) {
	next_indent_level := flags_base.indentation_level
	if !added_newline && flags_base.line_indent_level > next_indent_level {
		next_indent_level = flags_base.line_indent_level
	}

	self.parent = flags_base
	self.last_text = flags_base.last_text
	self.last_word = flags_base.last_word
	self.indentation_level = next_indent_level
}

func NewFlags(mode Mode) *Flags {
	flags := new(Flags)
	flags.mode = mode
	return flags
}
