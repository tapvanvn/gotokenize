package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

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
		stringifier.Content += " "
	}
	stringifier.Content += content
	stringifier.HasSpace = stroke.IsSpaceAfter
	stringifier.HasStrongBreak = stroke.IsStrongBreakEquivalent
}

func (stringifier *Stringifier) PutToken(token *gotokenize.Token) {

	if token.Type == TokenJSIfTrail {

		stringifier.PutIfTrail(token)

	} else if token.Type == TokenJSBlock {

		stringifier.PutBlock(token)

	} else if token.Type == TokenJSBracket {

		stringifier.PutBracket(token)

	} else if token.Type == TokenJSBracketSquare {

		stringifier.PutBracketSquare(token)

	} else if token.Type == TokenJSWord {

		stringifier.PutWord(token)

	} else if token.Type == TokenJSString {

		stringifier.PutString(token)

	} else if token.Type == TokenJSPhrase {

		stringifier.PutPhrase(token)

	} else if token.Type == TokenJSClass {

		stringifier.PutClass(token)

	} else if token.Type == TokenJSClassFunction {

		stringifier.PutClassFunction(token)

	} else if token.Type == TokenJSAssignVariable {

		stringifier.PutAssignVariable(token)

	} else if token.Type == TokenJSFunctionLambda {

		stringifier.PutLambda(token)

	} else if token.Type == TokenJSRegex {

		stringifier.PutRegex(token)

	} else {

		stringifier.put(token.Content, &DefaultStroke)
	}
}
func (stringifier *Stringifier) PutStream(iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.PutToken(token)
	}
}
func (stringifier *Stringifier) PutLambda(token *gotokenize.Token) {
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("=>", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &DefaultStroke)
}
func (stringifier *Stringifier) PutAssignVariable(token *gotokenize.Token) {

	stringifier.put("=", &BreakAfterStroke)
	stringifier.PutStream(token.Children.Iterator())
}

func (stringifier *Stringifier) PutPhrase(token *gotokenize.Token) {

	if token.Children.Length() == 0 {
		return
	}

	iter := token.Children.Iterator()
	stringifier.put("", &NeedBreakStroke)
	stringifier.PutStream(iter)
	stringifier.put("", &DefaultStroke)
}

func (stringifier *Stringifier) PutClass(token *gotokenize.Token) {

	stringifier.put("class", &DefaultStroke)
	stringifier.PutStream(token.Children.Iterator())
}
func (stringifier *Stringifier) PutClassFunction(token *gotokenize.Token) {

	stringifier.put("", &NeedBreakStroke)
	stringifier.PutStream(token.Children.Iterator())
	stringifier.put("", &DefaultStroke)
}

func (stringifier *Stringifier) PutBracket(token *gotokenize.Token) {

	stringifier.put("(", &BreakAfterStroke)

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		stringifier.PutToken(childToken)
	}
	stringifier.put(")", &BreakAfterStroke)
}
func (stringifier *Stringifier) PutBracketSquare(token *gotokenize.Token) {

	stringifier.put("[", &BreakAfterStroke)

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		stringifier.PutToken(childToken)
	}
	stringifier.put("]", &BreakAfterStroke)
}

func (stringifier *Stringifier) PutBlock(token *gotokenize.Token) {

	stringifier.put("{", &BreakAfterStroke)
	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		stringifier.PutToken(childToken)
	}
	stringifier.put("}", &BreakAfterStroke)
}

func (stringifier *Stringifier) PutIfTrail(token *gotokenize.Token) {

	stringifier.put("if", &NeedBreakStroke)

	conditionBracket := token.Children.GetTokenAt(0)
	bodyPhrase := token.Children.GetTokenAt(1)
	stringifier.PutToken(conditionBracket)
	stringifier.PutToken(bodyPhrase)

	nextIter := token.Children.Iterator()
	nextIter.Seek(2)
	for {
		token := nextIter.Get()
		if token == nil || (token.Type != TokenJSElseIf && token.Type != TokenJSElse) {
			break
		}
		nextIter.Read()
		if token.Type == TokenJSElseIf {
			stringifier.put("else if", &NeedBreakStroke)
			stringifier.PutToken(token.Children.GetTokenAt(0))
			bodyPhraseToken := nextIter.Read()
			if bodyPhraseToken == nil {
				break
			}
			stringifier.PutToken(bodyPhraseToken)

		} else if token.Type == TokenJSElse {
			stringifier.put("else", &NeedBreakStroke)
			bodyPhraseToken := nextIter.Read()
			if bodyPhraseToken == nil {
				break
			}
			stringifier.PutToken(bodyPhraseToken)
		}
	}
}
func (stringifier *Stringifier) PutRegex(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.put(token.Children.ConcatStringContent(), &DefaultStroke)
	stringifier.put(token.Content, &DefaultStroke)
}
func (stringifier *Stringifier) PutString(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.put(token.Children.ConcatStringContent(), &DefaultStroke)
	stringifier.put(token.Content, &DefaultStroke)
}
func (stringifier *Stringifier) PutWord(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultSpaceStroke)
}
