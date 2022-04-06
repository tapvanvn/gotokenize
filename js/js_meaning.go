package js

import (
	"github.com/tapvanvn/gotokenize"
)

func CreateJSRawMeaning() *JSRawMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{

		"#%^&*-+/!<>=?:@\"'` \\;\r\n\t{}[](),.|": {TokenType: TokenJSOperator, Separate: true},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	jsRawMeaning := NewJSRawMeaning()

	jsRawMeaning.SetSource(meaning)

	return jsRawMeaning
}

func CreateJSPhraseMeaning() *JSPhraseMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{

		"#%^&*-+/!<>=?:@\"'` \\;\r\n\t{}[](),.|": {TokenType: TokenJSOperator, Separate: true},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	jsRawMeaning := NewJSRawMeaning()

	jsRawMeaning.SetSource(meaning)

	jsPhraseMeaning := NewJSPhraseMeaning(jsRawMeaning)

	return jsPhraseMeaning
}

func CreateJSMeaning() *gotokenize.PatternMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{

		"#%^&*-+/!<>=?:@\"'` \\;\r\n\t{}[](),.|": {TokenType: TokenJSOperator, Separate: true},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	jsRawMeaning := NewJSRawMeaning()

	jsRawMeaning.SetSource(meaning)

	jsPhraseMeaning := NewJSPhraseMeaning(jsRawMeaning)

	return gotokenize.CreatePatternMeaning(jsPhraseMeaning, JSPatterns, gotokenize.NoTokens, JSGlobalNested)
}
