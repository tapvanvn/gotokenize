package js

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

const (
	TokenJSUnknown     = 0
	TokenJSWord        = 1
	TokenJSOperator    = 2
	TokenJSPhraseBreak = 3

	TokenJSKeyWord          = 100
	TokenJSSoftBreak        = 101 // ,
	TokenJSPhrase           = 102
	TokenJSBracket          = 103 //()
	TokenJSBlock            = 104 //{}
	TokenJSBracketSquare    = 105 //[]
	TokenJSUnaryOperator    = 106 // !, ~, ++, --
	TokenJSBinaryOperator   = 107 // <>+-*/%, <=, >=, <<, >>, >>>, ||, |, &&, &, ^, **, ==, ===, !=, !==
	TokenJSTreeDotOperator  = 108
	TokenJSQuestionOperator = 109 // ?
	TokenJSAssign           = 110 // =
	TokenJSRightArrow       = 111 // =>
	TokenJSLineComment      = 112 // //
	TokenJSBlockComment     = 113 // /**/
	TokenJSRegex            = 114 // //img

	TokenJSFunction       = 200
	TokenJSFunctionLambda = 201
	TokenJSVariable       = 202
	TokenJSString         = 203
	TokenJSFor            = 204
	TokenJSIf             = 205
	TokenJSElseIf         = 206
	TokenJSElse           = 207
	TokenJSInlineIfBody   = 208
	TokenJSSwitch         = 209
	TokenJSCase           = 210
	TokenJSDefault        = 211
	TokenJSBreak          = 212
	TokenJSWhile          = 213
	TokenJSDo             = 214
	TokenJSClass          = 215
	TokenJSClassFunction  = 216
	TokenJSIfTrail        = 217
	TokenJSObjectProperty = 218

	//TokenJSStrongBreak           = 300 //sure `;`
	//TokenJSArgumentList          = 301 //use in function declaration and function call
	//TokenJSStrongDeclareVariable = 302 //mark that the follow variable had been strong declared with `var` keyword
	//TokenJSStrongDeclareConstant = 303 //mark that the follow variable had been strong declared with `const` keyword
	//TokenJSStrongDeclareLet      = 304 //mark that the follow variable had been strong declared with `let` keyword
	TokenJSInlineIf       = 350
	TokenJSAssignVariable = 351 //declare a new or assign value to variable
	//TokenJSDeclareFunction = 351 //declare a new function (traditional or lambda)
	//TokenJSDeclareClass    = 352 //declare a class using class keyword
	//TokenJSDeclareObject   = 353 //declare an object
	//TokenJSDeclareArray    = 354 //declare an array using [] or new Array
	TokenJSVoidStatement = 355 //any void(...) void ...
	TokenJSOperatorTrail = 356
)

var JSTokenNameDictionary = map[int]string{
	TokenJSUnknown:          "unknown",
	TokenJSLineComment:      "line comment",
	TokenJSBlockComment:     "block comment",
	TokenJSWord:             "word",
	TokenJSKeyWord:          "key word",
	TokenJSPhrase:           "phrase",
	TokenJSRegex:            "regex",
	TokenJSOperator:         "operator",
	TokenJSUnaryOperator:    "unary operator",
	TokenJSBinaryOperator:   "binary operator",
	TokenJSTreeDotOperator:  "treedot operator",
	TokenJSQuestionOperator: "question operator",
	TokenJSAssign:           "assign",
	TokenJSRightArrow:       "right narrow",
	TokenJSPhraseBreak:      "phrase break",
	TokenJSSoftBreak:        "soft break",
	TokenJSBracket:          "bracket",
	TokenJSBracketSquare:    "square",
	TokenJSBlock:            "block",
	//
	TokenJSVariable:       "variable",
	TokenJSString:         "string",
	TokenJSFor:            "for",
	TokenJSIf:             "if",
	TokenJSElseIf:         "elseif",
	TokenJSElse:           "else",
	TokenJSInlineIfBody:   "inline if body",
	TokenJSSwitch:         "switch",
	TokenJSCase:           "case",
	TokenJSDefault:        "default",
	TokenJSBreak:          "break",
	TokenJSWhile:          "while",
	TokenJSDo:             "do",
	TokenJSFunction:       "function",
	TokenJSFunctionLambda: "lambda",
	TokenJSClass:          "class",
	TokenJSClassFunction:  "class function",

	TokenJSIfTrail:        "if trail",
	TokenJSOperatorTrail:  "operator trail",
	TokenJSObjectProperty: "property",
	TokenJSInlineIf:       "inline if",
	TokenJSAssignVariable: "assign variable",
	TokenJSVoidStatement:  "void statement",
}

