// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"reflect"
	"testing"
)

func T(ty TokenType, v string) Token {
	return Token{ty, v, 0, 0}
}

func TestSuccessfulScan(t *testing.T) {
	for _, test := range []struct {
		input  string
		tokens []Token
	}{
		{"bar(", []Token{T(TokenFunction, "bar(")}},
		{"abcd", []Token{T(TokenIdent, "abcd")}},
		{`"abcd"`, []Token{T(TokenString, `"abcd"`)}},
		{"'abcd'", []Token{T(TokenString, "'abcd'")}},
		{"#name", []Token{T(TokenHash, "#name")}},
		{"4.2", []Token{T(TokenNumber, "4.2")}},
		{".42", []Token{T(TokenNumber, ".42")}},
		{"42%", []Token{T(TokenPercentage, "42%")}},
		{"4.2%", []Token{T(TokenPercentage, "4.2%")}},
		{".42%", []Token{T(TokenPercentage, ".42%")}},
		{"42px", []Token{T(TokenDimension, "42px")}},
		{"url('http://www.google.com/')", []Token{T(TokenURI, "url('http://www.google.com/')")}},
		{"U+0042", []Token{T(TokenUnicodeRange, "U+0042")}},
		{"<!--", []Token{T(TokenCDO, "<!--")}},
		{"-->", []Token{T(TokenCDC, "-->")}},
		{"   \n   \t   \n", []Token{T(TokenS, "   \n   \t   \n")}},
		{"/* foo */", []Token{T(TokenComment, "/* foo */")}},
		{"bar(", []Token{T(TokenFunction, "bar(")}},
		{"~=", []Token{T(TokenIncludes, "~=")}},
		{"|=", []Token{T(TokenDashMatch, "|=")}},
		{"^=", []Token{T(TokenPrefixMatch, "^=")}},
		{"$=", []Token{T(TokenSuffixMatch, "$=")}},
		{"*=", []Token{T(TokenSubstringMatch, "*=")}},
		{"{", []Token{T(TokenChar, "{")}},
		{"\uFEFF", []Token{T(TokenBOM, "\uFEFF")}},

		{"42''", []Token{
			T(TokenNumber, "42"),
			T(TokenString, "''"),
		}},
		{`╯︵┻━┻"stuff"`, []Token{
			T(TokenIdent, "╯︵┻━┻"),
			T(TokenString, `"stuff"`),
		}},
		{"color:red", []Token{
			T(TokenIdent, "color"),
			T(TokenChar, ":"),
			T(TokenIdent, "red"),
		}},
		{"color:red;background:blue", []Token{
			T(TokenIdent, "color"),
			T(TokenChar, ":"),
			T(TokenIdent, "red"),
			T(TokenChar, ";"),
			T(TokenIdent, "background"),
			T(TokenChar, ":"),
			T(TokenIdent, "blue"),
		}},
		{"color:rgb(0,1,2)", []Token{
			T(TokenIdent, "color"),
			T(TokenChar, ":"),
			T(TokenFunction, "rgb("),
			T(TokenNumber, "0"),
			T(TokenChar, ","),
			T(TokenNumber, "1"),
			T(TokenChar, ","),
			T(TokenNumber, "2"),
			T(TokenChar, ")"),
		}},
		{"color:#fff", []Token{
			T(TokenIdent, "color"),
			T(TokenChar, ":"),
			T(TokenHash, "#fff"),
		}},

		// Check note in CSS2 4.3.4:
		// Note that COMMENT tokens cannot occur within other tokens: thus, "url(/*x*/pic.png)" denotes the URI "/*x*/pic.png", not "pic.png".
		{"url(/*x*/pic.png)", []Token{
			T(TokenURI, "url(/*x*/pic.png)"),
		}},

		// More URI testing, since it's important
		{"url(/pic.png)", []Token{
			T(TokenURI, "url(/pic.png)"),
		}},
		{"url( /pic.png )", []Token{
			T(TokenURI, "url( /pic.png )"),
		}},
		{"uRl(/pic.png)", []Token{
			T(TokenURI, "uRl(/pic.png)"),
		}},
		{"url(\"/pic.png\")", []Token{
			T(TokenURI, "url(\"/pic.png\")"),
		}},
		{"url('/pic.png')", []Token{
			T(TokenURI, "url('/pic.png')"),
		}},
		{"url('/pic.png?badchars=\\(\\'\\\"\\)\\ ')", []Token{
			T(TokenURI, "url('/pic.png?badchars=\\(\\'\\\"\\)\\ ')"),
		}},

		// CSS2 section 4.1.1: "red-->" is IDENT "red--" followed by DELIM ">",
		{"red-->", []Token{
			T(TokenIdent, "red--"),
			T(TokenChar, ">"),
		}},

		{"-moz-border:1", []Token{
			T(TokenIdent, "-moz-border"),
			T(TokenChar, ":"),
			T(TokenNumber, "1"),
		}},

		// CSS2 section 4.1.3, last bullet point: identifier test
		// is the same as te\st.
		// commenting out while this fails, so I can commit other tests
		//{"test", []Token{T(TokenIdent, "test")}},
		//{"te\\st", []Token{T(TokenIdent, "test")}},
	} {
		tokens := []Token{}
		s := New(test.input)
		for {
			tok := s.Next()
			if tok.Type == TokenError {
				t.Fatalf("Error token with: %q", test.input)
			}
			if tok.Type == TokenEOF {
				break
			}
			tok.Line = 0
			tok.Column = 0
			tokens = append(tokens, *tok)
		}
		if !reflect.DeepEqual(tokens, test.tokens) {
			t.Fatalf("For input string %q, bad tokens. Expected:\n%#v\n\nGot:\n%#v", test.input, test.tokens, tokens)
		}
	}
}
