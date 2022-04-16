package js

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

func NewJSPhraseMeaning(baseMeaning gotokenize.IMeaning) *JSPhraseMeaning {

	return &JSPhraseMeaning{

		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSPhraseMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSPhraseMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	if len(process.Context.AncestorTokens) == 0 && process.Iter.Offset == 0 {
		fmt.Print("\033[s") //save cursor the position
	}
	stackPhraseToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")

	token := meaning.getNextMeaningToken(&process.Context, process.Iter, stackPhraseToken)

	if token != nil {

		if token.Type == TokenJSPhrase {

			if token.Children.Length() == 1 { //avoid single element phrase

				*token = *token.Children.GetTokenAt(0)

			} else if token.Children.Length() == 0 {

				return nil
			}
		}

		process.Context.PreviousToken = token.Type
		process.Context.PreviousTokenContent = token.Content

	} else {
		fmt.Println("token nil")

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

func (meaning *JSPhraseMeaning) newPhraseToken(tokenType int) *gotokenize.Token {

	return gotokenize.NewToken(meaning.GetMeaningLevel(), tokenType, "")
}

func (meaning *JSPhraseMeaning) newStackToken() *gotokenize.Token {

	return gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
}

func (meaning *JSPhraseMeaning) processChild(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	if !gotokenize.IsContainToken(JSLevel2GlobalNested, parentToken.Type) {

		return
	}

	proc := gotokenize.NewMeaningProcessFromStream(append(context.AncestorTokens, parentToken.Type), &parentToken.Children)

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
func (meaning *JSPhraseMeaning) optimizePhrase(token *gotokenize.Token) *gotokenize.Token {

	if token.Type == TokenJSPhrase {
		if token.Children.Length() == 1 {

			return token.Children.GetTokenAt(0)
		}
	} else if token.Children.Length() == 1 && token.Children.GetTokenAt(0).Type == TokenJSPhrase {

		token.Children = token.Children.GetTokenAt(0).Children
	}
	return token
}

func (meaning *JSPhraseMeaning) getNextMeaningToken(context *gotokenize.MeaningContext,
	iter *gotokenize.Iterator, stackToken *gotokenize.Token) *gotokenize.Token {

	for {
		token := iter.Get()
		if token == nil {
			break
		}
		//fmt.Println(token.Content)
		if token.Content == "for" {

			if stackToken.Children.Length() > 0 {
				return meaning.optimizePhrase(stackToken)
			}

			_ = iter.Read()
			forPhrase := meaning.newPhraseToken(TokenJSPhraseFor)
			forPhrase.AddChild(*token)
			meaning.continuePassPhraseBreak(context, iter)
			if next := iter.Get(); next != nil && next.Type == TokenJSBracket {
				_ = iter.Read()
				meaning.processChild(context, next)
				forPhrase.AddChild(*meaning.optimizePhrase(next))

				meaning.continuePassPhraseBreak(context, iter)
				body := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())
				forPhrase.AddChild(*body)
			}

			return forPhrase

		} else if token.Content == "do" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			doPhrase := meaning.newPhraseToken(TokenJSPhraseDo)
			doPhrase.AddChild(*token)
			meaning.continuePassPhraseBreak(context, iter)
			body := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())
			doPhrase.AddChild(*body)
			meaning.continuePassPhraseBreak(context, iter)
			if next := iter.Get(); next != nil && next.Content == "while" {
				_ = iter.Read()
				doPhrase.AddChild(*next)
				meaning.continuePassPhraseBreak(context, iter)
				if next := iter.Get(); next != nil && next.Type == TokenJSBracket {
					_ = iter.Read()
					meaning.processChild(context, next)
					doPhrase.AddChild(*meaning.optimizePhrase(next))
				} else {
					//error
				}
			} else {
				//error
			}
			return doPhrase

		} else if token.Content == "while" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			whilePhrase := meaning.newPhraseToken(TokenJSPhraseWhile)
			whilePhrase.AddChild(*token)
			meaning.continuePassPhraseBreak(context, iter)
			if next := iter.Get(); next != nil && next.Type == TokenJSBracket {
				_ = iter.Read()
				meaning.processChild(context, next)
				whilePhrase.AddChild(*meaning.optimizePhrase(next))

				meaning.continuePassPhraseBreak(context, iter)
				body := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())
				//meaning.processChild(context, body)
				whilePhrase.AddChild(*body)
			} else {
				//error
			}
			return whilePhrase

		} else if token.Content == "if" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			ifPhrase := meaning.newPhraseToken(TokenJSPhraseIfTrail)
			ifPhrase.AddChild(*token)
			meaning.nextIfTrail(context, iter, ifPhrase)
			return ifPhrase

		} else if token.Content == "function" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			functionPhrase := meaning.newPhraseToken(TokenJSPhraseFunction)
			functionPhrase.AddChild(*token)
			meaning.nextFunction(context, iter, functionPhrase)
			return functionPhrase

		} else if token.Content == "class" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			classPhrase := meaning.newPhraseToken(TokenJSPhraseClass)
			classPhrase.AddChild(*token)
			meaning.nextClass(context, iter, classPhrase)
			return classPhrase

		} else if token.Content == "switch" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			switchPhrase := meaning.newPhraseToken(TokenJSPhraseSwitch)
			switchPhrase.AddChild(*token)
			meaning.nextSwitch(context, iter, switchPhrase)

			return switchPhrase

		} else if token.Content == "try" {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			tryPhrase := meaning.newPhraseToken(TokenJSPhraseTry)

			tryPhrase.AddChild(*token)
			meaning.nextTryCatch(context, iter, tryPhrase)
			return tryPhrase

		} else if token.Type == TokenJSPhraseBreak {

			if stackToken.Children.Length() > 0 {

				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			return token

		} else if token.Content == "=>" { //lambda
			_ = iter.Read()
			lambdaPhrase := meaning.newPhraseToken(TokenJSPhraseLambda)

			lambdaPhrase.AddChild(*stackToken)

			lambdaPhrase.AddChild(*token)

			rightPart := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())

			lambdaPhrase.AddChild(*rightPart)
			return lambdaPhrase

		} else if token.Type == TokenJSAssign {

			_ = iter.Read()
			assignPhrase := meaning.newPhraseToken(TokenJSPhraseAssign)

			assignPhrase.AddChild(*stackToken)

			assignPhrase.AddChild(*token)

			rightPart := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())
			assignPhrase.AddChild(*rightPart)

			return assignPhrase
		} else if token.Type == TokenJSQuestionOperator {
			_ = iter.Read()
			inlineIfPhrase := meaning.newPhraseToken(TokenJSPhraseInlineIf)

			inlineIfPhrase.AddChild(*stackToken)
			inlineIfPhrase.AddChild(*token)
			firstPart := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())
			//meaning.processChild(context, firstPart)
			inlineIfPhrase.AddChild(*firstPart)
			if next := iter.Get(); next != nil && next.Type == TokenJSColonOperator {
				_ = iter.Read()
				inlineIfPhrase.AddChild(*next)
				secondPart := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())
				//meaning.processChild(context, firstPart)
				inlineIfPhrase.AddChild(*secondPart)
			} else {
				//TODO: error
			}
			return inlineIfPhrase

		} else if token.Content == "return" {
			returnPhrase := meaning.newPhraseToken(TokenJSReturnStatement)
			_ = iter.Read()
			returnPhrase.AddChild(*token)
			meaning.nextReturnStatement(context, iter, returnPhrase)
			return returnPhrase
		} else if token.Type == TokenJSColonOperator {

			if stackToken.Children.Length() > 0 {
				return meaning.optimizePhrase(stackToken)
			}
			_ = iter.Read()
			return token
		} else {
			_ = iter.Read()
			meaning.processChild(context, token)
			stackToken.Children.AddToken(*meaning.optimizePhrase(token))

		}
	}

	return meaning.optimizePhrase(stackToken)
}

