package engine_tests

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/CS80-Team/Goolean/internal/engine"
)

var queries = func() []string {
	file, err := os.Open("queries.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	var queries []string
	for scanner.Scan() {
		queries = append(queries, scanner.Text())
	}

	return queries
}()

func executeQueries(b *testing.B, engine *engine.Engine) {
	for i := 0; i < b.N; i++ {
		for _, query := range queries {
			engine.QueryString(query)
		}
	}
}

func BenchmarkEngineWithOrderedSlice(b *testing.B) {
	b.StartTimer()

	engine := MockEngine(datasetPath, *OrderedSliceIndexManager)

	executeQueries(b, engine)

	b.StopTimer()
}

func BenchmarkEngineWithSkipPointerList(b *testing.B) {
	b.StartTimer()

	engine := MockEngine(datasetPath, *SkipPointerListIndexManager)

	executeQueries(b, engine)

	b.StopTimer()
}
