package js

import (
	"github.com/tapvanvn/gotokenize"
)

const (
	TokenJSUnknown           = 0
	TokenJSWord              = 1
	TokenJSOperator          = 2
	TokenJSPhraseBreak       = 3
	TokenJSPhraseStrongBreak = 4 //need ; after
	TokenJSScopeBegin        = 5
	TokenJSScopeEnd          = 6
	TokenJSWordBreak         = 7
	TokenJSGlueBegin         = 8
	TokenJSGlueEnd           = 9

	TokenJSBracket        = 100
	TokenJSBlock          = 101
	TokenJSBracketSquare  = 102
	TokenJSUnaryOperator  = 103 // !, ~, ++, --
	TokenJSBinaryOperator = 104 // <>+-*/%, <=, >=, <<, >>, >>>, ||, |, &&, &, ^, **, ==, ===, !=, !==
	TokenJSAssign         = 105 // =
	TokenJSRightArrow     = 106 // =>
	TokenJSLineComment    = 107
	TokenJSBlockComment   = 108
	TokenJSPhrase         = 109
	TokenJSRegex          = 110

	TokenJSFunction       = 200
	TokenJSFunctionLambda = 201
	TokenJSVariable       = 202
	TokenJSString         = 203
	TokenJSFor            = 204
	TokenJSIf             = 205
	TokenJSElseIf         = 206
	TokenJSElse           = 207
	TokenJSSwitch         = 208
	TokenJSWhile          = 209
	TokenJSDo             = 210

	//TokenJSCraft      = 300
	//TokenJSCraftDebug = 301

	//TokenJSPatchStream = 400
)

//JSTokenName return name from type of token
func JSTokenName(Type int) string {

	switch Type {

	case TokenJSUnknown:
		return "unknown"

	case TokenJSWord:
		return "word"

	case TokenJSOperator:
		return "operator"

	case TokenJSPhraseBreak:
		return "phrase break"

	case TokenJSBracket:
		return "bracket"

	case TokenJSBlock:
		return "block"

	case TokenJSBracketSquare:
		return "bracket square"

	case TokenJSUnaryOperator:
		return "unary operator"

	case TokenJSBinaryOperator:
		return "binary operator"

	case TokenJSAssign:
		return "assign"

	case TokenJSRightArrow:
		return "right arrow"

	case TokenJSLineComment:
		return "line comment"

	case TokenJSBlockComment:
		return "block comment"

	case TokenJSPhrase:
		return "phrase"

	case TokenJSFunction:
		return "function"

	case TokenJSFunctionLambda:
		return "lambda"

	case TokenJSVariable:
		return "variable"

	case TokenJSString:
		return "string"

	case TokenJSFor:
		return "for"

	case TokenJSIf:
		return "if"

	case TokenJSElseIf:
		return "else if"

	case TokenJSElse:
		return "else"

	case TokenJSSwitch:
		return "switch"

	case TokenJSWhile:
		return "while"

	case TokenJSDo:
		return "do"

	/*
		case TokenJSCraft:
				return "craft"

			case TokenJSCraftDebug:
				return "craft debug"
	*/
	default:
		return "unknown"
	}
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
,volatile,while,with,yield,`

//JSIgnores tokens that will be ignore
var JSIgnores = []int{

	TokenJSLineComment,

	TokenJSBlockComment,
}

//JSPatterns Patterns to detech structure of js
var JSPatterns = []gotokenize.Pattern{

	//pattern if block
	{
		Type:                 TokenJSIf,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "if", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern if phrase
	{
		Type:                 TokenJSIf,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "if", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{ExportType: TokenJSPhrase, IsPhraseUntil: true},
		},
	},

	//pattern else if block
	{
		Type:                 TokenJSElseIf,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "else", IsIgnoreInResult: true},
			{Content: "if", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern else if pharse
	{
		Type:                 TokenJSElseIf,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "else", IsIgnoreInResult: true},
			{Content: "if", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{ExportType: TokenJSPhrase, IsPhraseUntil: true},
		},
	},
	//pattern else block
	{
		Type:                 TokenJSElse,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "else", IsIgnoreInResult: true},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern else phrase
	{
		Type:                 TokenJSElse,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "else", IsIgnoreInResult: true},
			{ExportType: TokenJSPhrase, IsPhraseUntil: true},
		},
	},

	//pattern for
	{
		Type:                 TokenJSFor,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "for", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern function with keyword
	{
		Type:                 TokenJSFunction,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "function", IsIgnoreInResult: true},
			{Type: TokenJSWord},
			{Type: TokenJSBracket},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern lambda
	{
		Type:                 TokenJSFunctionLambda,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Type: TokenJSBracket},
			{Type: TokenJSRightArrow, IsIgnoreInResult: true},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern switch
	{
		Type:                 TokenJSSwitch,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "switch", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern while block
	{
		Type:                 TokenJSWhile,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "while", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{Type: TokenJSBlock, CanNested: true},
		},
	},
	//pattern while phrase
	{
		Type:                 TokenJSWhile,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "while", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
			{ExportType: TokenJSPhrase, IsPhraseUntil: true},
		},
	},

	//pattern do block
	{
		Type:                 TokenJSDo,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "do", IsIgnoreInResult: true},
			{Type: TokenJSBlock, CanNested: true},
			{Content: "while", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
		},
	},

	//pattern do phrase
	{
		Type:                 TokenJSDo,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "do", IsIgnoreInResult: true},
			{ExportType: TokenJSPhrase, IsPhraseUntil: true},
			{Content: "while", IsIgnoreInResult: true},
			{Type: TokenJSBracket},
		},
	},
	/*
		//pattern jscraft.require
		{
			Type:                 TokenJSCraft,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "jscraft", IsIgnoreInResult: true},
				{Content: ".", IsIgnoreInResult: true},
				{Content: "require"},
				{Type: TokenJSBracket},
			},
		},

		//pattern jscraft.template
		{
			Type:                 TokenJSCraft,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "jscraft", IsIgnoreInResult: true},
				{Content: ".", IsIgnoreInResult: true},
				{Content: "template"},
				{Type: TokenJSBracket, CanNested: true},
			},
		},

		//pattern jscraft.build
		{
			Type:                 TokenJSCraft,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "jscraft", IsIgnoreInResult: true},
				{Content: ".", IsIgnoreInResult: true},
				{Content: "build"},
				{Type: TokenJSBracket, CanNested: true},
			},
		},

		//pattern jscraft.conflict
		{
			Type:                 TokenJSCraft,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "jscraft", IsIgnoreInResult: true},
				{Content: ".", IsIgnoreInResult: true},
				{Content: "conflict"},
				{Type: TokenJSBracket},
			},
		},

		//pattern jscraft.conflict
		{
			Type:                 TokenJSCraft,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "jscraft", IsIgnoreInResult: true},
				{Content: ".", IsIgnoreInResult: true},
				{Content: "fetch"},
				{Type: TokenJSBracket},
			},
		},

		//pattern jscraft.conflict
		{
			Type:                 TokenJSCraftDebug,
			IsRemoveGlobalIgnore: true,
			Struct: []gotokenize.PatternToken{
				{Content: "jscraft_debug", IsIgnoreInResult: true},
				{Type: TokenJSBlock, CanNested: true},
			},
		},*/
}
