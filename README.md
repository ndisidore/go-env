# go-env

A minimal, type-safe library for loading configuration from environment variables using Go generics.

## Why go-env?

There are some excellent, mature systems for loading config in the golang ecosystem such as [cobra](https://cobra.dev/) and [urfave/cli](https://cli.urfave.org/). However, these do come with some extra weight (in the form of dependencies) and syntactic lock in. A simple [flag](https://pkg.go.dev/flag) is often more then suffient for a wide swath of applications. The ergonomics of flag begin to break down where loading from env with a default value is involved, and one ends up having to write a custom function for each type and can become a bit verbose where non-strings are involved.

This aims to bridge that gap. It was created with the intention of allowing environment and type variables to be used in parallel to configure an app.

### Features

- **Minimal**: Zero dependencies, small API surface
- **Type-safe**: Compile-time type checking with generics
- **Flexible**: Support for any custom type
- **Testable**: Easy to mock and test
- **Ergonomic**: Clean, readable code
- **Production-ready**: Used in production applications

Perfect for applications that need simple, reliable environment variable parsing without the overhead of larger configuration frameworks.


## Installation

```bash
go get github.com/ndisidore/go-env
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/ndisidore/go-env"
)

func main() {
    ctx := context.Background()

    // Basic usage with defaults
    port := env.MustFromEnvOrDefault(ctx, "PORT", 8080)
    debug := env.MustFromEnvOrDefault(ctx, "DEBUG", false)
    timeout := env.MustFromEnvOrDefault(ctx, "TIMEOUT", 30*time.Second)

    fmt.Printf("Port: %d, Debug: %t, Timeout: %v\n", port, debug, timeout)
}
```

## Basic Usage

### Simple Types

```go
ctx := context.Background()

// String
name := env.MustFromEnvOrDefault(ctx, "APP_NAME", "myapp")

// Numbers
port := env.MustFromEnvOrDefault(ctx, "PORT", 8080)
maxSize := env.MustFromEnvOrDefault(ctx, "MAX_SIZE", int64(1024))
ratio := env.MustFromEnvOrDefault(ctx, "RATIO", 0.75)

// Boolean
debug := env.MustFromEnvOrDefault(ctx, "DEBUG", false)

// Duration and Time
timeout := env.MustFromEnvOrDefault(ctx, "TIMEOUT", 30*time.Second)
startTime := env.MustFromEnvOrDefault(ctx, "START_TIME", time.Now())
```

### Arrays and Slices

```go
// Comma-separated values (default separator)
hosts := env.MustFromEnvOrDefault(ctx, "HOSTS", []string{"localhost"})
// HOSTS=api.example.com,db.example.com -> ["api.example.com", "db.example.com"]

ports := env.MustFromEnvOrDefault(ctx, "PORTS", []int{8080})
// PORTS=8080,8081,8082 -> [8080, 8081, 8082]

// Custom separator
tags := env.MustFromEnvOrDefault(ctx, "TAGS", []string{"default"}, 
    env.WithEnvParseSeparator("|"))
// TAGS=web|api|database -> ["web", "api", "database"]
```

### Error Handling

```go
// With error handling
config, err := env.FromEnvOrDefault(ctx, "CONFIG_FILE", "config.json")
if err != nil {
    log.Fatal(err)
}

// Fallback to default on parse error
port := env.MustFromEnvOrDefault(ctx, "PORT", 8080, 
    env.WithFallbackToDefaultOnError(true))
```

## Advanced Usage

### Custom Types

Extend the parser to support any custom type using marshallers:

```go
// Custom type
type DatabaseURL struct {
    Host     string
    Port     int
    Database string
    SSL      bool
}

// Custom parser function
func parseDBURL(value string) (DatabaseURL, error) {
    // Parse connection string: "host:port/db?ssl=true"
    // Implementation details...
    return DatabaseURL{}, nil
}

// Usage
dbURL := env.MustFromEnvOrDefault(ctx, "DATABASE_URL", DatabaseURL{}, 
    env.WithCustomMarshallerFunc[DatabaseURL](parseDBURL))
```

### JSON Configuration

```go
type Config struct {
    Database struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"database"`
    Redis struct {
        URL string `json:"url"`
    } `json:"redis"`
}

func parseJSONConfig(value string) (Config, error) {
    var config Config
    err := json.Unmarshal([]byte(value), &config)
    return config, err
}

// Usage
config := env.MustFromEnvOrDefault(ctx, "APP_CONFIG", Config{}, 
    env.WithCustomMarshallerFunc[Config](parseJSONConfig))
```

### Enums and Constants

```go
type LogLevel int

const (
    Debug LogLevel = iota
    Info
    Warn
    Error
)

