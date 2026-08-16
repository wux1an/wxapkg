package jsbeautifier

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/ditashi/jsbeautifier-go/optargs"
	"github.com/ditashi/jsbeautifier-go/tokenizer"
	"github.com/ditashi/jsbeautifier-go/unpackers"
	"github.com/ditashi/jsbeautifier-go/utils"
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

var default_options = map[string]interface{}{
	"indent_size":           4,
	"indent_char":           " ",
	"indent_with_tabs":      false,
	"preserve_newlines":     true,
	"max_preserve_newlines": 10,
	"space_in_paren":        false,
	"space_in_empty_paren":  false,
	"e4x":                       false,
	"jslint_happy":              false,
	"space_after_anon_function": false,
	"brace_style":               "collapse",
	"keep_array_indentation":    false,
	"keep_function_indentation": false,
	"eval_code":                 false,
	"unescape_strings":          false,
	"wrap_line_length":          0,
	"break_chained_methods":     false,
	"end_with_newline":          false,
}

var handlers = map[string]func(*jsbeautifier, *tokenizer.Token){
	"TK_START_EXPR":     (*jsbeautifier).handle_start_expr,
	"TK_END_EXPR":       (*jsbeautifier).handle_end_expr,
	"TK_START_BLOCK":    (*jsbeautifier).handle_start_block,
	"TK_END_BLOCK":      (*jsbeautifier).handle_end_block,
	"TK_WORD":           (*jsbeautifier).handle_word,
	"TK_RESERVED":       (*jsbeautifier).handle_word,
	"TK_SEMICOLON":      (*jsbeautifier).handle_semicolon,
	"TK_STRING":         (*jsbeautifier).handle_string,
	"TK_EQUALS":         (*jsbeautifier).handle_equals,
	"TK_OPERATOR":       (*jsbeautifier).handle_operator,
	"TK_COMMA":          (*jsbeautifier).handle_comma,
	"TK_BLOCK_COMMENT":  (*jsbeautifier).handle_block_comment,
	"TK_INLINE_COMMENT": (*jsbeautifier).handle_inline_comment,
	"TK_COMMENT":        (*jsbeautifier).handle_comment,
	"TK_DOT":            (*jsbeautifier).handle_dot,
	"TK_UNKNOWN":        (*jsbeautifier).handle_unknown,
	"TK_EOF":            (*jsbeautifier).handle_eof,
}

type jsbeautifier struct {
	flags            *Flags
	previous_flags   *Flags
	flag_store       []Flags
	tokens           []string
	token_pos        int
	options          optargs.MapType
	output           *output
	last_last_text   string
	last_type        string
	ternary_depth    int
	indent_string    string
	baseIndentString string
	token_queue      tokenizer.TokenStack
	tkch             chan tokenizer.Token
}

type Mode int

const (
	BlockStatement Mode = iota
	Statement
	ObjectLiteral
	ArrayLiteral
	ForInitializer
	Conditional
	Expression
)

func (self *jsbeautifier) unpack(s *string, eval_code bool) *string {
	return unpackers.Run(s)
}

func (self *jsbeautifier) beautify(s *string) (string, error) {

	if !utils.InStrArray(self.options["brace_style"].(string), []string{"expand", "collapse", "end-expand", "none"}) {
		return "", errors.New("opts.brace-style must be \"expand\", \"collapse\", \"end-exapnd\", or \"none\".")
	}

	s = self.blank_state(s)

	input := self.unpack(s, self.options["eval_code"].(bool))
	t := tokenizer.New(input, self.options, "    ")
	self.tkch = t.Tokenize()
	self.token_pos = 0

	for token := range self.tkch {
		self.parse_token(token)

		for !self.token_queue.Empty() {
			qtoken := self.token_queue.Shift()
			self.parse_token(*qtoken)
		}
	}

	sweet_code := self.output.get_code()
	if self.options["end_with_newline"].(bool) {
		sweet_code += "\n"
	}

	return sweet_code, nil
}

func (self *jsbeautifier) parse_token(t tokenizer.Token) {
	for _, comment_token := range t.CommentsBefore() {
		self.handle_token(&comment_token)
	}
	self.handle_token(&t)

	self.last_last_text = self.flags.last_text
	self.last_type = t.Type()
	self.flags.last_text = t.Text()
	self.token_pos++
}

