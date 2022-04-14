package js

import (
	"fmt"
	"strings"

	"github.com/tapvanvn/gotokenize/v2"
)

var notOperatorTrail = ",:,"

func NewJSOperatorMeaning(baseMeaning gotokenize.IMeaning) *JSOperatorMeaning {

	return &JSOperatorMeaning{

		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSOperatorMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSOperatorMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {
	if len(process.Context.AncestorTokens) == 0 && process.Iter.Offset == 0 {
		fmt.Print("\033[s") //save cursor the position
	}
	token := meaning.getNextMeaningToken(&process.Context, process.Iter)

	if token != nil {
		if gotokenize.IsContainToken(JSGlobalNested, token.Type) {

			meaning.processChild(&process.Context, token)
		}
		process.Context.PreviousToken = token.Type
		process.Context.PreviousTokenContent = token.Content
	} else {
		process.Context.PreviousToken = gotokenize.TokenNoType
		process.Context.PreviousTokenContent = ""
	}

	if len(process.Context.AncestorTokens) == 0 {
		fmt.Print("\033[u\033[K") //restore
		fmt.Printf("%s percent: %f%%\n", meaning.GetName(), process.GetPercent())
		fmt.Print("\033[A")
	}

	return token
}

func (meaning *JSOperatorMeaning) getNextMeaningToken(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) *gotokenize.Token {

	if token := iter.Read(); token != nil {

		if (token.Type == TokenJSWord ||
			token.Type == TokenJSPhrase ||
			token.Type == TokenJSString || token.Type == TokenJSBracket) && meaning.testOperatorPhrase(context, iter) {

			meaning.continueOperatorPhrase(context, iter, token)
		}
		return token
	}

	return nil
}
func (meaning *JSOperatorMeaning) processChild(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	proc := gotokenize.NewMeaningProcessFromStream(append(context.AncestorTokens, parentToken.Type), &parentToken.Children)

	newStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		if gotokenize.IsContainToken(JSGlobalNested, token.Type) {

			//meaning.processChild(context, token)
		} else {
			//fmt.Println("stop at", JSTokenName(token.Type))
		}
		newStream.AddToken(*token)
	}
	parentToken.Children = newStream
}
func (meaning *JSOperatorMeaning) testOperatorPhrase(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) bool {

	numOperator := 0
	i := 0
	for {
		token := iter.GetBy(i)

		if token == nil || strings.Contains(notOperatorTrail, ","+token.Content+",") {
			break
		}
		if token.Type == TokenJSBinaryOperator {
			numOperator++
		}
		i++
	}
	return numOperator > 0
}
func (meaning *JSOperatorMeaning) continueOperatorPhrase(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, parentToken *gotokenize.Token) {
	trail := gotokenize.CreateStream(meaning.GetMeaningLevel())
	trail.AddToken(*parentToken)
	numOperator := 0
	for {
		token := iter.Get()

		if token == nil || strings.Contains(notOperatorTrail, ","+token.Content+",") {
			break
		}
		if token.Type == TokenJSBinaryOperator {
			numOperator++
		}
		_ = iter.Read()
		trail.AddToken(*token)
		//trail.AddToken(*iter.Read())
	}
	if trail.Length() > 2 && numOperator > 0 {
		parentToken.Type = TokenJSOperatorTrail
		parentToken.Children = trail
	}
}

func (meaning *JSOperatorMeaning) GetName() string {

	return "JSOperatorMeaning"
}
