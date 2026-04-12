package counter

import (
	"bufio"
	"os"
	"strings"
)

func CountFile(filename string) (int, int, int64, error) {
	file, err := os.Open(filename)

	if err != nil {
		return 0, 0, 0, err

	}
	defer file.Close() //nolint:errcheck

	info, err := file.Stat()

	if err != nil {
		return 0, 0, 0, err
	}

	bytes := info.Size()

	scanner := bufio.NewScanner(file)

	lines, words := 0, 0

	for scanner.Scan() {
		lines++
		line := scanner.Text()
		fields := strings.Fields(line)
		words += len(fields)
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, 0, err
	}

	return lines, words, bytes, nil
}