func (self *jsbeautifier) handle_token(t *tokenizer.Token) {

	newlines := t.NewLines()
	keep_whitespace := self.options["keep_array_indentation"].(bool) && self.is_array(self.flags.mode)

	if keep_whitespace {
		for i := 0; i < newlines; i++ {
			self.print_newline(i > 0, false)
		}
	} else {
		if self.options["max_preserve_newlines"].(int) != 0 && newlines > self.options["max_preserve_newlines"].(int) {
			newlines = self.options["max_preserve_newlines"].(int)
		}
		if self.options["preserve_newlines"].(bool) && newlines > 1 {
			self.print_newline(false, false)
			for i := 1; i < newlines; i++ {
				self.print_newline(true, false)

			}
		}
	}

	handlers[t.Type()](self, t)

}

func (self *jsbeautifier) is_special_word(s string) bool {
	return utils.InStrArray(s, []string{"case", "return", "do", "if", "throw", "else"})
}

func (self *jsbeautifier) is_array(mode Mode) bool {
	return mode == ArrayLiteral
}

func (self *jsbeautifier) is_expression(mode Mode) bool {
	return mode == Expression || mode == ForInitializer || mode == Conditional
}

func (self *jsbeautifier) allow_wrap_or_preserved_newline(current_token tokenizer.Token, force_linewrap bool) {

	if self.output.just_added_newline() {
		return
	}

	if (self.options["preserve_newlines"].(bool) && current_token.WantedNewLine()) || force_linewrap {
		self.print_newline(false, true)
	} else if self.options["wrap_line_length"].(int) > 0 {
		proposed_line_length := self.output.current_line.get_character_count() + len(current_token.Text())

		if self.output.space_before_token {
			proposed_line_length += 1
		}

		if proposed_line_length >= self.options["wrap_line_length"].(int) {
			self.print_newline(false, true)
		}
	}
}

func (self *jsbeautifier) print_newline(force_newline bool, preserve_satatement_flags bool) {
	if !preserve_satatement_flags {

		if self.flags.last_text != ";" && self.flags.last_text != "," && self.flags.last_text != "=" && self.last_type != "TK_OPERATOR" {
			for self.flags.mode == Statement && !self.flags.if_block && !self.flags.do_block {
				self.restore_mode()
			}
		}
	}

	if self.output.add_new_line(force_newline) {
		self.flags.multiline_frame = true
	}
}

func (self *jsbeautifier) print_token_line_indentation(current_token tokenizer.Token) {

	if self.output.just_added_newline() {
		line := self.output.current_line
		if self.options["keep_array_indentation"].(bool) && self.is_array(self.flags.mode) && current_token.WantedNewLine() {
			line.push(current_token.WhitespaceBefore())
			self.output.space_before_token = false
		} else if self.output.set_indent(self.flags.indentation_level) {
			self.flags.line_indent_level = self.flags.indentation_level
		}
	}

}

func (self *jsbeautifier) print_token(current_token tokenizer.Token, s string) {
	if s == "" {
		s = current_token.Text()
	}
	self.print_token_line_indentation(current_token)
	self.output.add_token(s)
}

func (self *jsbeautifier) indent() {
	self.flags.indentation_level += 1
}

func (self *jsbeautifier) deindent() {
	allow_deindent := self.flags.indentation_level > 0 && ((self.flags.parent == nil) || self.flags.indentation_level > self.flags.parent.indentation_level)
	if allow_deindent {
		self.flags.indentation_level -= 1
	}
}

func (self *jsbeautifier) set_mode(mode Mode) {
	if self.flags != nil {
		self.flag_store = append(self.flag_store, *self.flags)
		self.previous_flags = self.flags
	} else {
		self.previous_flags = NewFlags(mode)
	}

	self.flags = NewFlags(mode)
	self.flags.apply_base(self.previous_flags, self.output.just_added_newline())

	self.flags.start_line_index = self.output.get_line_number()
}

func (self *jsbeautifier) restore_mode() {
	if len(self.flag_store) > 0 {
		self.previous_flags = self.flags

		self.flags = &self.flag_store[len(self.flag_store)-1]
		self.flag_store = self.flag_store[:len(self.flag_store)-1]

		if self.previous_flags.mode == Statement {
			self.output.remove_redundant_indentation(*self.previous_flags)
		}
	}
}

