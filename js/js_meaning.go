package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

var JSCharacterMap = map[string]gotokenize.RawTokenDefine{

	"#%^&*-+/!<>=?:@\"'` \\;\r\n\t{}[](),.|": {TokenType: TokenJSOperator, Separate: true},
}

func NewDefaultJSRawMeaning() *JSRawMeaning {

	meaning := gotokenize.CreateRawMeaning(JSCharacterMap, false)

	jsRawMeaning := NewJSRawMeaning()

	jsRawMeaning.SetSource(meaning)

	return jsRawMeaning
}

func NewDefaultJSMeaning() *gotokenize.PatternMeaning {

	jsRawMeaning := NewDefaultJSRawMeaning()

	return gotokenize.NewPatternMeaning(jsRawMeaning, JSPatterns, gotokenize.NoTokens, JSGlobalNested)
}

func NewDefaultJSInstructionMeaning() *JSInstructionMeaning {

	jsPatternMeaning := NewDefaultJSMeaning()

	return NewJSInstructionMeaning(jsPatternMeaning)
}
