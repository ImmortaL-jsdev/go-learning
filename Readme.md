```markdown
# Go CLI Tools

Набор консольных утилит на Go:
- `wc` – подсчёт строк, слов и байт в файле.
- `downloader` – параллельное скачивание файлов с прогрессом и отменой.

## Установка

```bash
git clone https://github.com/ImmortaL-jsdev/go-learning.git
cd go-learning/go-cli-tools
make build
```

## Использование

### wc

```bash
./bin/wc test.txt          # все статистики
./bin/wc -l test.txt       # только строки
./bin/wc -w test.txt       # только слова
./bin/wc -c test.txt       # только байты
```

### downloader

```bash
./bin/downloader -file urls.txt -workers 5
```

- `-file` – файл со списком URL (обязательный)
- `-workers` – количество параллельных загрузок (по умолчанию 5)

## Разработка

```bash
make test      # тесты
make bench     # бенчмарки
make lint      # линтер
make clean     # очистка
```

## Требования

- Go 1.24+
- `golangci-lint` (для линтинга)

## Структура проекта

```
go-cli-tools/
├── cmd/
│   ├── wc/main.go
│   └── downloader/main.go
├── internal/
│   └── counter/
│       ├── counter.go
│       └── counter_test.go
├── Makefile
├── go.mod
└── go.sum
```

## Лицензия

MIT
```