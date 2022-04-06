package gotokenize_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/gotokenize/css"
	"github.com/tapvanvn/gotokenize/js"
	"github.com/tapvanvn/gotokenize/json"
	"github.com/tapvanvn/gotokenize/xml"
)

func TestStream(t *testing.T) {
	content := "test=abc  ,\n  test2=def"
	stream := gotokenize.CreateStream()
	stream.Tokenize(content)
	checkContent := ""
	iter := stream.Iterator()
	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()
		checkContent += token.Content
	}
	if content != checkContent {
		t.Fail()
	}
}

func TestRawMeaning(t *testing.T) {
	content := "test=abc  test2=def"
	tokenMap := map[string]gotokenize.RawTokenDefine{
		"=": {
			TokenType: 1,
			Separate:  true,
		},
	}
	stream := gotokenize.CreateStream()
	stream.Tokenize(content)
	meaning := gotokenize.CreateRawMeaning(tokenMap, true)
	meaning.Prepare(&stream)
	for {
		token := meaning.Next()
		if token == nil {
			break
		}
		if token.Content == "=" && token.Type != 1 {
			t.Fail()
		}
	}
}

func TestPatternMeaning(t *testing.T) {
	content := "test=abc  test2=def"
	patterns := []gotokenize.Pattern{

		{
			Type: 100,
			Struct: []gotokenize.PatternToken{

				{Type: 0},
				{Content: "="},
				{Type: 0},
			},
			IsRemoveGlobalIgnore: true,
		},
	}
	stream := gotokenize.CreateStream()
	stream.Tokenize(content)

	tokenMap := map[string]gotokenize.RawTokenDefine{
		"=": {TokenType: 1, Separate: true},
		" ": {TokenType: 2, Separate: false},
	}

	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	patternMeaning := gotokenize.CreatePatternMeaning(&meaning, patterns, gotokenize.NoTokens, gotokenize.NoTokens)

	patternMeaning.Prepare(&stream)

	for {
		token := patternMeaning.Next()
		if token == nil {
			break
		}
		fmt.Println(token.Type, token.Content)
	}
}

func TestJSONMeaning(t *testing.T) {
	content := `{
		"user_name": "test",
		"age":30,
		"asset":["gold","silver","land"]
	}`

	stream := gotokenize.CreateStream()
	stream.Tokenize(content)

	meaning := json.CreateJSONMeaning()
	meaning.Prepare(&stream)

	token := meaning.Next()
	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, "[", json.JSONNaming(token.Type), "]")
		if token.Children.Length() > 0 {
			token.Children.Debug(1, json.JSONNaming)
		}
		token = meaning.Next()
	}

	//meaning.GetStream().Debug(0, json.JSONNaming)
}

func TestXMLRawMeaning(t *testing.T) {
	content := `<xml abc="def">
		<next/>
		<name>tapvanvn</name>
		<debug>{{ahshsdfkjlsdf}}</debug>
		<!--
			comment here
		-->
	</xml>`

	stream := gotokenize.CreateStream()
	stream.Tokenize(content)

	meaning := xml.CreateXMLRawMeaning()
	meaning.Prepare(&stream)

	token := meaning.Next()

	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, "[", xml.XMLNaming(token.Type), "]", token.Content)
		if token.Children.Length() > 0 {
			token.Children.Debug(1, xml.XMLNaming)
		}
		token = meaning.Next()
	}
}

func TestXMLMeaning(t *testing.T) {
	content := `<xml abc="def">
		<next/>
		<name>tapvanvn</name>
		<utf-8>kiểm tra</utf-8>
		<debug>{{ahshsdfkjlsdf}}</debug>
		<!--
			comment here
		-->
	</xml>`

	stream := gotokenize.CreateStream()
	stream.Tokenize(content)

	meaning := xml.CreateXMLMeaning()
	meaning.Prepare(&stream)

	token := meaning.Next()

	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, "[", xml.XMLNaming(token.Type), "]", token.Content)
		if token.Children.Length() > 0 {
			token.Children.Debug(1, xml.XMLNaming)
		}
		token = meaning.Next()
	}
}

func TestCSSMeaning(t *testing.T) {
	content := `
	@media only screen and (max-width: 900px) {
		.mobile_gone {
			display: none !important;
		}
	}
	
	[type="input"]{
		position: relative;
		border-bottom-width: 1px;
		border-bottom-style: dotted;
		padding-top: 5px;
		padding-bottom: 5px;
		border-bottom-color: gray;
	}`

	stream := gotokenize.CreateStream()
	stream.Tokenize(content)

	meaning := css.CreateCSSMeaning()
	meaning.Prepare(&stream)

	token := meaning.Next()

	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, "[", css.CSSNaming(token.Type), "]", token.Content)
		if token.Children.Length() > 0 {
			token.Children.Debug(1, css.CSSNaming)
		}
		token = meaning.Next()
	}
}

func TestJSMeaning(t *testing.T) {
	content := `
	var a = b //comment
	var c = d
	if(a==b)
		d
	function def() {
		()=>{bef}
		a = c
	}`

	stream := gotokenize.CreateStream()
	stream.Tokenize(content)

	meaning := js.CreateJSMeaning()
	meaning.Prepare(&stream)

	token := meaning.Next()

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName)
		token = meaning.Next()
	}
}