func (meaning *JSPhraseMeaning) nextReturnStatement(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, returnPhrase *gotokenize.Token) {

	meaning.continuePassPhraseBreak(context, iter)
	if next := iter.Get(); next != nil && next.Type == TokenJSPhraseBreak {
		return
	}
	returnBody := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())

	returnPhrase.AddChild(*returnBody)
}

func (meaning *JSPhraseMeaning) nextIfTrail(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, ifPhrase *gotokenize.Token) {
	meaning.continuePassPhraseBreak(context, iter)
	if next := iter.Get(); next != nil && next.Type == TokenJSBracket {
		_ = iter.Read()
		meaning.processChild(context, next)
		ifPhrase.AddChild(*meaning.optimizePhrase(next))

		meaning.continuePassPhraseBreak(context, iter)

		body := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())

		ifPhrase.AddChild(*body)

		for {
			meaning.continuePassPhraseBreak(context, iter)
			//read body
			if next := iter.Get(); next != nil && next.Content == "else" {
				_ = iter.Read()
				ifPhrase.AddChild(*next)

				meaning.continuePassPhraseBreak(context, iter)
				if next := iter.Get(); next != nil && next.Content == "if" {

					_ = iter.Read()
					ifPhrase.AddChild(*next)
					meaning.continuePassPhraseBreak(context, iter)
					if next := iter.Get(); next != nil && next.Type == TokenJSBracket {
						_ = iter.Read()
						meaning.processChild(context, next)
						ifPhrase.AddChild(*meaning.optimizePhrase(next))
					}
				}
				meaning.continuePassPhraseBreak(context, iter)
				body := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())

				ifPhrase.AddChild(*body)
			} else {
				break
			}
		}
	}
}

