package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/flynn/go-docopt"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
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

func main() {
	usage := `jsbeautifier.
Usage:
	jsbeautifier [options] PATH
Options:
	--indent-size=SIZE  Indentation Size [default: 4]
	--indent-char=CHAR  Indentation Char [default: space]
	--brace-style=STYLE  Brace Style (expand, collapse, end-exapnd or none) [default: collapse]
	--outfile=OUTPUTFILE  Output results to a file [default: stdout]
	-z --eval-code  Eval JS code (dangerous!)
	-f --keep-function-indentation  Keep original function indentation
	-t --indent-with-tabs  Use tabs to indent
	-d --disable-preserve-newlines  Don't preserve newlines
	-P --space-in-paren  Use spaces in parenthesis
	-E --space-in-empty-paren  Use spaces in empty parenthesis
	-j --jslint-happy  Keep JSLint happy :)
	-a --space_after_anon_function  Use spaces after anonymous functions
	-x --unescape-strings  Unescape strings
	-X --e4x  Parse XML by e4x standard
	-n --end-with-newline  End output with newline
	-w --wrap-line-length  Wrap line length
	-i --stdin  Use stdin as input
	-h --help  Show this screen
	--version  Show version
`
	arguments, _ := docopt.Parse(usage, nil, true, "jsbeautifier 0.3", false)

	file := arguments.String["PATH"]
	if arguments.Bool["--stdin"] {
		file = "-"
	}

	outfile := "stdout"

	options := jsbeautifier.DefaultOptions()

	if arguments.Bool["--keep-array-indentation"] {
		options["keep_array_indentation"] = true
	}
	if arguments.Bool["--keep-function-indentation"] {
		options["keep_function_indentation"] = true
	}
	if arguments.String["--outfile"] != "" {
		outfile = arguments.String["--outfile"]
	}
	if arguments.String["--indent-size"] != "" {
		options["indent_size"], _ = strconv.ParseInt(arguments.String["--indent-size"], 0, 8)
		options["indent_size"] = int(options["indent_size"].(int64))
	}
	if arguments.String["--indent-char"] != "" {
		if arguments.String["--indent-char"] == "space" {
			options["indent_char"] = " "
		} else {
			options["indent_char"] = arguments.String["--indent_char"]
		}
	}
	if arguments.Bool["--indent-with-tabs"] {
		options["indent_with_tabs"] = true
	}
	if arguments.Bool["--disable-preserve-newlines"] {
		options["preserve_newlines"] = false
	}
	if arguments.Bool["--space-in-paren"] {
		options["space_in_paren"] = true
	}
	if arguments.Bool["--space-in-empty-paren"] {
		options["space_in_empty_paren"] = true
	}
	if arguments.Bool["--jslint-happy"] {
		options["jslint_happy"] = true
	}
	if arguments.Bool["--space_after_anon_function"] {
		options["space_after_anon_function"] = true
	}
	if arguments.Bool["--eval-code"] {
		options["eval_code"] = true
	}
	if arguments.String["--brace-style"] != "" {
		options["brace_style"] = arguments.String["--brace-style"]
	}
	if arguments.Bool["--unescape-strings"] {
		options["unescape_strings"] = true
	}
	if arguments.Bool["--e4x"] {
		options["e4x"] = true
	}
	if arguments.Bool["--end-with-newline"] {
		options["end_with_newline"] = true
	}
	if arguments.String["--wrap-line-length"] != "" {
		options["wrap_line_length"], _ = strconv.ParseInt(arguments.String["--wrap-line-length"], 0, 8)

		options["wrap_line_length"] = int(options["wrap_line_length"].(int64))
	}

	code := jsbeautifier.BeautifyFile(file, options)

	if outfile == "stdout" {
		fmt.Print(*code)
	} else {
		ioutil.WriteFile(file, []byte(*code), 0777)
	}

}
