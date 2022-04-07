package gotokenize

type PatternMeaning struct {
	*AbstractMeaning
	Patterns       []Pattern
	IgnoreTokens   []int
	TokenCanNested []int
}

func NewPatternMeaning(parent IMeaning, patterns []Pattern, ignoreTokens []int, tokenCanNested []int) *PatternMeaning {

	pattern := &PatternMeaning{

		AbstractMeaning: NewAbtractMeaning(parent),
		Patterns:        patterns,
		IgnoreTokens:    ignoreTokens,
		TokenCanNested:  tokenCanNested,
	}

	return pattern
}

func (meaning *PatternMeaning) Next(process *MeaningProcess) *Token {

	token := meaning.getNextMeaningToken(process.Iter)
	if token != nil {

		if IsContainToken(meaning.TokenCanNested, token.Type) {

			childProcess := NewMeaningProcessFromStream(&token.Children)
			//TODO: fixbug Prepare
			//childMeaning.Prepare(childProcess)

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
	}
	return token
}

func (meaning *PatternMeaning) getNextMeaningToken(iter *Iterator) *Token {
	for {

		if iter.EOS() {
			break
		}

		marks := iter.FindPattern(meaning.Patterns, true, meaning.IgnoreTokens)

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

					childProcess := NewMeaningProcessFromStream(&childToken.Children)

					meaning.Prepare(childProcess)

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