//JSTokenName return name from type of token
func JSTokenName(tokenType int) string {

	if name, ok := JSTokenNameDictionary[tokenType]; ok {
		return name
	}

	return fmt.Sprintf("unknown-%d", tokenType)
}

//JSKeyWords keywords of javascript
var JSKeyWords string = `
,abstract,arguments,await,boolean,
,break,byte,case,catch,
,char,class,const,continue,
,debugger,default,delete,do,
,double,else,enum,eval,
,export,extends,false,final,finally,float,for,function,
,goto,if,implements,import,
,in,instanceof,int,interface,
,let,long,native,new,
,null,package,private,protected,
,public,return,short,static,
,super,switch,synchronized,this,
,throw,throws,transient,true,
,try,typeof,var,void,
,volatile,while,with,yield,constructor,`

var JSDebugOptions = &gotokenize.DebugOption{
	StringifyTokens: []int{
		TokenJSString,
		TokenJSRegex,
	},
}

//JSIgnores tokens that will be ignore
var JSIgnores = []int{

	TokenJSLineComment,
	TokenJSBlockComment,
	TokenJSPhraseBreak,
}

var JSPhraseAllow = []int{
	TokenJSWord,
	TokenJSOperator,
	//TokenJSAssign,
	TokenJSUnaryOperator,
	TokenJSBinaryOperator,
	TokenJSRegex,
	TokenJSString,
	TokenJSVariable,
	TokenJSBlock,
	TokenJSBracket,
	TokenJSBracketSquare,
	TokenJSRightArrow,
	TokenJSFunctionLambda,
	TokenJSFunction,
	TokenJSClass,
	TokenJSTreeDotOperator,
}
var JSGlobalNested = []int{
	TokenJSBlock,
	TokenJSBracket,
	TokenJSBracketSquare,
	TokenJSClass,
	TokenJSPhrase,
}

//JSPatternOperator the patterns to detect instruction structure
var JSPatternOperator = gotokenize.PatternMeaningDefine{
	IgnoreTokens:   gotokenize.NoTokens,
	TokenCanNested: JSGlobalNested,
	Patterns: []gotokenize.Pattern{
		//Inline If
		{
			Type:                 TokenJSPhrase,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Type: TokenJSBinaryOperator},
			},
		},
	},
}

var JSLevel2GlobalNested = append(
	JSGlobalNested,
	TokenJSOperatorTrail,
	TokenJSSwitch,
	TokenJSDo,
	TokenJSWhile,
	TokenJSIf,
	TokenJSElseIf,
	TokenJSElse,
	TokenJSFunction,
	TokenJSFunctionLambda,
	TokenJSPhrase,
	TokenJSClass,
	TokenJSClassFunction,
)

