## File synchronization service

Configuration file `niisva/internal/config/config.go`
```go
type Config struct {
GRPCAddress     string // server address
FileStoragePath string // tracked master directory
FileType        string // tracked file type (default: `.txt` file)
UsersPathStore  string // file with tracked slaves directories (default: `.csv` file)
}
```

### Usage:

In the directory: `cmd/watcher`
```
Execute: go run main.go
```
