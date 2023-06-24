package env_test

import (
	"math/rand"
	"reflect"
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
				ret, err := env.FromEnvOrDefault(tt.searchEnv, defaultVal, env.WithEnvLoader(loader))
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
