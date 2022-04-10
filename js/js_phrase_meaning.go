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
	TokenJSUnaryOperator,
	TokenJSBinaryOperator,
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
}
var jsPhraseNext = []int{
	TokenJSAssign,
	//TokenJSBinaryOperator,
	TokenJSRightArrow,
	TokenJSQuestionOperator,
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
,break,byte,case,catch,
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
,throw,throws,transient,true,
,var,void,try,
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
	return token
}
func (meaning *JSPhraseMeaning) processChild(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	if !gotokenize.IsContainToken(JSLevel2GlobalNested, parentToken.Type) {
		if parentToken.Type == TokenJSRightArrow {

		}
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

			fmt.Println("found lambda")
			meaning.nextLambdaBody(context, iter, token)

		} else if token.Content == "function" {

			meaning.nextFunction(context, iter, token)

		} else if token.Content == "class" {
			meaning.nextClass(context, iter, token)

		} else if token.Content == "switch" {

			meaning.nextSwitch(context, iter, token)
			/*
				} else if token.Content == "?" {

					meaning.nextInlineIfBody(context, iter, token, 0)
					//*/
		} else if token.Content == "try" {

			meaning.nextTryCatch(context, iter, token)

		} else if token.Type != TokenJSPhraseBreak &&
			( //token.Type == TokenJSWord ||
			token.Type != TokenJSAssign &&
				token.Type != TokenJSColonOperator &&
				token.Type != TokenJSQuestionOperator &&
				gotokenize.IsContainToken(jsPhraseAllow, token.Type) &&
				!gotokenize.IsContainToken(jsPhraseBreakAfter, token.Type) &&
				!(token.Type == TokenJSKeyWord && strings.Contains(jsPhraseBreakKeyWords, ","+token.Content+","))) {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			meaning.processChild(context, token)
			tmpToken.Children.AddToken(*token)
			meaning.continuePhrase(context, iter, tmpToken)

			return tmpToken

		} else if token.Type == TokenJSPhraseBreak {

			return meaning.getNextMeaningToken(context, iter)
		}
		return token
	}

	return nil
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
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)      //for
	meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
	if bracket := iter.Read(); bracket != nil {
		meaning.processChild(context, bracket)
		tmpToken.Children.AddToken(*bracket)           //bracket
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
	}
	*currentToken = *tmpToken
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
	tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	tmpToken.Children.AddToken(*currentToken)                                                                           //class
	meaning.continuePassPhraseBreak(context, iter)                                                                      //remove empty phrase break
	if nextToken := iter.Get(); nextToken != nil && (nextToken.Type == TokenJSWord || nextToken.Content == "extends") { //class name if existed
		tmpToken.Children.AddToken(*nextToken)
		_ = iter.Read()
		meaning.continuePassPhraseBreak(context, iter)

		if nextToken2 := iter.Get(); nextToken2 != nil && (nextToken2.Type == TokenJSWord || nextToken2.Content == "extends") {
			tmpToken.Children.AddToken(*nextToken2)
			_ = iter.Read()
			meaning.continuePassPhraseBreak(context, iter)
			tmpToken.Children.AddToken(*iter.Read())
		}
	}
	meaning.continuePassPhraseBreak(context, iter)
	if body := iter.Get(); body != nil && body.Type == TokenJSBlock {

		_ = iter.Read()
		meaning.nextClassBody(context, body.Children.Iterator(), body)
		tmpToken.Children.AddToken(*body)
	}

	*currentToken = *tmpToken
}
func (meaning *JSPhraseMeaning) nextClassBody(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	for {
		meaning.continuePassPhraseBreak(context, iter) //remove empty phrase break
		if funcName := iter.Get(); funcName != nil &&
			((funcName.Type == TokenJSKeyWord && funcName.Content == "constructor") ||
				funcName.Type == TokenJSWord) {
			_ = iter.Read()
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			tmpToken.Children.AddToken(*funcName)
			if bracket := iter.Read(); bracket != nil {
				meaning.processChild(context, bracket)
				tmpToken.Children.AddToken(*bracket) //bracket
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

	for {
		token := iter.Get()

		if token == nil {

			break
		}
		//Check if the phrase can continue

		if token.Type == TokenJSPhraseBreak {
			if lastToken := iter.GetBy(-1); lastToken == nil || lastToken.Type != TokenJSBinaryOperator {
				if nextToken := iter.GetBy(1); nextToken == nil || nextToken.Type != TokenJSBinaryOperator {

					break
				}
			}
			fmt.Println("here", token.Type, token.Content)
		} else if !gotokenize.IsContainToken(jsPhraseAllow, token.Type) {
			break
		}

		_ = iter.Read()

		meaning.processChild(context, token)

		if gotokenize.IsContainToken(jsPhraseNext, token.Type) {

			if currentToken.Type == TokenJSPhrase && currentToken.Children.Length() == 1 {

				*currentToken = *currentToken.Children.GetTokenAt(0)
			}

			tmpContainerPhrase := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			tmpContainerPhrase.Children.AddToken(*currentToken)
			tmpContainerPhrase.Children.AddToken(*token)

			if childToken := meaning.getNextMeaningToken(context, iter); childToken != nil {
				meaning.processChild(context, childToken)
				if childToken.Type == TokenJSPhrase && childToken.Children.Length() == 1 {

					*childToken = *childToken.Children.GetTokenAt(0)
				}
				tmpContainerPhrase.Children.AddToken(*childToken)
			}

			*currentToken = *tmpContainerPhrase

		} else {

			currentToken.Children.AddToken(*token)
		}
		if gotokenize.IsContainToken(jsPhraseBreakAfter, token.Type) ||
			(token.Type == TokenJSKeyWord && strings.Contains(jsPhraseBreakKeyWords, ","+token.Content+",")) {

			break
		}
	}
	if currentToken.Children.Length() == 1 {

		*currentToken = *currentToken.Children.GetTokenAt(0)
	}
}

func (meaning *JSPhraseMeaning) GetName() string {

	return "JSPhraseMeaning"
}
