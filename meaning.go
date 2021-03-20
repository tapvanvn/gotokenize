package gotokenize

//Meaning inteface for language meaning process
type IMeaning interface {
	Prepare(stream *TokenStream)
	Next() *Token
}

type Meaning struct {
	source IMeaning
	Stream TokenStream
	Iter   TokenStreamIterator
}

func CreateMeaning(source IMeaning) Meaning {

	meaning := &Meaning{
		Stream: CreateStream(),
	}

	if source != nil {

		meaning.source = source

		for {
			token := meaning.source.Next()
			if token == nil {
				break
			}
			meaning.Stream.AddToken(*token)
		}

	}
	meaning.Iter = meaning.Stream.Iterator()

	return *meaning
}

func (meaning *Meaning) Next() *Token {
	return meaning.Iter.Read()
}

func (meaning *Meaning) Prepare(stream *TokenStream) {

	if meaning.source != nil {
		meaning.source.Prepare(stream)
		meaning.Stream = CreateStream()
		for {
			token := meaning.source.Next()
			if token == nil {
				break
			}
			meaning.Stream.AddToken(*token)
		}
		meaning.Iter = meaning.Stream.Iterator()
	} else {
		iter := stream.Iterator()
		meaning.Stream = CreateStream()
		for {
			if iter.EOS() {
				break
			}
			token := iter.Read()
			meaning.Stream.AddToken(*token)
		}
		meaning.Iter = meaning.Stream.Iterator()
	}
}