func (self *jsbeautifier) start_of_object_property() bool {
	return self.flags.parent.mode == ObjectLiteral && self.flags.mode == Statement && ((self.flags.last_text == ":" && self.flags.ternary_depth == 0) || (self.last_type == "TK_RESERVED" && (self.flags.last_text == "get" || self.flags.last_text == "set")))
}

func (self *jsbeautifier) start_of_statement(current_token tokenizer.Token) bool {

	if (self.last_type == "TK_RESERVED" && utils.InStrArray(self.flags.last_text, []string{"var", "let", "const"}) && current_token.Type() == "TK_WORD") || (self.last_type == "TK_RESERVED" && self.flags.last_text == "do") || (self.last_type == "TK_RESERVED" && self.flags.last_text == "return" && !current_token.WantedNewLine()) || (self.last_type == "TK_RESERVED" && self.flags.last_text == "else" && !(current_token.Type() == "TK_RESERVED" && current_token.Text() == "if")) || (self.last_type == "TK_END_EXPR" && (self.previous_flags.mode == ForInitializer || self.previous_flags.mode == Conditional)) || (self.last_type == "TK_WORD" && self.flags.mode == BlockStatement && !self.flags.in_case && !(current_token.Text() == "--" || current_token.Text() == "++") && current_token.Type() != "TK_WORD" && current_token.Type() != "TK_RESERVED") || (self.flags.mode == ObjectLiteral && ((self.flags.last_text == ":" && self.flags.ternary_depth == 0) || (self.last_type == "TK_RESERVED" && utils.InStrArray(self.flags.last_text, []string{"get", "set"})))) {

		self.set_mode(Statement)
		self.indent()
		if self.last_type == "TK_RESERVED" && utils.InStrArray(self.flags.last_text, []string{"var", "let", "const"}) && current_token.Type() == "TK_WORD" {
			self.flags.declaration_statement = true
		}

		if !self.start_of_object_property() {
			self.allow_wrap_or_preserved_newline(current_token, (current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), []string{"do", "for", "if", "while"})))
		}

		return true
	} else {
		return false
	}
}

func (self *jsbeautifier) get_token() (tokenizer.Token, bool) {
	token, more := <-self.tkch

	if more {
		self.token_queue.Append(token)
	}

	return token, more
}

func (self *jsbeautifier) handle_start_expr(current_token *tokenizer.Token) {
	if self.start_of_statement(*current_token) {

	}
	next_mode := Expression
	if current_token.Text() == "[" {
		if self.last_type == "TK_WORD" || self.flags.last_text == ")" {
			if self.last_type == "TK_RESERVED" && utils.InStrArray(self.flags.last_text, tokenizer.GetLineStarters()) {
				self.output.space_before_token = true
			}
			self.set_mode(next_mode)
			self.print_token(*current_token, "")
			self.indent()
			if self.options["space_in_paren"].(bool) {
				self.output.space_before_token = true
			}
			return
		}

		next_mode = ArrayLiteral
		if self.is_array(self.flags.mode) {
			if self.flags.last_text == "[" || (self.flags.last_text == "," && (self.last_last_text == "]" || self.last_last_text == "}")) {
				if !self.options["keep_array_indentation"].(bool) {
					self.print_newline(false, false)
				}
			}
		}
	} else {
		if self.last_type == "TK_RESERVED" && self.flags.last_text == "for" {
			next_mode = ForInitializer
		} else if self.last_type == "TK_RESERVED" && (self.flags.last_text == "if" || self.flags.last_text == "while") {
			next_mode = Conditional
		} else {
			next_mode = Expression
		}
	}

	if self.flags.last_text == ";" || self.last_type == "TK_START_BLOCK" {
		self.print_newline(false, false)
	} else if utils.InStrArray(self.last_type, []string{"TK_END_EXPR", "TK_START_EXPR", "TK_END_BLOCK"}) || self.flags.last_text == "." {
		self.allow_wrap_or_preserved_newline(*current_token, current_token.WantedNewLine())
	} else if !(self.last_type == "TK_RESERVED" && current_token.Text() == "(") && !utils.InStrArray(self.last_type, []string{"TK_WORD", "TK_OPERATOR"}) {
		self.output.space_before_token = true
	} else if (self.last_type == "TK_RESERVED" && (self.flags.last_word == "function" || self.flags.last_word == "typeof")) || (self.flags.last_text == "*" && self.last_last_text == "function") {
		if self.options["space_after_anon_function"].(bool) {
			self.output.space_before_token = true
		}
	} else if self.last_type == "TK_RESERVED" && (utils.InStrArray(self.flags.last_text, tokenizer.GetLineStarters()) || self.flags.last_text == "catch") {
		self.output.space_before_token = true
	}

	if self.last_type == "TK_EQUALS" || self.last_type == "TK_OPERATOR" {
		if !self.start_of_object_property() {
			self.allow_wrap_or_preserved_newline(*current_token, false)
		}
	}

	self.set_mode(next_mode)
	self.print_token(*current_token, "")

	if self.options["space_in_paren"].(bool) {
		self.output.space_before_token = true
	}
	self.indent()
}

