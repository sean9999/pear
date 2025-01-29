package pear

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/DataDog/gostackparse"
)

func NicePanic(w io.Writer) {
	stack := debug.Stack()
	goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

	// frames := []*gostackparse.Frame{}

	// for _, frame := range goroutines[0].Stack {
	// 	if strings.HasPrefix(frame.File, runtime.GOROOT()) {
	// 		break
	// 	}
	// 	frames = append(frames, frame)
	// }

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err := enc.Encode(goroutines)
	if err != nil {
		fmt.Fprintf(w, "%s\n", stack)
	}
}
