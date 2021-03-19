package gotokenize

//Meaning inteface for language meaning process
type Meaning interface {
	GetNextMeaningToken() *Token
}