func (self *jsbeautifier) handle_start_block(current_token *tokenizer.Token) {
	next_token, nmore := self.get_token()
	second_token, smore := self.get_token()
	if smore && ((second_token.Text() == ":" && utils.InStrArray(next_token.Type(), []string{"TK_STRING", "TK_WORD", "TK_RESERVED"})) || (utils.InStrArray(next_token.Text(), []string{"get", "set"}) && utils.InStrArray(second_token.Type(), []string{"TK_WORD", "TK_RESE$RVED"}))) {
		if !utils.InStrArray(self.last_last_text, []string{"class", "interface"}) {
			self.set_mode(ObjectLiteral)
		} else {
			self.set_mode(BlockStatement)
		}
	} else {
		self.set_mode(BlockStatement)
	}

	empty_braces := (nmore) && len(next_token.CommentsBefore()) == 0 && next_token.Text() == "}"
	empty_anonymous_function := empty_braces && self.flags.last_word == "function" && self.last_type == "TK_END_EXPR"

	if self.options["brace_style"].(string) == "expand" || (self.options["brace_style"].(string) == "none" && current_token.WantedNewLine()) {
		if self.last_type != "TK_OPERATOR" && (empty_anonymous_function || self.last_type == "TK_EQUALS" || (self.last_type == "TK_RESERVED" && self.is_special_word(self.flags.last_text) && self.flags.last_text != "else")) {
			self.output.space_before_token = true
		} else {
			self.print_newline(false, true)
		}
	} else {
		if !utils.InStrArray(self.last_type, []string{"TK_OPERATOR", "TK_START_EXPR"}) {
			if self.last_type == "TK_START_BLOCK" {
				self.print_newline(false, false)
			} else {
				self.output.space_before_token = true
			}
		} else {
			if self.is_array(self.previous_flags.mode) && self.flags.last_text == "," {
				if self.last_last_text == "}" {
					self.output.space_before_token = true
				} else {
					self.print_newline(false, false)
				}
			}
		}
	}

	self.print_token(*current_token, "")
	self.indent()
}

func (self *jsbeautifier) handle_end_expr(current_token *tokenizer.Token) {
	for self.flags.mode == Statement {
		self.restore_mode()
	}

	if self.flags.multiline_frame {
		self.allow_wrap_or_preserved_newline(*current_token, current_token.Text() == "]" && self.is_array(self.flags.mode) && !self.options["keep_array_indentation"].(bool))
	}

	if self.options["space_in_paren"].(bool) {
		if self.last_type == "TK_START_EXPR" && !self.options["space_in_empty_paren"].(bool) {
			self.output.space_before_token = false
			self.output.trim(false)
		} else {
			self.output.space_before_token = true
		}
	}

	if current_token.Text() == "]" && self.options["keep_array_indentation"].(bool) {
		self.print_token(*current_token, "")
		self.restore_mode()
	} else {
		self.restore_mode()
		self.print_token(*current_token, "")
	}

	self.output.remove_redundant_indentation(*self.previous_flags)

	if self.flags.do_while && self.previous_flags.mode == Conditional {
		self.previous_flags.mode = Expression
		self.flags.do_block = false
		self.flags.do_block = false
	}
}

