package json

import (
	"github.com/tapvanvn/gotokenize"
)

const (
	TokenJSONNumber       = 1
	TokenJSONOperator     = 2
	TokenJSONString       = 3
	TokenJSONNumberString = 4

	TokenJSONBlock  = 10
	TokenJSONSquare = 11

	TokenJSONPair = 100
)

var JSONPatterns = []gotokenize.Pattern{

	{
		Type: TokenJSONPair,
		Struct: []gotokenize.PatternToken{
			{Type: TokenJSONString},
			{Content: ":", IsIgnoreInResult: true},
			{Type: TokenJSONBlock, CanNested: true},
		},
		IsRemoveGlobalIgnore: true,
	},

	{
		Type: TokenJSONPair,
		Struct: []gotokenize.PatternToken{

			{Type: TokenJSONString},
			{Content: ":", IsIgnoreInResult: true},
			{Type: TokenJSONNumberString},
		},
		IsRemoveGlobalIgnore: true,
	},

	{
		Type: TokenJSONPair,
		Struct: []gotokenize.PatternToken{

			{Type: TokenJSONString},
			{Content: ":", IsIgnoreInResult: true},
			{Type: TokenJSONString},
		},
		IsRemoveGlobalIgnore: true,
	},
	{
		Type: TokenJSONPair,
		Struct: []gotokenize.PatternToken{

			{Type: TokenJSONString},
			{Content: ":", IsIgnoreInResult: true},
			{Type: TokenJSONSquare},
		},
		IsRemoveGlobalIgnore: true,
	},
}
