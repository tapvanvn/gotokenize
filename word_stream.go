package gotokenize

const (
	TokenWord = iota
	TokenSpace
)

//WordStream word
type WordStream struct {
	Source string
	Tokens []Token
	Offset int
	runes  []rune
}

func (stream *WordStream) GetCurrentCharacter() rune {
	return stream.runes[stream.Offset]
}

func (stream *WordStream) ReadCharacter() rune {
	if !stream.EOS() {
		var tmpOffset = stream.Offset
		stream.Offset++
		return stream.runes[tmpOffset]
	}
	return rune(0)
}

func (stream *WordStream) NextIndexOf(ch rune) int {
	tmpOffset := stream.Offset + 1
	for {
		if tmpOffset == stream.Length() {
			break
		}
		tmpRune := stream.runes[tmpOffset]
		if tmpRune == ch {
			return tmpOffset
		}
		tmpOffset++
	}
	return -1
}

func (stream *WordStream) GetToCharacter(toRune rune) string {
	var rs string = ""
	if !stream.EOS() {
		var pos int = stream.NextIndexOf(toRune)
		if pos >= stream.Offset {
			rs = string(stream.runes[stream.Offset:pos])
		}
	}
	return rs
}

func (stream *WordStream) ReadToCharacter(toRune rune) string {
	var rs string = ""
	if !stream.EOS() {
		var pos int = stream.NextIndexOf(toRune)
		if pos >= stream.Offset {
			rs = string(stream.runes[stream.Offset:pos])
			stream.Offset = pos
		}
	}
	return rs
}

func (stream *WordStream) ReadWhileCharacterIn(filter string) string {

	var runeFilter []rune = []rune(filter)
	var rsRune []rune
	for {
		if stream.EOS() {
			break
		}
		ch := stream.GetCurrentCharacter()
		var found bool = false
		for _, runeCh := range runeFilter {
			if runeCh == ch {
				found = true
				rsRune = append(rsRune, stream.ReadCharacter())
				break
			}
		}
		if !found {
			break
		}
	}
	return string(rsRune)
}

//Tokenize tokenize a string
func (stream *WordStream) Tokenize(content string) {
	stream.Source = content
	stream.Offset = 0
	stream.runes = []rune(content)
}

//AddToken add token to stream
func (stream *WordStream) AddToken(token Token) {

}

//AddTokenByConntent AddTokenByConntent
func (stream *WordStream) AddTokenByConntent(content []rune, tokenType int) {

}

//ReadToken read token
func (stream *WordStream) ReadToken() Token {
	return Token{}
}

//ResetToBegin reset to begin
func (stream *WordStream) ResetToBegin() {
	stream.Offset = 0
}

//EOS is end of stream
func (stream *WordStream) EOS() bool {
	return stream.Offset >= len(stream.runes)
}

//Length get len of stream
func (stream *WordStream) Length() int {
	return len(stream.runes)
}
