package integration

import (
	"context"
	"io"
	"strings"
	"sync"
	"testing"

	gloo "github.com/gloo-foo/framework"
)

// TestParallelInstances runs separate command instances concurrently, each on
// its own input, and checks every result is correct and isolated.
func TestParallelInstances(t *testing.T) {
	cases := map[string]string{
		"alpha\nbeta\ngamma": "ALPHA\nBETA\nGAMMA",
		"one\ntwo\nthree":    "ONE\nTWO\nTHREE",
		"red\ngreen\nblue":   "RED\nGREEN\nBLUE",
	}
	var wg sync.WaitGroup
	for in, want := range cases {
		wg.Add(1)
		go func(in, want string) {
			defer wg.Done()
			got, err := runOn(upper(), in)
			if err != nil {
				t.Errorf("in=%q: %v", in, err)
				return
			}
			if strings.Join(got, "\n") != want {
				t.Errorf("in=%q got=%q want=%q", in, got, want)
			}
		}(in, want)
	}
	wg.Wait()
}

// TestReusableStatelessValue shares ONE stateless command value across many
// goroutines. With -race this proves stateless commands are reuse-safe.
func TestReusableStatelessValue(t *testing.T) {
	cmd := upper() // single shared value
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := runOn(cmd, "go\nrust")
			if err != nil {
				t.Errorf("execute: %v", err)
				return
			}
			if len(got) != 2 || got[0] != "GO" || got[1] != "RUST" {
				t.Errorf("got %q", got)
			}
		}()
	}
	wg.Wait()
}

// TestReusableStatefulValue shares ONE stateful command value across goroutines.
// Each Execute must get its own counter, so every goroutine sees 1,2,3 — no
// cross-goroutine state leak.
func TestReusableStatefulValue(t *testing.T) {
	cmd := numberLines() // single shared value, per-Execute state
	want := []string{"1:a", "2:b", "3:c"}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := runOn(cmd, "a\nb\nc")
			if err != nil {
				t.Errorf("execute: %v", err)
				return
			}
			if len(got) != len(want) {
				t.Errorf("got %q, want %q", got, want)
				return
			}
			for j := range want {
				if got[j] != want[j] {
					t.Errorf("got %q, want %q", got, want)
					return
				}
			}
		}()
	}
	wg.Wait()
}

// TestParallelPipelines composes and runs pipelines concurrently.
func TestParallelPipelines(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pipe := gloo.Pipe(upper(), numberLines())
			got, err := runOn(pipe, "x\ny")
			if err != nil {
				t.Errorf("execute: %v", err)
				return
			}
			if len(got) != 2 || got[0] != "1:X" || got[1] != "2:Y" {
				t.Errorf("got %q", got)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkSequential(b *testing.B) {
	input := strings.Repeat("test line\n", 1000)
	cmd := upper()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		src := gloo.ByteReaderSource([]io.Reader{strings.NewReader(input)})
		_, _ = gloo.Collect(ctx, cmd.Execute(ctx, src.Stream(ctx)))
	}
}

func BenchmarkParallel(b *testing.B) {
	input := strings.Repeat("test line\n", 1000)
	cmd := upper()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				src := gloo.ByteReaderSource([]io.Reader{strings.NewReader(input)})
				_, _ = gloo.Collect(ctx, cmd.Execute(ctx, src.Stream(ctx)))
			}()
		}
		wg.Wait()
	}
}