func (self *jsbeautifier) handle_end_block(current_token *tokenizer.Token) {
	for self.flags.mode == Statement {
		self.restore_mode()
	}

	empty_braces := self.last_type == "TK_START_BLOCK"

	if self.options["brace_style"].(string) == "expand" {
		if !empty_braces {
			self.print_newline(false, false)
		}
	} else {
		if !empty_braces {
			if self.is_array(self.flags.mode) && self.options["keep_array_indentation"].(bool) {
				self.options["keep_array_indentation"] = false
				self.print_newline(false, false)
				self.options["keep_array_indentation"] = true
			} else {
				self.print_newline(false, false)
			}
		}
	}

	self.restore_mode()
	self.print_token(*current_token, "")
}

func (self *jsbeautifier) handle_word(current_token *tokenizer.Token) {

	if current_token.Type() == "TK_RESERVED" && self.flags.mode != ObjectLiteral && (current_token.Text() == "set" || current_token.Text() == "get") {
		current_token.SetType("TK_WORD")
	}

	if current_token.Type() == "TK_RESERVED" && self.flags.mode == ObjectLiteral {
		next_token, _ := self.get_token()
		if next_token.Text() == ":" {
			current_token.SetType("TK_WORD")
		}
	}

	if self.start_of_statement(*current_token) {
	} else if current_token.WantedNewLine() && !self.is_expression(self.flags.mode) && (self.last_type != "TK_OPERATOR" || (self.flags.last_text == "--" || self.flags.last_text == "++")) && self.last_type != "TK_EQUALS" && (self.options["preserve_newlines"].(bool) || !(self.last_type == "TK_RESERVED" && utils.InStrArray(self.flags.last_text, []string{"var", "let", "const", "set", "get"}))) {
		self.print_newline(false, false)
	}

	if self.flags.do_block && !self.flags.do_while {
		if current_token.Type() == "TK_RESERVED" && current_token.Text() == "while" {
			self.output.space_before_token = true
			self.print_token(*current_token, "")
			self.output.space_before_token = true
			self.flags.do_while = true
			return
		} else {
			self.print_newline(false, false)
			self.flags.do_block = false
		}
	}

	if self.flags.if_block {
		if (!self.flags.else_block) && (current_token.Type() == "TK_RESERVED" && current_token.Text() == "else") {
			self.flags.else_block = true
		} else {
			for self.flags.mode == Statement {
				self.restore_mode()
			}
			self.flags.if_block = false
		}
	}

	if current_token.Type() == "TK_RESERVED" && (current_token.Text() == "case" || (current_token.Text() == "default" && self.flags.in_case_statement)) {
		self.print_newline(false, false)

		if self.flags.case_body || self.options["jslint_happy"].(bool) {
			self.flags.case_body = false
			self.deindent()
		}

		self.print_token(*current_token, "")
		self.flags.in_case = true
		self.flags.in_case_statement = true
		return
	}

	if current_token.Type() == "TK_RESERVED" && current_token.Text() == "function" {

		if (self.flags.last_text == "}" || self.flags.last_text == ";") || (self.output.just_added_newline() && !utils.InStrArray(self.flags.last_text, []string{"[", "{", ":", "=", ","})) {
			if !self.output.just_added_blankline() && len(current_token.CommentsBefore()) == 0 {
				self.print_newline(false, false)
				self.print_newline(true, false)
			}
		}

		if self.last_type == "TK_RESERVED" || self.last_type == "TK_WORD" {
			if self.last_type == "TK_RESERVED" && utils.InStrArray(self.flags.last_text, []string{"get", "set", "new", "return", "export"}) {
				self.output.space_before_token = true
			} else if self.last_type == "TK_RESERVED" && self.flags.last_text == "default" && self.last_last_text == "export" {
				self.output.space_before_token = true
			} else {
				self.print_newline(false, false)
			}
		} else if self.last_type == "TK_OPERATOR" || self.flags.last_text == "=" {
			self.output.space_before_token = true
		} else if !self.flags.multiline_frame && (self.is_expression(self.flags.mode) || self.is_array(self.flags.mode)) {

		} else {
			self.print_newline(false, false)
		}
	}

	if utils.InStrArray(self.last_type, []string{"TK_COMMA", "TK_START_EXPR", "TK_EQUALS", "TK_OPERATOR"}) {

		if !self.start_of_object_property() {
			self.allow_wrap_or_preserved_newline(*current_token, false)
		}
	}

	if current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), []string{"function", "get", "set"}) {
		self.print_token(*current_token, "")
		self.flags.last_word = current_token.Text()
		return
	}

	prefix := "NONE"
	if self.last_type == "TK_END_BLOCK" {
		if !(current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), []string{"else", "catch", "finally"})) {
			prefix = "NEWLINE"
		} else {
			if utils.InStrArray(self.options["brace_style"].(string), []string{"expand", "end-expand"}) || (self.options["brace_style"].(string) == "none" && current_token.WantedNewLine()) {
				prefix = "NEWLINE"
			} else {
				prefix = "SPACE"
				self.output.space_before_token = true
			}
		}
	} else if self.last_type == "TK_SEMICOLON" && self.flags.mode == BlockStatement {
		prefix = "NEWLINE"
	} else if self.last_type == "TK_SEMICOLON" && self.is_expression(self.flags.mode) {
		prefix = "SPACE"
	} else if self.last_type == "TK_STRING" {
		prefix = "NEWLINE"
	} else if self.last_type == "TK_RESERVED" || self.last_type == "TK_WORD" || (self.flags.last_text == "*" && self.last_last_text == "function") {
		prefix = "SPACE"
	} else if self.last_type == "TK_START_BLOCK" {
		prefix = "NEWLINE"
	} else if self.last_type == "TK_END_EXPR" {
		self.output.space_before_token = true
		prefix = "NEWLINE"
	}

	if current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), tokenizer.GetLineStarters()) && self.flags.last_text != ")" {

		if self.flags.last_text == "else" || self.flags.last_text == "export" {
			prefix = "SPACE"
		} else {
			prefix = "NEWLINE"
		}
	}

	if current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), []string{"else", "catch", "finally"}) {
		if self.last_type != "TK_END_BLOCK" || self.options["brace_style"].(string) == "expand" || self.options["brace_style"].(string) == "end-expand" || (self.options["brace_style"].(string) == "none" && current_token.WantedNewLine()) {
			self.print_newline(false, false)
		} else {
			self.output.trim(true)

			if self.output.current_line.last() != "}" {
				self.print_newline(false, false)
			}

			self.output.space_before_token = true
		}
	} else if prefix == "NEWLINE" {
		if self.last_type == "TK_RESERVED" && self.is_special_word(self.flags.last_text) {
			self.output.space_before_token = true
		} else if self.last_type != "TK_END_EXPR" {
			if (self.last_type != "TK_START_EXPR" || !(current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), []string{"var", "let", "const"}))) && self.flags.last_text != ":" {
				if current_token.Type() == "TK_RESERVED" && current_token.Text() == "if" && self.flags.last_text == "else" {
					self.output.space_before_token = true
				} else {
					self.print_newline(false, false)
				}
			}
		} else if current_token.Type() == "TK_RESERVED" && utils.InStrArray(current_token.Text(), tokenizer.GetLineStarters()) && self.flags.last_text != ")" {
			self.print_newline(false, false)
		}
	} else if self.flags.multiline_frame && self.is_array(self.flags.mode) && self.flags.last_text == "," && self.last_last_text == "}" {
		self.print_newline(false, false)
	} else if prefix == "SPACE" {
		self.output.space_before_token = true
	}

	self.print_token(*current_token, "")
	self.flags.last_word = current_token.Text()

	if current_token.Type() == "TK_RESERVED" && current_token.Text() == "do" {
		self.flags.do_block = true
	}

	if current_token.Type() == "TK_RESERVED" && current_token.Text() == "if" {
		self.flags.if_block = true
	}

}

