package js

import (
	"fmt"
	"strings"

	"github.com/tapvanvn/gotokenize/v2"
)

var jsPhraseAllow = []int{
	TokenJSWord,
	TokenJSKeyWord,
	TokenJSOperator,

	TokenJSTreeDotOperator,
	TokenJSAssign,
	//TokenJSUnaryOperator,
	//TokenJSBinaryOperator,

	TokenJSQuestionOperator,
	TokenJSColonOperator,

	TokenJSRegex,
	TokenJSString,

	TokenJSVariable,
	TokenJSBlock,
	TokenJSBracket,
	TokenJSBracketSquare,
	TokenJSRightArrow,
	TokenJSFunctionLambda,
	TokenJSFunction,
	TokenJSClass,
	TokenJSReturnStatement,
}

var jsPhraseNext = []int{
	TokenJSAssign,
	//TokenJSQuestionOperator,
	//TokenJSColonOperator,
	//TokenJSBinaryOperator,
	//TokenJSRightArrow,
	//TokenJSQuestionOperator,
	TokenJSColonOperator,
}

var jsPhraseBreakAfter = []int{
	//TokenJSBracket,
	//TokenJSBracketSquare,
	//TokenJSBlock,
}

type JSPhraseContext struct {
	InContextOfToken int
}

var jsPhraseBreakKeyWords string = `
,abstract,arguments,await,boolean,
,byte,case,catch,
,char,class,const,continue,
,debugger,default,delete,do,
,double,else,enum,eval,
,export,false,final,finally,float,for,function,
,goto,if,implements,import,
,int,interface,
,let,long,native,
,package,private,protected,
,public,short,static,
,super,switch,synchronized,
,transient,true,
,var,try,return,
,volatile,while,with,yield,constructor,`