func parseLogLevel(value string) (LogLevel, error) {
    switch strings.ToLower(value) {
    case "debug": return Debug, nil
    case "info":  return Info, nil
    case "warn":  return Warn, nil
    case "error": return Error, nil
    default:
        return Info, fmt.Errorf("invalid log level: %s", value)
    }
}

logLevel := env.MustFromEnvOrDefault(ctx, "LOG_LEVEL", Info, 
    env.WithCustomMarshallerFunc[LogLevel](parseLogLevel))
```

## Integration Examples

### With Standard Library Flags

```go
import (
    "flag"
    "github.com/ndisidore/go-env"
)

var (
    port = flag.Int("port", env.MustFromEnvOrDefault(ctx, "PORT", 8080), "server port")
    host = flag.String("host", env.MustFromEnvOrDefault(ctx, "HOST", "localhost"), "server host")
)
```

### Application Configuration

```go
type AppConfig struct {
    Port        int
    Host        string
    DatabaseURL string
    Debug       bool
    Timeout     time.Duration
}

func LoadConfig() AppConfig {
    ctx := context.Background()
    
    return AppConfig{
        Port:        env.MustFromEnvOrDefault(ctx, "PORT", 8080),
        Host:        env.MustFromEnvOrDefault(ctx, "HOST", "localhost"),
        DatabaseURL: env.MustFromEnvOrDefault(ctx, "DATABASE_URL", "postgres://localhost/myapp"),
        Debug:       env.MustFromEnvOrDefault(ctx, "DEBUG", false),
        Timeout:     env.MustFromEnvOrDefault(ctx, "TIMEOUT", 30*time.Second),
    }
}
```

## Configuration Options

All functions accept optional configuration parameters:

| Option | Description | Default |
|--------|-------------|---------|
| `WithEnvLoader(loader)` | Override environment variable loading | `os.Getenv` |
| `WithEnvParseSeparator(sep)` | Array/slice separator | `","` |
| `WithFallbackToDefaultOnError(bool)` | Return default on parse error | `false` |
| `WithTimeLayout(layout)` | Time parsing format | `time.RFC3339` |
| `WithSensitive(bool)` | Mark value as sensitive for logging | `false` |
| `WithCustomMarshaller[T](marshaller)` | Register custom type marshaller | - |
| `WithCustomMarshallerFunc[T](func)` | Register custom parser function | - |

### Examples with Options

```go
// Custom separator for arrays
hosts := env.MustFromEnvOrDefault(ctx, "HOSTS", []string{}, 
    env.WithEnvParseSeparator(";"))

// Graceful error handling
timeout := env.MustFromEnvOrDefault(ctx, "TIMEOUT", 30*time.Second,
    env.WithFallbackToDefaultOnError(true))

// Custom time format
startTime := env.MustFromEnvOrDefault(ctx, "START_TIME", time.Now(),
    env.WithTimeLayout("2006-01-02 15:04:05"))

// Multiple options
config := env.MustFromEnvOrDefault(ctx, "DATABASE_CONFIG", defaultConfig,
    env.WithCustomMarshallerFunc[DatabaseConfig](parseDBConfig),
    env.WithFallbackToDefaultOnError(true),
    env.WithSensitive(true))
```

## Testing

The library is designed to be testable by allowing you to override the environment loader:

```go
func TestConfig(t *testing.T) {
    // Create a mock environment
    mockEnv := map[string]string{
        "PORT":  "9000",
        "DEBUG": "true",
        "HOSTS": "api.test.com,db.test.com",
    }
    
    loader := func(key string) string {
        return mockEnv[key]
    }
    
    ctx := context.Background()
    
    // Test with mock environment
    port := env.MustFromEnvOrDefault(ctx, "PORT", 8080, 
        env.WithEnvLoader(loader))
    
    assert.Equal(t, 9000, port)
    
    hosts := env.MustFromEnvOrDefault(ctx, "HOSTS", []string{}, 
        env.WithEnvLoader(loader))
    
    assert.Equal(t, []string{"api.test.com", "db.test.com"}, hosts)
}
```

## Supported Types

### Built-in Types
- `string`
- `bool` 
- `int`, `uint`, `int64`, `uint64`
- `float64`
- `time.Duration`
- `time.Time`
- `url.URL`
- Slices of all above types: `[]string`, `[]int`, etc.

### Custom Types
Any type can be supported by implementing a custom marshaller function or interface.

## API Reference

### Core Functions

```go
// Parse with error handling
func FromEnvOrDefault[T any](ctx context.Context, envVar string, defaultVal T, opts ...EnvParseOption) (T, error)

// Parse with panic on error  
func MustFromEnvOrDefault[T any](ctx context.Context, envVar string, defaultVal T, opts ...EnvParseOption) T
```

### Custom Marshaller Interface

```go
type CustomMarshaller interface {
    UnmarshalEnv(value string) (any, error)
}
```