func (self *jsbeautifier) handle_semicolon(current_token *tokenizer.Token) {
	if self.start_of_statement(*current_token) {
		self.output.space_before_token = false
	}
	for self.flags.mode == Statement && !self.flags.if_block && !self.flags.do_block {
		self.restore_mode()
	}

	self.print_token(*current_token, "")
}

func (self *jsbeautifier) handle_string(current_token *tokenizer.Token) {
	if self.start_of_statement(*current_token) {
		self.output.space_before_token = true
	} else if self.last_type == "TK_RESERVED" || self.last_type == "TK_WORD" {
		self.output.space_before_token = true
	} else if utils.InStrArray(self.last_type, []string{"TK_COMMA", "TK_START_EXPR", "TK_EQUALS", "TK_OPERATOR"}) {
		if !self.start_of_object_property() {
			self.allow_wrap_or_preserved_newline(*current_token, false)
		}
	} else {
		self.print_newline(false, false)
	}

	self.print_token(*current_token, "")
}

func (self *jsbeautifier) handle_equals(current_token *tokenizer.Token) {
	if self.start_of_statement(*current_token) {

	}

	if self.flags.declaration_statement {
		self.flags.declaration_assignment = true
	}

	self.output.space_before_token = true
	self.print_token(*current_token, "")
	self.output.space_before_token = true
}

