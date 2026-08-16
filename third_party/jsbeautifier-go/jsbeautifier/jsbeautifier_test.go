package jsbeautifier

import (
	"testing"

	"github.com/ditashi/jsbeautifier-go/optargs"
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

var ts *testing.T

var test_options optargs.MapType

func TestBeautifier(t *testing.T) {
	ts = t

	test_options = default_options

	test_options["indent_size"] = 4
	test_options["indent_char"] = " "
	test_options["preserve_newlines"] = true
	test_options["jslint_happy"] = false
	test_options["keep_array_indentation"] = false
	test_options["brace_style"] = "collapse"

	// Unicode Support
	test("var \u0CA0_\u0CA0 = \"hi\";")
	test("var \u00e4 x = {\n    \u00e4rgerlich: true\n};")

	// End With Newline - (eof = "\n")
	test_options["end_with_newline"] = true
	test("", "\n")
	test("   return .5", "   return .5\n")
	test("   \n\nreturn .5\n\n\n\n", "   return .5\n")
	test("\n")

	// End With Newline - (eof = "")
	test_options["end_with_newline"] = false
	test("")
	test("   return .5")
	test("   \n\nreturn .5\n\n\n\n", "   return .5")
	test("\n", "")

	// New Test Suite

	// Old tests
	test("")
	test("   return .5")
	test("   return .5;\n   a();")
	test("    return .5;\n    a();")
	test("     return .5;\n     a();")
	test("   < div")
	test("a        =          1", "a = 1")
	test("a=1", "a = 1")
	test("(3) / 2")
	test("[\"a\", \"b\"].join(\"\")")
	test("a();\n\nb();")
	test("var a = 1 var b = 2", "var a = 1\nvar b = 2")
	test("var a=1, b=c[d], e=6;", "var a = 1,\n    b = c[d],\n    e = 6;")
	test("var a,\n    b,\n    c;")
	test("let a = 1 let b = 2", "let a = 1\nlet b = 2")
	test("let a=1, b=c[d], e=6;", "let a = 1,\n    b = c[d],\n    e = 6;")
	test("let a,\n    b,\n    c;")
	test("const a = 1 const b = 2", "const a = 1\nconst b = 2")
	test("const a=1, b=c[d], e=6;", "const a = 1,\n    b = c[d],\n    e = 6;")
	test("const a,\n    b,\n    c;")
	test("a = \" 12345 \"")
	test("a = ' 12345 '")
	test("if (a == 1) b = 2;")
	test("if(1){2}else{3}", "if (1) {\n    2\n} else {\n    3\n}")
	test("if(1||2);", "if (1 || 2);")
	test("(a==1)||(b==2)", "(a == 1) || (b == 2)")
	test("var a = 1 if (2) 3;", "var a = 1\nif (2) 3;")
	test("a = a + 1")
	test("a = a == 1")
	test("/12345[^678]*9+/.match(a)")
	test("a /= 5")
	test("a = 0.5 * 3")
	test("a *= 10.55")
	test("a < .5")
	test("a <= .5")
	test("a<.5", "a < .5")
	test("a<=.5", "a <= .5")
	test("a = 0xff;")
	test("a=0xff+4", "a = 0xff + 4")
	test("a = [1, 2, 3, 4]")
	test("F*(g/=f)*g+b", "F * (g /= f) * g + b")
	test("a.b({c:d})", "a.b({\n    c: d\n})")
	test("a.b\n(\n{\nc:\nd\n}\n)", "a.b({\n    c: d\n})")
	test("a.b({c:\"d\"})", "a.b({\n    c: \"d\"\n})")
	test("a.b\n(\n{\nc:\n\"d\"\n}\n)", "a.b({\n    c: \"d\"\n})")
	test("a=!b", "a = !b")
	test("a=!!b", "a = !!b")
	test("a?b:c", "a ? b : c")
	test("a?1:2", "a ? 1 : 2")
	test("a?(b):c", "a ? (b) : c")
	test("x={a:1,b:w==\"foo\"?x:y,c:z}", "x = {\n    a: 1,\n    b: w == \"foo\" ? x : y,\n    c: z\n}")
	test("x=a?b?c?d:e:f:g;", "x = a ? b ? c ? d : e : f : g;")
	test("x=a?b?c?d:{e1:1,e2:2}:f:g;", "x = a ? b ? c ? d : {\n    e1: 1,\n    e2: 2\n} : f : g;")
	test("function void(void) {}")
	test("if(!a)foo();", "if (!a) foo();")
	test("a=~a", "a = ~a")
	test("a;/*comment*/b;", "a; /*comment*/\nb;")
	test("a;/* comment */b;", "a; /* comment */\nb;")

	// simple comments don't get touched at all
	test("a;/*\ncomment\n*/b;", "a;\n/*\ncomment\n*/\nb;")
	test("a;/**\n* javadoc\n*/b;", "a;\n/**\n * javadoc\n */\nb;")
	test("a;/**\n\nno javadoc\n*/b;", "a;\n/**\n\nno javadoc\n*/\nb;")

	// comment blocks detected and reindented even w/o javadoc starter
	test("a;/*\n* javadoc\n*/b;", "a;\n/*\n * javadoc\n */\nb;")
	test("if(a)break;", "if (a) break;")
	test("if(a){break}", "if (a) {\n    break\n}")
	test("if((a))foo();", "if ((a)) foo();")
	test("for(var i=0;;) a", "for (var i = 0;;) a")
	test("for(var i=0;;)\na", "for (var i = 0;;)\n    a")
	test("a++;")
	test("for(;;i++)a()", "for (;; i++) a()")
	test("for(;;i++)\na()", "for (;; i++)\n    a()")
	test("for(;;++i)a", "for (;; ++i) a")
	test("return(1)", "return (1)")
	test("try{a();}catch(b){c();}finally{d();}", "try {\n    a();\n} catch (b) {\n    c();\n} finally {\n    d();\n}")

	//  magic function call
	test("(xx)()")

	// another magic function call
	test("a[1]()")
	test("if(a){b();}else if(c) foo();", "if (a) {\n    b();\n} else if (c) foo();")
	test("switch(x) {case 0: case 1: a(); break; default: break}", "switch (x) {\n    case 0:\n    case 1:\n        a();\n        break;\n    default:\n        break\n}")
	test("switch(x){case -1:break;case !y:break;}", "switch (x) {\n    case -1:\n        break;\n    case !y:\n        break;\n}")
	test("a !== b")
	test("if (a) b(); else c();", "if (a) b();\nelse c();")

	// typical greasemonkey start
	test("// comment\n(function something() {})")

	// duplicating newlines
	test("{\n\n    x();\n\n}")
	test("if (a in b) foo();")
	test("if(X)if(Y)a();else b();else c();", "if (X)\n    if (Y) a();\n    else b();\nelse c();")
	test("if (foo) bar();\nelse break")
	test("var a, b;")
	test("var a = new function();")
	test("new function")
	test("var a, b")
	test("{a:1, b:2}", "{\n    a: 1,\n    b: 2\n}")
	test("a={1:[-1],2:[+1]}", "a = {\n    1: [-1],\n    2: [+1]\n}")
	test("var l = {'a':'1', 'b':'2'}", "var l = {\n    'a': '1',\n    'b': '2'\n}")
	test("if (template.user[n] in bk) foo();")
	test("return 45")
	test("return this.prevObject ||\n\n    this.constructor(null);")
	test("If[1]")
	test("Then[1]")
	test("a = 1e10")
	test("a = 1.3e10")
	test("a = 1.3e-10")
	test("a = -1.3e-10")
	test("a = 1e-10")
	test("a = e - 10")
	test("a = 11-10", "a = 11 - 10")
	test("a = 1;// comment", "a = 1; // comment")
	test("a = 1; // comment")
	test("a = 1;\n // comment", "a = 1;\n// comment")
	test("a = [-1, -1, -1]")

	// The exact formatting these should have is open for discussion, but they are at least reasonable
	test("a = [ // comment\n    -1, -1, -1\n]")
	test("var a = [ // comment\n    -1, -1, -1\n]")
	test("a = [ // comment\n    -1, // comment\n    -1, -1\n]")
	test("var a = [ // comment\n    -1, // comment\n    -1, -1\n]")
	test("o = [{a:b},{c:d}]", "o = [{\n    a: b\n}, {\n    c: d\n}]")

	// was: extra space appended
	test("if (a) {\n    do();\n}")

	// if/else statement with empty body
	test("if (a) {\n// comment\n}else{\n// comment\n}", "if (a) {\n    // comment\n} else {\n    // comment\n}")

	// multiple comments indentation
	test("if (a) {\n// comment\n// comment\n}", "if (a) {\n    // comment\n    // comment\n}")
	test("if (a) b() else c();", "if (a) b()\nelse c();")
	test("if (a) b() else if c() d();", "if (a) b()\nelse if c() d();")
	test("{}")
	test("{\n\n}")
	test("do { a(); } while ( 1 );", "do {\n    a();\n} while (1);")
	test("do {} while (1);")
	test("do {\n} while (1);", "do {} while (1);")
	test("do {\n\n} while (1);")
	test("var a = x(a, b, c)")
	test("delete x if (a) b();", "delete x\nif (a) b();")
	test("delete x[x] if (a) b();", "delete x[x]\nif (a) b();")
	test("for(var a=1,b=2)d", "for (var a = 1, b = 2) d")
	test("for(var a=1,b=2,c=3) d", "for (var a = 1, b = 2, c = 3) d")
	test("for(var a=1,b=2,c=3;d<3;d++)\ne", "for (var a = 1, b = 2, c = 3; d < 3; d++)\n    e")
	test("function x(){(a||b).c()}", "function x() {\n    (a || b).c()\n}")
	test("function x(){return - 1}", "function x() {\n    return -1\n}")
	test("function x(){return ! a}", "function x() {\n    return !a\n}")
	test("x => x")
	test("(x) => x")
	test("x => { x }", "x => {\n    x\n}")
	test("(x) => { x }", "(x) => {\n    x\n}")

	// a common snippet in jQuery plugins
	test("settings = $.extend({},defaults,settings);", "settings = $.extend({}, defaults, settings);")
	test("$http().then().finally().default()")
	test("$http()\n.then()\n.finally()\n.default()", "$http()\n    .then()\n    .finally()\n    .default()")
	test("$http().when.in.new.catch().throw()")
	test("$http()\n.when\n.in\n.new\n.catch()\n.throw()", "$http()\n    .when\n    .in\n    .new\n    .catch()\n    .throw()")
	test("{xxx;}()", "{\n    xxx;\n}()")
	test("a = 'a'\nb = 'b'")
	test("a = /reg/exp")
	test("a = /reg/")
	test("/abc/.test()")
	test("/abc/i.test()")
	test("{/abc/i.test()}", "{\n    /abc/i.test()\n}")
	test("var x=(a)/a;", "var x = (a) / a;")
	test("x != -1")
	test("for (; s-->0;)t", "for (; s-- > 0;) t")
	test("for (; s++>0;)u", "for (; s++ > 0;) u")
	test("a = s++>s--;", "a = s++ > s--;")
	test("a = s++>--s;", "a = s++ > --s;")
	test("{x=#1=[]}", "{\n    x = #1=[]\n}")
	test("{a:#1={}}", "{\n    a: #1={}\n}")
	test("{a:#1#}", "{\n    a: #1#\n}")
	test("\"incomplete-string")
	test("'incomplete-string")
	test("/incomplete-regex")
	test("`incomplete-template-string")
	test("{a:1},{a:2}", "{\n    a: 1\n}, {\n    a: 2\n}")
	test("var ary=[{a:1}, {a:2}];", "var ary = [{\n    a: 1\n}, {\n    a: 2\n}];")
	// incomplete
	test("{a:#1", "{\n    a: #1")

	// incomplete
	test("{a:#", "{\n    a: #")

	// incomplete
	test("}}}", "}\n}\n}")
	test("<!--\nvoid();\n// -->")

	// incomplete regexp
	test("a=/regexp", "a = /regexp")
	test("{a:#1=[],b:#1#,c:#999999#}", "{\n    a: #1=[],\n    b: #1#,\n    c: #999999#\n}")
	test("a = 1e+2")
	test("a = 1e-2")
	test("do{x()}while(a>1)", "do {\n    x()\n} while (a > 1)")
	test("x(); /reg/exp.match(something)", "x();\n/reg/exp.match(something)")
	test("something();(", "something();\n(")
	test("#!she/bangs, she bangs\nf=1", "#!she/bangs, she bangs\n\nf = 1")
	test("#!she/bangs, she bangs\n\nf=1", "#!she/bangs, she bangs\n\nf = 1")
	test("#!she/bangs, she bangs\n\n/* comment */")
	test("#!she/bangs, she bangs\n\n\n/* comment */")
	test("#")
	test("#!")
	test("function namespace::something()")
	test("<!--\nsomething();\n-->")
	test("<!--\nif(i<0){bla();}\n-->", "<!--\nif (i < 0) {\n    bla();\n}\n-->")
	test("{foo();--bar;}", "{\n    foo();\n    --bar;\n}")
	test("{foo();++bar;}", "{\n    foo();\n    ++bar;\n}")
	test("{--bar;}", "{\n    --bar;\n}")
	test("{++bar;}", "{\n    ++bar;\n}")
	test("if(true)++a;", "if (true) ++a;")
	test("if(true)\n++a;", "if (true)\n    ++a;")
	test("if(true)--a;", "if (true) --a;")
	test("if(true)\n--a;", "if (true)\n    --a;")

	// Handling of newlines around unary ++ and -- operators
	test("{foo\n++bar;}", "{\n    foo\n    ++bar;\n}")
	test("{foo++\nbar;}", "{\n    foo++\n    bar;\n}")

	// This is invalid, but harder to guard against. Issue #203.
	test("{foo\n++\nbar;}", "{\n    foo\n    ++\n    bar;\n}")

	// regexps
	test("a(/abc\\/\\/def/);b()", "a(/abc\\/\\/def/);\nb()")
	test("a(/a[b\\[\\]c]d/);b()", "a(/a[b\\[\\]c]d/);\nb()")

	// incomplete char class
	test("a(/a[b\\[")

	// allow unescaped / in char classes
	test("a(/[a/b]/);b()", "a(/[a/b]/);\nb()")
	test("typeof /foo\\//;")
	test("yield /foo\\//;")
	test("throw /foo\\//;")
	test("do /foo\\//;")
	test("return /foo\\//;")
	test("switch (a) {\n    case /foo\\//:\n        b\n}")
	test("if (a) /foo\\//\nelse /foo\\//;")
	test("if (foo) /regex/.test();")
	test("function foo() {\n    return [\n        \"one\",\n        \"two\"\n    ];\n}")
	test("a=[[1,2],[4,5],[7,8]]", "a = [\n    [1, 2],\n    [4, 5],\n    [7, 8]\n]")
	test("a=[[1,2],[4,5],function(){},[7,8]]", "a = [\n    [1, 2],\n    [4, 5],\n    function() {},\n    [7, 8]\n]")
	test("a=[[1,2],[4,5],function(){},function(){},[7,8]]", "a = [\n    [1, 2],\n    [4, 5],\n    function() {},\n    function() {},\n    [7, 8]\n]")
	test("a=[[1,2],[4,5],function(){},[7,8]]", "a = [\n    [1, 2],\n    [4, 5],\n    function() {},\n    [7, 8]\n]")
	test("a=[b,c,function(){},function(){},d]", "a = [b, c, function() {}, function() {}, d]")
	test("a=[b,c,\nfunction(){},function(){},d]", "a = [b, c,\n    function() {},\n    function() {},\n    d\n]")
	test("a=[a[1],b[4],c[d[7]]]", "a = [a[1], b[4], c[d[7]]]")
	test("[1,2,[3,4,[5,6],7],8]", "[1, 2, [3, 4, [5, 6], 7], 8]")

	test("[[[\"1\",\"2\"],[\"3\",\"4\"]],[[\"5\",\"6\",\"7\"],[\"8\",\"9\",\"0\"]],[[\"1\",\"2\",\"3\"],[\"4\",\"5\",\"6\",\"7\"],[\"8\",\"9\",\"0\"]]]", "[\n    [\n        [\"1\", \"2\"],\n        [\"3\", \"4\"]\n    ],\n    [\n        [\"5\", \"6\", \"7\"],\n        [\"8\", \"9\", \"0\"]\n    ],\n    [\n        [\"1\", \"2\", \"3\"],\n        [\"4\", \"5\", \"6\", \"7\"],\n        [\"8\", \"9\", \"0\"]\n    ]\n]")
	test("{[x()[0]];indent;}", "{\n    [x()[0]];\n    indent;\n}")

	test("{{}/z/}", "{\n    {}\n    /z/\n}")
	test("return ++i", "return ++i")
	test("return !!x", "return !!x")
	test("return !x", "return !x")
	test("return [1,2]", "return [1, 2]")
	test("return;", "return;")
	test("return\nfunc", "return\nfunc")
	test("catch(e)", "catch (e)")
	test("yield [1, 2]")

	test("var a=1,b={foo:2,bar:3},{baz:4,wham:5},c=4;",
		"var a = 1,\n    b = {\n        foo: 2,\n        bar: 3\n    },\n    {\n        baz: 4,\n        wham: 5\n    }, c = 4;")
	test("var a=1,b={foo:2,bar:3},{baz:4,wham:5},\nc=4;",
		"var a = 1,\n    b = {\n        foo: 2,\n        bar: 3\n    },\n    {\n        baz: 4,\n        wham: 5\n    },\n    c = 4;")

	// inline comment
	test("function x(/*int*/ start, /*string*/ foo)", "function x( /*int*/ start, /*string*/ foo)")

	// javadoc comment
	test("/**\n* foo\n*/", "/**\n * foo\n */")
	test("{\n/**\n* foo\n*/\n}", "{\n    /**\n     * foo\n     */\n}")

	// starless block comment
	test("/**\nfoo\n*/")
	test("/**\nfoo\n**/")
	test("/**\nfoo\nbar\n**/")
	test("/**\nfoo\n\nbar\n**/")
	test("/**\nfoo\n    bar\n**/")
	test("{\n/**\nfoo\n*/\n}", "{\n    /**\n    foo\n    */\n}")
	test("{\n/**\nfoo\n**/\n}", "{\n    /**\n    foo\n    **/\n}")
	test("{\n/**\nfoo\nbar\n**/\n}", "{\n    /**\n    foo\n    bar\n    **/\n}")
	test("{\n/**\nfoo\n\nbar\n**/\n}", "{\n    /**\n    foo\n\n    bar\n    **/\n}")
	test("{\n/**\nfoo\n    bar\n**/\n}", "{\n    /**\n    foo\n        bar\n    **/\n}")
	test("{\n    /**\n    foo\nbar\n    **/\n}")

	test("var a,b,c=1,d,e,f=2;", "var a, b, c = 1,\n    d, e, f = 2;")
	test("var a,b,c=[],d,e,f=2;", "var a, b, c = [],\n    d, e, f = 2;")
	test("function() {\n    var a, b, c, d, e = [],\n        f;\n}")

	test("do/regexp/;\nwhile(1);", "do /regexp/;\nwhile (1);")

	test("var a = a,\na;\nb = {\nb\n}", "var a = a,\n    a;\nb = {\n    b\n}")

	test("var a = a,\n    /* c */\n    b;")
	test("var a = a,\n    // c\n    b;")

	test("foo.(\"bar\");")

	test("if (a) a()\nelse b()\nnewline()")
	test("if (a) a()\nnewline()")
	test("a=typeof(x)", "a = typeof(x)")

	test("var a = function() {\n        return null;\n    },\n    b = false;")

	test("var a = function() {\n    func1()\n}")
	test("var a = function() {\n    func1()\n}\nvar b = function() {\n    func2()\n}")

	// Code with and without semicolons
	test("var whatever = require(\"whatever\");\nfunction() {\n    a = 6;\n}",
		"var whatever = require(\"whatever\");\n\nfunction() {\n    a = 6;\n}")
	test("var whatever = require(\"whatever\")\nfunction() {\n    a = 6\n}", "var whatever = require(\"whatever\")\n\nfunction() {\n    a = 6\n}")

	test_options["space_after_anon_function"] = true

	test("switch(x) {case 0: case 1: a(); break; default: break}",
		"switch (x) {\n    case 0:\n    case 1:\n        a();\n        break;\n    default:\n        break\n}")
	test("switch(x){case -1:break;case !y:break;}", "switch (x) {\n    case -1:\n        break;\n    case !y:\n        break;\n}")
	test("a=typeof(x)", "a = typeof (x)")
	test("x();\n\nfunction(){}", "x();\n\nfunction () {}")
	test("function () {\n    var a, b, c, d, e = [],\n        f;\n}")
	test("// comment 1\n(function()", "// comment 1\n(function ()")
	test("var o1=$.extend(a);function(){alert(x);}", "var o1 = $.extend(a);\n\nfunction () {\n    alert(x);\n}")
	test("function* () {\n    yield 1;\n}")

	test_options["space_after_anon_function"] = false

	test_options["jslint_happy"] = true

	test("x();\n\nfunction(){}", "x();\n\nfunction () {}")
	test("function () {\n    var a, b, c, d, e = [],\n        f;\n}")
	test("switch(x) {case 0: case 1: a(); break; default: break}", "switch (x) {\ncase 0:\ncase 1:\n    a();\n    break;\ndefault:\n    break\n}")
	test("switch(x){case -1:break;case !y:break;}", "switch (x) {\ncase -1:\n    break;\ncase !y:\n    break;\n}")
	test("// comment 1\n(function()", "// comment 1\n(function ()")
	test("var o1=$.extend(a);function(){alert(x);}", "var o1 = $.extend(a);\n\nfunction () {\n    alert(x);\n}")
	test("a=typeof(x)", "a = typeof (x)")
	test("function* () {\n    yield 1;\n}")

	test_options["jslint_happy"] = false
	test("switch(x) {case 0: case 1: a(); break; default: break}", "switch (x) {\n    case 0:\n    case 1:\n        a();\n        break;\n    default:\n        break\n}")
	test("switch(x){case -1:break;case !y:break;}", "switch (x) {\n    case -1:\n        break;\n    case !y:\n        break;\n}")

	test("// comment 2\n(function()", "// comment 2\n(function()")

	test("var a2, b2, c2, d2 = 0, c = function() {}, d = '';", "var a2, b2, c2, d2 = 0,\n    c = function() {},\n    d = '';")
	test("var a2, b2, c2, d2 = 0, c = function() {},\nd = '';", "var a2, b2, c2, d2 = 0,\n    c = function() {},\n    d = '';")
	test("var o2=$.extend(a);function(){alert(x);}", "var o2 = $.extend(a);\n\nfunction() {\n    alert(x);\n}")
	test("function*() {\n    yield 1;\n}")
	test("function* x() {\n    yield 1;\n}")

	test("{\"x\":[{\"a\":1,\"b\":3},\n7,8,8,8,8,{\"b\":99},{\"a\":11}]}", "{\n    \"x\": [{\n            \"a\": 1,\n            \"b\": 3\n        },\n        7, 8, 8, 8, 8, {\n            \"b\": 99\n        }, {\n            \"a\": 11\n        }\n    ]\n}")
	test("{\"x\":[{\"a\":1,\"b\":3},7,8,8,8,8,{\"b\":99},{\"a\":11}]}", "{\n    \"x\": [{\n        \"a\": 1,\n        \"b\": 3\n    }, 7, 8, 8, 8, 8, {\n        \"b\": 99\n    }, {\n        \"a\": 11\n    }]\n}")

	test("{\"1\":{\"1a\":\"1b\"},\"2\"}", "{\n    \"1\": {\n        \"1a\": \"1b\"\n    },\n    \"2\"\n}")
	test("{a:{a:b},c}", "{\n    a: {\n        a: b\n    },\n    c\n}")

	test("{[y[a]];keep_indent;}", "{\n    [y[a]];\n    keep_indent;\n}")

	test("if (x) {y} else { if (x) {y}}", "if (x) {\n    y\n} else {\n    if (x) {\n        y\n    }\n}")

	test("if (foo) one()\ntwo()\nthree()")
	test("if (1 + foo() && bar(baz()) / 2) one()\ntwo()\nthree()")
	test("if (1 + foo() && bar(baz()) / 2) one();\ntwo();\nthree();")

	test_options["indent_size"] = 1
	test_options["indent_char"] = " "
	test("{ one_char() }", "{\n one_char()\n}")

	test("var a,b=1,c=2", "var a, b = 1,\n c = 2")

	test_options["indent_size"] = 4
	test_options["indent_char"] = " "
	test("{ one_char() }", "{\n    one_char()\n}")

	test_options["indent_size"] = 1
	test_options["indent_char"] = "\t"
	test("{ one_char() }", "{\n\tone_char()\n}")
	test("x = a ? b : c; x;", "x = a ? b : c;\nx;")

	test_options["indent_size"] = 5
	test_options["indent_char"] = " "
	test_options["indent_with_tabs"] = true

	test("{ one_char() }", "{\n\tone_char()\n}")
	test("x = a ? b : c; x;", "x = a ? b : c;\nx;")

	test_options["indent_size"] = 4
	test_options["indent_char"] = " "
	test_options["indent_with_tabs"] = false

	test_options["preserve_newlines"] = false
	test("var\na=dont_preserve_newlines;", "var a = dont_preserve_newlines;")

	test("function foo() {\n    return 1;\n}\n\nfunction foo() {\n    return 1;\n}")
	test("function foo() {\n    return 1;\n}\nfunction foo() {\n    return 1;\n}", "function foo() {\n    return 1;\n}\n\nfunction foo() {\n    return 1;\n}")
	test("function foo() {\n    return 1;\n}\n\n\nfunction foo() {\n    return 1;\n}", "function foo() {\n    return 1;\n}\n\nfunction foo() {\n    return 1;\n}")

	test_options["preserve_newlines"] = true
	test("var\na=do_preserve_newlines;", "var\n    a = do_preserve_newlines;")
	test("// a\n// b\n\n// c\n// d")
	test("if (foo) //  comment\n{\n    bar();\n}")

	test_options["keep_array_indentation"] = false
	test("a = ['a', 'b', 'c',\n    'd', 'e', 'f']",
		"a = ['a', 'b', 'c',\n    'd', 'e', 'f'\n]")
	test("a = ['a', 'b', 'c',\n    'd', 'e', 'f',\n        'g', 'h', 'i']",
		"a = ['a', 'b', 'c',\n    'd', 'e', 'f',\n    'g', 'h', 'i'\n]")
	test("a = ['a', 'b', 'c',\n        'd', 'e', 'f',\n            'g', 'h', 'i']",
		"a = ['a', 'b', 'c',\n    'd', 'e', 'f',\n    'g', 'h', 'i'\n]")
	test("var x = [{}\n]", "var x = [{}]")
	test("var x = [{foo:bar}\n]", "var x = [{\n    foo: bar\n}]")
	test("a = ['something',\n    'completely',\n    'different'];\nif (x);",
		"a = ['something',\n    'completely',\n    'different'\n];\nif (x);")
	test("a = ['a','b','c']", "a = ['a', 'b', 'c']")
	test("a = ['a',   'b','c']", "a = ['a', 'b', 'c']")
	test("x = [{'a':0}]",
		"x = [{\n    'a': 0\n}]")
	test("{a([[a1]], {b;});}", "{\n    a([\n        [a1]\n    ], {\n        b;\n    });\n}")
	test("a();\n   [\n   ['sdfsdfsd'],\n        ['sdfsdfsdf']\n   ].toString();",
		"a();\n[\n    ['sdfsdfsd'],\n    ['sdfsdfsdf']\n].toString();")
	test("a();\na = [\n   ['sdfsdfsd'],\n        ['sdfsdfsdf']\n   ].toString();",
		"a();\na = [\n    ['sdfsdfsd'],\n    ['sdfsdfsdf']\n].toString();")
	test("function() {\n    Foo([\n        ['sdfsdfsd'],\n        ['sdfsdfsdf']\n    ]);\n}",
		"function() {\n    Foo([\n        ['sdfsdfsd'],\n        ['sdfsdfsdf']\n    ]);\n}")
	test("function foo() {\n    return [\n        \"one\",\n        \"two\"\n    ];\n}")
	// 4 spaces per indent input, processed with 4-spaces per indent
	test("function foo() {\n"+
		"    return [\n"+
		"        {\n"+
		"            one: 'x',\n"+
		"            two: [\n"+
		"                {\n"+
		"                    id: 'a',\n"+
		"                    name: 'apple'\n"+
		"                }, {\n"+
		"                    id: 'b',\n"+
		"                    name: 'banana'\n"+
		"                }\n"+
		"            ]\n"+
		"        }\n"+
		"    ];\n"+
		"}",
		"function foo() {\n"+
			"    return [{\n"+
			"        one: 'x',\n"+
			"        two: [{\n"+
			"            id: 'a',\n"+
			"            name: 'apple'\n"+
			"        }, {\n"+
			"            id: 'b',\n"+
			"            name: 'banana'\n"+
			"        }]\n"+
			"    }];\n"+
			"}")

	test("function foo() {\n"+
		"   return [\n"+
		"      {\n"+
		"         one: 'x',\n"+
		"         two: [\n"+
		"            {\n"+
		"               id: 'a',\n"+
		"               name: 'apple'\n"+
		"            }, {\n"+
		"               id: 'b',\n"+
		"               name: 'banana'\n"+
		"            }\n"+
		"         ]\n"+
		"      }\n"+
		"   ];\n"+
		"}",
		"function foo() {\n"+
			"    return [{\n"+
			"        one: 'x',\n"+
			"        two: [{\n"+
			"            id: 'a',\n"+
			"            name: 'apple'\n"+
			"        }, {\n"+
			"            id: 'b',\n"+
			"            name: 'banana'\n"+
			"        }]\n"+
			"    }];\n"+
			"}")

	test_options["keep_array_indentation"] = true
	test("a = ['a', 'b', 'c',\n    'd', 'e', 'f']")
	test("a = ['a', 'b', 'c',\n    'd', 'e', 'f',\n        'g', 'h', 'i']")
	test("a = ['a', 'b', 'c',\n        'd', 'e', 'f',\n            'g', 'h', 'i']")
	test("var x = [{}\n]", "var x = [{}\n]")
	test("var x = [{foo:bar}\n]", "var x = [{\n        foo: bar\n    }\n]")
	test("a = ['something',\n    'completely',\n    'different'];\nif (x);")
	test("a = ['a','b','c']", "a = ['a', 'b', 'c']")
	test("a = ['a',   'b','c']", "a = ['a', 'b', 'c']")
	test("x = [{'a':0}]",
		"x = [{\n    'a': 0\n}]")
	test("{a([[a1]], {b;});}", "{\n    a([[a1]], {\n        b;\n    });\n}")
	test("a();\n   [\n   ['sdfsdfsd'],\n        ['sdfsdfsdf']\n   ].toString();",
		"a();\n   [\n   ['sdfsdfsd'],\n        ['sdfsdfsdf']\n   ].toString();")
	test("a();\na = [\n   ['sdfsdfsd'],\n        ['sdfsdfsdf']\n   ].toString();",
		"a();\na = [\n   ['sdfsdfsd'],\n        ['sdfsdfsdf']\n   ].toString();")
	test("function() {\n    Foo([\n        ['sdfsdfsd'],\n        ['sdfsdfsdf']\n    ]);\n}",
		"function() {\n    Foo([\n        ['sdfsdfsd'],\n        ['sdfsdfsdf']\n    ]);\n}")
	test("function foo() {\n    return [\n        \"one\",\n        \"two\"\n    ];\n}")

	test("function foo() {\n" +
		"    return [\n" +
		"        {\n" +
		"            one: 'x',\n" +
		"            two: [\n" +
		"                {\n" +
		"                    id: 'a',\n" +
		"                    name: 'apple'\n" +
		"                }, {\n" +
		"                    id: 'b',\n" +
		"                    name: 'banana'\n" +
		"                }\n" +
		"            ]\n" +
		"        }\n" +
		"    ];\n" +
		"}")

	test_options["keep_array_indentation"] = false

	test("a = //comment\n    /regex/;")

	test("/*\n * X\n */")
	test("/*\r\n * X\r\n */", "/*\n * X\n */")

	test("if (a)\n{\nb;\n}\nelse\n{\nc;\n}", "if (a) {\n    b;\n} else {\n    c;\n}")

	test("var a = new function();")
	test("new function")

	test_options["brace_style"] = "expand"

	test("//case 1\nif (a == 1)\n{}\n//case 2\nelse if (a == 2)\n{}")
	test("if(1){2}else{3}", "if (1)\n{\n    2\n}\nelse\n{\n    3\n}")
	test("try{a();}catch(b){c();}catch(d){}finally{e();}", "try\n{\n    a();\n}\ncatch (b)\n{\n    c();\n}\ncatch (d)\n{}\nfinally\n{\n    e();\n}")
	test("if(a){b();}else if(c) foo();", "if (a)\n{\n    b();\n}\nelse if (c) foo();")
	test("if (a) {\n// comment\n}else{\n// comment\n}",
		"if (a)\n{\n    // comment\n}\nelse\n{\n    // comment\n}")
	test("if (x) {y} else { if (x) {y}}", "if (x)\n{\n    y\n}\nelse\n{\n    if (x)\n    {\n        y\n    }\n}")
	test("if (a)\n{\nb;\n}\nelse\n{\nc;\n}", "if (a)\n{\n    b;\n}\nelse\n{\n    c;\n}")
	test("    /*\n* xx\n*/\n// xx\nif (foo) {\n    bar();\n}", "    /*\n     * xx\n     */\n    // xx\n    if (foo)\n    {\n        bar();\n    }")
	test("if (foo)\n{}\nelse /regex/.test();")
	test("if (foo) {", "if (foo)\n{")
	test("foo {", "foo\n{")
	test("return {", "return {")
	test("return /* inline */ {", "return /* inline */ {")
	test("return;\n{", "return;\n{")
	test("throw {}")
	test("throw {\n    foo;\n}")
	test("var foo = {}")
	test("function x() {\n    foo();\n}zzz", "function x()\n{\n    foo();\n}\nzzz")
	test("a: do {} while (); xxx", "a: do {} while ();\nxxx")
	test("{a: do {} while (); xxx}", "{\n    a: do {} while ();xxx\n}")
	test("var a = new function() {};")
	test("var a = new function a() {};", "var a = new function a()\n{};")
	test("var a = new function()\n{};", "var a = new function() {};")
	test("var a = new function a()\n{};")
	test("var a = new function a()\n    {},\n    b = new function b()\n    {};")
	test("foo({\n    'a': 1\n},\n10);",
		"foo(\n    {\n        'a': 1\n    },\n    10);")
	test("([\"foo\",\"bar\"]).each(function(i) {return i;});", "([\"foo\", \"bar\"]).each(function(i)\n{\n    return i;\n});")
	test("(function(i) {return i;})();", "(function(i)\n{\n    return i;\n})();")
	test("test( /*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test(\n"+
		"/*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"},\n"+
		"/*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test(\n"+
			"    /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test( /*Argument 1*/\n"+
		"{\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */\n"+
		"{\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")

	test_options["brace_style"] = "collapse"

	test("//case 1\nif (a == 1) {}\n//case 2\nelse if (a == 2) {}")
	test("if(1){2}else{3}", "if (1) {\n    2\n} else {\n    3\n}")
	test("try{a();}catch(b){c();}catch(d){}finally{e();}", "try {\n    a();\n} catch (b) {\n    c();\n} catch (d) {} finally {\n    e();\n}")
	test("if(a){b();}else if(c) foo();", "if (a) {\n    b();\n} else if (c) foo();")
	test("if (a) {\n// comment\n}else{\n// comment\n}",
		"if (a) {\n    // comment\n} else {\n    // comment\n}")
	test("if (x) {y} else { if (x) {y}}", "if (x) {\n    y\n} else {\n    if (x) {\n        y\n    }\n}")
	test("if (a)\n{\nb;\n}\nelse\n{\nc;\n}", "if (a) {\n    b;\n} else {\n    c;\n}")
	test("    /*\n* xx\n*/\n// xx\nif (foo) {\n    bar();\n}", "    /*\n     * xx\n     */\n    // xx\n    if (foo) {\n        bar();\n    }")
	test("if (foo) {} else /regex/.test();")
	test("if (foo) {", "if (foo) {")
	test("foo {", "foo {")
	test("return {", "return {")
	test("return /* inline */ {", "return /* inline */ {")
	test("return;\n{", "return; {")
	test("throw {}")
	test("throw {\n    foo;\n}")
	test("var foo = {}")
	test("function x() {\n    foo();\n}zzz", "function x() {\n    foo();\n}\nzzz")
	test("a: do {} while (); xxx", "a: do {} while ();\nxxx")
	test("{a: do {} while (); xxx}", "{\n    a: do {} while ();xxx\n}")
	test("var a = new function() {};")
	test("var a = new function a() {};")
	test("var a = new function()\n{};", "var a = new function() {};")
	test("var a = new function a()\n{};", "var a = new function a() {};")
	test("var a = new function a()\n    {},\n    b = new function b()\n    {};", "var a = new function a() {},\n    b = new function b() {};")
	test("foo({\n    'a': 1\n},\n10);",
		"foo({\n        'a': 1\n    },\n    10);")
	test("([\"foo\",\"bar\"]).each(function(i) {return i;});", "([\"foo\", \"bar\"]).each(function(i) {\n    return i;\n});")
	test("(function(i) {return i;})();", "(function(i) {\n    return i;\n})();")
	test("test( /*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/ {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test(\n"+
		"/*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"},\n"+
		"/*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test(\n"+
			"    /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test( /*Argument 1*/\n"+
		"{\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */\n"+
		"{\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/ {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")

	test_options["brace_style"] = "end-expand"

	test("//case 1\nif (a == 1) {}\n//case 2\nelse if (a == 2) {}")
	test("if(1){2}else{3}", "if (1) {\n    2\n}\nelse {\n    3\n}")
	test("try{a();}catch(b){c();}catch(d){}finally{e();}", "try {\n    a();\n}\ncatch (b) {\n    c();\n}\ncatch (d) {}\nfinally {\n    e();\n}")
	test("if(a){b();}else if(c) foo();", "if (a) {\n    b();\n}\nelse if (c) foo();")
	test("if (a) {\n// comment\n}else{\n// comment\n}",
		"if (a) {\n    // comment\n}\nelse {\n    // comment\n}")
	test("if (x) {y} else { if (x) {y}}", "if (x) {\n    y\n}\nelse {\n    if (x) {\n        y\n    }\n}")
	test("if (a)\n{\nb;\n}\nelse\n{\nc;\n}", "if (a) {\n    b;\n}\nelse {\n    c;\n}")
	test("    /*\n* xx\n*/\n// xx\nif (foo) {\n    bar();\n}", "    /*\n     * xx\n     */\n    // xx\n    if (foo) {\n        bar();\n    }")
	test("if (foo) {}\nelse /regex/.test();")
	test("if (foo) {", "if (foo) {")
	test("foo {", "foo {")
	test("return {", "return {")
	test("return /* inline */ {", "return /* inline */ {")
	test("return;\n{", "return; {")
	test("throw {}")
	test("throw {\n    foo;\n}")
	test("var foo = {}")
	test("function x() {\n    foo();\n}zzz", "function x() {\n    foo();\n}\nzzz")
	test("a: do {} while (); xxx", "a: do {} while ();\nxxx")
	test("{a: do {} while (); xxx}", "{\n    a: do {} while ();xxx\n}")
	test("var a = new function() {};")
	test("var a = new function a() {};")
	test("var a = new function()\n{};", "var a = new function() {};")
	test("var a = new function a()\n{};", "var a = new function a() {};")
	test("var a = new function a()\n    {},\n    b = new function b()\n    {};", "var a = new function a() {},\n    b = new function b() {};")
	test("foo({\n    'a': 1\n},\n10);",
		"foo({\n        'a': 1\n    },\n    10);")
	test("([\"foo\",\"bar\"]).each(function(i) {return i;});", "([\"foo\", \"bar\"]).each(function(i) {\n    return i;\n});")
	test("(function(i) {return i;})();", "(function(i) {\n    return i;\n})();")
	test("test( /*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/ {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test(\n"+
		"/*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"},\n"+
		"/*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test(\n"+
			"    /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test( /*Argument 1*/\n"+
		"{\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */\n"+
		"{\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/ {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")

	test_options["brace_style"] = "none"

	test("//case 1\nif (a == 1)\n{}\n//case 2\nelse if (a == 2)\n{}")
	test("if(1){2}else{3}", "if (1) {\n    2\n} else {\n    3\n}")
	test("try{a();}catch(b){c();}catch(d){}finally{e();}", "try {\n    a();\n} catch (b) {\n    c();\n} catch (d) {} finally {\n    e();\n}")
	test("if(a){b();}else if(c) foo();", "if (a) {\n    b();\n} else if (c) foo();")
	test("if (a) {\n// comment\n}else{\n// comment\n}",
		"if (a) {\n    // comment\n} else {\n    // comment\n}")
	test("if (x) {y} else { if (x) {y}}", "if (x) {\n    y\n} else {\n    if (x) {\n        y\n    }\n}")
	test("if (a)\n{\nb;\n}\nelse\n{\nc;\n}", "if (a)\n{\n    b;\n}\nelse\n{\n    c;\n}")
	test("    /*\n* xx\n*/\n// xx\nif (foo) {\n    bar();\n}", "    /*\n     * xx\n     */\n    // xx\n    if (foo) {\n        bar();\n    }")
	test("if (foo)\n{}\nelse /regex/.test();")
	test("if (foo) {")
	test("foo {")
	test("return {")
	test("return /* inline */ {")
	test("return;\n{")
	test("throw {}")
	test("throw {\n    foo;\n}")
	test("var foo = {}")
	test("function x() {\n    foo();\n}zzz", "function x() {\n    foo();\n}\nzzz")
	test("a: do {} while (); xxx", "a: do {} while ();\nxxx")
	test("{a: do {} while (); xxx}", "{\n    a: do {} while ();xxx\n}")
	test("var a = new function() {};")
	test("var a = new function a() {};")
	test("var a = new function()\n{};", "var a = new function() {};")
	test("var a = new function a()\n{};")
	test("var a = new function a()\n    {},\n    b = new function b()\n    {};")
	test("foo({\n    'a': 1\n},\n10);",
		"foo({\n        'a': 1\n    },\n    10);")
	test("([\"foo\",\"bar\"]).each(function(i) {return i;});", "([\"foo\", \"bar\"]).each(function(i) {\n    return i;\n});")
	test("(function(i) {return i;})();", "(function(i) {\n    return i;\n})();")
	test("test( /*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/ {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test(\n"+
		"/*Argument 1*/ {\n"+
		"    'Value1': '1'\n"+
		"},\n"+
		"/*Argument 2\n"+
		" */ {\n"+
		"    'Value2': '2'\n"+
		"});",

		"test(\n"+
			"    /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")
	test("test( /*Argument 1*/\n"+
		"{\n"+
		"    'Value1': '1'\n"+
		"}, /*Argument 2\n"+
		" */\n"+
		"{\n"+
		"    'Value2': '2'\n"+
		"});",

		"test( /*Argument 1*/\n"+
			"    {\n"+
			"        'Value1': '1'\n"+
			"    },\n"+
			"    /*Argument 2\n"+
			"     */\n"+
			"    {\n"+
			"        'Value2': '2'\n"+
			"    });")

	test_options["brace_style"] = "collapse"

	test("a = <?= external() ?> ;")
	test("a = <%= external() %> ;")

	test("roo = {\n    /*\n    ****\n      FOO\n    ****\n    */\n    BAR: 0\n};")
	test("if (zz) {\n    // ....\n}\n(function")

	test_options["preserve_newlines"] = true
	test("var a = 42; // foo\n\nvar b;")
	test("var a = 42; // foo\n\n\nvar b;")
	test("var a = 'foo' +\n    'bar';")
	test("var a = \"foo\" +\n    \"bar\";")

	test("\"foo\"\"bar\"\"baz\"", "\"foo\"\n\"bar\"\n\"baz\"")
	test("'foo''bar''baz'", "'foo'\n'bar'\n'baz'")
	test("{\n    get foo() {}\n}")
	test("{\n    var a = get\n    foo();\n}")
	test("{\n    set foo() {}\n}")
	test("{\n    var a = set\n    foo();\n}")
	test("var x = {\n    get function()\n}")
	test("var x = {\n    set function()\n}")

	test("var x = set\n\na() {}", "var x = set\n\na() {}")
	test("var x = set\n\nfunction() {}", "var x = set\n\nfunction() {}")

	test("<!-- foo\nbar();\n-->")
	test("<!-- dont crash")
	test("for () /abc/.test()")
	test("if (k) /aaa/m.test(v) && l();")
	test("switch (true) {\n    case /swf/i.test(foo):\n        bar();\n}")
	test("createdAt = {\n    type: Date,\n    default: Date.now\n}")
	test("switch (createdAt) {\n    case a:\n        Date,\n    default:\n        Date.now\n}")

	test("return function();")
	test("var a = function();")
	test("var a = 5 + function();")

	test("{\n    foo // something\n    ,\n    bar // something\n    baz\n}")
	test("function a(a) {} function b(b) {} function c(c) {}", "function a(a) {}\n\nfunction b(b) {}\n\nfunction c(c) {}")

	test("3.*7;", "3. * 7;")
	test("a = 1.e-64 * 0.5e+4 / 6e-23;")
	test("import foo.*;", "import foo.*;")
	test("function f(a: a, b: b)")
	test("foo(a, function() {})")
	test("foo(a, /regex/)")

	test("/* foo */\n\"x\"")

	test_options["break_chained_methods"] = false
	test_options["preserve_newlines"] = false
	test("foo\n.bar()\n.baz().cucumber(fat)", "foo.bar().baz().cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat); foo.bar().baz().cucumber(fat)", "foo.bar().baz().cucumber(fat);\nfoo.bar().baz().cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat)\n foo.bar().baz().cucumber(fat)", "foo.bar().baz().cucumber(fat)\nfoo.bar().baz().cucumber(fat)")
	test("this\n.something = foo.bar()\n.baz().cucumber(fat)", "this.something = foo.bar().baz().cucumber(fat)")
	test("this.something.xxx = foo.moo.bar()")
	test("this\n.something\n.xxx = foo.moo\n.bar()", "this.something.xxx = foo.moo.bar()")

	test_options["break_chained_methods"] = false
	test_options["preserve_newlines"] = true
	test("foo\n.bar()\n.baz().cucumber(fat)", "foo\n    .bar()\n    .baz().cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat); foo.bar().baz().cucumber(fat)", "foo\n    .bar()\n    .baz().cucumber(fat);\nfoo.bar().baz().cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat)\n foo.bar().baz().cucumber(fat)", "foo\n    .bar()\n    .baz().cucumber(fat)\nfoo.bar().baz().cucumber(fat)")
	test("this\n.something = foo.bar()\n.baz().cucumber(fat)", "this\n    .something = foo.bar()\n    .baz().cucumber(fat)")
	test("this.something.xxx = foo.moo.bar()")
	test("this\n.something\n.xxx = foo.moo\n.bar()", "this\n    .something\n    .xxx = foo.moo\n    .bar()")

	test_options["break_chained_methods"] = true
	test_options["preserve_newlines"] = false
	test("foo\n.bar()\n.baz().cucumber(fat)", "foo.bar()\n    .baz()\n    .cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat); foo.bar().baz().cucumber(fat)", "foo.bar()\n    .baz()\n    .cucumber(fat);\nfoo.bar()\n    .baz()\n    .cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat)\n foo.bar().baz().cucumber(fat)", "foo.bar()\n    .baz()\n    .cucumber(fat)\nfoo.bar()\n    .baz()\n    .cucumber(fat)")
	test("this\n.something = foo.bar()\n.baz().cucumber(fat)", "this.something = foo.bar()\n    .baz()\n    .cucumber(fat)")
	test("this.something.xxx = foo.moo.bar()")
	test("this\n.something\n.xxx = foo.moo\n.bar()", "this.something.xxx = foo.moo.bar()")

	test_options["break_chained_methods"] = true
	test_options["preserve_newlines"] = true
	test("foo\n.bar()\n.baz().cucumber(fat)", "foo\n    .bar()\n    .baz()\n    .cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat); foo.bar().baz().cucumber(fat)", "foo\n    .bar()\n    .baz()\n    .cucumber(fat);\nfoo.bar()\n    .baz()\n    .cucumber(fat)")
	test("foo\n.bar()\n.baz().cucumber(fat)\n foo.bar().baz().cucumber(fat)", "foo\n    .bar()\n    .baz()\n    .cucumber(fat)\nfoo.bar()\n    .baz()\n    .cucumber(fat)")
	test("this\n.something = foo.bar()\n.baz().cucumber(fat)", "this\n    .something = foo.bar()\n    .baz()\n    .cucumber(fat)")
	test("this.something.xxx = foo.moo.bar()")
	test("this\n.something\n.xxx = foo.moo\n.bar()", "this\n    .something\n    .xxx = foo.moo\n    .bar()")
	test_options["break_chained_methods"] = false
	test_options["preserve_newlines"] = false

	wrap_input_1 := ("foo.bar().baz().cucumber((fat && \"sassy\") || (leans\n&& mean));\n" +
		"Test_very_long_variable_name_this_should_never_wrap\n.but_this_can\n" +
		"if (wraps_can_occur && inside_an_if_block) that_is_\n.okay();\n" +
		"object_literal = {\n" +
		"    propertx: first_token + 12345678.99999E-6,\n" +
		"    property: first_token_should_never_wrap + but_this_can,\n" +
		"    propertz: first_token_should_never_wrap + !but_this_can,\n" +
		"    proper: \"first_token_should_never_wrap\" + \"but_this_can\"\n" +
		"}")

	wrap_input_2 := ("{\n" +
		"    foo.bar().baz().cucumber((fat && \"sassy\") || (leans\n&& mean));\n" +
		"    Test_very_long_variable_name_this_should_never_wrap\n.but_this_can\n" +
		"    if (wraps_can_occur && inside_an_if_block) that_is_\n.okay();\n" +
		"    object_literal = {\n" +
		"        propertx: first_token + 12345678.99999E-6,\n" +
		"        property: first_token_should_never_wrap + but_this_can,\n" +
		"        propertz: first_token_should_never_wrap + !but_this_can,\n" +
		"        proper: \"first_token_should_never_wrap\" + \"but_this_can\"\n" +
		"    }" +
		"}")

	test_options["preserve_newlines"] = false
	test_options["wrap_line_length"] = 0

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat && \"sassy\") || (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap.but_this_can\n"+
			"if (wraps_can_occur && inside_an_if_block) that_is_.okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token + 12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap + but_this_can,\n"+
			"    propertz: first_token_should_never_wrap + !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" + \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 70

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat && \"sassy\") || (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap.but_this_can\n"+
			"if (wraps_can_occur && inside_an_if_block) that_is_.okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token + 12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap + but_this_can,\n"+
			"    propertz: first_token_should_never_wrap + !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" + \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 40

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat &&\n"+
			"    \"sassy\") || (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap\n"+
			"    .but_this_can\n"+
			"if (wraps_can_occur &&\n"+
			"    inside_an_if_block) that_is_.okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token +\n"+
			"        12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap +\n"+
			"        but_this_can,\n"+
			"    propertz: first_token_should_never_wrap +\n"+
			"        !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" +\n"+
			"        \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 41

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat && \"sassy\") ||\n"+
			"    (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap\n"+
			"    .but_this_can\n"+
			"if (wraps_can_occur &&\n"+
			"    inside_an_if_block) that_is_.okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token +\n"+
			"        12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap +\n"+
			"        but_this_can,\n"+
			"    propertz: first_token_should_never_wrap +\n"+
			"        !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" +\n"+
			"        \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 45

	test(wrap_input_2,

		"{\n"+
			"    foo.bar().baz().cucumber((fat && \"sassy\") ||\n"+
			"        (leans && mean));\n"+
			"    Test_very_long_variable_name_this_should_never_wrap\n"+
			"        .but_this_can\n"+
			"    if (wraps_can_occur &&\n"+
			"        inside_an_if_block) that_is_.okay();\n"+
			"    object_literal = {\n"+
			"        propertx: first_token +\n"+
			"            12345678.99999E-6,\n"+
			"        property: first_token_should_never_wrap +\n"+
			"            but_this_can,\n"+
			"        propertz: first_token_should_never_wrap +\n"+
			"            !but_this_can,\n"+
			"        proper: \"first_token_should_never_wrap\" +\n"+
			"            \"but_this_can\"\n"+
			"    }\n"+
			"}")

	test_options["preserve_newlines"] = true
	test_options["wrap_line_length"] = 0

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat && \"sassy\") || (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap\n"+
			"    .but_this_can\n"+
			"if (wraps_can_occur && inside_an_if_block) that_is_\n"+
			"    .okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token + 12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap + but_this_can,\n"+
			"    propertz: first_token_should_never_wrap + !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" + \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 70

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat && \"sassy\") || (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap\n"+
			"    .but_this_can\n"+
			"if (wraps_can_occur && inside_an_if_block) that_is_\n"+
			"    .okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token + 12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap + but_this_can,\n"+
			"    propertz: first_token_should_never_wrap + !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" + \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 40

	test(wrap_input_1,

		"foo.bar().baz().cucumber((fat &&\n"+
			"    \"sassy\") || (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap\n"+
			"    .but_this_can\n"+
			"if (wraps_can_occur &&\n"+
			"    inside_an_if_block) that_is_\n"+
			"    .okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token +\n"+
			"        12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap +\n"+
			"        but_this_can,\n"+
			"    propertz: first_token_should_never_wrap +\n"+
			"        !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" +\n"+
			"        \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 41

	test(wrap_input_1,
		"foo.bar().baz().cucumber((fat && \"sassy\") ||\n"+
			"    (leans && mean));\n"+
			"Test_very_long_variable_name_this_should_never_wrap\n"+
			"    .but_this_can\n"+
			"if (wraps_can_occur &&\n"+
			"    inside_an_if_block) that_is_\n"+
			"    .okay();\n"+
			"object_literal = {\n"+
			"    propertx: first_token +\n"+
			"        12345678.99999E-6,\n"+
			"    property: first_token_should_never_wrap +\n"+
			"        but_this_can,\n"+
			"    propertz: first_token_should_never_wrap +\n"+
			"        !but_this_can,\n"+
			"    proper: \"first_token_should_never_wrap\" +\n"+
			"        \"but_this_can\"\n"+
			"}")

	test_options["wrap_line_length"] = 45

	test(wrap_input_2,

		"{\n"+
			"    foo.bar().baz().cucumber((fat && \"sassy\") ||\n"+
			"        (leans && mean));\n"+
			"    Test_very_long_variable_name_this_should_never_wrap\n"+
			"        .but_this_can\n"+
			"    if (wraps_can_occur &&\n"+
			"        inside_an_if_block) that_is_\n"+
			"        .okay();\n"+
			"    object_literal = {\n"+
			"        propertx: first_token +\n"+
			"            12345678.99999E-6,\n"+
			"        property: first_token_should_never_wrap +\n"+
			"            but_this_can,\n"+
			"        propertz: first_token_should_never_wrap +\n"+
			"            !but_this_can,\n"+
			"        proper: \"first_token_should_never_wrap\" +\n"+
			"            \"but_this_can\"\n"+
			"    }\n"+
			"}")

	test_options["wrap_line_length"] = 0

	test_options["preserve_newlines"] = false
	test("if (foo) // comment\n    bar();")
	test("if (foo) // comment\n    (bar());")
	test("if (foo) // comment\n    (bar());")
	test("if (foo) // comment\n    /asdf/;")
	test("this.oa = new OAuth(\n"+
		"    _requestToken,\n"+
		"    _accessToken,\n"+
		"    consumer_key\n"+
		");", "this.oa = new OAuth(_requestToken, _accessToken, consumer_key);")
	test("foo = {\n    x: y, // #44\n    w: z // #44\n}")
	test("switch (x) {\n    case \"a\":\n        // comment on newline\n        break;\n    case \"b\": // comment on same line\n        break;\n}")
	test("this.type =\n    this.options =\n    // comment\n    this.enabled null;", "this.type = this.options =\n    // comment\n    this.enabled null;")
	test("someObj\n    .someFunc1()\n    // This comment should not break the indent\n    .someFunc2();", "someObj.someFunc1()\n    // This comment should not break the indent\n    .someFunc2();")

	test("if (true ||\n!true) return;", "if (true || !true) return;")

	test("if\n(foo)\nif\n(bar)\nif\n(baz)\nwhee();\na();", "if (foo)\n    if (bar)\n        if (baz) whee();\na();")
	test("if\n(foo)\nif\n(bar)\nif\n(baz)\nwhee();\nelse\na();", "if (foo)\n    if (bar)\n        if (baz) whee();\n        else a();")
	test("if (foo)\nbar();\nelse\ncar();", "if (foo) bar();\nelse car();")

	test("if (foo) if (bar) if (baz);\na();", "if (foo)\n    if (bar)\n        if (baz);\na();")
	test("if (foo) if (bar) if (baz) whee();\na();", "if (foo)\n    if (bar)\n        if (baz) whee();\na();")
	test("if (foo) a()\nif (bar) if (baz) whee();\na();", "if (foo) a()\nif (bar)\n    if (baz) whee();\na();")
	test("if (foo);\nif (bar) if (baz) whee();\na();", "if (foo);\nif (bar)\n    if (baz) whee();\na();")
	test("if (options)\n"+"    for (var p in options)\n"+"        this[p] = options[p];", "if (options)\n"+"    for (var p in options) this[p] = options[p];")
	test("if (options) for (var p in options) this[p] = options[p];", "if (options)\n    for (var p in options) this[p] = options[p];")

	test("if (options) do q(); while (b());", "if (options)\n    do q(); while (b());")
	test("if (options) while (b()) q();", "if (options)\n    while (b()) q();")
	test("if (options) do while (b()) q(); while (a());", "if (options)\n    do\n        while (b()) q(); while (a());")

	test("function f(a, b, c,\nd, e) {}", "function f(a, b, c, d, e) {}")

	test("function f(a,b) {if(a) b()}function g(a,b) {if(!a) b()}", "function f(a, b) {\n    if (a) b()\n}\n\nfunction g(a, b) {\n    if (!a) b()\n}")
	test("function f(a,b) {if(a) b()}\n\n\n\nfunction g(a,b) {if(!a) b()}", "function f(a, b) {\n    if (a) b()\n}\n\nfunction g(a, b) {\n    if (!a) b()\n}")

	test("(if(a) b())(if(a) b())", "(\n    if (a) b())(\n    if (a) b())")
	test("(if(a) b())\n\n\n(if(a) b())", "(\n    if (a) b())\n(\n    if (a) b())")

	test("/*\n * foo\n */\nfunction foo() {}")
	test("// a nice function\nfunction foo() {}")
	test("function foo() {}\nfunction foo() {}", "function foo() {}\n\nfunction foo() {}")

	test("[\n    function() {}\n]")

	test("if\n(a)\nb();", "if (a) b();")
	test("var a =\nfoo", "var a = foo")
	test("var a = {\n\"a\":1,\n\"b\":2}", "var a = {\n    \"a\": 1,\n    \"b\": 2\n}")
	test("var a = {\n'a':1,\n'b':2}", "var a = {\n    'a': 1,\n    'b': 2\n}")
	test("var a = /*i*/ \"b\";")
	test("var a = /*i*/\n\"b\";", "var a = /*i*/ \"b\";")
	test("var a = /*i*/\nb;", "var a = /*i*/ b;")
	test("{\n\n\n\"x\"\n}", "{\n    \"x\"\n}")
	test("if(a &&\nb\n||\nc\n||d\n&&\ne) e = f", "if (a && b || c || d && e) e = f")
	test("if(a &&\n(b\n||\nc\n||d)\n&&\ne) e = f", "if (a && (b || c || d) && e) e = f")
	test("\n\n\"x\"", "\"x\"")
	test("a = 1;\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nb = 2;", "a = 1;\nb = 2;")

	test_options["preserve_newlines"] = true
	test("if (foo) // comment\n    bar();")
	test("if (foo) // comment\n    (bar());")
	test("if (foo) // comment\n    (bar());")
	test("if (foo) // comment\n    /asdf/;")
	test("this.oa = new OAuth(\n" + "    _requestToken,\n" + "    _accessToken,\n" + "    consumer_key\n" + ");")
	test("foo = {\n    x: y, // #44\n    w: z // #44\n}")
	test("switch (x) {\n    case \"a\":\n        // comment on newline\n        break;\n    case \"b\": // comment on same line\n        break;\n}")
	test("this.type =\n    this.options =\n    // comment\n    this.enabled null;")
	test("someObj\n    .someFunc1()\n    // This comment should not break the indent\n    .someFunc2();")

	test("if (true ||\n!true) return;", "if (true ||\n    !true) return;")

	test("if\n(foo)\nif\n(bar)\nif\n(baz)\nwhee();\na();", "if (foo)\n    if (bar)\n        if (baz)\n            whee();\na();")
	test("if\n(foo)\nif\n(bar)\nif\n(baz)\nwhee();\nelse\na();", "if (foo)\n    if (bar)\n        if (baz)\n            whee();\n        else\n            a();")
	test("if (foo) bar();\nelse\ncar();", "if (foo) bar();\nelse\n    car();")

	test("if (foo) if (bar) if (baz);\na();", "if (foo)\n    if (bar)\n        if (baz);\na();")
	test("if (foo) if (bar) if (baz) whee();\na();", "if (foo)\n    if (bar)\n        if (baz) whee();\na();")
	test("if (foo) a()\nif (bar) if (baz) whee();\na();", "if (foo) a()\nif (bar)\n    if (baz) whee();\na();")
	test("if (foo);\nif (bar) if (baz) whee();\na();", "if (foo);\nif (bar)\n    if (baz) whee();\na();")
	test("if (options)\n" + "    for (var p in options)\n" + "        this[p] = options[p];")
	test("if (options) for (var p in options) this[p] = options[p];", "if (options)\n    for (var p in options) this[p] = options[p];")

	test("if (options) do q(); while (b());", "if (options)\n    do q(); while (b());")
	test("if (options) do; while (b());", "if (options)\n    do; while (b());")
	test("if (options) while (b()) q();", "if (options)\n    while (b()) q();")
	test("if (options) do while (b()) q(); while (a());", "if (options)\n    do\n        while (b()) q(); while (a());")

	test("function f(a, b, c,\nd, e) {}", "function f(a, b, c,\n    d, e) {}")

	test("function f(a,b) {if(a) b()}function g(a,b) {if(!a) b()}", "function f(a, b) {\n    if (a) b()\n}\n\nfunction g(a, b) {\n    if (!a) b()\n}")
	test("function f(a,b) {if(a) b()}\n\n\n\nfunction g(a,b) {if(!a) b()}", "function f(a, b) {\n    if (a) b()\n}\n\n\n\nfunction g(a, b) {\n    if (!a) b()\n}")

	test("(if(a) b())(if(a) b())", "(\n    if (a) b())(\n    if (a) b())")
	test("(if(a) b())\n\n\n(if(a) b())", "(\n    if (a) b())\n\n\n(\n    if (a) b())")

	test("if\n(a)\nb();", "if (a)\n    b();")
	test("var a =\nfoo", "var a =\n    foo")
	test("var a = {\n\"a\":1,\n\"b\":2}", "var a = {\n    \"a\": 1,\n    \"b\": 2\n}")
	test("var a = {\n'a':1,\n'b':2}", "var a = {\n    'a': 1,\n    'b': 2\n}")
	test("var a = /*i*/ \"b\";")
	test("var a = /*i*/\n\"b\";", "var a = /*i*/\n    \"b\";")
	test("var a = /*i*/\nb;", "var a = /*i*/\n    b;")
	test("{\n\n\n\"x\"\n}", "{\n\n\n    \"x\"\n}")
	test("if(a &&\nb\n||\nc\n||d\n&&\ne) e = f", "if (a &&\n    b ||\n    c || d &&\n    e) e = f")
	test("if(a &&\n(b\n||\nc\n||d)\n&&\ne) e = f", "if (a &&\n    (b ||\n        c || d) &&\n    e) e = f")
	test("\n\n\"x\"", "\"x\"")

	test("a = 1;\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nb = 2;", "a = 1;\n\n\n\n\n\n\n\n\n\nb = 2;")
	test_options["max_preserve_newlines"] = 8
	test("a = 1;\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nb = 2;", "a = 1;\n\n\n\n\n\n\n\nb = 2;")

	test_options["space_in_paren"] = false
	test_options["space_in_empty_paren"] = false
	test("if(p) foo(a,b)", "if (p) foo(a, b)")
	test("try{while(true){willThrow()}}catch(result)switch(result){case 1:++result }", "try {\n    while (true) {\n        willThrow()\n    }\n} catch (result) switch (result) {\n    case 1:\n        ++result\n}")
	test("((e/((a+(b)*c)-d))^2)*5;", "((e / ((a + (b) * c) - d)) ^ 2) * 5;")
	test("function f(a,b) {if(a) b()}function g(a,b) {if(!a) b()}", "function f(a, b) {\n    if (a) b()\n}\n\nfunction g(a, b) {\n    if (!a) b()\n}")
	test("a=[];", "a = [];")
	test("a=[b,c,d];", "a = [b, c, d];")
	test("a= f[b];", "a = f[b];")

	test_options["space_in_paren"] = true
	test("if(p) foo(a,b)", "if ( p ) foo( a, b )")
	test("try{while(true){willThrow()}}catch(result)switch(result){case 1:++result }", "try {\n    while ( true ) {\n        willThrow()\n    }\n} catch ( result ) switch ( result ) {\n    case 1:\n        ++result\n}")
	test("((e/((a+(b)*c)-d))^2)*5;", "( ( e / ( ( a + ( b ) * c ) - d ) ) ^ 2 ) * 5;")
	test("function f(a,b) {if(a) b()}function g(a,b) {if(!a) b()}", "function f( a, b ) {\n    if ( a ) b()\n}\n\nfunction g( a, b ) {\n    if ( !a ) b()\n}")
	test("a=[ ];", "a = [];")
	test("a=[b,c,d];", "a = [ b, c, d ];")
	test("a= f[b];", "a = f[ b ];")

	test_options["space_in_empty_paren"] = true
	test("if(p) foo(a,b)", "if ( p ) foo( a, b )")
	test("try{while(true){willThrow()}}catch(result)switch(result){case 1:++result }", "try {\n    while ( true ) {\n        willThrow( )\n    }\n} catch ( result ) switch ( result ) {\n    case 1:\n        ++result\n}")
	test("((e/((a+(b)*c)-d))^2)*5;", "( ( e / ( ( a + ( b ) * c ) - d ) ) ^ 2 ) * 5;")
	test("function f(a,b) {if(a) b()}function g(a,b) {if(!a) b()}", "function f( a, b ) {\n    if ( a ) b( )\n}\n\nfunction g( a, b ) {\n    if ( !a ) b( )\n}")
	test("a=[ ];", "a = [ ];")
	test("a=[b,c,d];", "a = [ b, c, d ];")
	test("a= f[b];", "a = f[ b ];")

	test_options["space_in_paren"] = false
	test_options["space_in_empty_paren"] = false

	test("`This is a ${template} string.`", "`This is a ${template} string.`")
	test("`This\n  is\n  a\n  ${template}\n  string.`", "`This\n  is\n  a\n  ${template}\n  string.`")

	test("xml=<a b=\"c\"><d/><e>\n foo</e>x</a>;", "xml = < a b = \"c\" > < d / > < e >\n    foo < /e>x</a > ;")
	test_options["e4x"] = true

	test("xml=<a b=\"c\"><d/><e>\n foo</e>x</a>;", "xml = <a b=\"c\"><d/><e>\n foo</e>x</a>;")
	test("<a b='This is a quoted \"c\".'/>", "<a b='This is a quoted \"c\".'/>")
	test("<a b=\"This is a quoted 'c'.\"/>", "<a b=\"This is a quoted 'c'.\"/>")
	test("<a b=\"A quote ' inside string.\"/>", "<a b=\"A quote ' inside string.\"/>")
	test("<a b='A quote \" inside string.'/>", "<a b='A quote \" inside string.'/>")
	test("<a b='Some \"\"\" quotes \"\"  inside string.'/>", "<a b='Some \"\"\" quotes \"\"  inside string.'/>")

	test("xml=<{a} b=\"c\"><d/><e v={z}>\n foo</e>x</{a}>;", "xml = <{a} b=\"c\"><d/><e v={z}>\n foo</e>x</{a}>;")

	test("xml = <_:.valid.xml- _:.valid.xml-=\"123\"/>;", "xml = <_:.valid.xml- _:.valid.xml-=\"123\"/>;")

	test("xml=<a b=\"c\"><![CDATA[d/>\n</a></{}]]></a>;", "xml = <a b=\"c\"><![CDATA[d/>\n</a></{}]]></a>;")
	test("xml=<![CDATA[]]>;", "xml = <![CDATA[]]>;")
	test("xml=<![CDATA[ b=\"c\"><d/><e v={z}>\n foo</e>x/]]>;", "xml = <![CDATA[ b=\"c\"><d/><e v={z}>\n foo</e>x/]]>;")

	test("xml=<a x=\"jn\"><c></b></f><a><d jnj=\"jnn\"><f></a ></nj></a>;", "xml = <a x=\"jn\"><c></b></f><a><d jnj=\"jnn\"><f></a ></nj></a>;")

	test("xml=<a></b>\nc<b;", "xml = <a></b>\nc<b;")
	test_options["e4x"] = false

	test("obj\n" + "    .last({\n" + "        foo: 1,\n" + "        bar: 2\n" + "    });\n" + "var test = 1;")

	test("obj\n" + "    .last(a, function() {\n" + "        var test;\n" + "    });\n" + "var test = 1;")

	test("obj.first()\n" + "    .second()\n" + "    .last(function(err, response) {\n" + "        console.log(err);\n" + "    });")

	test("obj.last(a, function() {\n" + "    var test;\n" + "});\n" + "var test = 1;")

	test("obj.last(a,\n" + "    function() {\n" + "        var test;\n" + "    });\n" + "var test = 1;")

	test("(function() {if (!window.FOO) window.FOO || (window.FOO = function() {var b = {bar: \"zort\"};});})();", "(function() {\n"+"    if (!window.FOO) window.FOO || (window.FOO = function() {\n"+"        var b = {\n"+"            bar: \"zort\"\n"+"        };\n"+"    });\n"+"})();")

	test("define([\"dojo/_base/declare\", \"my/Employee\", \"dijit/form/Button\",\n" + "    \"dojo/_base/lang\", \"dojo/Deferred\"\n" + "], function(declare, Employee, Button, lang, Deferred) {\n" + "    return declare(Employee, {\n" + "        constructor: function() {\n" + "            new Button({\n" + "                onClick: lang.hitch(this, function() {\n" + "                    new Deferred().then(lang.hitch(this, function() {\n" + "                        this.salary * 0.25;\n" + "                    }));\n" + "                })\n" + "            });\n" + "        }\n" + "    });\n" + "});")
	test("define([\"dojo/_base/declare\", \"my/Employee\", \"dijit/form/Button\",\n" + "        \"dojo/_base/lang\", \"dojo/Deferred\"\n" + "    ],\n" + "    function(declare, Employee, Button, lang, Deferred) {\n" + "        return declare(Employee, {\n" + "            constructor: function() {\n" + "                new Button({\n" + "                    onClick: lang.hitch(this, function() {\n" + "                        new Deferred().then(lang.hitch(this, function() {\n" + "                            this.salary * 0.25;\n" + "                        }));\n" + "                    })\n" + "                });\n" + "            }\n" + "        });\n" + "    });")

	test("(function() {\n" +
		"    return {\n" +
		"        foo: function() {\n" +
		"            return \"bar\";\n" +
		"        },\n" +
		"        bar: [\"bar\"]\n" +
		"    };\n" +
		"}());")

	test("var name = \"a;\n" + "name = \"b\";")
	test("var name = \"a; \\\n" + "    name = b\";")

	test("var c = \"_ACTION_TO_NATIVEAPI_\" + ++g++ + +new Date;")
	test("var c = \"_ACTION_TO_NATIVEAPI_\" - --g-- - -new Date;")

	test("a = {\n" + "    function: {},\n" + "    \"function\": {},\n" + "    throw: {},\n" + "    \"throw\": {},\n" + "    var: {},\n" + "    \"var\": {},\n" + "    set: {},\n" + "    \"set\": {},\n" + "    get: {},\n" + "    \"get\": {},\n" + "    if: {},\n" + "    \"if\": {},\n" + "    then: {},\n" + "    \"then\": {},\n" + "    else: {},\n" + "    \"else\": {},\n" + "    yay: {}\n" + "};")

	test("if(x){a();}else{b();}if(y){c();}", "if (x) {\n"+"    a();\n"+"} else {\n"+"    b();\n"+"}\n"+"if (y) {\n"+"    c();\n"+"}")

	test("var v = [\"a\",\n" + "    function() {\n" + "        return;\n" + "    }, {\n" + "        id: 1\n" + "    }\n" + "];")
	test("var v = [\"a\", function() {\n" + "    return;\n" + "}, {\n" + "    id: 1\n" + "}];")

	test("module \"Even\" {\n" + "    import odd from \"Odd\";\n" + "    export function sum(x, y) {\n" + "        return x + y;\n" + "    }\n" + "    export var pi = 3.141593;\n" + "    export default moduleName;\n" + "}")
	test("module \"Even\" {\n" + "    export default function div(x, y) {}\n" + "}")

	test("set[\"name\"]")
	test("get[\"name\"]")
	test("a = {\n" + "    set b(x) {},\n" + "    c: 1,\n" + "    d: function() {}\n" + "};")
	test("a = {\n" + "    get b() {\n" + "        retun 0;\n" + "    },\n" + "    c: 1,\n" + "    d: function() {}\n" + "};")

	test("'use strict';\n" +
		"if ([].some(function() {\n" +
		"        return false;\n" +
		"    })) {\n" +
		"    console.log('hello');\n" +
		"}")

	test("class Test {\n" +
		"    blah: string[];\n" +
		"    foo(): number {\n" +
		"        return 0;\n" +
		"    }\n" +
		"    bar(): number {\n" +
		"        return 0;\n" +
		"    }\n" +
		"}")
	test("interface Test {\n" +
		"    blah: string[];\n" +
		"    foo(): number {\n" +
		"        return 0;\n" +
		"    }\n" +
		"    bar(): number {\n" +
		"        return 0;\n" +
		"    }\n" +
		"}")

	test("var a=1,b={bang:2},c=3;", "var a = 1,\n    b = {\n        bang: 2\n    },\n    c = 3;")
	test("var a={bing:1},b=2,c=3;", "var a = {\n        bing: 1\n    },\n    b = 2,\n    c = 3;")
}

func test(options ...string) {
	if len(options) == 1 {
		assertMatch(options[0], options[0])
	} else if len(options) == 2 {
		assertMatch(options[0], options[1])
	} else {
		ts.Error("Cannot test for nothing or more than 1 input")
	}
}

func assertMatch(input, expect string) {
	result, _ := Beautify(&input, test_options)

	if result != expect {
		ts.Error("Input", input, "Result: ", result, " did not match ", expect)
	}

}
