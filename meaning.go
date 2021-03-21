package gotokenize

//Meaning inteface for language meaning process
type IMeaning interface {
	Prepare(stream *TokenStream)
	Next() *Token
	GetIter() *Iterator
	GetStream() *TokenStream
	SetStream(stream TokenStream)
	Clone() IMeaning
	GetSource() IMeaning
}

type Meaning struct {
	source IMeaning
	Stream TokenStream
	Iter   Iterator
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

func (meaning *Meaning) GetIter() *Iterator {
	return &meaning.Iter
}

func (meaning *Meaning) GetStream() *TokenStream {
	return &meaning.Stream
}

func (meaning *Meaning) SetStream(stream TokenStream) {
	meaning.Stream = stream
	meaning.Iter = meaning.Stream.Iterator()
}

func (meaning *Meaning) Clone() IMeaning {
	clone := &Meaning{
		Stream: CreateStream(),
	}
	if meaning.source != nil {
		clone.source = meaning.source.Clone()
	}
	clone.Iter = clone.Stream.Iterator()
	return clone
}

func (meaning *Meaning) GetSource() IMeaning {
	return meaning.source
}
