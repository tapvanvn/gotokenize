package gotokenize

import (
	"fmt"
	"strconv"
	"strings"
)

//TokenStreamIterator struct use to access token stream
type TokenStreamIterator struct {
	Stream *TokenStream
	Offset int
	Level  int
}

//DebugMark debug mark
func (iterator *TokenStreamIterator) DebugMark(level int, mark *Mark, ignores *[]int, fnName func(int) string) {

	length := mark.End - mark.Begin

	iter := 0

	for {
		if length <= 0 || iterator.EOS() {
			break
		}

		token := iterator.GetAt(mark.Begin + iter)
		fmt.Printf("%s", ColorOffset(mark.Begin+iter))
		if token != nil {

			for i := 0; i <= level; i++ {

				if i == 0 {

					fmt.Printf("|%s ", ColorType(token.Type))

				} else {

					fmt.Print("| ")
				}
			}

			if !isIgnoreInMark(mark.Begin+iter, ignores) {

				trimContent := strings.Trim(token.Content, " \n\r")

				if len(trimContent) > 0 {

					fmt.Printf("%s", ColorContent(token.Content))

				} else {

					fmt.Print("")
				}

				fmt.Printf("-%s\n", ColorName(fnName(token.Type)))

			} else {

				fmt.Printf("%s", ColorIgnore())
			}

		} else {

			fmt.Printf("%s", "nil")
		}

		fmt.Println("")

		length--

		iter++
	}
}

//Get read token but not move pointer
func (iterator *TokenStreamIterator) Get() *Token {

	if iterator.Offset <= len(iterator.Stream.Tokens)-1 {

		off := iterator.Offset

		return &iterator.Stream.Tokens[off]
	}
	return nil
}

//GetBy get token at (offset + iterator) position
func (iterator *TokenStreamIterator) GetBy(iter int) *Token {

	if iterator.Offset+iter <= len(iterator.Stream.Tokens)-1 {

		off := iterator.Offset + iter

		return &iterator.Stream.Tokens[off]
	}
	return nil
}

//GetAt get token at offset
func (iterator *TokenStreamIterator) GetAt(offset int) *Token {

	if offset <= len(iterator.Stream.Tokens)-1 {

		return &iterator.Stream.Tokens[offset]
	}
	return nil
}

func IsContainToken(ignoreTokens []int, tokenType int) bool {
	for _, i := range ignoreTokens {
		if i == tokenType {
			return true
		}
	}
	return false
}

