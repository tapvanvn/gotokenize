package gotokenize

import (
	"fmt"
	"math"
)

type PatternMeaning struct {
	*AbstractMeaning
	Patterns          []Pattern
	IgnoreTokens      []int
	TokenCanNested    []int
	PreventLoopTokens []int
}

type PatternMeaningDefine struct {
	Patterns          []Pattern
	IgnoreTokens      []int
	TokenCanNested    []int
	PreventLoopTokens []int
}

func NewPatternMeaning(parent IMeaning, patterns []Pattern, ignoreTokens []int, tokenCanNested []int, preventLoopTokens []int) *PatternMeaning {

	pattern := &PatternMeaning{

		AbstractMeaning: NewAbtractMeaning(parent),
		Patterns:        patterns,
		IgnoreTokens:    ignoreTokens,
		TokenCanNested:  tokenCanNested,
	}

	return pattern
}

func (meaning *PatternMeaning) Next(process *MeaningProcess) *Token {

	if len(process.Context.AncestorTokens) == 0 && process.Iter.Offset == 0 {
		fmt.Print("\033[s") //save cursor the position
	}
	token := meaning.getNextMeaningToken(process.Iter, &process.Context)
	if token != nil {

		if IsContainToken(meaning.TokenCanNested, token.Type) {

			childProcess := NewMeaningProcessFromStream(append(process.Context.AncestorTokens, token.Type), &token.Children)

			subStream := CreateStream(meaning.GetMeaningLevel())

			for {

				nestedToken := meaning.Next(childProcess)

				if nestedToken == nil {
					break
				}
				subStream.AddToken(*nestedToken)
			}
			token.Children = subStream
		}
		process.Context.PreviousToken = token.Type
		process.Context.PreviousTokenContent = token.Content

		fmt.Print("\033[u\033[K") //restore

		percent := int64(math.Floor(float64(process.Iter.Offset)/float64(process.Stream.Length())) * 100)
		fmt.Printf("percent: %d\n", percent)

	} else {
		process.Context.PreviousToken = TokenNoType
		process.Context.PreviousTokenContent = ""
	}
	if len(process.Context.AncestorTokens) == 0 {
		fmt.Print("\033[u\033[K") //restore
		fmt.Printf("%s percent: %f%%\n", meaning.GetName(), process.GetPercent())
		fmt.Print("\033[A")
	}
	return token
}

func (meaning *PatternMeaning) getNextMeaningToken(iter *Iterator, context *MeaningContext) *Token {
	for {

		if iter.EOS() {
			break
		}

		marks := iter.FindPattern(meaning.Patterns, true, meaning.IgnoreTokens, context)

		if len(marks) > 0 {

			mark := marks[0]

			curToken := Token{
				Type: mark.Type,
			}

			for _, childMark := range mark.Children {

				if childMark.IsIgnoreInResult {

					continue
				}

				childToken := iter.GetMaskedToken(childMark, &mark.Ignores)

				if childToken != nil && childMark.CanNested {

					childProcess := NewMeaningProcessFromStream(append(context.AncestorTokens, childToken.Type), &childToken.Children)

					subStream := CreateStream(meaning.GetMeaningLevel())

					for {

						nestedToken := meaning.Next(childProcess)

						if nestedToken == nil {
							break
						}
						subStream.AddToken(*nestedToken)
					}

					childToken.Children = subStream
				}

				if childToken != nil {

					curToken.Children.AddToken(*childToken)

				}
			}
			iter.Seek(mark.End)
			return &curToken
		}
		return iter.Read()
	}
	return nil
}

func (meaning *PatternMeaning) GetName() string {

	return "PatternMeaning"
}
