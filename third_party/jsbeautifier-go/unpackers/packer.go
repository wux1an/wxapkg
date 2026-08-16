package unpackers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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

type DotPacker struct{}

func (self *DotPacker) unpack(source *string) (*string, error) {
	payload, symtab, radix, count, _ := self.getargs(source)

	if count != len(symtab) {
		return nil, errors.New("Malformed p.a.c.k.e.r symtab")
	}
	unbaser, _ := newunbaser(radix)

	lookup := func(match string) string {
		res := symtab[unbaser.unbase(&match)]
		if res != "" {
			return res
		}

		return match
	}

	*source = regexp.MustCompile(`\b\w+\b`).ReplaceAllStringFunc(payload, lookup)

	return self.replacestrings(source), nil
}

func (self *DotPacker) replacestrings(source *string) *string {
	match := regexp.MustCompile(`(?s)var *(_\w+)\=\["(.*?)"\];`).FindStringSubmatch(*source)

	if match != nil {
		varname, string_values := match[1], match[2]

		startpoint := len(match[0])
		lookup := strings.Split(string_values, `","`)
		variable := fmt.Sprintf("%s[%%d]", varname)
		for index, value := range lookup {
			*source = strings.Replace(*source, fmt.Sprintf(variable, index), `"`+value+`"`, -1)
		}

		*source = (*source)[startpoint:]
		return source
	}
	return source
}
func (self *DotPacker) detect(source *string) bool {
	return strings.HasPrefix(strings.Replace(*source, " ", "", -1), "eval(function(p,a,c,k,e,")
}

func (self *DotPacker) getargs(source *string) (string, []string, int, int, error) {
	juicers := []*regexp.Regexp{regexp.MustCompile(`(?s)}\('(.*)', *(\d+), *(\d+), *'(.*)'\.split\('\|'\), *(\d+), *(.*)\)\)`),
		regexp.MustCompile(`(?s)}\('(.*)', *(\d+), *(\d+), *'(.*)'\.split\('\|'\)`)}

	for _, juicer := range juicers {
		args := juicer.FindStringSubmatch(*source)

		if args != nil {
			arg1, _ := strconv.Atoi(args[2])
			arg2, _ := strconv.Atoi(args[3])
			return args[1], strings.Split(args[4], "|"), arg1, arg2, nil
		} else {
			return "", []string{}, 0, 0, errors.New("Corrupted p.a.c.k.e.r data")
		}
	}
	return "", []string{}, 0, 0, errors.New("Could not make sense of p.a.c.k.e.r data (unexpected code structure)")
}

type unbaser struct {
	base   int
	unbase func(s *string) int
	dict   map[string]int
}

func newunbaser(base int) (*unbaser, error) {
	ALPHABET := map[int]string{62: `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`,
		95: ` !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ'
              '[\]^_` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`,
	}

	unbaser := new(unbaser)

	if 2 <= base && base <= 36 {
		unbaser.unbase = func(s *string) int { val, _ := strconv.ParseInt(*s, base, 0); return int(val) }
	} else {
		if _, ok := ALPHABET[base]; !ok {
			return nil, errors.New("Unsupported base encoding")
		}

		unbaser.dict = make(map[string]int)
		for index, cipher := range ALPHABET[base] {
			unbaser.dict[string(cipher)] = index
		}

		unbaser.unbase = unbaser.dictunbaser

	}

	unbaser.base = base
	return unbaser, nil
}

func (self *unbaser) dictunbaser(s *string) int {
	ret := 0
	for index, cipher := range reverse(*s) {
		ret += pow(self.base, index) * self.dict[string(cipher)]
	}

	return ret
}

func reverse(s string) string {
	o := make([]rune, utf8.RuneCountInString(s))
	i := len(o)
	for _, c := range s {
		i--
		o[i] = c
	}
	return string(o)
}

func pow(a, b int) int {
	var result int = 1

	for 0 != b {
		if 0 != (b & 1) {
			result *= a
		}
		b >>= 1
		a *= a
	}

	return result
}