var JSPatternLevel1 = gotokenize.PatternMeaningDefine{
	IgnoreTokens:   gotokenize.NoTokens,
	TokenCanNested: JSGlobalNested,
	Patterns: []gotokenize.Pattern{

		//inline if
		{
			Type:                 TokenJSInlineIf,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{IsAny: true, CanNested: true},
				{Type: TokenJSQuestionOperator, IsIgnoreInResult: true},
				{IsAny: true, CanNested: true},
				{Content: ":", IsIgnoreInResult: true},
				{IsAny: true, CanNested: true},
			},
		},
		//Function call
		//pattern if block
		{
			Type:                 TokenJSIf,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "if", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
				{IsAny: true, CanNested: true},
			},
		},
		//pattern else if block
		{
			Type:                 TokenJSElseIf,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "else", IsIgnoreInResult: true},
				{Content: "if", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
				{IsAny: true, CanNested: true},
			},
		},
		//pattern else block
		{
			Type:                 TokenJSElse,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "else", IsIgnoreInResult: true},
				{IsAny: true, CanNested: true},
			},
		},
		//pattern for
		{
			Type:                 TokenJSFor,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "for", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
				{IsAny: true, CanNested: true},
			},
		},
		//pattern switch
		{
			Type:                 TokenJSSwitch,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "switch", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
				{Type: TokenJSBlock, CanNested: true},
			},
		},
		//pattern do block
		{
			Type:                 TokenJSDo,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "do", IsIgnoreInResult: true},
				{IsAny: true, CanNested: true},
				{Content: "while", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
			},
		},
		//pattern while block
		{
			Type:                 TokenJSWhile,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "while", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
				{IsAny: true, CanNested: true},
			},
		},
		/*
			// Object Properties with word
			{
				Type:                 TokenJSObjectProperty,
				IsRemoveGlobalIgnore: true,
				Struct: []gotokenize.PatternToken{
					{Type: TokenJSWord},
					{Content: ":", IsIgnoreInResult: true},
					{IsAny: true, CanNested: true},
					{Content: ",", IsIgnoreInResult: true},
				},
				LivingContext: []int{TokenJSBlock},
			},
			// Object Properties with string
			{
				Type:                 TokenJSObjectProperty,
				IsRemoveGlobalIgnore: true,
				Struct: []gotokenize.PatternToken{
					{Type: TokenJSString},
					{Content: ":", IsIgnoreInResult: true},
					{IsAny: true, CanNested: true},
				},
				LivingContext: []int{TokenJSBlock},
			},
			// Object last Properties with word
			{
				Type:                 TokenJSObjectProperty,
				IsRemoveGlobalIgnore: true,
				Struct: []gotokenize.PatternToken{
					{Type: TokenJSWord},
					{Content: ":", IsIgnoreInResult: true},
					{IsAny: true, CanNested: true},
				},
				LivingContext: []int{TokenJSBlock},
			},
			// Object last Properties with string
			{
				Type:                 TokenJSObjectProperty,
				IsRemoveGlobalIgnore: true,
				Struct: []gotokenize.PatternToken{
					{Type: TokenJSString},
					{Content: ":", IsIgnoreInResult: true},
					{IsAny: true, CanNested: true},
				},
				LivingContext: []int{TokenJSBlock},
			},*/
		//pattern function with keyword and name
		{
			Type:                 TokenJSFunction,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "function", IsIgnoreInResult: true},
				{Type: TokenJSPhrase, CanNested: true, Nested: []gotokenize.PatternToken{
					{Type: TokenJSWord},
					{Type: TokenJSBracket, CanNested: true},
				}},
				{Type: TokenJSBlock, CanNested: true},
			},
		},
		//pattern function with keyword without name
		{
			Type:                 TokenJSFunction,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "function", IsIgnoreInResult: true},
				{Type: TokenJSBracket, CanNested: true},
				{Type: TokenJSBlock, CanNested: true},
			},
		},
		//class
		{
			Type:                 TokenJSClass,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "class", IsIgnoreInResult: true},
				{IsAny: true, CanNested: true},
			},
		},
		//pattern function without keyword : using to detech class's function
		{
			Type:                 TokenJSClassFunction,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "constructor"},
				{Type: TokenJSBracket, CanNested: true},
				{Type: TokenJSBlock, CanNested: true},
			},
			LivingContext: []int{TokenJSClass},
		},
		//pattern function without keyword : using to detech class's function
		{
			Type:                 TokenJSClassFunction,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Type: TokenJSPhrase, CanNested: true, Nested: []gotokenize.PatternToken{
					{Type: TokenJSWord},
					{Type: TokenJSBracket, CanNested: true},
				}},
				{Type: TokenJSBlock, CanNested: true},
			},
			LivingContext: []int{TokenJSClass},
		},
		//pattern lambda single word
		{
			Type:                 TokenJSFunctionLambda,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{IsAny: true, CanNested: true},
				{Type: TokenJSPhrase, CanNested: true, Nested: []gotokenize.PatternToken{
					{Type: TokenJSRightArrow},
					{IsAny: true},
				}},
			},
		},
	},
}

var JSPatternLevel2 = gotokenize.PatternMeaningDefine{
	IgnoreTokens:   gotokenize.NoTokens,
	TokenCanNested: JSLevel2GlobalNested,
	Patterns: []gotokenize.Pattern{
		{
			Type:                 TokenJSAssignVariable,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{IsAny: true, CanNested: true},
				{Type: TokenJSAssign},
				{IsAny: true, CanNested: true},
			},
		},
	},
}

var JSPatterns = []gotokenize.PatternMeaningDefine{

	JSPatternLevel1,
	JSPatternLevel2,
	//JSPatternLevel3,
}

var JSStrongBreakEquivalent = []int{
	TokenJSBlock,
}
