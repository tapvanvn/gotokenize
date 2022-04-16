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

func (stroke *TokenStroke) ShouldSpace(current bool) bool {
	return stroke.NeedSpace && !current
}
func (stroke *TokenStroke) ShouldStrongBreak(current bool) bool {
	//fmt.Println("shoud:", stroke.NeedStrongBreak, current)
	return stroke.NeedStrongBreak && !current
}

var DefaultStroke TokenStroke = TokenStroke{}
var DefaultSpaceStroke TokenStroke = TokenStroke{
	NeedSpace: true,
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
		HasSpace:       true,
		HasStrongBreak: true,
	}
}

type Stringifier struct {
	Content        string
	HasSpace       bool
	HasStrongBreak bool
}

func (stringifier *Stringifier) put(content string, stroke *TokenStroke) {
	if stroke.ShouldStrongBreak(stringifier.HasStrongBreak) {

		stringifier.Content += ";"
	} else if stroke.ShouldSpace(stringifier.HasSpace) {
		stringifier.Content += ""
	}
	stringifier.Content += content
	stringifier.HasSpace = stroke.IsSpaceAfter
	stringifier.HasStrongBreak = stroke.IsStrongBreakEquivalent
}

func (stringifier *Stringifier) PutToken(token *gotokenize.Token) {
	stringifier.Content += token.Content
}

func (stringifier *Stringifier) GetContent() string {
	return stringifier.Content
}
