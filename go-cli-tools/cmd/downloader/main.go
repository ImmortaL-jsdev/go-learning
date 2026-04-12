package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func downloadFile(ctx context.Context, url string, index int) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {

		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	filename := fmt.Sprintf("file_%d", index)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	filesFlag := flag.String("file", "", "print file name")
	workersFlag := flag.Int("workers", 5, "print count of workers")

	flag.Parse()

	if *filesFlag == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -file <filename> [-workers N]\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("Файл со списком URL: %s, количество воркеров: %d\n", *filesFlag, *workersFlag)

	file, err := os.Open(*filesFlag)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
		os.Exit(1)
	}

	defer file.Close() //nolint:errcheck

	scanner := bufio.NewScanner(file)

	var urls []string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		urls = append(urls, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: файл не содержит ни одного URL")
		os.Exit(1)
	}

	fmt.Printf("Найдено URL: %d\n", len(urls))

	for i, u := range urls {
		fmt.Printf("%d: %s\n", i+1, u)
	}

	downloaded := make(map[string]bool)

	downFile, err := os.Open("downloaded.txt")
	if err != nil {
		if os.IsNotExist(err) {

		} else {
			fmt.Fprintf(os.Stderr, "Ошибка при открытии downloaded.txt: %v\n", err)
			os.Exit(1)
		}
	} else {
		defer downFile.Close() //nolint:errcheck
		scanner := bufio.NewScanner(downFile)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			downloaded[line] = true
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения downloaded.txt: %v\n", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Получен сигнал завершения. Отмена...")
		cancel()
	}()

	jobs := make(chan struct {
		url string
		id  int
	}, len(urls))

	success := make(chan string, len(urls))

	var wg sync.WaitGroup

	for w := 0; w < *workersFlag; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					} else {
						err := downloadFile(ctx, job.url, job.id)
						if err != nil {
							fmt.Printf("Error %s: %v\n", job.url, err)
						} else {
							success <- job.url
							fmt.Printf("Success %s\n", job.url)
						}
					}
				}
			}
		}()

	}
	for i, url := range urls {
		if downloaded[url] {
			continue
		} else {
			jobs <- struct {
				url string
				id  int
			}{url, i}
		}
	}
	close(jobs)
	wg.Wait()
	close(success)

	outFile, err := os.OpenFile("downloaded.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия downloaded.txt для записи: %v\n", err)
	} else {
		defer outFile.Close() //nolint:errcheck
		for url := range success {
			if _, err := fmt.Fprintf(outFile, "%s\n", url); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка записи URL %s: %v\n", url, err)
			}
		}
	}
	fmt.Println("All downloads are complete!")
}
