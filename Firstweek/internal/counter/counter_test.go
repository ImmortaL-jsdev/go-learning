package counter

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type table struct {
	name      string
	content   string
	wantLine  int
	wantWords int
	wantBytes int
	wantErr   bool
}

func SumSlice(s []int) int {
	var sum int
	for _, v := range s {
		sum += v
	}
	return sum
}

func BenchmarkSumSlice(b *testing.B) {

	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SumSlice(s) // передаём слайс
	}
}

func TestCountFile(t *testing.T) {
	tests := []table{
		{
			name:      "empty file",
			content:   "",
			wantLine:  0,
			wantWords: 0,
			wantBytes: 0,
			wantErr:   false,
		},
		{
			name:      "test1",
			content:   "hello world",
			wantLine:  1,
			wantWords: 2,
			wantBytes: 11,
			wantErr:   false,
		},
		{
			name:      "test2",
			content:   "hello\nworld",
			wantLine:  2,
			wantWords: 2,
			wantBytes: 11,
			wantErr:   false,
		},
		{
			name:      "test3",
			content:   "test3",
			wantLine:  0,
			wantWords: 0,
			wantBytes: 0,
			wantErr:   true,
		},
	}
	for _, tc := range tests {
		if tc.wantErr {
			_, _, _, err := CountFile("/nonexistent/file.txt")
			if err == nil {
				t.Errorf("expected error but got nil for test %s", tc.name)
			}
			continue
		} else {
			fileName := fmt.Sprintf("test_%s.txt", strings.ReplaceAll(tc.name, " ", "_"))
			err := os.WriteFile(fileName, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(fileName)
			lines, words, bytes, err := CountFile(fileName)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if lines != tc.wantLine {
				t.Errorf("unexpected lines: %v", lines)
			}

			if words != tc.wantWords {
				t.Errorf("unexpected words: %v", words)
			}

			if bytes != int64(tc.wantBytes) {
				t.Errorf("unexpected bytes: %v", bytes)
			}

		}
	}
}

func BenchmarkCountFile(b *testing.B) {
	content := strings.Repeat("hello world\n", 1000)
	fileName := "benchmark_file.txt"
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		b.Fatalf("failed to create file: %v", err)
	}
	defer os.Remove(fileName)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _, err := CountFile(fileName)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
