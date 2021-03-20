package gotokenize

type PatternMeaning struct {
	Meaning
	Patterns       []Pattern
	IgnoreTokens   []int
	TokenCanNested []int
}

func CreatePatternMeaning(stream *TokenStream, patterns []Pattern, ignoreTokens []int, tokenCanNested []int) PatternMeaning {
	pattern := &PatternMeaning{
		Patterns:       patterns,
		IgnoreTokens:   ignoreTokens,
		TokenCanNested: tokenCanNested,
	}
	pattern.Prepare(stream)
	return *pattern
}

func (meaning *PatternMeaning) Prepare(stream *TokenStream) {

	meaning.Meaning.Prepare(stream)
}

func (meaning *PatternMeaning) Next() *Token {

	for {

		if meaning.Meaning.Iter.EOS() {
			break
		}

		marks := meaning.Meaning.Iter.FindPattern(meaning.Patterns, true, -1, meaning.IgnoreTokens)

		if len(marks) > 0 {

			mark := marks[0]

			curToken := Token{
				Type: mark.Type,
			}

			for _, childMark := range mark.Children {

				if childMark.IsIgnoreInResult {

					continue
				}

				childToken := meaning.Meaning.Iter.GetMaskedToken(childMark, &mark.Ignores)

				if childToken != nil && childMark.CanNested {

					childMeaning := CreatePatternMeaning(&childToken.Children, meaning.Patterns, mark.Ignores, meaning.TokenCanNested)
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

			meaning.Meaning.Iter.Seek(mark.End)

			return &curToken

		} else { // case no mark found

			if normalToken := meaning.Meaning.Iter.Read(); normalToken != nil {

				if normalToken.Children.Length() > 0 && IsContainToken(meaning.TokenCanNested, normalToken.Type) {

					childMeaning := CreatePatternMeaning(&normalToken.Children, meaning.Patterns, meaning.IgnoreTokens, meaning.TokenCanNested)

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
