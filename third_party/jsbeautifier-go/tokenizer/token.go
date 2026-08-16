package tokenizer

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

type Token struct {
	text                    string
	tktype                  string
	n_newlines              int
	wanted_newline          bool
	whitespace_before_token string
	comments_before         []Token
	parent                  *Token
}

func (self *Token) Type() string {
	return self.tktype
}

func (self *Token) SetType(tktype string) {
	self.tktype = tktype
}

func (self *Token) CommentsBefore() []Token {
	return self.comments_before
}

func (self *Token) Text() string {
	return self.text
}

func (self *Token) NewLines() int {
	return self.n_newlines
}

func (self *Token) WhitespaceBefore() string {
	return self.whitespace_before_token
}
func (self *Token) WantedNewLine() bool {
	return self.wanted_newline
}

func (self *Token) Print() {
	println("Token Type", self.tktype, " - Value :", self.text)
}

func NewSimpleToken(text string, tktype string, n_newlines int, whitespace_before_token string) Token {
	return Token{text, tktype, n_newlines, n_newlines > 0, whitespace_before_token, make([]Token, 0), nil}
}
