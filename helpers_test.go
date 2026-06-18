package integration

import (
	"context"
	"io"
	"strconv"
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// upper is a stateless command that upper-cases each line. Stateless commands
// are safe to reuse across goroutines.
func upper() gloo.Command[[]byte, []byte] {
	return patterns.Map(func(line []byte) ([]byte, error) {
		return []byte(strings.ToUpper(string(line))), nil
	})
}

// numberLines is a stateful command: a fresh "n:line" counter per Execute. The
// per-Execute state is what makes a single value reusable across goroutines.
func numberLines() gloo.Command[[]byte, []byte] {
	return patterns.StatefulMap(func() func([]byte) ([]byte, error) {
		n := 0
		return func(line []byte) ([]byte, error) {
			n++
			return []byte(strconv.Itoa(n) + ":" + string(line)), nil
		}
	})
}

// runOn executes cmd over input (one item per line) and returns the output
// lines. It is safe to call from multiple goroutines: each call builds its own
// source and stream.
func runOn(cmd gloo.Command[[]byte, []byte], input string) ([]string, error) {
	ctx := context.Background()
	source := gloo.ByteReaderSource([]io.Reader{strings.NewReader(input)})
	items, err := gloo.Collect(ctx, cmd.Execute(ctx, source.Stream(ctx)))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(items))
	for i, b := range items {
		out[i] = string(b)
	}
	return out, nil
}
