package util

import (
	"bytes"
	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/tidwall/pretty"
	"github.com/yosssi/gohtml"
	"regexp"
	"strings"
)

var regScriptInHtml = regexp.MustCompile(`(?s) *<script.*?>(.*?)</script>`)
var jsOptions = jsbeautifier.DefaultOptions()

func PrettyJson(data []byte) []byte {
	return pretty.Pretty(data)
}

func PrettyHtml(data []byte) []byte {
	data = gohtml.FormatBytes(bytes.TrimSpace(data)) // use `TrimSpace` to remove leading whitespace
	data = regScriptInHtml.ReplaceAllFunc(data, func(script []byte) []byte {
		var space = countLeadingSpaces(script)

		var jsCode = regScriptInHtml.FindSubmatch(script)[1]
		var jsStr = strings.Repeat(" ", space+2) + string(bytes.TrimSpace(jsCode))

		beautify, err := jsbeautifier.Beautify(&jsStr, jsOptions)
		if err == nil {
			return bytes.Replace(script, jsCode, []byte("\n"+beautify+"\n"+strings.Repeat(" ", space)), 1)
		}
		return script
	})
	return data
}

func PrettyJavaScript(data []byte) []byte {
	var code = string(bytes.TrimSpace(data)) // use `TrimSpace` to remove leading whitespace
	beautify, err := jsbeautifier.Beautify(&code, jsOptions)
	if err != nil {
		return data
	}

	return []byte(beautify)
}

func countLeadingSpaces(data []byte) int {
	var result = 0
	for _, c := range data {
		if c == ' ' {
			result++
		} else {
			return result
		}
	}

	return result
}
