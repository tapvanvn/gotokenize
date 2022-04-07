package main

import (
	"os"

	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/js"
)

func main() {
	data, err := os.ReadFile("complex.js")
	if err != nil {
		panic(err)
	}
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(string(data))

	meaning := js.NewDefaultJSInstructionMeaning()

	proc := gotokenize.NewMeaningProcessFromStream(&stream)

	meaning.Prepare(proc)

	token := meaning.Next(proc)

	for {
		if token == nil {
			break
		}
		token.Debug(0, js.JSTokenName)
		token = meaning.Next(proc)
	}
	gotokenize.DebugMeaning(meaning)
}
