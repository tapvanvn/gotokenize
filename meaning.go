package gotokenize

import "fmt"

//Meaning inteface for language meaning process
type IMeaning interface {
	GetName() string
	Prepare(process *MeaningProcess)
	Next(process *MeaningProcess) *Token
	GetSource() IMeaning
	SetSource(IMeaning)
	GetMeaningLevel() int
	Clone() IMeaning
	Propagate(func(meaning IMeaning))
}

func NewAbtractMeaning(base IMeaning) *AbstractMeaning {

	return &AbstractMeaning{

		BaseMeaning: base,
	}
}

//AbstractMeaning this struct actions like a two-side middle layer between the owner struct and the real base struct
//incase BaseMeaning is nil, the struct actions like base struct for owner struct
//incase BaseMeaning is not nil, the struct actions like the wrapper for the base meaning struct.
type AbstractMeaning struct {
	BaseMeaning IMeaning
}

func (meaning *AbstractMeaning) Prepare(process *MeaningProcess) {

	if meaning.BaseMeaning != nil {

		if process.Stream.MeaningLevel < meaning.GetMeaningLevel() {

			meaning.BaseMeaning.Prepare(process)
			//fmt.Printf("begin do prepare %s numToken:%d\n", meaning.BaseMeaning.GetName(), process.Stream.Length())
			tmpStream := CreateStream(meaning.GetMeaningLevel())

			for {
				token := meaning.BaseMeaning.Next(process)
				if token == nil {
					break
				}
				tmpStream.AddToken(*token)
			}
			//fmt.Printf("after do prepare %s numToken:%d\n", meaning.BaseMeaning.GetName(), tmpStream.Length())
			process.SetStream(process.ParentTokens, &tmpStream)

		}
	}
}

func (meaning *AbstractMeaning) Next(process *MeaningProcess) *Token {

	if meaning.BaseMeaning != nil {

		return meaning.BaseMeaning.Next(process)
	}
	return process.Iter.Read()
}
func (meaning *AbstractMeaning) GetSource() IMeaning {

	return meaning.BaseMeaning
}

func (meaning *AbstractMeaning) GetMeaningLevel() int {

	if meaning.BaseMeaning != nil {

		return meaning.BaseMeaning.GetMeaningLevel() + 1
	}
	return 0
}

func (meaning *AbstractMeaning) SetSource(baseMeaning IMeaning) {

	if baseMeaning != nil {
		if test, ok := baseMeaning.(*AbstractMeaning); ok {
			meaning.BaseMeaning = test.BaseMeaning
		} else {
			meaning.BaseMeaning = baseMeaning
		}
		return
	}
	meaning.BaseMeaning = nil
}

func (meaning *AbstractMeaning) Clone() IMeaning {

	if meaning.BaseMeaning != nil {

		return meaning.BaseMeaning.Clone()
	}
	return nil
}

func (meaning *AbstractMeaning) Propagate(fn func(meaning IMeaning)) {

	fn(meaning)

	if meaning.BaseMeaning != nil {

		meaning.BaseMeaning.Propagate(fn)

	}
}
func (meaning *AbstractMeaning) GetName() string {

	if meaning.BaseMeaning != nil {

		return meaning.BaseMeaning.GetName()
	}
	return "AbstractMeaning"
}

func DebugMeaning(meaning IMeaning) {

	if meaning != nil {

		p := func(mean IMeaning) {

			fmt.Printf("level:%d\t%s\n", mean.GetMeaningLevel(), mean.GetName())
		}
		meaning.Propagate(p)
	}
}
