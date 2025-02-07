# PhotoCifu Pocketbase extended with Go

## Install

```bash
go mod tidy
```

## Run Dev

```bash
go run . serve
```

## Build to publish

```bash
$env:GOOS = "linux"
go generate ./...
go build

./photo-cifu serve
```

## Module creation

`go mod init github.com/shujink0/photo-cifu`
