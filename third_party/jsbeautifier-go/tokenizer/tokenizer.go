package tokenizer

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ditashi/jsbeautifier-go/optargs"
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

type tokenizer struct {
	input                   *string
	options                 optargs.MapType
	last_token              Token
	tokens_parsed           int
	indent_string           string
	parser_pos              int
	n_newlines              int
	whitespace_before_token string
	in_html_comment         bool
	acorn                   *acorn
}

var whitespace = [4]string{"\n", "\r", "\t", " "}

var digit, _ = regexp.Compile("[0-9]")
var hexPrefix = regexp.MustCompile("[Xx]")
var hexDigit = regexp.MustCompile("[0123456789abcdefABCDEF]")
var exponentMarker = regexp.MustCompile("[Ee]")
var exponentSign = regexp.MustCompile("[+-]")
var e4xTag = regexp.MustCompile(`^<(!\[CDATA\[[\s\S]*?\]\]|[-a-zA-Z:0-9_.]+|\{[^{}]*\})\s*([-a-zA-Z:0-9_.]+=(\{[^{}]*\}|"[^"]*"|'[^']*')\s*)*\/?\s*>`)
var e4xXML = regexp.MustCompile(`<(\/?)(!\[CDATA\[[\s\S]*?\]\]|[-a-zA-Z:0-9_.]+|\{[^{}]*\})\s*([-a-zA-Z:0-9_.]+=(\{[^{}]*\}|"[^"]*"|\'[^\']*\')\s*)*(\/?)\s*>`)

var punct = strings.Split("+ - * / % & ++ -- = += -= *= /= %= == === != !== > < >= <= >> << >>> >>>= >>= <<= && &= | || ! ~ , : ? ^ ^= |= :: => + ' <?= <? ?> <%= <% %>", " ")

var line_starters = strings.Split("continue,try,throw,return,var,let,const,if,switch,case,default,for,while,break,function,yield,import,export", ",")

var reserved_words = append(line_starters, []string{"do", "in", "else", "get", "set", "new", "catch", "finally", "typeof"}...)

func (self *tokenizer) Tokenize() chan Token {
	tkch := make(chan Token)

	go func() {
		var open *Token

		comments := make([]Token, 0)
		open_stack := new(TokenStack)

		var last Token
		// while not (not last == None and last.type == 'TK_EOF'):
		for last.Type() == "" || last.tktype != "TK_EOF" {
			token_value, tktype := self.getNextToken()

			token := NewSimpleToken(token_value, tktype, self.n_newlines, self.whitespace_before_token)
			for token.tktype == "TK_INLINE_COMMENT" || token.tktype == "TK_COMMENT" || token.tktype == "TK_BLOCK_COMMENT" || token.tktype == "TK_UNKNOWN" {
				comments = append(comments, token)
				token_value, tktype := self.getNextToken()
				token = NewSimpleToken(token_value, tktype, self.n_newlines, self.whitespace_before_token)
			}
			if len(comments) > 0 {
				token.comments_before = comments
				comments = make([]Token, 0)
			}

			if token.tktype == "TK_START_BLOCK" || token.tktype == "TK_START_EXPR" {

				token.parent = new(Token)
				*token.parent = last

				open = new(Token)
				*open = token

				open_stack.Append(token)
			} else if (token.tktype == "TK_END_BLOCK" || token.tktype == "TK_END_EXPR") && (open != nil && ((token.text == "]" && open.text == "[") ||
				(token.text == ")" && open.text == "(") ||
				(token.text == "}" && open.text == "{"))) {

				token.parent = open.parent
				open = open_stack.Pop()
			}
			tkch <- token
			last = token
			self.last_token = token
		}
		close(tkch)
	}()
	return tkch
}

func (self *tokenizer) GetCharSlice(backoffset int, nextoffset int) (string, bool) {
	back_pos, next_pos := 0, 0
	for i := 0; i < backoffset; i++ {
		_, width := utf8.DecodeLastRuneInString((*self.input)[:self.parser_pos])
		back_pos += width
	}
	for i := 0; i < nextoffset; i++ {
		_, width := utf8.DecodeRuneInString((*self.input)[self.parser_pos:])
		next_pos += width
	}
	if self.parser_pos+next_pos > len(*self.input) || self.parser_pos-back_pos < 0 {
		return "", true
	}
	return (*self.input)[self.parser_pos-back_pos : self.parser_pos+next_pos], false
}

