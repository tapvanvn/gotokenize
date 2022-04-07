package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewJSInstructionMeaning(baseMeaning gotokenize.IMeaning) *JSInstructionMeaning {

	//abtract := gotokenize.NewAbtractMeaning(baseMeaning)
	//patternMeaning := gotokenize.NewPatternMeaning(abtract, JSInstructionPatterns, gotokenize.NoTokens, JSInstructionGlobalNested)

	return &JSInstructionMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSInstructionMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSInstructionMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := meaning.getNextMeaningToken(process.Iter)
	return token
}

func (meaning *JSInstructionMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	if token := iter.Read(); token != nil {

		if token.Type == TokenJSFor {

			meaning.continueFor(iter, token)

		} else if token.Type == TokenJSIf {

			meaning.continueIfTrail(iter, token)

		} else if token.Type == TokenJSDo {

			meaning.continueDo(iter, token)

		} else if token.Type == TokenJSWhile {

			meaning.continueWhite(iter, token)
		}

		return token
	}

	return nil
}
func (meaning *JSInstructionMeaning) continuePassPhraseBreak(iter *gotokenize.Iterator) {
	for {
		token := iter.Get()
		if token.Type != TokenJSPhraseBreak {
			break
		}
		_ = iter.Read()
	}
}

func (meaning *JSInstructionMeaning) continuePhrase(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	for {
		token := iter.Get()
		if token == nil || token.Type == TokenJSPhraseBreak || !gotokenize.IsContainToken(JSPhraseAllow, token.Type) {
			break
		}
		_ = iter.Read()
		currentToken.Children.AddToken(*token)
	}
}

func (meaning *JSInstructionMeaning) continueFor(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	meaning.continuePassPhraseBreak(iter)
	nextToken := iter.Read()
	if nextToken != nil {
		currentToken.Children.AddToken(*nextToken)
	}
}

func (meaning *JSInstructionMeaning) continueIfTrail(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	meaning.continuePassPhraseBreak(iter)

	ifBodyPhrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(iter, ifBodyPhrase)
	currentToken.Children.AddToken(*ifBodyPhrase)

	for {
		testElseIf := iter.Get()
		if testElseIf == nil || testElseIf.Type != TokenJSElseIf {
			break
		}
		iter.Read()
		meaning.continuePassPhraseBreak(iter)
		elseIfBodyPhrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(iter, elseIfBodyPhrase)
		currentToken.Children.AddToken(*testElseIf)
		currentToken.Children.AddToken(*elseIfBodyPhrase)
	}
	testElse := iter.Get()
	if testElse != nil && testElse.Type == TokenJSElse {
		iter.Read()
		meaning.continuePassPhraseBreak(iter)
		elseBodyPhrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(iter, elseBodyPhrase)
		currentToken.Children.AddToken(*testElse)
		currentToken.Children.AddToken(*elseBodyPhrase)
	}
}

func (meaning *JSInstructionMeaning) continueDo(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	meaning.continuePassPhraseBreak(iter)
	if nextToken := iter.Read(); nextToken != nil {
		currentToken.Children.AddToken(*nextToken)
		meaning.continuePassPhraseBreak(iter)
		testWhileToken := iter.Get()
		if testWhileToken != nil && testWhileToken.Type == TokenJSWhile {
			iter.Read()
			meaning.continuePassPhraseBreak(iter)
			if testWhileBodyToken := iter.Read(); testWhileBodyToken != nil && testWhileBodyToken.Type == TokenJSBracket {
				testWhileToken.Children.AddToken(*testWhileBodyToken)
			}
			currentToken.Children.AddToken(*testWhileToken)
		}
		//else error
	}

}

func (meaning *JSInstructionMeaning) continueWhite(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	meaning.continuePassPhraseBreak(iter)
	if nextToken := iter.Read(); nextToken != nil && nextToken.Type == TokenJSBracket {
		currentToken.Children.AddToken(*nextToken)
		meaning.continuePassPhraseBreak(iter)
		if bodyToken := iter.Read(); bodyToken != nil {
			currentToken.Children.AddToken(*bodyToken)
		}
	}
	//else error
}

func (meaning *JSInstructionMeaning) processPhrase(currentToken *gotokenize.Token) {

	/*process := gotokenize.NewMeaningProcessFromStream(&currentToken.Children)
	newStream := gotokenize.CreateStream(0)

	for {
		token := meaning.patternMeaning.Next(process)

		if token == nil {

			break
		}
		if token.Type == TokenJSAssign {
			fmt.Println("found assign")
			fmt.Println(meaning.patternMeaning.Patterns)
		}
		if gotokenize.IsContainToken(JSInstructionGlobalNested, token.Type) {

			chilProc := gotokenize.NewMeaningProcessFromStream(&token.Children)
			chilStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

			for {
				childToken := meaning.Next(chilProc)

				if childToken == nil {

					break
				}

				chilStream.AddToken(*childToken)
			}
			token.Children = chilStream
		}
		newStream.AddToken(*token)
	}
	currentToken.Children = newStream*/
}

func (meaning *JSInstructionMeaning) GetName() string {

	return "JSInstructionMeaning"
}
