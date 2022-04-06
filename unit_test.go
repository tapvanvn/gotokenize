package gotokenize_test

import (
	"testing"

	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/css"
	"github.com/tapvanvn/gotokenize/v2/js"
	"github.com/tapvanvn/gotokenize/v2/json"
	"github.com/tapvanvn/gotokenize/v2/xml"
)

func TestStream(t *testing.T) {
	content := "test=abc  ,\n  test2=def"
	stream := gotokenize.CreateStream(0)
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
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(&stream)
	meaning := gotokenize.CreateRawMeaning(tokenMap, true)
	meaning.Prepare(proc)
	for {
		token := meaning.Next(proc)
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
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	tokenMap := map[string]gotokenize.RawTokenDefine{
		"=": {TokenType: 1, Separate: true},
		" ": {TokenType: 2, Separate: false},
	}

	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	patternMeaning := gotokenize.CreatePatternMeaning(meaning, patterns, []int{2}, gotokenize.NoTokens)

	proc := gotokenize.NewMeaningProcessFromStream(&stream)

	patternMeaning.Prepare(proc)

	for {
		token := patternMeaning.Next(proc)
		if token == nil {
			break
		}
		token.Debug(0, nil)

		if token.Type != 100 {
			t.Fail()
		}
	}
	gotokenize.DebugMeaning(patternMeaning)
}

func TestJSONMeaning(t *testing.T) {
	content := `{
		"user_name": "test",
		"age":30,
		"asset":["gold","silver","land"]
	}`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(&stream)
	meaning := json.CreateJSONMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		token.Children.Debug(0, json.JSONNaming)
	}
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

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(&stream)
	meaning := xml.NewXMLHighMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}

		token.Children.Debug(0, xml.XMLNaming)

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

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(&stream)
	meaning := xml.NewXMLHighMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}

		token.Children.Debug(0, xml.XMLNaming)
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

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(&stream)
	meaning := css.CreateCSSMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		token.Children.Debug(0, css.CSSNaming)
	}
}
func TestJSRawMeaning(t *testing.T) {
	content := `

	function  {
		()=>{bef}
		a = c
	}`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	meaning := js.CreateJSRawMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(&stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}

func TestJSPhraseMeaning(t *testing.T) {
	content := `

	function def() {
		()=>{bef}
		a = c
		for(var a = 0; a< 10; a++) {
			b()
		}
	}`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	meaning := js.CreateJSPhraseMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(&stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}

func TestJSMeaning(t *testing.T) {
	content := `

	var IntSet = class {
		constructor() {
		  this.array = new Array();
		}
		add(value) {
		  let contains = this.contains(value);
		  this.array[value | 0] = value | 0;
		  return !contains;
		}
		contains(value) {
		  return this.array[value | 0] != void 0;
		}
		remove(value) {
		  this.array[value | 0] = void 0;
		}
		clear() {
		  this.array.length = 0;
		}
	};`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	meaning := js.CreateJSMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(&stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}
