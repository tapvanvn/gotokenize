package gotokenize

import (
	"fmt"
)

//TokenStream token stream
type TokenStream struct {
	Tokens       []Token
	MeaningLevel int //the level of meaning the stream got on a meaning process
}

func CreateStream(meaningLevel int) TokenStream {

	return TokenStream{

		Tokens:       []Token{},
		MeaningLevel: meaningLevel,
	}
}

//Iterator make iterator of stream
func (stream *TokenStream) Iterator() *Iterator {

	return &Iterator{Stream: stream, Offset: 0, Level: 0}
}

//Tokenize tokenize a string
func (stream *TokenStream) Tokenize(content string) {

	runes := []rune(content)

	for _, rune := range runes {

		token := Token{Content: fmt.Sprintf("%c", rune)}

		stream.AddToken(token)
	}
}

//AddToken add token to stream
func (stream *TokenStream) AddToken(token Token) {

	stream.Tokens = append(stream.Tokens, token)
}

//AddTokenFromString split string to character and add each character as a token with type is providing type.
func (stream *TokenStream) AddTokenFromString(tokenType int, str string) {

	for r := range []rune(str) {

		stream.AddToken(Token{Type: tokenType, Content: string(rune(r))})
	}
}

//AddTokenByContent add token
func (stream *TokenStream) AddTokenByContent(content []rune, tokenType int) {

	stream.Tokens = append(stream.Tokens, Token{Content: string(content), Type: tokenType})
}

//Debug print debug tree
func (stream *TokenStream) Debug(level int, fnName func(int) string, options *DebugOption) {

	for _, token := range stream.Tokens {

		token.Debug(level, fnName, options)
	}
}

//GetTokenAt get token at offset
func (stream *TokenStream) GetTokenAt(offset int) *Token {

	if offset <= len(stream.Tokens)-1 {

		return &stream.Tokens[offset]
	}
	return nil
}

func isIgnoreInMark(iterator int, ignores *[]int) bool {

	for _, i := range *ignores {

		if i == iterator {

			return true
		}
	}
	return false
}

//Length get len of stream
func (stream *TokenStream) Length() int {

	return len(stream.Tokens)
}

//ConcatStringContent concat content of tokens
func (stream *TokenStream) ConcatStringContent() string {

	var iterator = stream.Iterator()

	iterator.Reset()

	content := ""

	for {
		if iterator.EOS() {

			break
		}
		token := iterator.Read()

		content += string(token.Content)
	}

	return content
}

//ToArray get array of tokens
func (stream *TokenStream) ToArray() []Token {

	var rs []Token

	var iterator = stream.Iterator()

	iterator.Reset()

	for {
		if iterator.EOS() {

			break
		}
		token := iterator.Read()

		rs = append(rs, *token)
	}

	return rs
}

//ReadFirstTokenType read first token of type
func (stream *TokenStream) ReadFirstTokenType(tokenType int) *Token {

	var iterator = stream.Iterator()

	iterator.Reset()

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