func (self *tokenizer) getNextToken() (string, string) {
	defer func() { self.tokens_parsed++ }()

	whitespace_on_this_line := make([]string, 0)
	self.n_newlines = 0
	self.whitespace_before_token = ""

	if self.parser_pos >= len(*self.input) {
		return "", "TK_EOF"
	}

	if self.tokens_parsed == 0 {
		self.last_token = NewSimpleToken("{", "TK_START_BLOCK", self.n_newlines, self.whitespace_before_token)
	}

	c := self.AdvanceNextChar()

	for utils.InStrArray(c, whitespace[:]) {
		if c == "\n" {
			self.n_newlines += 1
			// whitespace_on_this_line = []
		} else if c == self.indent_string {
			whitespace_on_this_line = append(whitespace_on_this_line, self.indent_string)
		} else if c != "\r" {
			whitespace_on_this_line = append(whitespace_on_this_line, " ")
		}

		if self.parser_pos >= len(*self.input) {
			return "", "TK_EOF"
		}

		c = self.AdvanceNextChar()
	}

	if len(whitespace_on_this_line) != 0 {
		self.whitespace_before_token = strings.Join(whitespace_on_this_line, "")
	}

	if digit.Match([]byte(c)) {

		allow_decimal := true
		allow_e := true
		local_digit := digit

		if c == "0" && self.parser_pos < len(*self.input) && hexPrefix.Match([]byte(self.GetNextChar())) {
			allow_decimal = false
			allow_e = false
			c += self.AdvanceNextChar()
			local_digit = hexDigit
		} else {
			c = ""
			self.BackupChar()
		}

		for self.parser_pos < len(*self.input) && local_digit.Match([]byte(self.GetNextChar())) {
			c += self.AdvanceNextChar()

			if allow_decimal && self.parser_pos < len(*self.input) && self.GetNextChar() == "." {

				c += self.AdvanceNextChar()
				allow_decimal = false
			}

			if allow_e && self.parser_pos < len(*self.input) && exponentMarker.Match([]byte(self.GetNextChar())) {
				c += self.AdvanceNextChar()

				if self.parser_pos < len(*self.input) && exponentSign.Match([]byte(self.GetNextChar())) {
					c += self.AdvanceNextChar()
				}

				allow_e = false
				allow_decimal = false
			}
		}

		return c, "TK_WORD"
	}

	if self.acorn.IsIdentifierStart(self.GetLastRune()) {
		if self.parser_pos < len(*self.input) {
			for self.acorn.IsIdentifierChar(self.GetNextRune()) {
				c += self.AdvanceNextChar()
				if self.parser_pos == len(*self.input) {
					break
				}
			}
		}

		if !(self.last_token.tktype == "TK_DOT" || (self.last_token.tktype == "TK_RESERVED" && (self.last_token.text == "set" || self.last_token.text == "get"))) && utils.InStrArray(c, reserved_words) {
			if c == "in" {
				return c, "TK_OPERATOR"
			}

			return c, "TK_RESERVED"
		}

		return c, "TK_WORD"
	}

	if c == "(" || c == "[" {
		return c, "TK_START_EXPR"
	}

	if c == ")" || c == "]" {
		return c, "TK_END_EXPR"
	}

	if c == "{" {
		return c, "TK_START_BLOCK"
	}

	if c == "}" {
		return c, "TK_END_BLOCK"
	}

	if c == ";" {
		return c, "TK_SEMICOLON"
	}

	if c == "/" {
		comment := ""
		inline_comment := true
		nextCh, nextWidth := self.GetNextCharWithWidth()
		if nextCh == "*" {
			self.parser_pos += nextWidth

			nextCh, nextWidth = self.GetNextCharWithWidth()
			if self.parser_pos < len(*self.input) {
				for !(nextCh == "*" && self.parser_pos+nextWidth < len(*self.input) && string((*self.input)[self.parser_pos+nextWidth]) == "/") && self.parser_pos < len(*self.input) {
					c = self.AdvanceNextChar()
					comment += c
					if c == "\r" || c == "\n" {
						inline_comment = false
					}

					if self.parser_pos >= len(*self.input) {
						break
					}
					nextCh, nextWidth = self.GetNextCharWithWidth()
				}
			}
			self.parser_pos += 2

			if inline_comment && self.n_newlines == 0 {
				return "/*" + comment + "*/", "TK_INLINE_COMMENT"
			} else {
				return "/*" + comment + "*/", "TK_BLOCK_COMMENT"
			}
		}

		if self.GetNextChar() == "/" {
			comment = c
			for self.GetNextChar() != "\r" && self.GetNextChar() != "\n" {
				comment += self.AdvanceNextChar()
				if self.parser_pos >= len(*self.input) {
					break
				}
			}

			return comment, "TK_COMMENT"
		}
	}

	tempr := e4xTag

	if c == "`" || c == "'" || c == "\"" || ((c == "/") || (self.options["e4x"].(bool) && c == "<" && tempr.Match([]byte(string((*self.input)[self.parser_pos-1:]))))) && ((self.last_token.tktype == "TK_RESERVED" && utils.InStrArray(self.last_token.text, []string{"return", "case", "throw", "else", "do", "typeof", "yield"})) || (self.last_token.tktype == "TK_END_EXPR" && self.last_token.text == ")" && self.last_token.parent != nil && self.last_token.parent.tktype == "TK_RESERVED" && utils.InStrArray(self.last_token.parent.text, []string{"if", "while", "for"})) || (utils.InStrArray(self.last_token.tktype, []string{"TK_COMMENT", "TK_START_EXPR", "TK_START_BLOCK", "TK_END_BLOCK", "TK_OPERATOR", "TK_EQUALS", "TK_EOF", "TK_SEMICOLON", "TK_COMMA"}))) {

		sep := c
		esc := false
		esc1 := 0
		esc2 := 0
		resulting_string := c
		in_char_class := false

		if sep == "/" { //regexp
			in_char_class = false
			for self.parser_pos < len(*self.input) && (esc || in_char_class || self.GetNextChar() != sep) && !self.acorn.GetNewline().Match([]byte(self.GetNextChar())) {
				resulting_string += self.AdvanceNextChar()
				if !esc {
					esc = self.GetLastChar() == "\\"
					if self.GetLastChar() == "[" {
						in_char_class = true
					} else if self.GetLastChar() == "]" {
						in_char_class = false
					}
				} else {
					esc = false
				}

			}
		} else if self.options["e4x"].(bool) && sep == "<" { // xml
			xmlRegExp := e4xXML
			xmlStr := (*self.input)[self.parser_pos-1:]
			match := xmlRegExp.FindStringSubmatch((xmlStr))

			if match != nil {
				rootTag := match[2]
				depth := 0
				lastPos := 0
				for match != nil {
					isEndTag := match[1]
					tagName := match[2]
					isSingeltonTag := (match[len(match)-1] != "") || (strings.HasPrefix(match[2], "![CDATA["))
					if tagName == rootTag && !isSingeltonTag {
						if isEndTag != "" {
							depth -= 1
						} else {
							depth += 1
						}
					}

					if depth <= 0 {
						break
					}

					indices := xmlRegExp.FindStringSubmatchIndex(xmlStr[lastPos:])
					lastPos += indices[1]

					match = xmlRegExp.FindStringSubmatch(xmlStr[lastPos:])
				}
				xmlLength := 0
				if match != nil {
					indices := xmlRegExp.FindStringSubmatchIndex(xmlStr[lastPos:])
					xmlLength = lastPos + indices[1]
				} else {
					xmlLength = len(xmlStr)
				}

				self.parser_pos += xmlLength - 1

				return xmlStr[:xmlLength], "TK_STRING"
			}

		} else { // string
			for self.parser_pos < len(*self.input) && (esc || (self.GetNextChar() != sep && (sep == "`" || !self.acorn.GetNewline().Match([]byte(self.GetNextChar()))))) {
				resulting_string += self.AdvanceNextChar()
				if esc1 > 0 && esc1 >= esc2 {
					esc1, ok := strconv.ParseUint(resulting_string[esc2:], 16, 0)

					if ok != nil {
						esc1 = 0
					}
					if esc1 >= 0x20 && esc1 <= 0x7e {
						esc1c := string(esc1)
						resulting_string = resulting_string[:len(resulting_string)-2-esc2]
						if esc1c == sep || esc1c == "\\" {
							resulting_string += "\\"
						}
						resulting_string += esc1c
					}
					esc1 = 0
				}
				if esc1 > 0 {
					esc1 += 1
				} else if !esc {
					esc = self.GetLastChar() == "\\"
				} else {
					esc = false
					if self.options["unescape_strings"].(bool) {
						if self.GetLastChar() == "x" {
							esc1 += 1
							esc2 = 2
						} else if self.GetLastChar() == "u" {
							esc1 += 1
							esc2 = 4
						}
					}
				}
			}
		}

		if self.parser_pos < len(*self.input) && self.GetNextChar() == sep {
			resulting_string += self.AdvanceNextChar()

			if sep == "/" {
				for self.parser_pos < len(*self.input) && self.acorn.IsIdentifierStart(self.GetNextRune()) {
					resulting_string += self.AdvanceNextChar()
				}
			}
		}
		return resulting_string, "TK_STRING"
	}

	if c == "#" {

		if self.tokens_parsed == 0 && len(*self.input) > self.parser_pos && self.GetNextChar() == "!" {
			resulting_string := c
			for self.parser_pos < len(*self.input) && c != "\n" {
				c = self.AdvanceNextChar()
				resulting_string += c
			}
			return strings.TrimSpace(resulting_string) + "\n", "TK_UNKNOWN"
		}

		sharp := "#"
		if match := digit.Match([]byte(self.GetNextChar())); self.parser_pos < len(*self.input) && match {
			for {
				c = self.AdvanceNextChar()
				sharp += c
				if self.parser_pos >= len(*self.input) || c == "#" || c == "=" {
					break
				}
			}
		}

		nextCh, nextWidth := self.GetNextCharWithWidth()
		if c == "#" || self.parser_pos >= len(*self.input) {

		} else if nextCh == "[" && string((*self.input)[self.parser_pos+nextWidth]) == "]" {
			sharp += "[]"
			self.parser_pos += 2
		} else if nextCh == "{" && string((*self.input)[self.parser_pos+nextWidth]) == "}" {
			sharp += "{}"
			self.parser_pos += 2
		}
		return sharp, "TK_WORD"
	}

	slice, out_of_bounds := self.GetCharSlice(1, 3)

	if c == "<" && !out_of_bounds && slice == "<!--" {
		for self.parser_pos < len(*self.input) && self.GetNextChar() != "\n" {
			c += self.AdvanceNextChar()
		}
		self.in_html_comment = true
		return c, "TK_COMMENT"
	}

	slice, out_of_bounds = self.GetCharSlice(1, 2)

	if c == "-" && self.in_html_comment && !out_of_bounds && slice == "-->" {
		self.in_html_comment = false
		self.parser_pos += 2
		return "-->", "TK_COMMENT"
	}

	if c == "." {
		return c, "TK_DOT"
	}

	if utils.InStrArray(c, punct) {
		for self.parser_pos < len(*self.input) && utils.InStrArray(c+self.GetNextChar(), punct) {
			c += self.AdvanceNextChar()
			if self.parser_pos >= len(*self.input) {
				break
			}
		}

		if c == "," {
			return c, "TK_COMMA"
		}

		if c == "=" {
			return c, "TK_EQUALS"
		}
		return c, "TK_OPERATOR"
	}

	return c, "TK_UNKNOWN"
}

func GetLineStarters() []string {
	return line_starters
}

func New(s *string, options optargs.MapType, indent_string string) *tokenizer {
	t := new(tokenizer)
	t.input = s
	t.options = options
	t.indent_string = indent_string
	t.acorn = NewAcorn()
	return t
}
