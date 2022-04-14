package js

import (
	"strings"

	"github.com/tapvanvn/gotokenize/v2"
)

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
	//fmt.Println("+ ", content, stroke.IsStrongBreakEquivalent)
}

func (stringifier *Stringifier) PutToken(token *gotokenize.Token) {

	if token == nil {

		return
	}
	switch token.Type {
	case TokenJSBlock:
		stringifier.PutBlock(token)

	case TokenJSBracket:
		stringifier.PutBracket(token)

	case TokenJSBracketSquare:
		stringifier.PutBracketSquare(token)

	case TokenJSInlineIf:
		stringifier.PutInlineIf(token)

	case TokenJSIf:
		stringifier.PutIf(token)

	case TokenJSElseIf:
		stringifier.PutElseIf(token)

	case TokenJSElse:
		stringifier.PutElse(token)

	case TokenJSDo:
		stringifier.PutDo(token)
	case TokenJSWhile:
		stringifier.PutWhile(token)

	case TokenJSFor:
		stringifier.PutFor(token)

	case TokenJSSwitch:
		stringifier.PutSwitch(token)

	case TokenJSFunction:
		stringifier.PutFunction(token)

	case TokenJSWord:
		stringifier.PutWord(token)

	case TokenJSString:
		stringifier.PutString(token)

	case TokenJSPhrase, TokenJSMinusPhrase:
		stringifier.PutPhrase(token)

	case TokenJSClass:
		stringifier.PutClass(token)
	case TokenJSClassFunction:
		stringifier.PutClassFunction(token)
	case TokenJSAssignVariable:
		stringifier.PutAssignVariable(token)
	case TokenJSOperatorTrail:
		stringifier.PutOperatorTrail(token)
	case TokenJSFunctionLambda:
		stringifier.PutLambda(token)
	case TokenJSRegex:
		stringifier.PutRegexToken(token)
	case TokenJSObjectProperty, TokenJSObjectLastProperty:

		stringifier.PutObjectProperty(token)

	case TokenJSAssign:
		stringifier.put(token.Content, &BreakAfterStroke)

	case TokenJSKeyWord:
		if strings.Contains(requireBreakKeyWords, ","+token.Content+",") {

			stringifier.put(token.Content+" ", &NeedAndHasBreakStroke)

		} else if strings.Contains(requireSpaceKeyWords, ","+token.Content+",") {

			stringifier.put(" "+token.Content+" ", &BreakAfterStroke)
		} else if strings.Contains(requireBreakNormalAfter, ","+token.Content+",") {
			stringifier.put(" "+token.Content+" ", &NeedBreakStroke)
		} else {
			stringifier.put(" "+token.Content+" ", &BreakAfterStroke)
		}
	case TokenJSBreak:
		stringifier.put("break ", &NeedBreakStroke)
	case TokenJSColonOperator:
		stringifier.put(":", &BreakAfterStroke)
	case TokenJSOperator:
		stringifier.put(token.Content, &BreakAfterStroke)
	case TokenJSUnaryOperator:
		stringifier.put(token.Content, &BreakAfterStroke)
	case TokenJSPhraseBreak:
		break
	case TokenJSStrongBreak:

		stringifier.put(";", &BreakAfterStroke)
	case TokenJSReturnStatement:
		stringifier.PutReturnStatement(token)
	case TokenJSLabel:
		stringifier.PutLabel(token)
	default:
		if token.Content == "," {
			stringifier.put(",", &BreakAfterStroke)
		} else if token.Type == TokenJSLineComment || token.Type == TokenJSBlockComment {
		} else {
			stringifier.put(token.Content, &DefaultStroke)
		}
	}
}

func (stringifier *Stringifier) PutLabel(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))

	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))

	stringifier.put("", &BreakAfterStroke)
}
func (stringifier *Stringifier) PutPhrase(token *gotokenize.Token) {

	if token.Children.Length() == 0 {
		return
	}
	stringifier.put("", &NeedAndHasBreakStroke)
	if token.Type == TokenJSMinusPhrase {
		stringifier.put("-", &BreakAfterStroke)
		stringifier.PutStream(token.Children.GetTokenAt(0).Children.Iterator())
	} else {
		iter := token.Children.Iterator()
		stringifier.put("", &BreakAfterStroke)
		stringifier.PutStream(iter)
	}

	stringifier.put("", &DefaultStroke)
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
func (stringifier *Stringifier) PutStreamNoBreak(iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put("", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
}
func (stringifier *Stringifier) PutStreamSpace(iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put(" ", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
}

func (stringifier *Stringifier) PutAssignVariable(token *gotokenize.Token) {

	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))

	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))

	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))

	stringifier.put("", &DefaultStroke)
}
func (stringifier *Stringifier) PutOperatorTrail(token *gotokenize.Token) {

	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutStream(token.Children.Iterator())
	stringifier.put("", &DefaultStroke)
}
func (stringifier *Stringifier) PutReturnStatement(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutStreamSpace(token.Children.Iterator())
	stringifier.put("", &DefaultStroke)
}
func (stringifier *Stringifier) PutFor(token *gotokenize.Token) {
	stringifier.put("for", &NeedAndHasBreakStroke)
	stringifier.PutForBracket(token.Children.GetTokenAt(0))
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &DefaultStroke)
}

//TODO: improve this
func (stringifier *Stringifier) PutForBracket(token *gotokenize.Token) {

	stringifier.put("(", &BreakAfterStroke)

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		stringifier.put("", &BreakAfterStroke)
		stringifier.PutToken(childToken)
		stringifier.put("", &BreakAfterStroke)
	}
	stringifier.put(")", &DefaultStroke)
}