func (meaning *JSPhraseMeaning) nextFunction(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, funcPhrase *gotokenize.Token) {
	//phrase may contain name
	meaning.continuePassPhraseBreak(context, iter)
	next := iter.Get()
	if next == nil {
		return
	}
	if next.Type == TokenJSWord {
		_ = iter.Read()
		funcPhrase.AddChild(*next)
		meaning.continuePassPhraseBreak(context, iter)
		next = iter.Get()
		if next == nil {
			return
		}
	}

	if next.Type != TokenJSBracket {
		return
	}
	_ = iter.Read()
	meaning.processChild(context, next)
	funcPhrase.AddChild(*meaning.optimizePhrase(next))
	meaning.continuePassPhraseBreak(context, iter)
	body := meaning.getNextMeaningToken(context, iter, meaning.newStackToken())

	funcPhrase.AddChild(*body)
}

func (meaning *JSPhraseMeaning) nextClass(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, classPhrase *gotokenize.Token) {
	meaning.continuePassPhraseBreak(context, iter)
	//class may have class name
	next := iter.Get()
	if next == nil {
		return
	}
	if next.Type == TokenJSWord || (next.Type == TokenJSKeyWord && next.Content == "extends") {
		_ = iter.Read()
		classPhrase.AddChild(*next)
		meaning.continuePassPhraseBreak(context, iter)
		if next.Type == TokenJSWord {
			if next2 := iter.Get(); next2 != nil && next2.Content == "extends" {
				_ = iter.Read()
				classPhrase.AddChild(*next)
				meaning.continuePassPhraseBreak(context, iter)
				classPhrase.AddChild(*iter.Read()) //base class
				meaning.continuePassPhraseBreak(context, iter)
			}
		} else {
			classPhrase.AddChild(*iter.Read()) //base class
			meaning.continuePassPhraseBreak(context, iter)
		}
		next = iter.Get()
		if next == nil {
			return
		}
	}
	if next.Type == TokenJSBlock {
		_ = iter.Read()

		meaning.nextClassBody(context, next)

		classPhrase.AddChild(*next)
	}
}

