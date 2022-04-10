package main

import (
	"flag"
	"os"

	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/js"
)

var file = flag.String("f", "complex.js", "input file")

func main() {
	flag.Parse()

	data, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(string(data))
	var meaning gotokenize.IMeaning = js.NewDefaultJSInstructionMeaning()

	if flag.NArg() > 0 {
		if flag.Arg(0) == "meaning" {
			meaning = js.NewDefaultJSMeaning()
		} else if flag.Arg(0) == "pattern" {
			meaning = js.NewDefaultJSPatternMeaning()
		} else if flag.Arg(0) == "phrase" {
			meaning = js.NewDefaultJSPhraseMeaning()
		} else if flag.Arg(0) == "operator" {
			meaning = js.NewDefaultJSOperatorMeaning()
		} else if flag.Arg(0) == "raw" {
			meaning = js.NewDefaultJSRawMeaning()
		}
	}

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	stringifer := js.NewStringfier()
	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName, js.JSDebugOptions)
		stringifer.PutToken(token)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)

	os.WriteFile("out.js", []byte(stringifer.Content), 0644)
}