func NewJSPhraseMeaning(baseMeaning gotokenize.IMeaning) *JSPhraseMeaning {

	return &JSPhraseMeaning{

		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSPhraseMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSPhraseMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	/*if len(process.Context.AncestorTokens) == 0 && process.Iter.Offset == 0 {
		fmt.Print("\033[s") //save cursor the position
	}*/
	token := meaning.getNextMeaningToken(&process.Context, process.Iter)

	if token != nil {

		if token.Type != TokenJSPhrase {

			meaning.processChild(&process.Context, token)

			if token.Children.Length() == 1 {
				//get the first if phrase have only one blocking typed child
				if firstToken := token.Children.GetTokenAt(0); firstToken.Type == TokenJSPhrase {

					token.Children = firstToken.Children
				}
			}
		} else if token.Children.Length() == 1 { //avoid single element phrase

			*token = *token.Children.GetTokenAt(0)

		} else if token.Children.Length() == 0 {

			return meaning.getNextMeaningToken(&process.Context, process.Iter) //remove empty phrase
		}
		process.Context.PreviousToken = token.Type
		process.Context.PreviousTokenContent = token.Content

	} else {

		process.Context.PreviousToken = gotokenize.TokenNoType
		process.Context.PreviousTokenContent = ""
	}
	/*if len(process.Context.AncestorTokens) == 0 {
		fmt.Print("\033[u\033[K") //restore
		fmt.Printf("%s percent: %f%%\n", meaning.GetName(), process.GetPercent())
		fmt.Print("\033[A")
	}*/
	if token != nil {

		//fmt.Printf("found:[%s]\tcontent:%s\tlevel:%d\n", JSTokenName(token.Type), token.Content, len(process.Context.AncestorTokens))

	}
	return token
}
func (meaning *JSPhraseMeaning) processChild(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	if !gotokenize.IsContainToken(JSLevel2GlobalNested, parentToken.Type) {

		return
	}

	proc := gotokenize.NewMeaningProcessFromStream(append(context.AncestorTokens, parentToken.Type), &parentToken.Children)

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
func (meaning *JSPhraseMeaning) getNextMeaningToken(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) *gotokenize.Token {

	meaning.continuePassPhraseBreak(context, iter)

	if token := iter.Read(); token != nil {

		if token.Content == "for" {

			meaning.nextFor(context, iter, token)

		} else if token.Content == "do" {

			meaning.nextDo(context, iter, token)

		} else if token.Content == "while" {

			meaning.nextWhile(context, iter, token)

		} else if token.Content == "if" {

			meaning.nextIfTrail(context, iter, token)

		} else if token.Content == "=>" { //lambda

			meaning.nextLambdaBody(context, iter, token)

		} else if token.Content == "function" {

			meaning.nextFunction(context, iter, token)

		} else if token.Content == "class" {

			meaning.nextClass(context, iter, token)

		} else if token.Content == "switch" {

			meaning.nextSwitch(context, iter, token)

		} else if token.Content == "try" {

			meaning.nextTryCatch(context, iter, token)

		} else if token.Content == "return" {

			meaning.nextReturnStatement(context, iter, token)

		} else if token.Type != TokenJSPhraseBreak &&
			( //token.Type == TokenJSWord ||
			token.Type != TokenJSAssign &&
				token.Type != TokenJSColonOperator &&
				token.Type != TokenJSQuestionOperator &&

				gotokenize.IsContainToken(jsPhraseAllow, token.Type) &&
				!gotokenize.IsContainToken(jsPhraseBreakAfter, token.Type) &&
				!(token.Type == TokenJSKeyWord && strings.Contains(jsPhraseBreakKeyWords, ","+token.Content+","))) {

			//fmt.Println("start phrase at", iter.Offset)
			//fmt.Println("begin ", JSTokenName(token.Type), token.Content)
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			meaning.processChild(context, token)
			tmpToken.Children.AddToken(*token)
			meaning.continuePhrase(context, iter, tmpToken)

			//fmt.Println("\t end phrase at", iter.Offset)
			if tmpToken.Type == TokenJSPhrase {

				if tmpToken.Children.Length() == 1 {

					*tmpToken = *tmpToken.Children.GetTokenAt(0)

				} else if tmpToken.Children.Length() == 0 {

					return meaning.getNextMeaningToken(context, iter)
				}
			}
			return tmpToken

		} else if token.Type == TokenJSPhraseBreak && token.Content != ";" {

			return meaning.getNextMeaningToken(context, iter)

		} else {

			//fmt.Println("not in ", JSTokenName(token.Type), token.Content)
		}
		if token.Type == TokenJSPhrase && token.Children.Length() == 1 {

			*token = *token.Children.GetTokenAt(0)
		}
		return token
	}

	return nil
}

func (meaning *JSPhraseMeaning) getNextNonPhraseToken(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) (*gotokenize.Token, bool, int) {
	processed := true
	backupOffset := iter.Offset
	fmt.Printf("start %s at %d\n", gotokenize.ColorContent("nonphrase"), iter.Offset)
	if token := iter.Read(); token != nil {

		if token.Content == "for" {

			meaning.nextFor(context, iter, token)

		} else if token.Content == "do" {

			meaning.nextDo(context, iter, token)

		} else if token.Content == "while" {

			meaning.nextWhile(context, iter, token)

		} else if token.Content == "if" {

			meaning.nextIfTrail(context, iter, token)

		} else if token.Content == "=>" { //lambda

			meaning.nextLambdaBody(context, iter, token)

		} else if token.Type == TokenJSAssign {

			meaning.nextAssign(context, iter, token)

		} else if token.Content == "function" {

			meaning.nextFunction(context, iter, token)

		} else if token.Content == "class" {

			meaning.nextClass(context, iter, token)

		} else if token.Content == "switch" {

			meaning.nextSwitch(context, iter, token)

		} else if token.Content == "try" {

			meaning.nextTryCatch(context, iter, token)

		} else if token.Content == "return" {

			meaning.nextReturnStatement(context, iter, token)

		} else {

			processed = false
		}
		if token.Type == TokenJSPhrase && token.Children.Length() == 1 {

			*token = *token.Children.GetTokenAt(0)
		}
		fmt.Printf("\tend %s at %d %s:[%s]\n", gotokenize.ColorContent("nonphrase"), iter.Offset, JSTokenName(token.Type), token.Content)
		return token, processed, backupOffset
	}

	return nil, false, backupOffset
}

func (meaning *JSPhraseMeaning) nextReturnStatement(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSReturnStatement, "")
	tmpToken.Children.AddToken(*currentToken) //return

	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBracket {
		meaning.processChild(context, nextToken)
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	} else {

		phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(context, iter, phrase)
		tmpToken.Children.AddToken(*phrase)

	}
	fmt.Printf("\tend %s at %d\n", gotokenize.ColorContent("return"), iter.Offset)
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextTryCatch(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //try
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if block := iter.Read(); block != nil {        //try body
		meaning.processChild(context, block)
		tmpToken.Children.AddToken(*block)
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		if catchToken := iter.Get(); catchToken != nil && catchToken.Content == "catch" {
			tmpToken.Children.AddToken(*catchToken)
			_ = iter.Read()
			meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
			if bracket := iter.Read(); bracket != nil {
				meaning.processChild(context, bracket)
				tmpToken.Children.AddToken(*bracket)
				meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
				if nextToken := iter.Read(); nextToken != nil {
					meaning.processChild(context, nextToken)
					tmpToken.Children.AddToken(*nextToken)
				}
			}
		}
	}
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextFor(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	fmt.Printf("found %s at %d\n", gotokenize.ColorContent("for"), iter.Offset)
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //for
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if bracket := iter.Read(); bracket != nil {
		meaning.processForBracket(context, bracket)

		tmpToken.Children.AddToken(*bracket)           //bracket
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBlock {
			meaning.processChild(context, nextToken)
			tmpToken.Children.AddToken(*nextToken)
			_ = iter.Read()
		} else {
			fmt.Printf("\t\tstart %s at %d\n", gotokenize.ColorContent("for body"), iter.Offset)
			phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			meaning.continuePhrase(context, iter, phrase)
			tmpToken.Children.AddToken(*phrase)
			fmt.Printf("\t\t\tend %s at %d\n", gotokenize.ColorContent("for body"), iter.Offset)
		}
	}
	fmt.Printf("\tend %s at %d\n", gotokenize.ColorContent("for"), iter.Offset)
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) processForBracket(context *gotokenize.MeaningContext, bracket *gotokenize.Token) {

	iter := bracket.Children.Iterator()
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		if token.Type == TokenJSPhraseBreak && token.Content == ";" {
			token.Type = TokenJSStrongBreak
		}
	}
	meaning.processChild(context, bracket)
}

func (meaning *JSPhraseMeaning) nextDo(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //do
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBlock {
		meaning.processChild(context, nextToken)
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	} else {
		phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(context, iter, phrase)
		tmpToken.Children.AddToken(*phrase)
	}
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	tmpToken.Children.AddToken(*iter.Read())       //while
	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBracket {
		meaning.processChild(context, nextToken)
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	} else {
		phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(context, iter, phrase)
		tmpToken.Children.AddToken(*phrase)
	}
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextWhile(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //while
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBracket {
		meaning.processChild(context, nextToken)
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	} else {
		phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(context, iter, phrase)
		tmpToken.Children.AddToken(*phrase)
	}
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBlock {
		meaning.processChild(context, nextToken)
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	} else {
		phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(context, iter, phrase)
		tmpToken.Children.AddToken(*phrase)
	}
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextIfTrail(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //if
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	tmpToken.Children.AddToken(*iter.Read())       //bracket
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil { //if body
		if nextToken.Type == TokenJSBlock {
			meaning.processChild(context, nextToken)
			tmpToken.Children.AddToken(*nextToken)
			_ = iter.Read()
		} else {
			phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			meaning.continuePhrase(context, iter, phrase)
			tmpToken.Children.AddToken(*phrase)
		}
	}
	//test else if
	for {
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		if nextToken := iter.Get(); nextToken != nil && nextToken.Content == "else" {
			tmpToken.Children.AddToken(*nextToken)
			_ = iter.Read()
			meaning.continuePassPhraseBreak(context, iter)
			nextToken2 := iter.Get()
			if nextToken2 != nil && nextToken2.Content == "if" {
				tmpToken.Children.AddToken(*nextToken2)
				_ = iter.Read()
				//else if
				meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
				if bracket := iter.Read(); bracket != nil {
					meaning.processChild(context, bracket)
					tmpToken.Children.AddToken(*bracket) //bracket
				}
				meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
				nextToken2 = iter.Get()
			}
			if nextToken2 != nil { //elseif block or else block
				if nextToken2.Type == TokenJSBlock {
					meaning.processChild(context, nextToken2)
					tmpToken.Children.AddToken(*nextToken2)
					_ = iter.Read()
				} else {
					phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
					meaning.continuePhrase(context, iter, phrase)
					tmpToken.Children.AddToken(*phrase)
				}
			}
		} else {
			break
		}
	}
	*currentToken = *tmpToken

}

func (meaning *JSPhraseMeaning) nextLambdaBody(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)                                        //=>
	meaning.continuePassPhraseBreak(context, iter)                                   //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSBlock { //body
		meaning.processChild(context, nextToken)
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	} else {
		phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
		meaning.continuePhrase(context, iter, phrase)
		tmpToken.Children.AddToken(*phrase)
	}

	*currentToken = *tmpToken

}

func (meaning *JSPhraseMeaning) nextAssign(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //assign
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break

	phrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(context, iter, phrase)
	tmpToken.Children.AddToken(*phrase)

	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextFunction(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)                                       //function
	meaning.continuePassPhraseBreak(context, iter)                                  //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil && nextToken.Type == TokenJSWord { //function name if existed
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
	}
	if bracket := iter.Read(); bracket != nil {
		meaning.processChild(context, bracket)
		tmpToken.Children.AddToken(*bracket)
		if body := iter.Read(); body != nil {
			meaning.processChild(context, body)
			tmpToken.Children.AddToken(*body)
		}
	}
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextClass(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "class")
	tmpToken.Children.AddToken(*currentToken) //class
	meaning.continuePassPhraseBreak(context, iter)
	if nextToken := iter.Get(); nextToken != nil && (nextToken.Type == TokenJSWord || nextToken.Content == "extends") { //class name if existed
		tmpToken.Children.AddToken(*nextToken) //name or extends
		_ = iter.Read()

		meaning.continuePassPhraseBreak(context, iter)

		if nextToken2 := iter.Get(); nextToken2 != nil {

			tmpToken.Children.AddToken(*nextToken2)
			_ = iter.Read()

			if nextToken2.Content == "extends" {
				meaning.continuePassPhraseBreak(context, iter)
				base := iter.Read()
				tmpToken.Children.AddToken(*base)
			}
		}
	}
	meaning.continuePassPhraseBreak(context, iter)
	body := iter.Get()
	if body != nil && body.Type == TokenJSBlock {
		_ = iter.Read()

		meaning.nextClassBody(context, body)

		tmpToken.Children.AddToken(*body)
	}

	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextClassBody(context *gotokenize.MeaningContext, currentToken *gotokenize.Token) {
	iter := currentToken.Children.Iterator()
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	for {
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		funcName := iter.Read()
		if funcName != nil &&
			((funcName.Type == TokenJSKeyWord && funcName.Content == "constructor") ||
				funcName.Type == TokenJSWord) {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			tmpToken.Children.AddToken(*funcName)

			meaning.continuePassPhraseBreak(context, iter)
			if bracket := iter.Read(); bracket != nil {
				meaning.processChild(context, bracket)
				tmpToken.Children.AddToken(*bracket) //bracket

				meaning.continuePassPhraseBreak(context, iter)
				if body := iter.Read(); body != nil {

					meaning.processChild(context, body)
					tmpToken.Children.AddToken(*body) //body
				}
			}
			tmpStream.AddToken(*tmpToken)
		} else {
			break
		}
	}
	currentToken.Children = tmpStream

}
func (meaning *JSPhraseMeaning) nextSwitch(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //switch
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	tmpToken.Children.AddToken(*iter.Read())       //bracket
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break

	if body := iter.Get(); body != nil && body.Type == TokenJSBlock {
		_ = iter.Read()
		meaning.nextSwitchBody(context, body.Children.Iterator(), body)

		tmpToken.Children.AddToken(*body)
	}
	*currentToken = *tmpToken
}

func (meaning *JSPhraseMeaning) nextSwitchBody(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	var tmpPhrase *gotokenize.Token = nil
	for {
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		if token := meaning.getNextMeaningToken(context, iter); token != nil {

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
	currentToken.Children = tmpStream
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

//continuePhrase read until the phrase end or meet a block
func (meaning *JSPhraseMeaning) continuePhrase(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	token, processed, backupOffset := meaning.getNextNonPhraseToken(context, iter)

	_ = backupOffset
	if processed && token != nil {
		currentToken.Children.AddToken(*token)

		return
	}
	iter.Seek(backupOffset)
	for {

		token := iter.Get()

		if token == nil {
			//*currentToken = *token
			//return
			//fmt.Println("\t\t\t\t restore1", backupOffset)
			//iter.Seek(backup)
			break
		}

		//Check if the phrase can continue when reach a phrase break
		if token.Type == TokenJSPhraseBreak {
			if lastToken := iter.GetBy(-1); lastToken == nil || lastToken.Type != TokenJSBinaryOperator {
				if nextToken := iter.GetBy(1); nextToken == nil || nextToken.Type != TokenJSBinaryOperator {

					//fmt.Println("\t token break", id)
					//iter.Seek(backup)
					//fmt.Println("\t\t\t\t restore2", backupOffset)
					break
				}
			}
			fmt.Println("====continue====")
		} else if !gotokenize.IsContainToken(jsPhraseAllow, token.Type) {
			//fmt.Println("\t token not allow", id)
			//iter.Seek(backup)
			//*currentToken = *token
			//fmt.Println("\t\t\t\t restore3", backupOffset)
			break
		}
		_ = iter.Read()
		meaning.processChild(context, token)

		//fmt.Println("\ton phrase:", id, JSTokenName(token.Type), token.Content)

		if gotokenize.IsContainToken(jsPhraseNext, token.Type) {

			if currentToken.Type == TokenJSPhrase && currentToken.Children.Length() == 1 {

				*currentToken = *currentToken.Children.GetTokenAt(0)
			}

			tmpContainerPhrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			tmpContainerPhrase.Children.AddToken(*currentToken)
			tmpContainerPhrase.Children.AddToken(*token)

			if childToken := meaning.getNextMeaningToken(context, iter); childToken != nil {

				if childToken.Type != TokenJSPhrase { //prevent loop

					meaning.processChild(context, childToken)
				}
				if childToken.Type == TokenJSPhrase && childToken.Children.Length() == 1 {

					*childToken = *childToken.Children.GetTokenAt(0)
				}
				tmpContainerPhrase.Children.AddToken(*childToken)
			}
			if tmpContainerPhrase.Children.Length() == 1 {

				*tmpContainerPhrase = *tmpContainerPhrase.Children.GetTokenAt(0)
			}
			//fmt.Println("\tpharese next", id)
			*currentToken = *tmpContainerPhrase
		} else {

			currentToken.Children.AddToken(*token)
		}
		if gotokenize.IsContainToken(jsPhraseBreakAfter, token.Type) ||
			(token.Type == TokenJSKeyWord && strings.Contains(jsPhraseBreakKeyWords, ","+token.Content+",")) {
			//fmt.Println("\t break after", id)
			break
		}
	}
	if currentToken.Type == TokenJSPhrase && currentToken.Children.Length() == 1 {

		*currentToken = *currentToken.Children.GetTokenAt(0)
		//fmt.Println("\tminimize token", id, JSTokenName(currentToken.Type), currentToken.Content)
	}
}

func (meaning *JSPhraseMeaning) GetName() string {

	return "JSPhraseMeaning"
}
