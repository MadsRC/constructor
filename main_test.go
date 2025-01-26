package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTitleFunc(t *testing.T) {
	input := "fiskePinde"
	expected := "FiskePinde"
	actual := titleFunc(input)
	require.Equal(t, expected, actual)

}

func BenchmarkTitleFunc(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		titleFunc("fiskePinde")
	}
}
