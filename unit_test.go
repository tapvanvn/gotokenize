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
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)
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

	patternMeaning := gotokenize.NewPatternMeaning(meaning, patterns, []int{2}, gotokenize.NoTokens)

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)

	patternMeaning.Prepare(proc)

	for {
		token := patternMeaning.Next(proc)
		if token == nil {
			break
		}
		token.Debug(0, nil, nil)

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
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)
	meaning := json.CreateJSONMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		token.Children.Debug(0, json.JSONNaming, gotokenize.EmptyDebugOption)
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
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)
	meaning := xml.NewXMLHighMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}

		token.Children.Debug(0, xml.XMLNaming, gotokenize.EmptyDebugOption)

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
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)
	meaning := xml.NewXMLHighMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}

		token.Children.Debug(0, xml.XMLNaming, gotokenize.EmptyDebugOption)
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
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)
	meaning := css.CreateCSSMeaning()
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		token.Children.Debug(0, css.CSSNaming, gotokenize.EmptyDebugOption)
	}
}
func TestJSRawMeaning(t *testing.T) {
	content := `

	var spine = (()=>{
		var __defProp = Object.defineProperty;
		var __markAsModule = (target) => __defProp(target, "__esModule", { value: true });
		var __export = (target, all) => {
		  __markAsModule(target);
		  for (var name in all)
			__defProp(target, name, { get: all[name], enumerable: true });
		};
		})()`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	meaning := js.NewDefaultJSRawMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName, js.JSDebugOptions)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}

func TestJSMeaning(t *testing.T) {
	content := `
	var spine = (()=>{
		var __defProp = Object.defineProperty;
		var __markAsModule = (target) => __defProp(target, "__esModule", { value: true });
		var __export = (target, all) => {
		  __markAsModule(target);
		  for (var name in all)
			__defProp(target, name, { get: all[name], enumerable: true });
		};
		})()
	`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	meaning := js.NewDefaultJSMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName, js.JSDebugOptions)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}

func TestJSInstructionMeaning(t *testing.T) {
	content := `
	for (var name in all)
		__defProp(target, name, { get: all[name], enumerable: true });`

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)

	meaning := js.NewDefaultJSInstructionMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName, js.JSDebugOptions)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}
