package gotokenize

import "fmt"

var __root_meaning = RootMeaning{}

func CreateRootMeaning() *RootMeaning {

	return &__root_meaning
}

type RootMeaning struct{}

func (meaning *RootMeaning) GetMeaningLevel() int {
	return 0
}
func (meaning *RootMeaning) Clone() IMeaning {
	return &RootMeaning{}
}
func (meaning *RootMeaning) Next(process *MeaningProcess) *Token {
	return process.Iter.Read()
}

func (meaning *RootMeaning) Prepare(proc *MeaningProcess) {}
func (meaning *RootMeaning) GetSource() IMeaning {
	return nil
}
func (meaning *RootMeaning) SetSource(baseMeaning IMeaning) {}

func (meaning *RootMeaning) Propagate(fn func(meaning IMeaning)) {
	fmt.Println("root propagate")
	fn(meaning)
}
func (meaning *RootMeaning) GetName() string {
	return "RootMeaning"
}
