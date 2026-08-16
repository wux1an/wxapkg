package tokenizer

import "unicode/utf8"

func (self *tokenizer) GetNextRune() rune {
	c, _ := utf8.DecodeRuneInString((*self.input)[self.parser_pos:])
	return c
}

func (self *tokenizer) GetNextChar() string {
	c, _ := utf8.DecodeRuneInString((*self.input)[self.parser_pos:])
	return string(c)
}

func (self *tokenizer) GetNextCharWithWidth() (string, int) {
	c, width := utf8.DecodeRuneInString((*self.input)[self.parser_pos:])
	return string(c), width
}

func (self *tokenizer) AdvanceNextChar() string {
	r, width := self.GetNextCharWithWidth()
	self.parser_pos += width
	return r
}

func (self *tokenizer) GetLastRune() rune {
	c, _ := utf8.DecodeLastRuneInString((*self.input)[:self.parser_pos])
	return c
}

func (self *tokenizer) GetLastChar() string {
	c, _ := utf8.DecodeLastRuneInString((*self.input)[:self.parser_pos])
	return string(c)
}

func (self *tokenizer) GetLastCharWithWidth() (string, int) {
	c, width := utf8.DecodeLastRuneInString((*self.input)[:self.parser_pos])
	return string(c), width
}

func (self *tokenizer) BackupChar() string {
	c, width := self.GetLastCharWithWidth()
	self.parser_pos -= width
	return string(c)
}
