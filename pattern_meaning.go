package gotokenize

type PatternMeaning struct {
	IMeaning
	Patterns       []Pattern
	IgnoreTokens   []int
	TokenCanNested []int
}

func CreatePatternMeaning(parent IMeaning, patterns []Pattern, ignoreTokens []int, tokenCanNested []int) PatternMeaning {

	pattern := &PatternMeaning{
		IMeaning:       parent,
		Patterns:       patterns,
		IgnoreTokens:   ignoreTokens,
		TokenCanNested: tokenCanNested,
	}

	return *pattern
}

func (meaning *PatternMeaning) Prepare(stream *TokenStream) {
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

func (meaning *PatternMeaning) Next() *Token {

	iter := meaning.GetIter()

	for {

		if iter.EOS() {
			break
		}

		marks := iter.FindPattern(meaning.Patterns, true, -1, meaning.IgnoreTokens)

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

					childMeaning := CreatePatternMeaning(meaning.IMeaning.Clone(), meaning.Patterns, mark.Ignores, meaning.TokenCanNested)
					childMeaning.Prepare(&childToken.Children)
					//childMeaning = PatternMeaning.init(stream: child_token!.children, pattern_groups: self.pattern_groups, is_ignore_func: self.is_ignore_func)

					subStream := CreateStream()

					for {

						nestedToken := childMeaning.Next()

						if nestedToken == nil {
							break
						}
						subStream.AddToken(*nestedToken)
					}

					childToken.Children = subStream
				}

				if childToken != nil {

					curToken.Children.AddToken(*childToken)

				} else {

					//meaning.Iterator.DebugMark(0, &patternMark, &patternMark.Ignores, js.TokenName)

					//meaning.Iterator.DebugMark(1, m, &patternMark.Ignores, js.TokenName)
				}
			}

			iter.Seek(mark.End)

			return &curToken

		} else { // case no mark found

			if normalToken := iter.Read(); normalToken != nil {

				if normalToken.Children.Length() > 0 && IsContainToken(meaning.TokenCanNested, normalToken.Type) {

					childMeaning := CreatePatternMeaning(meaning.IMeaning.Clone(), meaning.Patterns, meaning.IgnoreTokens, meaning.TokenCanNested)
					childMeaning.Prepare(&normalToken.Children)
					subStream := CreateStream()

					for {

						nestedToken := childMeaning.Next()

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
	}
	return nil
}
