package main

import (
	"flag"
	"os"

	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/js"
)

var file = flag.String("f", "complex.js", "input file")
var printDebug = flag.Bool("d", false, "print debug")
var printMeaningDebug = flag.Bool("m", false, "print meaning debug")

func main() {
	flag.Parse()

	data, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(string(data))
	var meaning gotokenize.IMeaning = js.NewDefaultJSInstructionMeaning()
	var stringifier gotokenize.IStringifier = js.NewDefaultInstructionStrigifier()

	if flag.NArg() > 0 {
		if flag.Arg(0) == "meaning" {
			meaning = js.NewDefaultJSMeaning()

		} else if flag.Arg(0) == "pattern" {
			meaning = js.NewDefaultJSPatternMeaning()
		} else if flag.Arg(0) == "phrase" {
			meaning = js.NewDefaultJSPhraseMeaning()
			stringifier = js.NewDefaultPhraseStringifier()
		} else if flag.Arg(0) == "operator" {
			meaning = js.NewDefaultJSOperatorMeaning()
			stringifier = js.NewDefaultOperatorStringifier()
		} else if flag.Arg(0) == "raw" {
			meaning = js.NewDefaultJSRawMeaning()
			stringifier = js.NewDefaultRawStringifier()
		}
	}

	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.EmptyParentTokens, &stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		if *printDebug {
			token.Debug(0, js.JSTokenName, js.JSDebugOptions)
		}
		stringifier.PutToken(token)
		token = meaning.Next(proc)
	}
	if *printMeaningDebug {
		gotokenize.DebugMeaning(meaning)
	}

	if err := os.WriteFile("out.js", []byte(stringifier.GetContent()), 0644); err != nil {
		panic(err)
	}

}
