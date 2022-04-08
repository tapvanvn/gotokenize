package js

import (
	"fmt"
	"strings"

	"github.com/tapvanvn/gotokenize/v2"
)

func NewStringifyContext() *StringifyContext {
	return &StringifyContext{}
}

type StringifyContext struct {
	LastTokenType int
}

func StringifyStream(iter *gotokenize.Iterator, ctx *StringifyContext) string {
	contents := []string{}
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		contents = append(contents, StringifyToken(token, ctx))
	}
	return strings.Join(contents, " ")
}

func StringifyToken(token *gotokenize.Token, ctx *StringifyContext) string {

	switch token.Type {

	case TokenJSWord:
		return token.Content
	case TokenJSString:
		return fmt.Sprintf("%s%s%s", token.Content, token.Children.ConcatStringContent(), token.Content)
	case TokenJSFor:
		return stringifyFor(token, ctx)
	case TokenJSBlockComment, TokenJSLineComment:
		return ""
	case TokenJSBlock:
		return stringifyBlock(token, ctx)
	case TokenJSBracket:
		return stringifyBracket(token, ctx)
	case TokenJSBracketSquare:
		return stringifyBracketSquare(token, ctx)
	case TokenJSPhraseBreak:
		return ""
	case TokenJSStrongBreak:
		return ";"
	case TokenJSIf:
		return fmt.Sprintf("if%s", StringifyStream(token.Children.Iterator(), ctx))
	case TokenJSElseIf:
		return fmt.Sprintf("else if%s", StringifyStream(token.Children.Iterator(), ctx))
	case TokenJSElse:
		return fmt.Sprintf("else%s", StringifyStream(token.Children.Iterator(), ctx))
	case TokenJSIfTrail:
		return stringifyIfTrail(token, ctx)
	case TokenJSPhrase:
		return stringifyPhrase(token, ctx)
	case TokenJSAssignVariable:
		return fmt.Sprintf("=%s", StringifyStream(token.Children.Iterator(), ctx))
	case TokenJSFunctionLambda:
		return stringifyFunctionLambda(token, ctx)
	case TokenJSClass:
		return stringifyClass(token, ctx)
	case TokenJSFunction:
		return stringifyFunction(token, ctx)
	case TokenJSClassFunction:
		return stringifyClassFunction(token, ctx)
	case TokenJSRegex:
		return token.Children.ConcatStringContent()
	case TokenJSSwitch:
		return fmt.Sprintf("switch%s", StringifyStream(token.Children.Iterator(), ctx))
	case TokenJSCase:
		return fmt.Sprintf("case %s:", StringifyStream(token.Children.Iterator(), ctx))
	case TokenJSDefault:
		return "default:"
	case TokenJSBreak:
		return "break"
	default:
		//fmt.Println("unknown ", token.Type)
		ctx.LastTokenType = token.Type
		return token.Content
	}
	return ""
}

func stringifyPhrase(token *gotokenize.Token, ctx *StringifyContext) string {

	return StringifyStream(token.Children.Iterator(), ctx)
}

func stringifyFor(token *gotokenize.Token, ctx *StringifyContext) string {

	if token.Children.Length() == 2 {

		return "for" + StringifyStream(token.Children.Iterator(), ctx)
	}
	return ""
}

func stringifyBracket(token *gotokenize.Token, ctx *StringifyContext) string {

	return fmt.Sprintf("%s%s%s", "(", StringifyStream(token.Children.Iterator(), ctx), ")")
}
func stringifyBracketSquare(token *gotokenize.Token, ctx *StringifyContext) string {

	return fmt.Sprintf("%s%s%s", "[", StringifyStream(token.Children.Iterator(), ctx), "]")
}

func stringifyBlock(token *gotokenize.Token, ctx *StringifyContext) string {

	return fmt.Sprintf("%s%s%s", "{", StringifyStream(token.Children.Iterator(), ctx), "}")
}

func stringifyIfTrail(token *gotokenize.Token, ctx *StringifyContext) string {
	contents := []string{
		"if",
	}
	conditionBracket := token.Children.GetTokenAt(0)
	bodyPhrase := token.Children.GetTokenAt(1)
	contents = append(contents, StringifyToken(conditionBracket, ctx), stringifyPhrase(bodyPhrase, ctx))
	nextIter := token.Children.Iterator()
	nextIter.Seek(2)
	for {
		token := nextIter.Get()
		if token == nil || (token.Type != TokenJSElseIf && token.Type != TokenJSElse) {
			break
		}
		nextIter.Read()
		if token.Type == TokenJSElseIf {
			contents = append(contents, "else if")
			contents = append(contents, StringifyToken(token.Children.GetTokenAt(0), ctx))
			bodyPhraseToken := nextIter.Read()
			if bodyPhraseToken == nil {
				break
			}
			contents = append(contents, stringifyPhrase(bodyPhraseToken, ctx))
		} else if token.Type == TokenJSElse {
			contents = append(contents, "else")
			bodyPhraseToken := nextIter.Read()
			if bodyPhraseToken == nil {
				break
			}
			contents = append(contents, stringifyPhrase(bodyPhraseToken, ctx))
		}
	}
	return strings.Join(contents, "")
}

func stringifyFunctionLambda(token *gotokenize.Token, ctx *StringifyContext) string {

	if token.Children.Length() >= 2 {

		return fmt.Sprintf("%s=>%s",
			StringifyToken(token.Children.GetTokenAt(0), ctx),
			StringifyToken(token.Children.GetTokenAt(1), ctx))
	} else {
		fmt.Println("wrong lambda", token.Children.Length())
	}
	return ""
}
func stringifyClassFunction(token *gotokenize.Token, ctx *StringifyContext) string {

	return StringifyStream(token.Children.Iterator(), ctx)
}
func stringifyFunction(token *gotokenize.Token, ctx *StringifyContext) string {

	return "function " + StringifyStream(token.Children.Iterator(), ctx)
}
func stringifyClass(token *gotokenize.Token, ctx *StringifyContext) string {

	return fmt.Sprintf("class%s", StringifyStream(token.Children.Iterator(), ctx))
}
