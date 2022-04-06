package gotokenize

func NewMeaningProcessFromStream(stream *TokenStream) *MeaningProcess {
	proc := &MeaningProcess{}
	proc.SetStream(stream)
	return proc
}

type MeaningProcess struct {
	Stream *TokenStream
	Iter   *Iterator
}

func (proc *MeaningProcess) SetStream(stream *TokenStream) {
	proc.Stream = stream
	proc.Iter = proc.Stream.Iterator()
}
