package snowflake

import (
	"testing"
)

var worker = New(1001, 1002)

func Test_Generate(t *testing.T) {
	_, err := worker.Generate()
	if err != nil {
		t.Fatal("err")
	}
}

func Test_Generate_Repetition(t *testing.T) {
	data := make(map[int64]bool)
	dataOfRepetition := make([]int64, 0)
	for i := 0; i < 1000000; i++ {
		id, err := worker.Generate()
		if err != nil {
			t.Fatal("err")
		}
		if _, ok := data[id]; ok {
			dataOfRepetition = append(dataOfRepetition, id)
		}
	}
	if len(dataOfRepetition) > 0 {
		t.Fatal("err", len(dataOfRepetition))
	}
}

func Benchmark_Generate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := worker.Generate()
		if err != nil {
			continue
		}
	}
}
