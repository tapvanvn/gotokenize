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

	TokenJSBracket         = 100 //()
	TokenJSBlock           = 101 //{}
	TokenJSBracketSquare   = 102 //[]
	TokenJSUnaryOperator   = 103 // !, ~, ++, --
	TokenJSBinaryOperator  = 104 // <>+-*/%, <=, >=, <<, >>, >>>, ||, |, &&, &, ^, **, ==, ===, !=, !==
	TokenJSTreeDotOperator = 105
	TokenJSAssign          = 106 // =
	TokenJSRightArrow      = 107 // =>
	TokenJSLineComment     = 108 // //
	TokenJSBlockComment    = 109 // /**/
	TokenJSPhrase          = 110
	TokenJSRegex           = 111 // //img

	TokenJSFunction       = 200
	TokenJSFunctionLambda = 201
	TokenJSVariable       = 202
	TokenJSString         = 203
	TokenJSFor            = 204
	TokenJSIf             = 205
	TokenJSElseIf         = 206
	TokenJSElse           = 207
	TokenJSSwitch         = 208
	TokenJSCase           = 209
	TokenJSDefault        = 210
	TokenJSBreak          = 211
	TokenJSWhile          = 212
	TokenJSDo             = 213
	TokenJSClass          = 214
	TokenJSClassFunction  = 215
	TokenJSIfTrail        = 216

	TokenJSStrongBreak           = 300 //sure `;`
	TokenJSArgumentList          = 301 //use in function declaration and function call
	TokenJSStrongDeclareVariable = 302 //mark that the follow variable had been strong declared with `var` keyword
	TokenJSStrongDeclareConstant = 303 //mark that the follow variable had been strong declared with `const` keyword
	TokenJSStrongDeclareLet      = 304 //mark that the follow variable had been strong declared with `let` keyword
	TokenJSAssignVariable        = 350 //declare a new or assign value to variable
	TokenJSDeclareFunction       = 351 //declare a new function (traditional or lambda)
	TokenJSDeclareClass          = 352 //declare a class using class keyword
	TokenJSDeclareObject         = 353 //declare an object
	TokenJSDeclareArray          = 354 //declare an array using [] or new Array
	TokenJSVoidStatement         = 355 //any void(...) void ...
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

	case TokenJSTreeDotOperator:
		return "treedot operator"

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

	case TokenJSCase:
		return "case"

	case TokenJSDefault:
		return "default"

	case TokenJSBreak:
		return "break"

	case TokenJSWhile:
		return "while"

	case TokenJSDo:
		return "do"

	case TokenJSClass:
		return "class"

	case TokenJSClassFunction:
		return "class function"

	case TokenJSRegex:
		return "regex"

	case TokenJSIfTrail:
		return "iftrail"
		//Instruction
	case TokenJSArgumentList:
		return "argument list"
	case TokenJSStrongDeclareVariable:
		return "var"
	case TokenJSStrongDeclareConstant:
		return "const"
	case TokenJSStrongDeclareLet:
		return "let"
	case TokenJSAssignVariable:
		return "assign variable"
	case TokenJSDeclareFunction:
		return "declare function"
	case TokenJSDeclareClass:
		return "declare class"
	case TokenJSDeclareObject:
		return "declare object"
	case TokenJSDeclareArray:
		return "declare array"
	case TokenJSVoidStatement:
		return "void statement"
	default:
		return fmt.Sprintf("unknown-%d", Type)
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

//JSPatterns Patterns to detech structure of js
var JSPatterns = []gotokenize.Pattern{

	//pattern if block
	{
		Type:                 TokenJSIf,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "if", IsIgnoreInResult: true},
			{Type: TokenJSBracket, CanNested: true},
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
		},
	},
	//pattern else block
	{
		Type:                 TokenJSElse,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "else", IsIgnoreInResult: true},
		},
	},
	//pattern for
	{
		Type:                 TokenJSFor,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "for", IsIgnoreInResult: true},
			{Type: TokenJSBracket, CanNested: true},
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
		},
	},
	//pattern while block
	{
		Type:                 TokenJSWhile,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "while", IsIgnoreInResult: true},
		},
	},
	//pattern function with keyword and name
	{
		Type:                 TokenJSFunction,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "function", IsIgnoreInResult: true},
			{Type: TokenJSWord},
			{Type: TokenJSBracket, CanNested: true},
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

	//pattern function without keyword : using to detech class's function
	{
		Type:                 TokenJSClassFunction,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Type: TokenJSWord},
			{Type: TokenJSBracket, CanNested: true},
			{Type: TokenJSBlock, CanNested: true},
		},
	},

	//pattern lambda
	{
		Type:                 TokenJSFunctionLambda,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Type: TokenJSBracket, CanNested: true},
			{Type: TokenJSRightArrow, IsIgnoreInResult: true},
		},
	},

	//pattern lambda single word
	{
		Type:                 TokenJSFunctionLambda,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Type: TokenJSBracket, CanNested: true},
			{Type: TokenJSRightArrow, IsIgnoreInResult: true},
		},
	},
	//class without name
	{
		Type:                 TokenJSClass,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "class", IsIgnoreInResult: true},
			{Type: TokenJSBlock, CanNested: true},
		},
	},
	//class with name
	{
		Type:                 TokenJSClass,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "class", IsIgnoreInResult: true},
			{Type: TokenJSWord},
			{Type: TokenJSBlock, CanNested: true},
		},
	},
}

var JSGlobalNested = []int{
	TokenJSBlock,
	TokenJSBracket,
	TokenJSBracketSquare,
	TokenJSClass,
	TokenJSPhrase,
}

//JSInstructionPatterns the patterns to detect instruction structure
var JSInstructionPatterns = []gotokenize.Pattern{
	{
		Type:                 TokenJSStrongDeclareVariable,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "var", IsIgnoreInResult: true},
		},
	},
	//Assign variable
	{
		Type:                 TokenJSAssignVariable,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{IsAny: true, CanNested: true},
			{Type: TokenJSAssign, IsIgnoreInResult: true},
			{IsAny: true, CanNested: true},
		},
	},
	//Void
	{
		Type:                 TokenJSVoidStatement,
		IsRemoveGlobalIgnore: true,
		Struct: []gotokenize.PatternToken{
			{Content: "void", IsIgnoreInResult: true},
			{IsAny: true, CanNested: true},
		},
	},
}

var JSInstructionGlobalNested = append(
	JSGlobalNested,
	TokenJSDeclareClass,
	TokenJSAssignVariable,
	TokenJSFunction,
	TokenJSFunctionLambda,
	TokenJSClassFunction,
	TokenJSSwitch,
)

var JSStrongBreakEquivalent = []int{
	TokenJSBlock,
}
