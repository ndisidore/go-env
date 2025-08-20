package env_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ndisidore/go-env"
)

func TestParsesParseable(t *testing.T) {
	t.Parallel()

	// TODO: table driven tests, but generics makes this difficult
	var makeLoader func(envs map[string]string) env.EnvLoader = func(envs map[string]string) env.EnvLoader {
		return func(key string) string {
			return envs[key]
		}
	}

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		const defaultVal = "default"
		var (
			loader = makeLoader(map[string]string{"KNOWN_STRING": "a string"})
			cases  = []struct {
				searchEnv string
				expected  string
			}{
				{searchEnv: "KNOWN_STRING", expected: "a string"},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				if err != nil {
					t.Logf("unexpected error: %v", err)
					t.Fail()
				}
				if ret != tt.expected {
					t.Logf("return value (%s) does not match expected (%s)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		const defaultVal = false
		var (
			loader = makeLoader(map[string]string{"KNOWN_BOOL": "true", "NOT_BOOL": "abcd"})
			cases  = []struct {
				searchEnv           string
				expected            bool
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_BOOL", expected: true},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_BOOL", expected: false, expectedErrContains: "invalid syntax"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%t) does not match expected (%t)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()
		var defaultVal = rand.Int()
		var (
			loader = makeLoader(map[string]string{"KNOWN_INT": "123", "NOT_INT": "abcd"})
			cases  = []struct {
				searchEnv           string
				expected            int
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_INT", expected: 123},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_INT", expectedErrContains: "invalid syntax"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%d) does not match expected (%d)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("uint", func(t *testing.T) {
		t.Parallel()
		const defaultVal = uint(555)
		var (
			loader = makeLoader(map[string]string{"KNOWN_UINT": "123", "NOT_UINT": "abcd"})
			cases  = []struct {
				searchEnv           string
				expected            uint
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_UINT", expected: 123},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_UINT", expectedErrContains: "invalid syntax"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%d) does not match expected (%d)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		var (
			defaultVal = rand.Int63()
			loader     = makeLoader(map[string]string{"KNOWN_INT": "8675309", "NOT_INT": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            int64
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_INT", expected: 8675309},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_INT", expectedErrContains: "invalid syntax"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%d) does not match expected (%d)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("uint64", func(t *testing.T) {
		t.Parallel()
		var (
			defaultVal = rand.Uint64()
			loader     = makeLoader(map[string]string{"KNOWN_UINT": "5555555", "NOT_UINT": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            uint64
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_UINT", expected: 5555555},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_UINT", expectedErrContains: "invalid syntax"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%d) does not match expected (%d)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		t.Parallel()
		var (
			defaultVal = rand.Float64()
			loader     = makeLoader(map[string]string{"KNOWN_FLOAT": "69.69", "NOT_FLOAT": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            float64
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_FLOAT", expected: 69.69},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_FLOAT", expectedErrContains: "invalid syntax"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%.4f) does not match expected (%.4f)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("time.Duration", func(t *testing.T) {
		t.Parallel()
		var (
			defaultVal = time.Minute * 5
			loader     = makeLoader(map[string]string{"KNOWN_DURATION": "10s", "NOT_DURATION": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            time.Duration
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_DURATION", expected: time.Second * 10},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_DURATION", expectedErrContains: "invalid duration"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case ret != tt.expected:
					t.Logf("return value (%s) does not match expected (%s)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("time.Time", func(t *testing.T) {
		t.Parallel()
		var (
			defaultVal = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
			loader     = makeLoader(map[string]string{"KNOWN_TIME": "2021-01-01T00:00:00Z", "NOT_TIME": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            time.Time
				expectedErrContains string
				options             []env.EnvParseOption
			}{
				{searchEnv: "KNOWN_TIME", expected: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{searchEnv: "KNOWN_TIME", expectedErrContains: "parsing time", options: []env.EnvParseOption{env.WithTimeLayout(time.RFC1123)}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				{searchEnv: "NOT_TIME", expectedErrContains: "parsing time"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, append(tt.options, env.WithEnvLoader(loader))...)
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !ret.Equal(tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]string", func(t *testing.T) {
		t.Parallel()
		var (
			defaultVal = []string{"hello", "world"}
			loader     = makeLoader(map[string]string{"KNOWN_STR_ARRAY": "testy,mctesterson,jr", "NOT_STR_ARRAY": "abcd"})
			cases      = []struct {
				searchEnv string
				expected  []string
			}{
				{searchEnv: "KNOWN_STR_ARRAY", expected: []string{"testy", "mctesterson", "jr"}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_STR_ARRAY", expected: []string{"abcd"}},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]bool", func(t *testing.T) {
		var (
			defaultVal = []bool{true, false, true}
			loader     = makeLoader(map[string]string{"KNOWN_BOOL_ARRAY": "true, true,false", "NOT_BOOL_ARRAY": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            []bool
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_BOOL_ARRAY", expected: []bool{true, true, false}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_BOOL_ARRAY", expectedErrContains: "item abcd (pos: 0) failed to parse"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]int", func(t *testing.T) {
		var (
			defaultVal = []int{12, 8, 263, -6}
			loader     = makeLoader(map[string]string{"KNOWN_INT_ARRAY": "63, 52,-8,285", "NOT_INT_ARRAY": "abcd"})
			cases      = []struct {
				searchEnv           string
				expected            []int
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_INT_ARRAY", expected: []int{63, 52, -8, 285}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_INT_ARRAY", expectedErrContains: "item abcd (pos: 0) failed to parse"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]uint", func(t *testing.T) {
		var (
			defaultVal = []uint{12, 8, 263, 481}
			loader     = makeLoader(map[string]string{"KNOWN_UINT_ARRAY": "63, 52,0,285", "NOT_UINT_ARRAY": "32,-2,abcd"})
			cases      = []struct {
				searchEnv           string
				expected            []uint
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_UINT_ARRAY", expected: []uint{63, 52, 0, 285}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_UINT_ARRAY", expectedErrContains: "item -2 (pos: 1) failed to parse"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]int64", func(t *testing.T) {
		var (
			defaultVal = []int64{rand.Int63(), rand.Int63(), rand.Int63(), rand.Int63()}
			loader     = makeLoader(map[string]string{"KNOWN_INT_ARRAY": "616515641, 52,0,-6115122", "NOT_INT_ARRAY": "32,-2,abcd"})
			cases      = []struct {
				searchEnv           string
				expected            []int64
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_INT_ARRAY", expected: []int64{616515641, 52, 0, -6115122}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_INT_ARRAY", expectedErrContains: "item abcd (pos: 2) failed to parse"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]uint64", func(t *testing.T) {
		var (
			defaultVal = []uint64{rand.Uint64(), rand.Uint64(), rand.Uint64(), rand.Uint64()}
			loader     = makeLoader(map[string]string{"KNOWN_UINT_ARRAY": "616515641, 52,0,6115122", "NOT_UINT_ARRAY": "32,-2,abcd"})
			cases      = []struct {
				searchEnv           string
				expected            []uint64
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_UINT_ARRAY", expected: []uint64{616515641, 52, 0, 6115122}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_UINT_ARRAY", expectedErrContains: "item -2 (pos: 1) failed to parse"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})

	t.Run("[]float64", func(t *testing.T) {
		var (
			defaultVal = []float64{rand.Float64(), rand.Float64(), rand.Float64()}
			loader     = makeLoader(map[string]string{"KNOWN_FLOAT_ARRAY": "845.15, -52.3,0.0,666.5154, 7", "NOT_FLOAT_ARRAY": "32.22,-2.151,abcd"})
			cases      = []struct {
				searchEnv           string
				expected            []float64
				expectedErrContains string
			}{
				{searchEnv: "KNOWN_FLOAT_ARRAY", expected: []float64{845.15, -52.3, 0.0, 666.5154, 7}},
				{searchEnv: "UNKNOWN_ENV", expected: defaultVal},
				// since even single values are a string, we don't expect an error here, just a single length array
				{searchEnv: "NOT_FLOAT_ARRAY", expectedErrContains: "item abcd (pos: 2) failed to parse"},
			}
		)
		for _, tt := range cases {
			t.Run("", func(t *testing.T) {
				ret, err := env.FromEnvOrDefault(context.Background(), tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
				switch {
				case err != nil && tt.expectedErrContains != "":
					if !strings.Contains(err.Error(), tt.expectedErrContains) {
						t.Logf("unexpected error: %v", err)
						t.Fail()
					}
				case err != nil:
					t.Logf("unexpected error: %v", err)
					t.Fail()
				case !reflect.DeepEqual(ret, tt.expected):
					t.Logf("return value (%v) does not match expected (%v)", ret, tt.expected)
					t.Fail()
				}
			})
		}
	})
}

// Custom types for testing custom marshallers
type IPAddress net.IP

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

// Custom marshaller implementations
func parseIPAddress(value string) (IPAddress, error) {
	ip := net.ParseIP(value)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", value)
	}
	return IPAddress(ip), nil
}

func parseConfig(value string) (Config, error) {
	var config Config
	if err := json.Unmarshal([]byte(value), &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config JSON: %w", err)
	}
	return config, nil
}

func parseLogLevel(value string) (LogLevel, error) {
	switch strings.ToLower(value) {
	case "debug":
		return LogLevelDebug, nil
	case "info":
		return LogLevelInfo, nil
	case "warn":
		return LogLevelWarn, nil
	case "error":
		return LogLevelError, nil
	default:
		if level, err := strconv.Atoi(value); err == nil && level >= 0 && level <= 3 {
			return LogLevel(level), nil
		}
		return LogLevelError, fmt.Errorf("invalid log level: %s", value)
	}
}

// Custom marshaller for testing interface-based marshalling
type customMarshaller struct{}

func (cm customMarshaller) UnmarshalEnv(value string) (any, error) {
	return parseIPAddress(value)
}

// JSON marshaller for testing JSON config parsing
type jsonMarshaller struct{}

func (jm jsonMarshaller) UnmarshalEnv(value string) (any, error) {
	var config Config
	err := json.Unmarshal([]byte(value), &config)
	return config, err
}

func TestCustomMarshallerFunc(t *testing.T) {
	ctx := context.Background()
	
	// Test IP address parsing
	os.Setenv("TEST_IP", "192.168.1.1")
	defer os.Unsetenv("TEST_IP")
	
	ip, err := env.FromEnvOrDefault(ctx, "TEST_IP", IPAddress(nil), 
		env.WithCustomMarshallerFunc[IPAddress](parseIPAddress))
	
	if err != nil {
		t.Fatalf("Failed to parse IP address: %v", err)
	}
	
	expected := net.ParseIP("192.168.1.1")
	if !net.IP(ip).Equal(expected) {
		t.Errorf("Expected IP %v, got %v", expected, net.IP(ip))
	}
}

func TestCustomMarshallerJSON(t *testing.T) {
	ctx := context.Background()
	
	// Test JSON config parsing
	os.Setenv("TEST_CONFIG", `{"host":"localhost","port":8080}`)
	defer os.Unsetenv("TEST_CONFIG")
	
	config, err := env.FromEnvOrDefault(ctx, "TEST_CONFIG", Config{}, 
		env.WithCustomMarshallerFunc[Config](parseConfig))
	
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}
	
	if config.Host != "localhost" || config.Port != 8080 {
		t.Errorf("Expected config {localhost 8080}, got %+v", config)
	}
}

func TestCustomMarshallerEnum(t *testing.T) {
	ctx := context.Background()
	
	// Test enum parsing by string
	os.Setenv("TEST_LOG_LEVEL", "warn")
	defer os.Unsetenv("TEST_LOG_LEVEL")
	
	level, err := env.FromEnvOrDefault(ctx, "TEST_LOG_LEVEL", LogLevelInfo, 
		env.WithCustomMarshallerFunc[LogLevel](parseLogLevel))
	
	if err != nil {
		t.Fatalf("Failed to parse log level: %v", err)
	}
	
	if level != LogLevelWarn {
		t.Errorf("Expected log level %v, got %v", LogLevelWarn, level)
	}
	
	// Test enum parsing by number
	os.Setenv("TEST_LOG_LEVEL", "3")
	
	level, err = env.FromEnvOrDefault(ctx, "TEST_LOG_LEVEL", LogLevelInfo, 
		env.WithCustomMarshallerFunc[LogLevel](parseLogLevel))
	
	if err != nil {
		t.Fatalf("Failed to parse log level: %v", err)
	}
	
	if level != LogLevelError {
		t.Errorf("Expected log level %v, got %v", LogLevelError, level)
	}
}

func TestCustomMarshallerInterface(t *testing.T) {
	ctx := context.Background()
	
	os.Setenv("TEST_IP_INTERFACE", "10.0.0.1")
	defer os.Unsetenv("TEST_IP_INTERFACE")
	
	ip, err := env.FromEnvOrDefault(ctx, "TEST_IP_INTERFACE", IPAddress(nil), 
		env.WithCustomMarshaller[IPAddress](customMarshaller{}))
	
	if err != nil {
		t.Fatalf("Failed to parse IP address with interface: %v", err)
	}
	
	expected := net.ParseIP("10.0.0.1")
	if !net.IP(ip).Equal(expected) {
		t.Errorf("Expected IP %v, got %v", expected, net.IP(ip))
	}
}

func TestCustomMarshallerError(t *testing.T) {
	ctx := context.Background()
	
	// Test error handling
	os.Setenv("TEST_INVALID_IP", "not-an-ip")
	defer os.Unsetenv("TEST_INVALID_IP")
	
	_, err := env.FromEnvOrDefault(ctx, "TEST_INVALID_IP", IPAddress(nil), 
		env.WithCustomMarshallerFunc[IPAddress](parseIPAddress))
	
	if err == nil {
		t.Fatal("Expected error for invalid IP address")
	}
	
	if !strings.Contains(err.Error(), "invalid IP address") {
		t.Errorf("Expected error to contain 'invalid IP address', got: %v", err)
	}
}

func TestCustomMarshallerFallback(t *testing.T) {
	ctx := context.Background()
	
	// Test fallback to default on error
	os.Setenv("TEST_INVALID_IP_FALLBACK", "not-an-ip")
	defer os.Unsetenv("TEST_INVALID_IP_FALLBACK")
	
	defaultIP := IPAddress(net.ParseIP("127.0.0.1"))
	ip, err := env.FromEnvOrDefault(ctx, "TEST_INVALID_IP_FALLBACK", defaultIP, 
		env.WithCustomMarshallerFunc[IPAddress](parseIPAddress),
		env.WithFallbackToDefaultOnError(true))
	
	if err != nil {
		t.Fatalf("Expected no error with fallback enabled, got: %v", err)
	}
	
	if !net.IP(ip).Equal(net.ParseIP("127.0.0.1")) {
		t.Errorf("Expected fallback to default IP 127.0.0.1, got %v", net.IP(ip))
	}
}