func (stringifier *Stringifier) PutDo(token *gotokenize.Token) {
	stringifier.put("do ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("while ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &DefaultStroke)
}

func (stringifier *Stringifier) PutWhile(token *gotokenize.Token) {
	stringifier.put("while ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &DefaultStroke)
}
func (stringifier *Stringifier) PutSwitch(token *gotokenize.Token) {
	stringifier.put("switch ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("{", &BreakAfterStroke)
	body := token.Children.GetTokenAt(1)
	if body != nil {
		iter := body.Children.Iterator()
		for {
			childToken := iter.Read()
			if childToken == nil {
				break
			}
			if childToken.Type == TokenJSPhrase {
				iter := childToken.Children.Iterator()
				identity := iter.Read()
				if identity.Type == TokenJSCase {
					stringifier.put("case ", &NeedAndHasBreakStroke)
					stringifier.PutStream(identity.Children.Iterator())
					stringifier.put(":", &BreakAfterStroke)
				} else if identity.Type == TokenJSDefault {
					stringifier.put("default:", &NeedAndHasBreakStroke)
				}
				stringifier.PutStream(iter)
			}
		}
	}
	stringifier.put("}", &BreakAfterStroke)
}
func (stringifier *Stringifier) PutObjectProperty(token *gotokenize.Token) {
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put(":", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	if token.Type == TokenJSObjectProperty {
		stringifier.put(",", &DefaultStroke)
	}
}
func (stringifier *Stringifier) PutLambda(token *gotokenize.Token) {
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("=>", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &DefaultStroke)
}

func (stringifier *Stringifier) PutFunction(token *gotokenize.Token) {

	stringifier.put("function ", &NeedAndHasBreakStroke)
	stringifier.PutStream(token.Children.Iterator())
	stringifier.put("", &BreakAfterStroke)
}
func (stringifier *Stringifier) PutClass(token *gotokenize.Token) {

	stringifier.put("class ", &NeedAndHasBreakStroke)
	stringifier.PutStreamSpace(token.Children.Iterator())
}
func (stringifier *Stringifier) PutClassFunction(token *gotokenize.Token) {

	stringifier.put("", &NeedBreakStroke)
	stringifier.PutStream(token.Children.Iterator())
	stringifier.put("", &BreakAfterStroke)
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
	stringifier.put(")", &DefaultStroke)
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
	stringifier.put("]", &DefaultStroke)
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

func (stringifier *Stringifier) PutInlineIf(token *gotokenize.Token) {

	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("?", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
	stringifier.put(":", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(4))
	stringifier.put("", &DefaultStroke)
	/*body := token.Children.GetTokenAt(2)
	stringifier.PutToken(body.Children.GetTokenAt(0))
	stringifier.put(":", &BreakAfterStroke)
	stringifier.PutToken(body.Children.GetTokenAt(2))
	stringifier.put("", &DefaultStroke)*/
}

func (stringifier *Stringifier) PutIf(token *gotokenize.Token) {

	stringifier.put("if ", &NeedAndHasBreakStroke)
	conditionBracket := token.Children.GetTokenAt(0)
	bodyPhrase := token.Children.GetTokenAt(1)

	stringifier.PutToken(conditionBracket)

	stringifier.put("", &BreakAfterStroke)

	if bodyPhrase != nil {

		stringifier.PutToken(bodyPhrase)
		if bodyPhrase.Type != TokenJSBlock {

			stringifier.put(";", &BreakAfterStroke)
		} else {
			stringifier.put("", &BreakAfterStroke)
		}
	}
}

func (stringifier *Stringifier) PutElseIf(token *gotokenize.Token) {
	stringifier.put("else if ", &NeedBreakStroke)

	conditionBracket := token.Children.GetTokenAt(0)
	bodyPhrase := token.Children.GetTokenAt(1)
	stringifier.PutToken(conditionBracket)
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(bodyPhrase)
}

func (stringifier *Stringifier) PutElse(token *gotokenize.Token) {
	stringifier.put("else ", &NeedAndHasBreakStroke)
	bodyPhrase := token.Children.GetTokenAt(0)
	stringifier.PutToken(bodyPhrase)
}

func (stringifier *Stringifier) PutRegexToken(token *gotokenize.Token) {

	if token.Type == TokenJSBracket ||
		token.Type == TokenJSBlock ||
		token.Type == TokenJSBracketSquare {
		stringifier.PutRegexBlockStyle(token)
	} else if token.Type == TokenJSRegex {
		stringifier.PutRegexStream(token)
	} else {
		stringifier.put(token.Content, &DefaultStroke)
	}
}
func (stringifier *Stringifier) PutRegexStream(parentToken *gotokenize.Token) {

	iter := parentToken.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenJSBracket ||
			childToken.Type == TokenJSBlock ||
			childToken.Type == TokenJSBracketSquare {
			stringifier.PutRegexBlockStyle(childToken)
		} else {
			stringifier.PutRegexToken(childToken)
		}
	}
}
func (stringifier *Stringifier) PutRegexBlockStyle(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.PutRegexStream(token)
	if token.Type == TokenJSBracket {
		stringifier.put(")", &DefaultStroke)
	} else if token.Type == TokenJSBlock {
		stringifier.put("}", &DefaultStroke)
	} else {
		stringifier.put("]", &DefaultStroke)
	}

}

func (stringifier *Stringifier) PutString(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.put(token.Children.ConcatStringContent(), &DefaultStroke)
	stringifier.put(token.Content, &DefaultStroke)
}
func (stringifier *Stringifier) PutWord(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultSpaceStroke)
}
