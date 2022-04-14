package gotokenize

func NewMeaningProcessFromStream(parentTokens []int, stream *TokenStream) *MeaningProcess {
	proc := &MeaningProcess{}
	proc.SetStream(parentTokens, stream)
	return proc
}

var EmptyParentTokens = []int{}

type MeaningProcess struct {
	Stream  *TokenStream
	Iter    *Iterator
	Context MeaningContext
}

func (proc *MeaningProcess) GetPercent() float64 {
	if proc.Stream == nil || proc.Stream.Length() == 0 {
		return 0
	}
	return float64(proc.Iter.Offset) / float64(proc.Stream.Length()) * 100
}

func (proc *MeaningProcess) SetStream(ancestors []int, stream *TokenStream) {
	proc.Stream = stream
	proc.Iter = proc.Stream.Iterator()
	proc.Context.AncestorTokens = ancestors
	proc.Context.PreviousToken = TokenNoType
}
