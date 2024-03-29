package gotokenize

import "fmt"

//PatternToken pattern
type PatternToken struct {
	Type             int
	Content          string
	IsPhraseUntil    bool
	IsIgnoreInResult bool
	CanNested        bool
	ExportType       int
	IsAny            bool           //if true match any next token
	NumNextMoreAny   int            // if > 0 && isAny = true, check if there is next NumNextMoreAny + 1 tokens
	Nested           []PatternToken //check if contain these tokens
}

//Pattern define a pattern is a array of token type
type Pattern struct {
	Type                 int
	Struct               []PatternToken
	IsRemoveGlobalIgnore bool
	LivingContext        []int //only available when the pattern token is found inside one of the livingContext token array.
	RequireLastToken     int
}

//Mark define a result of finding parttern process
type Mark struct {
	Type             int
	Begin            int
	End              int
	Ignores          []int //iterator that should be ignore
	CanNested        bool
	Children         []*Mark
	IsIgnoreInResult bool
	IsTokenStream    bool
}

//Debug print debug
func (mark *Mark) Debug(level int, fnName func(int) string) {

	for i := 0; i <= level; i++ {

		if i == 0 {

			fmt.Printf("|%s ", ColorType(mark.Type))

		} else {

			fmt.Print("| ")
		}
	}

	fmt.Printf("-%s", ColorName(fnName(mark.Type)))
	fmt.Printf(" ignore:%t nested:%t stream:%t", mark.IsIgnoreInResult, mark.CanNested, mark.IsTokenStream)
	fmt.Printf(" begin:%d", mark.Begin)
	fmt.Printf(" end:%d", mark.End)
	fmt.Printf(" child:%d\n", len(mark.Children))

	for _, m := range mark.Children {

		m.Debug(level+1, fnName)
	}
}
