package xml

import "github.com/tapvanvn/gotokenize"

const (
	TokenXMLSpace    = 1
	TokenXMLString   = 2
	TokenXMLOperator = 3

	TokenXMLAttribute  = 100
	TokenXMLTagBegin   = 101
	TokenXMLTagEnd     = 102
	TokenXMLTagUnknown = 103
	TokenXMLComment    = 104

	TokenXMLElement           = 1000
	TokenXMLElementAttributes = 1001
	TokenXMLElementBody       = 1002
)

var XMLGlobalNested = []int{
	TokenXMLTagBegin,
	TokenXMLTagUnknown,
}
var XMLIgnores = []int{
	TokenXMLComment,
}

func XMLNaming(tokenType int) string {
	switch tokenType {
	case 0:
		return "word"
	case TokenXMLSpace:
		return "space"
	case TokenXMLString:
		return "string"
	case TokenXMLOperator:
		return "operator"
	case TokenXMLComment:
		return "comment"
	case TokenXMLAttribute:
		return "attribute"
	case TokenXMLTagBegin:
		return "tag begin"
	case TokenXMLTagEnd:
		return "tag end"
	case TokenXMLTagUnknown:
		return "single tag"
	case TokenXMLElement:
		return "element"
	case TokenXMLElementAttributes:
		return "attribute set"
	case TokenXMLElementBody:
		return "body"
	}
	return "unknown"
}

var XMLPatterns = []gotokenize.Pattern{
	//pattern attribute "key"="value"
	{
		Type: TokenXMLAttribute,
		Struct: []gotokenize.PatternToken{
			{Type: TokenXMLString},
			{Content: "=", IsIgnoreInResult: true},
			{Type: TokenXMLString},
		},
		IsRemoveGlobalIgnore: true,
	},
	//pattern attbute key="value"
	{
		Type: TokenXMLAttribute,
		Struct: []gotokenize.PatternToken{
			{Type: 0},
			{Content: "=", IsIgnoreInResult: true},
			{Type: TokenXMLString},
		},
		IsRemoveGlobalIgnore: true,
	},
	//pattern attbute key=value
	{
		Type: TokenXMLAttribute,
		Struct: []gotokenize.PatternToken{
			{Type: 0},
			{Content: "=", IsIgnoreInResult: true},
			{Type: 0},
		},
		IsRemoveGlobalIgnore: true,
	},
	//pattern attbute "key"=value
	{
		Type: TokenXMLAttribute,
		Struct: []gotokenize.PatternToken{
			{Type: TokenXMLString},
			{Content: "=", IsIgnoreInResult: true},
			{Type: 0},
		},
		IsRemoveGlobalIgnore: true,
	},
}
