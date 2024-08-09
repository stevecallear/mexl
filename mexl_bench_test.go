package mexl_test

import (
	"testing"

	"github.com/stevecallear/mexl"
)

type benchmark struct {
	expr string
	env  map[string]any
}

var (
	basic = benchmark{
		expr: `orders gt 2`,
		env: map[string]any{
			"orders": 10,
		},
	}

	extended = benchmark{
		expr: `lower(email) ew "@mail.com" or (scope in ["internal", "beta"] and orders gt 0)`,
		env: map[string]any{
			"email":  "Test@Company.com",
			"scope":  "internal",
			"orders": 10,
		},
	}
)

func BenchmarkEval_Basic(b *testing.B) {
	benchmarkEval(b, basic)
}

func BenchmarkEval_Extended(b *testing.B) {
	benchmarkEval(b, extended)
}

func benchmarkEval(b *testing.B, bm benchmark) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mexl.Eval(bm.expr, bm.env)
		if err != nil {
			b.Fatalf("got %v, expected nil", err)
		}
	}
	b.StopTimer()
}

func BenchmarkCompile_Basic(b *testing.B) {
	benchmarkCompile(b, basic)
}

func BenchmarkCompile_Extended(b *testing.B) {
	benchmarkCompile(b, extended)
}

func benchmarkCompile(b *testing.B, bm benchmark) {
	b.Helper()

	for i := 0; i < b.N; i++ {
		_, err := mexl.Compile(bm.expr)
		if err != nil {
			b.Fatalf("got %v, expected nil", err)
		}
	}
}

func BenchmarkRun_Basic(b *testing.B) {
	benchmarkRun(b, basic)
}

func BenchmarkRun_Extended(b *testing.B) {
	benchmarkRun(b, extended)
}

func benchmarkRun(b *testing.B, bm benchmark) {
	b.Helper()

	p, err := mexl.Compile(bm.expr)
	if err != nil {
		b.Fatalf("got %v, expected nil", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		out, err := mexl.Run(p, bm.env)
		if err != nil {
			b.Fatalf("got %v, expected nil", err)
		}
		if out != true {
			b.Fatalf("got %v, expected true", out)
		}
	}
}
