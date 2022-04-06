package gotokenize

type PatternMeaning struct {
	*AbstractMeaning
	Patterns       []Pattern
	IgnoreTokens   []int
	TokenCanNested []int
}

func CreatePatternMeaning(parent IMeaning, patterns []Pattern, ignoreTokens []int, tokenCanNested []int) *PatternMeaning {

	pattern := &PatternMeaning{

		AbstractMeaning: NewAbtractMeaning(parent),
		Patterns:        patterns,
		IgnoreTokens:    ignoreTokens,
		TokenCanNested:  tokenCanNested,
	}

	return pattern
}

/*
func (meaning *PatternMeaning) Prepare(stream *TokenStream) {
	//fmt.Println("pattern prepare")
	meaning.IMeaning.Prepare(stream)
	tmpStream := CreateStream()
	token := meaning.IMeaning.Next()
	for {
		if token == nil {
			break
		}
		tmpStream.AddToken(*token)
		token = meaning.IMeaning.Next()
	}
	meaning.IMeaning.SetStream(tmpStream)
}
*/

func (meaning *PatternMeaning) Next(process *MeaningProcess) *Token {

	iter := process.Iter

	for {

		if iter.EOS() {
			break
		}
		beforeOffset := iter.Offset

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

		} else { // case no mark found

			if normalToken := iter.Read(); normalToken != nil {

				if normalToken.Children.Length() > 0 && IsContainToken(meaning.TokenCanNested, normalToken.Type) {

					childMeaning := CreatePatternMeaning(meaning.AbstractMeaning.Clone(), meaning.Patterns, meaning.IgnoreTokens, meaning.TokenCanNested)
					childProcess := &MeaningProcess{}

					childProcess.SetStream(&normalToken.Children)

					childMeaning.Prepare(childProcess)

					subStream := CreateStream(meaning.GetMeaningLevel())

					for {

						nestedToken := childMeaning.Next(childProcess)

						if nestedToken == nil {
							break
						}
						subStream.AddToken(*nestedToken)
					}
					normalToken.Children = subStream
				}
				return normalToken
			}
		}
		if iter.Offset == beforeOffset {
			break
		}
	}
	return nil
}
func (meaning *PatternMeaning) GetName() string {

	return "PatternMeaning"
}