func (self *jsbeautifier) handle_comma(current_token *tokenizer.Token) {
	if self.flags.declaration_statement {
		if self.is_expression(self.flags.parent.mode) {
			self.flags.declaration_assignment = false
		}
		self.print_token(*current_token, "")

		if self.flags.declaration_assignment {
			self.flags.declaration_assignment = false
			self.print_newline(false, true)
		} else {
			self.output.space_before_token = true
		}

		return
	}

	self.print_token(*current_token, "")

	if self.flags.mode == ObjectLiteral || (self.flags.mode == Statement && self.flags.parent.mode == ObjectLiteral) {
		if self.flags.mode == Statement {
			self.restore_mode()
		}

		self.print_newline(false, false)
	} else {
		self.output.space_before_token = true
	}
}

func (self *jsbeautifier) handle_operator(current_token *tokenizer.Token) {
	if self.start_of_statement(*current_token) {

	}

	if self.last_type == "TK_RESERVED" && self.is_special_word(self.flags.last_text) {
		self.output.space_before_token = true
		self.print_token(*current_token, "")
		return
	}

	if current_token.Text() == "*" && self.last_type == "TK_DOT" {
		self.print_token(*current_token, "")
		return
	}

	if current_token.Text() == ":" && self.flags.in_case {
		self.flags.case_body = true
		self.indent()
		self.print_token(*current_token, "")
		self.print_newline(false, false)
		self.flags.in_case = false
		return
	}

	if current_token.Text() == "::" {
		self.print_token(*current_token, "")
		return
	}

	if current_token.WantedNewLine() && (current_token.Text() == "--" || current_token.Text() == "++") {
		self.print_newline(false, true)
	}

	if self.last_type == "TK_OPERATOR" {
		self.allow_wrap_or_preserved_newline(*current_token, false)
	}

	space_before := true
	space_after := true

	if utils.InStrArray(current_token.Text(), []string{"--", "++", "!", "~"}) || ((current_token.Text() == "+" || current_token.Text() == "-") && (utils.InStrArray(self.last_type, []string{"TK_START_BLOCK", "TK_START_EXPR", "TK_EQUALS", "TK_OPERATOR"}) || utils.InStrArray(self.flags.last_text, tokenizer.GetLineStarters()) || self.flags.last_text == ",")) {
		space_after = false
		space_before = false

		if self.flags.last_text == ";" && self.is_expression(self.flags.mode) {
			space_before = true
		}

		if self.last_type == "TK_RESERVED" || self.last_type == "TK_END_EXPR" {
			space_before = true
		} else if self.last_type == "TK_OPERATOR" {
			space_before = (utils.InStrArray(current_token.Text(), []string{"--", "-"}) && utils.InStrArray(self.flags.last_text, []string{"--", "-"})) || (utils.InStrArray(current_token.Text(), []string{"++", "+"}) && utils.InStrArray(self.flags.last_text, []string{"++", "+"}))
		}

		if self.flags.mode == BlockStatement && (self.flags.last_text == "{" || self.flags.last_text == ";") {
			self.print_newline(false, false)
		}
	} else if current_token.Text() == ":" {
		if self.flags.ternary_depth == 0 {
			space_before = false
		} else {
			self.flags.ternary_depth -= 1
		}
	} else if current_token.Text() == "?" {
		self.flags.ternary_depth += 1
	} else if current_token.Text() == "*" && self.last_type == "TK_RESERVED" && self.flags.last_text == "function" {
		space_before = false
		space_after = false
	}

	if space_before {
		self.output.space_before_token = true
	}

	self.print_token(*current_token, "")

	if space_after {
		self.output.space_before_token = true
	}
}

