package css

import "github.com/tapvanvn/gotokenize/v2"

const (
	TokenCSSOperator     = 1
	TokenCSSString       = 2
	TokenCSSBlockComment = 3
	TokenCSSSpace        = 4

	TokenCSSAttribute          = 100
	TokenCSSAttributeImportain = 101
	TokenCSSBlock              = 102
	TokenCSSSquare             = 103
	TokenCSSHeader             = 104

	TokenCSSDefine = 1000
)

var CSSGlobalNested = []int{
	TokenCSSBlock,
}

var CSSIgnores = []int{}

func CSSNaming(tokenType int) string {
	switch tokenType {
	case TokenCSSOperator:
		return "operator"
	case TokenCSSString:
		return "string"
	case TokenCSSBlockComment:
		return "comment"
	case TokenCSSSpace:
		return "space"

	case TokenCSSAttribute:
		return "attribute"
	case TokenCSSAttributeImportain:
		return "attribute!"
	case TokenCSSBlock:
		return "block"
	case TokenCSSSquare:
		return "square"
	case TokenCSSHeader:
		return "header"

	case TokenCSSDefine:
		return "define"
	}
	return "unknown"
}

var CSSPatterns = []gotokenize.Pattern{}