func (meaning *JSPhraseMeaning) nextClassBody(context *gotokenize.MeaningContext, classBlock *gotokenize.Token) {
	iter := classBlock.Children.Iterator()
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	for {
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		funcName := iter.Read()
		if funcName != nil &&
			((funcName.Type == TokenJSKeyWord && funcName.Content == "constructor") ||
				funcName.Type == TokenJSWord) {

			tmpToken := meaning.newPhraseToken(TokenJSPhraseClassFunction)
			tmpToken.Children.AddToken(*funcName)

			meaning.continuePassPhraseBreak(context, iter)
			if bracket := iter.Read(); bracket != nil {
				meaning.processChild(context, bracket)
				tmpToken.Children.AddToken(*meaning.optimizePhrase(bracket)) //bracket

				meaning.continuePassPhraseBreak(context, iter)
				if body := iter.Read(); body != nil {

					meaning.processChild(context, body)
					tmpToken.Children.AddToken(*meaning.optimizePhrase(body)) //body
				}
			}
			tmpStream.AddToken(*tmpToken)
		} else {
			break
		}
	}
	classBlock.Children = tmpStream
}

func (meaning *JSPhraseMeaning) nextTryCatch(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, tryPhrase *gotokenize.Token) {

	meaning.continuePassPhraseBreak(context, iter)
	if next := iter.Get(); next != nil && next.Type == TokenJSBlock {

		_ = iter.Read()
		meaning.processChild(context, next)
		tryPhrase.AddChild(*next)
		if next := iter.Get(); next != nil && (next.Content == "catch" || next.Content == "finally") {
			_ = iter.Read()
			tryPhrase.AddChild(*next)

			if next.Content == "catch" {
				if next := iter.Get(); next != nil && next.Type == TokenJSBracket {
					_ = iter.Read()
					meaning.processChild(context, next)
					tryPhrase.AddChild(*meaning.optimizePhrase(next))
					if next := iter.Get(); next != nil && next.Type == TokenJSBlock {
						_ = iter.Read()
						meaning.processChild(context, next)
						tryPhrase.AddChild(*meaning.optimizePhrase(next))
					}
				}
			} else {
				if next := iter.Get(); next != nil && next.Type == TokenJSBlock {
					_ = iter.Read()
					meaning.processChild(context, next)
					tryPhrase.AddChild(*meaning.optimizePhrase(next))
				}
				return
			}
			if next := iter.Get(); next != nil && next.Content == "finally" {
				_ = iter.Read()
				tryPhrase.AddChild(*next)
				if next := iter.Get(); next != nil && next.Type == TokenJSBlock {
					_ = iter.Read()
					meaning.processChild(context, next)
					tryPhrase.AddChild(*meaning.optimizePhrase(next))
				}
			}
		}
	}
}

func (meaning *JSPhraseMeaning) nextSwitch(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, switchPhrase *gotokenize.Token) {

	meaning.continuePassPhraseBreak(context, iter)
	if next := iter.Get(); next != nil && next.Type == TokenJSBracket {

		_ = iter.Read()
		meaning.processChild(context, next)
		switchPhrase.AddChild(*meaning.optimizePhrase(next))

		meaning.continuePassPhraseBreak(context, iter)
		if next := iter.Get(); next != nil && next.Type == TokenJSBlock {
			_ = iter.Read()

			meaning.nextSwitchBody(context, next)
			switchPhrase.AddChild(*next)
		}
	}
}

func (meaning *JSPhraseMeaning) nextSwitchBody(context *gotokenize.MeaningContext, switchBlock *gotokenize.Token) {
	iter := switchBlock.Children.Iterator()
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	var tmpPhrase *gotokenize.Token = nil
	for {

		if token := iter.Read(); token != nil {

			if token.Type == TokenJSCase || token.Type == TokenJSDefault {

				if tmpPhrase != nil {

					meaning.processChild(context, tmpPhrase)
					tmpStream.AddToken(*tmpPhrase)
				}
				tmpPhrase = gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")

			}
			if tmpPhrase != nil {

				tmpPhrase.Children.AddToken(*token)
			}
		} else {
			break
		}
	}
	if tmpPhrase != nil {

		meaning.processChild(context, tmpPhrase)
		tmpStream.AddToken(*tmpPhrase)
	}
	switchBlock.Children = tmpStream
}

func (meaning *JSPhraseMeaning) continuePassPhraseBreak(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) {
	for {
		token := iter.Get()
		if token == nil || token.Type != TokenJSPhraseBreak {
			break
		}
		_ = iter.Read()
	}
}

func (meaning *JSPhraseMeaning) GetName() string {

	return "JSPhraseMeaning"
}
