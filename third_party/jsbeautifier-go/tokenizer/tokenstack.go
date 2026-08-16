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

type TokenStack struct {
	tokens []Token
	size   int
}

func (self *TokenStack) Append(t Token) {
	self.tokens = append(self.tokens, t)
	self.size++
}

func (self *TokenStack) Pop() *Token {
	if self.size == 0 {
		return nil
	}
	t := self.tokens[self.size-1]
	self.tokens = self.tokens[:self.size-1]
	self.size--
	return &t
}

func (self *TokenStack) Shift() *Token {
	if self.size == 0 {
		return nil
	}

	t := self.tokens[0]
	if self.size == 1 {
		self.tokens = nil
	} else {
		self.tokens = self.tokens[1:]
	}
	self.size--
	return &t
}

func (self *TokenStack) Empty() bool {
	return self.size == 0
}
