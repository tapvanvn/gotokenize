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
func NewDefaultJSPhraseMeaning() *JSPhraseMeaning {

	return NewJSPhraseMeaning(NewDefaultJSRawMeaning())
}

func NewDefaultJSPatternMeaning() gotokenize.IMeaning {

	var last gotokenize.IMeaning = NewDefaultJSOperatorMeaning()

	for _, pattern := range JSPatterns {

		last = gotokenize.NewPatternMeaning(last, pattern.Patterns, pattern.IgnoreTokens, pattern.TokenCanNested)
	}
	return last
}

func NewDefaultJSMeaning() gotokenize.IMeaning {

	//return NewJSPhraseMeaning(NewDefaultJSPatternMeaning())
	return NewDefaultJSPatternMeaning()
}

func NewDefaultJSOperatorMeaning() *JSOperatorMeaning {

	return NewJSOperatorMeaning(NewDefaultJSPhraseMeaning())
}

func NewDefaultJSInstructionMeaning() *JSInstructionMeaning {

	jsPatternMeaning := NewDefaultJSMeaning()

	return NewJSInstructionMeaning(jsPatternMeaning)
}