//FindPattern search pattern
func (iterator *TokenStreamIterator) FindPattern(patterns []Pattern, stopWhenFound bool, phraseBreak int, ignoreTokens []int) []Mark {

	marks := []Mark{}

	log := &Log{}

	defer log.Print()

	for _, pattern := range patterns {

		iter := 0

		iterToken := 0

		traceIterToken := -1

		patternTokenNum := len(pattern.Struct)

		ignores := []int{}

		children := []*Mark{}

		var patternToken PatternToken

		var childMark *Mark = nil

		for {
			if iterToken >= patternTokenNum {

				mark := Mark{Type: pattern.Type, Begin: iterator.Offset, End: iterator.Offset + iter, Ignores: ignores, Children: children}

				marks = append(marks, mark)

				log.Append(fmt.Sprintf("=>[%s] \n", ColorSuccess()))

				if stopWhenFound {

					return marks
				}
				break
			}
			if iterToken > traceIterToken {

				traceIterToken = iterToken

				patternToken = pattern.Struct[iterToken]

				childMark = &Mark{
					Type:             patternToken.Type,
					CanNested:        patternToken.CanNested,
					IsIgnoreInResult: patternToken.IsIgnoreInResult,
					IsTokenStream:    patternToken.IsPhraseUntil,
				}
				if patternToken.ExportType > 0 {

					childMark.Type = patternToken.ExportType
				}

				children = append(children, childMark)

				childMark.Begin = iterator.Offset + iter

				log.Append(fmt.Sprintf("\n\t[%s %s] %s %t", ColorType(patternToken.Type), ColorName(strconv.Itoa(patternToken.Type)), ColorContent(patternToken.Content), patternToken.IsPhraseUntil))
			}
			var match bool = true

			var moveIter int = 0

			nextToken := iterator.GetBy(iter)

			if nextToken == nil {
				break
			}

			if nextToken.Type == phraseBreak || IsContainToken(ignoreTokens, nextToken.Type) {

				if pattern.IsRemoveGlobalIgnore || patternToken.IsIgnoreInResult {

					ignores = append(ignores, iterator.Offset+iter)
				}
				iter++

				log.Append(fmt.Sprintf("\n"))

				continue
			}
			if patternToken.Content != "" {

				var currToken = iterator.GetBy(iter)

				if currToken == nil || currToken.Content != patternToken.Content {

					match = false

					log.Append(fmt.Sprintf("=>[%s %s %s]", ColorFail(), ColorType(currToken.Type), ColorContent(currToken.Content)))
				}
				if patternToken.IsIgnoreInResult {

					ignores = append(ignores, iterator.Offset+iter+moveIter)
				}

				childMark.Begin = iterator.Offset + iter

				moveIter = 1

			} else if patternToken.Type > -1 {

				var currToken = iterator.GetBy(iter)

				if currToken == nil || (currToken.Type != phraseBreak && currToken.Type != patternToken.Type) {

					match = false

					log.Append(fmt.Sprintf("=>[%s %s %s]", ColorFail(), ColorType(currToken.Type), ColorContent(currToken.Content)))
				}

				if patternToken.IsIgnoreInResult {

					ignores = append(ignores, iterator.Offset+iter+moveIter)
				}

				if currToken.Type == patternToken.Type {

					childMark.Begin = iterator.Offset + iter
				}

				moveIter = 1

			} else if patternToken.IsPhraseUntil {

				isWordFound := false

				for {
					var currToken = iterator.GetBy(iter + moveIter)

					if currToken == nil {

						match = false

						//log.Append(fmt.Sprintf("=>[%s]", ColorFail()))

						break
					}
					if IsContainToken(ignoreTokens, currToken.Type) {

						if pattern.IsRemoveGlobalIgnore || patternToken.IsIgnoreInResult {

							ignores = append(ignores, iterator.Offset+iter+moveIter)
						}
						moveIter++

						continue
					}
					if currToken.Type == phraseBreak && isWordFound {

						if pattern.IsRemoveGlobalIgnore || patternToken.IsIgnoreInResult {

							ignores = append(ignores, iterator.Offset+iter+moveIter)
						}
						moveIter++

						break

					} else if currToken.Type != phraseBreak && len(currToken.Content) > 0 {

						isWordFound = true
					}

					if patternToken.IsIgnoreInResult {

						ignores = append(ignores, iterator.Offset+iter+moveIter)
					}

					moveIter++
				}
			}
			if !match {

				break
			}

			iter += moveIter

			childMark.End = iterator.Offset + iter

			iterToken++
			log.Append(fmt.Sprintf("\n"))
		}
	}
	return marks
}

//GetMaskedToken get token from mask
func (iterator *TokenStreamIterator) GetMaskedToken(mark *Mark, ignores *[]int) *Token {

	if mark.IsTokenStream {

		token := Token{Type: mark.Type}

		len := mark.End - mark.Begin

		iter := 0

		for {
			if len <= 0 || iterator.EOS() {

				break
			}
			nextToken := iterator.GetAt(mark.Begin + iter)

			if !isIgnoreInMark(mark.Begin+iter, ignores) {

				token.Children.AddToken(*nextToken)

			}
			len--

			iter++
		}

		return &token

	} else {

		len := mark.End - mark.Begin

		iter := 0

		for {
			if len <= 0 || iterator.EOS() {

				break
			}
			nextToken := iterator.GetAt(mark.Begin + iter)

			if !isIgnoreInMark(mark.Begin+iter, ignores) {

				return nextToken

			}
			len--

			iter++
		}
	}
	return nil
}

//Read read token
func (iterator *TokenStreamIterator) Read() *Token {

	if iterator.Offset <= len(iterator.Stream.Tokens)-1 {

		off := iterator.Offset

		iterator.Offset++

		return &iterator.Stream.Tokens[off]
	}
	return nil
}

//Reset reset to begin
func (iterator *TokenStreamIterator) Reset() {

	iterator.Offset = 0
}

//EOS is end of stream
func (iterator *TokenStreamIterator) EOS() bool {

	return iterator.Offset >= len(iterator.Stream.Tokens)
}

//FirstType get first token of type
func (iterator *TokenStreamIterator) FirstType(tokenType int) *Token {
	iterator.Reset()
	return iterator.NextType(tokenType)
}

//NextType read from current position to next match of token type
func (iterator *TokenStreamIterator) NextType(tokenType int) *Token {

	for {
		if iterator.EOS() {

			break
		}
		token := iterator.Read()

		if token.Type == tokenType {

			return token
		}
	}
	return nil
}

//Seek move to offset
func (iterator *TokenStreamIterator) Seek(offset int) {
	iterator.Offset = offset
}
