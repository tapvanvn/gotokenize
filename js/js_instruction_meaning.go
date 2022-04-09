package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

var optimizePhraseTokens []int = []int{
	TokenJSPhrase,
	TokenJSClass,
	TokenJSFunction,
	TokenJSFunctionLambda,
	TokenJSClassFunction,
	TokenJSBracket,
	TokenJSBracketSquare,
	TokenJSSwitch,
	TokenJSFor,
	TokenJSDo,
	TokenJSWhile,
	TokenJSIfTrail,
	TokenJSBlock,
	TokenJSAssignVariable,
}

func NewJSInstructionMeaning(baseMeaning gotokenize.IMeaning) *JSInstructionMeaning {

	return &JSInstructionMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSInstructionMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSInstructionMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := meaning.getNextMeaningToken(process.Iter)
	if token != nil {
		if token.Type == TokenJSClassFunction {
			//fmt.Println("process token clas function")
			iter := token.Children.Iterator()
			iter.Seek(1)
			meaning.processIter(iter)
		} else if token.Type != TokenJSPhrase {

			meaning.processChild(token)
		} else if token.Children.Length() == 1 {

			if firstToken := token.Children.GetTokenAt(0); gotokenize.IsContainToken(optimizePhraseTokens, firstToken.Type) {
				return firstToken
			}
		}
	}
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

		} else if token.Type == TokenJSFunctionLambda {

			meaning.continueLambda(iter, token)

		} else if token.Type == TokenJSAssign {

			meaning.continueAssignValue(iter, token)

		} else if token.Type == TokenJSBlock {

			meaning.processBlock(token)

		} else if token.Type != TokenJSPhraseBreak && !gotokenize.IsContainToken(optimizePhraseTokens, token.Type) /*token.Type != TokenJSPhrase*/ {

			//fmt.Println("begin phrase", JSTokenName(token.Type))
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			meaning.processChild(token)
			tmpToken.Children.AddToken(*token)
			meaning.continuePhrase(iter, tmpToken)
			return tmpToken

		} else if token.Type == TokenJSPhraseBreak {

			return meaning.getNextMeaningToken(iter)
		}
		return token
	}

	return nil
}

func (meaning *JSInstructionMeaning) continuePassPhraseBreak(iter *gotokenize.Iterator) {
	for {
		token := iter.Get()
		if token == nil || token.Type != TokenJSPhraseBreak {
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
		//fmt.Println("in phrase", token.Type)
		_ = iter.Read()
		if token.Type == TokenJSFunctionLambda {

			meaning.continueLambda(iter, token)
		}
		meaning.processChild(token)
		currentToken.Children.AddToken(*token)
	}

	if currentToken.Children.Length() == 1 {

		if firstToken := currentToken.Children.GetTokenAt(0); gotokenize.IsContainToken(optimizePhraseTokens, firstToken.Type) {
			currentToken.Type = firstToken.Type
			currentToken.Content = firstToken.Content
			currentToken.Children = firstToken.Children
		}
	}
}

func (meaning *JSInstructionMeaning) continueFor(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	meaning.continuePassPhraseBreak(iter)
	phraseToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(iter, phraseToken)
	currentToken.Children.AddToken(*phraseToken)
}
func (meaning *JSInstructionMeaning) processIter(iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		meaning.processChild(token)
	}
}

func (meaning *JSInstructionMeaning) processChild(parentToken *gotokenize.Token) {

	if !gotokenize.IsContainToken(JSInstructionGlobalNested, parentToken.Type) {
		return
	}
	proc := gotokenize.NewMeaningProcessFromStream(&parentToken.Children)

	//meaning.Prepare(proc)
	newStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		newStream.AddToken(*token)
	}
	parentToken.Children = newStream
}

func (meaning *JSInstructionMeaning) continueIfTrail(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	meaning.continuePassPhraseBreak(iter)

	ifBodyPhrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(iter, ifBodyPhrase)
	currentToken.Children.AddToken(*ifBodyPhrase)
	currentToken.Type = TokenJSIfTrail

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
		meaning.processChild(nextToken)
		currentToken.Children.AddToken(*nextToken)
		meaning.continuePassPhraseBreak(iter)
		testWhileToken := iter.Get()
		if testWhileToken != nil && testWhileToken.Type == TokenJSWhile {
			iter.Read()
			meaning.continuePassPhraseBreak(iter)
			if testWhileBodyToken := iter.Read(); testWhileBodyToken != nil && testWhileBodyToken.Type == TokenJSBracket {
				meaning.processChild(testWhileBodyToken)
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
		meaning.processChild(nextToken)
		currentToken.Children.AddToken(*nextToken)
		meaning.continuePassPhraseBreak(iter)
		if bodyToken := iter.Read(); bodyToken != nil {
			meaning.processChild(bodyToken)
			currentToken.Children.AddToken(*bodyToken)
		}
	}
	//else error
}

func (meaning *JSInstructionMeaning) continueLambda(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	meaning.continuePassPhraseBreak(iter)
	phraseToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(iter, phraseToken)
	currentToken.Children.AddToken(*phraseToken)

}

func (meaning *JSInstructionMeaning) continueAssignValue(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	meaning.continuePassPhraseBreak(iter)
	currentToken.Type = TokenJSAssignVariable
	phraseToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(iter, phraseToken)
	currentToken.Children.AddToken(*phraseToken)
}

func (meaning *JSInstructionMeaning) processBlock(currentToken *gotokenize.Token) {

	newStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	iter := currentToken.Children.Iterator()
	meaning.continuePassPhraseBreak(iter)
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		newStream.AddToken(*token)
	}
	currentToken.Children = newStream
}

func (meaning *JSInstructionMeaning) GetName() string {

	return "JSInstructionMeaning"
}
