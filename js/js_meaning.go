package js

import (
	"github.com/tapvanvn/gotokenize"
)

func CreateJSMeaning() *gotokenize.PatternMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{

		"#%^&*-+/!<>=?:@\"'` \\;\r\n\t{}[](),.|": {TokenType: TokenJSOperator, Separate: true},
		//"0123456789":   {TokenType: TokenJSNumber, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	jsRawMeaning := NewJSRawMeaning()

	jsRawMeaning.SetSource(meaning)

	jsPhraseMeaning := NewJSPhraseMeaning(jsRawMeaning)

	return gotokenize.CreatePatternMeaning(jsPhraseMeaning, JSPatterns, gotokenize.NoTokens, JSGlobalNested)
}
