package js

import (
	"strings"

	"github.com/tapvanvn/gotokenize/v2"
)

var jsPhraseAllow = []int{
	TokenJSWord,
	TokenJSKeyWord,
	TokenJSOperator,

	TokenJSTreeDotOperator,
	//TokenJSAssign,
	TokenJSUnaryOperator,
	//TokenJSBinaryOperator,

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

var jsPhraseBreakAfter = []int{
	TokenJSBracket,
	TokenJSBracketSquare,
	TokenJSBlock,
}

var jsPhraseJoin = []int{
	TokenJSFunctionLambda,
}

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

var jsPhraseBreakKeyWords string = `
,abstract,arguments,await,boolean,
,break,byte,case,catch,
,char,class,const,continue,
,debugger,default,delete,do,
,double,else,enum,eval,
,export,extends,false,final,finally,float,for,function,
,goto,if,implements,import,
,in,instanceof,int,interface,
,let,long,native,new,
,null,package,private,protected,
,public,return,short,static,
,super,switch,synchronized,
,throw,throws,transient,true,
,try,typeof,var,void,
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

	token := meaning.getNextMeaningToken(process.ParentTokens, process.Iter)

	if token != nil {

		if token.Type != TokenJSPhrase {

			meaning.processChild(process.ParentTokens, token)

			if token.Children.Length() == 1 {
				//get the first if phrase have only one blocking typed child
				if firstToken := token.Children.GetTokenAt(0); firstToken.Type == TokenJSPhrase {

					token.Children = firstToken.Children
				}
			}
		} else if token.Children.Length() == 0 {

			return meaning.getNextMeaningToken(process.ParentTokens, process.Iter) //remove empty phrase
		}
	}

	return token
}

func (meaning *JSPhraseMeaning) getNextMeaningToken(ancesstorTokens []int, iter *gotokenize.Iterator) *gotokenize.Token {

	meaning.continuePassPhraseBreak(iter)
	if token := iter.Read(); token != nil {

		if gotokenize.IsContainToken(jsPhraseJoin, token.Type) {

			meaning.continueJoin(ancesstorTokens, iter, token)

		} else if token.Type != TokenJSPhraseBreak &&
			gotokenize.IsContainToken(jsPhraseAllow, token.Type) &&
			!gotokenize.IsContainToken(jsPhraseBreakAfter, token.Type) &&
			!(token.Type == TokenJSKeyWord && strings.Index(jsPhraseBreakKeyWords, token.Content) > 0) {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
			meaning.processChild(ancesstorTokens, token)
			tmpToken.Children.AddToken(*token)
			meaning.continuePhrase(ancesstorTokens, iter, tmpToken)
			return tmpToken

		} else if token.Type == TokenJSPhraseBreak {

			return meaning.getNextMeaningToken(ancesstorTokens, iter)
		}
		return token
	}

	return nil
}

func (meaning *JSPhraseMeaning) continuePassPhraseBreak(iter *gotokenize.Iterator) {
	for {
		token := iter.Get()
		if token == nil || token.Type != TokenJSPhraseBreak {
			break
		}
		_ = iter.Read()
	}
}

func (meaning *JSPhraseMeaning) continuePhrase(ancesstorTokens []int, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	for {
		token := iter.Get()

		if token == nil || token.Type == TokenJSPhraseBreak || !gotokenize.IsContainToken(jsPhraseAllow, token.Type) {

			break
		}
		_ = iter.Read()
		if gotokenize.IsContainToken(jsPhraseJoin, token.Type) {

			meaning.continueJoin(ancesstorTokens, iter, token)
		}
		meaning.processChild(ancesstorTokens, token)

		currentToken.Children.AddToken(*token)

		if gotokenize.IsContainToken(jsPhraseBreakAfter, token.Type) || (token.Type == TokenJSKeyWord && strings.Index(jsPhraseBreakKeyWords, token.Content) > 0) {

			break
		}
	}
	if currentToken.Children.Length() == 1 {

		firstToken := currentToken.Children.GetTokenAt(0)
		currentToken.Type = firstToken.Type
		currentToken.Content = firstToken.Content
		currentToken.Children = firstToken.Children
	}
}
func (meaning *JSPhraseMeaning) continueJoin(ancesstorTokens []int, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	meaning.continuePassPhraseBreak(iter)
	phraseToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, "")
	meaning.continuePhrase(ancesstorTokens, iter, phraseToken)

	if phraseToken.Type != TokenJSPhrase {
		currentToken.Children.AddToken(*phraseToken)
	} else {
		childIter := phraseToken.Children.Iterator()
		for {
			childToken := childIter.Read()
			if childToken == nil {
				break
			}
			currentToken.Children.AddToken(*childToken)
		}
	}
}
func (meaning *JSPhraseMeaning) processChild(ancestorTokens []int, parentToken *gotokenize.Token) {

	if !gotokenize.IsContainToken(JSLevel2GlobalNested, parentToken.Type) {
		return
	}
	proc := gotokenize.NewMeaningProcessFromStream(append(ancestorTokens, parentToken.Type), &parentToken.Children)

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
func (meaning *JSPhraseMeaning) GetName() string {

	return "JSPhraseMeaning"
}
