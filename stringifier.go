package gotokenize

type IStringifier interface {
	PutToken(token *Token)
	GetContent() string
}

func NewAbtractStringifier(base IStringifier) *AbstractStringifier {

	return &AbstractStringifier{

		BaseStringifier: base,
	}
}

//AbstractStrigifier this struct actions like a two-side middle layer between the owner struct and the real base struct
//incase BaseMeaning is nil, the struct actions like base struct for owner struct
//incase BaseMeaning is not nil, the struct actions like the wrapper for the base meaning struct.
type AbstractStringifier struct {
	BaseStringifier IStringifier
}

func (meaning *AbstractStringifier) AddToken(token *Token) {

	if meaning.BaseStringifier != nil {

		meaning.BaseStringifier.PutToken(token)
	}
}
func (meaning *AbstractStringifier) GetContent() string {
	if meaning.BaseStringifier != nil {
		return meaning.BaseStringifier.GetContent()
	}
	return ""
}
