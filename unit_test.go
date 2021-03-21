package gotokenize_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/gotokenize/json"
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
