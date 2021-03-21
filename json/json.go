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

func JSONNaming(tokenType int) string {
	switch tokenType {
	case 0:
		return "word"
	case TokenJSONNumber:
		return "number"
	case TokenJSONNumberString:
		return "number_string"
	case TokenJSONBlock:
		return "object"
	case TokenJSONSquare:
		return "array"
	case TokenJSONPair:
		return "pair"
	case TokenJSONString:
		return "string"
	}
	return "unknow"
}

var JSONGlobalNested = []int{TokenJSONBlock}
