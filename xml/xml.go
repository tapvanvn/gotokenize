package xml

import "github.com/tapvanvn/gotokenize/v2"

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
	TokenXMLSingleElement     = 1001
	TokenXMLEndElement        = 1002
	TokenXMLElementAttributes = 1003
	TokenXMLElementBody       = 1004
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
		return "xml_space"
	case TokenXMLString:
		return "xml_string"
	case TokenXMLOperator:
		return "xml_operator"
	case TokenXMLComment:
		return "xml_comment"
	case TokenXMLAttribute:
		return "xml_attribute"
	case TokenXMLTagBegin:
		return "xml_tag begin"
	case TokenXMLTagEnd:
		return "xml_tag end"
	case TokenXMLTagUnknown:
		return "xml_single tag"
	case TokenXMLEndElement:
		return "xml_end element"
	case TokenXMLElement:
		return "xml_element"
	case TokenXMLElementAttributes:
		return "xml_attribute set"
	case TokenXMLElementBody:
		return "xml_body"
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
