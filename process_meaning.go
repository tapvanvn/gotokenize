package gotokenize

func NewMeaningProcessFromStream(parentTokens []int, stream *TokenStream) *MeaningProcess {
	proc := &MeaningProcess{}
	proc.SetStream(parentTokens, stream)
	return proc
}

var EmptyParentTokens = []int{}

type MeaningProcess struct {
	Stream          *TokenStream
	Iter            *Iterator
	ParentTokens    []int
	PassedTokenType int
}

func (proc *MeaningProcess) SetStream(parentTokens []int, stream *TokenStream) {
	proc.Stream = stream
	proc.Iter = proc.Stream.Iterator()
	proc.ParentTokens = parentTokens
	proc.PassedTokenType = TokenNoType
}