func (self *jsbeautifier) handle_block_comment(current_token *tokenizer.Token) {
	lines := strings.Split(strings.Replace(current_token.Text(), "\x0d", "", -1), "\x0a")

	javadoc := true
	starless := true
	last_indent := current_token.WhitespaceBefore()
	last_indent_length := len(last_indent)

	self.print_newline(false, true)

	if len(lines) > 1 {
		for _, l := range lines[1:] {
			trims := strings.TrimSpace(l)
			if trims == "" || trims[0] != '*' {
				javadoc = false
				break
			}
		}

		if !javadoc {
			for _, l := range lines[1:] {
				trims := strings.TrimSpace(l)
				if trims != "" && !strings.HasPrefix(l, last_indent) {
					starless = false
					break
				}
			}
		}
	} else {
		javadoc = false
		starless = false
	}

	self.print_token(*current_token, lines[0])
	for _, l := range lines[1:] {
		self.print_newline(false, true)
		if javadoc {
			self.print_token(*current_token, " "+strings.TrimSpace(l))
		} else if starless && len(l) > last_indent_length {
			self.print_token(*current_token, l[last_indent_length:])
		} else {
			self.output.add_token(l)
		}
	}
	self.print_newline(false, true)
}

func (self *jsbeautifier) handle_inline_comment(current_token *tokenizer.Token) {
	self.output.space_before_token = true
	self.print_token(*current_token, "")
	self.output.space_before_token = true
}

func (self *jsbeautifier) handle_comment(current_token *tokenizer.Token) {
	if current_token.WantedNewLine() {
		self.print_newline(false, true)
	}
	if !current_token.WantedNewLine() {
		self.output.trim(true)
	}

	self.output.space_before_token = true
	self.print_token(*current_token, "")
	self.print_newline(false, true)
}

func (self *jsbeautifier) handle_dot(current_token *tokenizer.Token) {
	if self.start_of_statement(*current_token) {
	}

	if self.last_type == "TK_RESERVED" && self.is_special_word(self.flags.last_text) {
		self.output.space_before_token = true
	} else {
		self.allow_wrap_or_preserved_newline(*current_token, self.flags.last_text == ")" && self.options["break_chained_methods"].(bool))
	}
	self.print_token(*current_token, "")
}

func (self *jsbeautifier) handle_unknown(current_token *tokenizer.Token) {
	self.print_token(*current_token, "")
	if current_token.Text()[len(current_token.Text())-1:] == "\n" {
		self.print_newline(false, false)
	}
}

func (self *jsbeautifier) handle_eof(current_token *tokenizer.Token) {
	for self.flags.mode == Statement {
		self.restore_mode()
	}
}

func (self *jsbeautifier) blank_state(s *string) *string {
	if self.options["jslint_happy"].(bool) {
		self.options["space_after_anon_function"] = true
	}

	if self.options["indent_with_tabs"].(bool) {
		self.options["indent_char"] = "\t"
		self.options["indent_size"] = 1
	}

	self.indent_string = ""

	for i := 0; i < self.options["indent_size"].(int); i++ {
		self.indent_string += self.options["indent_char"].(string)
	}

	self.baseIndentString = ""
	self.last_type = "TK_START_BLOCK"
	self.last_last_text = ""

	preindent_index := 0

	if s != nil && len(*s) > 0 {
		for preindent_index < len(*s) && (string((*s)[preindent_index]) == " " || string((*s)[preindent_index]) == "\t") {
			self.baseIndentString += string((*s)[preindent_index])
			preindent_index++
		}
		temps := string((*s)[preindent_index:])
		s = &temps
	}

	self.output = NewOutput(self.indent_string, self.baseIndentString)

	self.set_mode(BlockStatement)
	return s
}

func DefaultOptions() map[string]interface{} {
	return default_options
}
func New(options optargs.MapType) *jsbeautifier {
	b := new(jsbeautifier)

	b.options.Copy(options)
	b.blank_state(nil)
	return b
}

func Beautify(data *string, options optargs.MapType) (string, error) {
	b := New(options)
	return b.beautify(data)
}

func BeautifyFile(file string, options optargs.MapType) *string {
	if file == "-" {
		panic("Reading stdin not implemented yet")
	}

	data, _ := ioutil.ReadFile(file)
	sdata := string(data)
	val, _ := Beautify(&sdata, options)
	return &val
}
