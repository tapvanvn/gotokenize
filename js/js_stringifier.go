package js

import "github.com/tapvanvn/gotokenize/v2"

var requireBreakKeyWords = ",var,const,let,delete,throw,continue,"
var requireSpaceKeyWords = ",typeof,new,catch,extends,instanceof,"
var requireBreakNormalAfter = ",return,"

type TokenStroke struct {
	NeedSpace               bool //need space before
	NeedStrongBreak         bool //need break with before context
	IsSpaceAfter            bool //is token contain a space (or equivalent) after
	IsStrongBreakEquivalent bool //is token also mean strong break
}

func (stroke *TokenStroke) ShouldSpace(current bool, currentBreak bool) bool {
	return !currentBreak && stroke.NeedSpace && !current
}
func (stroke *TokenStroke) ShouldStrongBreak(current bool) bool {
	//fmt.Println("shoud:", stroke.NeedStrongBreak, current)
	return stroke.NeedStrongBreak && !current
}

var DefaultStroke TokenStroke = TokenStroke{}
var NeedSpaceStroke TokenStroke = TokenStroke{
	NeedSpace: true,
}
var SpaceAfterStroke TokenStroke = TokenStroke{

	IsSpaceAfter: true,
}
var NeedAndHasSpaceStroke TokenStroke = TokenStroke{
	NeedSpace:    true,
	IsSpaceAfter: true,
}
var NeedBreakStroke TokenStroke = TokenStroke{
	NeedStrongBreak: true,
}
var BreakAfterStroke TokenStroke = TokenStroke{
	IsStrongBreakEquivalent: true,
}
var NeedAndHasBreakStroke TokenStroke = TokenStroke{
	IsStrongBreakEquivalent: true,
	NeedStrongBreak:         true,
}

func NewStringfier() *Stringifier {
	return &Stringifier{
		HasSpace:        true,
		HasStrongBreak:  true,
		TokenProcessors: map[int]*TokenStringifier{},
	}
}

type TokenStringifier struct {
	ProcessFunction func(*Stringifier, *gotokenize.Token)
}

func (tokenStringifier *TokenStringifier) ProcessToken(stringifier *Stringifier, token *gotokenize.Token) {

	if tokenStringifier.ProcessFunction != nil {

		tokenStringifier.ProcessFunction(stringifier, token)
	}
}

type Stringifier struct {
	Content           string
	HasSpace          bool
	HasStrongBreak    bool
	TokenProcessors   map[int]*TokenStringifier
	NonTokenProcessor *TokenStringifier
}

func (stringifier *Stringifier) SetProcessor(tokenType int, tokenStringifier func(*Stringifier, *gotokenize.Token)) {

	stringifier.TokenProcessors[tokenType] = &TokenStringifier{
		ProcessFunction: tokenStringifier,
	}
}
func (stringifier *Stringifier) SetNonTokenProcessor(tokenStringifier func(*Stringifier, *gotokenize.Token)) {

	stringifier.NonTokenProcessor = &TokenStringifier{
		ProcessFunction: tokenStringifier,
	}
}

func (stringifier *Stringifier) put(content string, stroke *TokenStroke) {

	if stroke.ShouldStrongBreak(stringifier.HasStrongBreak) {

		stringifier.Content += ";"

	} else if stroke.ShouldSpace(stringifier.HasSpace, stringifier.HasStrongBreak) {

		stringifier.Content += " "
	}
	stringifier.Content += content
	stringifier.HasSpace = stroke.IsSpaceAfter
	stringifier.HasStrongBreak = stroke.IsStrongBreakEquivalent
}

func (stringifier *Stringifier) PutToken(token *gotokenize.Token) {
	if token == nil {
		return
	}
	if processor, ok := stringifier.TokenProcessors[token.Type]; ok {

		processor.ProcessToken(stringifier, token)

	} else if stringifier.NonTokenProcessor != nil {

		stringifier.NonTokenProcessor.ProcessToken(stringifier, token)

	} else {

		stringifier.Content += token.Content
	}
}

func (stringifier *Stringifier) GetContent() string {

	return stringifier.Content
}
