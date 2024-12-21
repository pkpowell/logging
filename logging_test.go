package logging

import (
	"io"
	"log"
	"os"
	"testing"
)

// SETUP
// you need to call Run() once you've done what you need
func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func BenchmarkInfo(b *testing.B) {
	b.ReportAllocs()
	Init(bp(false), bp(false))

	for range 50 {
		Info("Hello", "World")
		Error("Hello", "World")
		Warn("Hello", "World")
	}
}

func BenchmarkInfof(b *testing.B) {
	b.ReportAllocs()
	Init(bp(false), bp(false))

	for range 50 {
		Info("Hello", "World")
		Error("Hello", "World")
		Warn("Hello", "World")
	}
}

func bp(b bool) *bool {
	return &b
}
