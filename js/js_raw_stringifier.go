package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewDefaultRawStringifier() *Stringifier {

	stringifier := NewStringfier()
	stringifier.SetProcessor(TokenJSBracket, ProcessRawBracket)
	stringifier.SetProcessor(TokenJSBlock, ProcessRawBlock)
	stringifier.SetProcessor(TokenJSBracketSquare, ProcessRawBracketSquare)
	stringifier.SetProcessor(TokenJSRegex, ProcessRawRegexToken)
	stringifier.SetProcessor(TokenJSString, ProcessRawString)
	stringifier.SetProcessor(TokenJSWord, ProcessRawWord)
	stringifier.SetProcessor(TokenJSKeyWord, ProcessRawKeyWord)
	stringifier.SetProcessor(TokenJSOperator, ProcessRawStrongOperator)
	stringifier.SetProcessor(TokenJSUnaryOperator, ProcessRawStrongOperator)
	stringifier.SetProcessor(TokenJSBinaryOperator, ProcessRawStrongOperator)
	stringifier.SetProcessor(TokenJSQuestionOperator, ProcessRawStrongOperator)
	stringifier.SetProcessor(TokenJSColonOperator, ProcessRawStrongOperator)
	stringifier.SetProcessor(TokenJSSoftBreak, ProcessRawSoftOperator)
	stringifier.SetProcessor(TokenJSPhraseBreak, ProcessRawPhraseBreak)
	stringifier.SetProcessor(TokenJSStrongBreak, ProcessRawStrongBreak)
	stringifier.SetNonTokenProcessor(ProcessNonToken)
	return stringifier
}

func ProcessNonToken(stringifier *Stringifier, token *gotokenize.Token) {
	if token.Content == "," {
		stringifier.put(",", &BreakAfterStroke)
	} else if token.Type == TokenJSLineComment || token.Type == TokenJSBlockComment {
	} else {
		stringifier.put(token.Content, &DefaultStroke)
	}
}

func ProcessRawBlock(stringifier *Stringifier, token *gotokenize.Token) {

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
func ProcessRawBracket(stringifier *Stringifier, token *gotokenize.Token) {

	stringifier.put("(", &BreakAfterStroke)

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		//fmt.Println("-----", childToken.Content)
		stringifier.PutToken(childToken)
	}
	stringifier.put(")", &DefaultStroke)
}
func ProcessRawBracketSquare(stringifier *Stringifier, token *gotokenize.Token) {

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

func ProcessRawRegexToken(stringifier *Stringifier, token *gotokenize.Token) {

	if token.Type == TokenJSBracket ||
		token.Type == TokenJSBlock ||
		token.Type == TokenJSBracketSquare {
		ProcessRawRegexBlockStyle(stringifier, token)
	} else if token.Type == TokenJSRegex {
		ProcessRawRegexStream(stringifier, token)
	} else {
		stringifier.put(token.Content, &DefaultStroke)
	}
}
func ProcessRawRegexStream(stringifier *Stringifier, parentToken *gotokenize.Token) {

	iter := parentToken.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenJSBracket ||
			childToken.Type == TokenJSBlock ||
			childToken.Type == TokenJSBracketSquare {
			ProcessRawRegexBlockStyle(stringifier, childToken)
		} else {
			ProcessRawRegexToken(stringifier, childToken)
		}
	}
}
func ProcessRawRegexBlockStyle(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	ProcessRawRegexStream(stringifier, token)
	if token.Type == TokenJSBracket {
		stringifier.put(")", &DefaultStroke)
	} else if token.Type == TokenJSBlock {
		stringifier.put("}", &DefaultStroke)
	} else {
		stringifier.put("]", &DefaultStroke)
	}
}
func ProcessRawString(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.put(token.Children.ConcatStringContent(), &DefaultStroke)
	stringifier.put(token.Content, &DefaultStroke)
}
func ProcessRawKeyWord(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put(token.Content, &NeedSpaceStroke)
}
func ProcessRawWord(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put(token.Content, &NeedSpaceStroke)
}
func ProcessRawOperator(stringifier *Stringifier, token *gotokenize.Token) {

	stringifier.put(token.Content, &SpaceAfterStroke)
}
func ProcessRawStrongOperator(stringifier *Stringifier, token *gotokenize.Token) {

	stringifier.put(token.Content, &BreakAfterStroke)
}
func ProcessRawSoftOperator(stringifier *Stringifier, token *gotokenize.Token) {

	stringifier.put(token.Content, &BreakAfterStroke)
}

func ProcessRawPhraseBreak(stringifier *Stringifier, token *gotokenize.Token) {
	if token.Content == ";" {
		stringifier.put(";", &BreakAfterStroke)
	} else {
		stringifier.put("", &BreakAfterStroke)
	}
}
func ProcessRawStrongBreak(stringifier *Stringifier, token *gotokenize.Token) {

	stringifier.put(";", &BreakAfterStroke)
}
