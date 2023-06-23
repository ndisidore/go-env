package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// Parseable represents the types the parser is capable of handling.
	Parseable interface {
		string | bool | int | uint | int64 | uint64 | float64 | time.Duration | []string | []bool | []int | []uint | []int64 | []uint64 | []float64
	}

	envParseOpts struct {
		separator      string
		defaultOnError bool
	}

	// EnvParseOption is a means to customize parse options via variadic parameters.
	EnvParseOption func(o *envParseOpts)
)

var (
	defaultParseOptions = envParseOpts{
		separator:      ",",
		defaultOnError: false,
	}
)

// WithEnvParseSeparator allows overriding the separated used to parse arrays/slices of a given type.
func WithEnvParseSeparator(sep string) EnvParseOption {
	return func(o *envParseOpts) {
		o.separator = sep
	}
}

// WithFallbackToDefaultOnError informs the parser that if an error is encountered during parsing, it should fallback to the default value.
func WithFallbackToDefaultOnError(fallback bool) EnvParseOption {
	return func(o *envParseOpts) {
		o.defaultOnError = fallback
	}
}

// MustFromEnvOrDefault attempts to parse the environment variable provided. If it is empty or missing, the default value is used.
//
// If an error is encountered, depending on whether the `WithFallbackToDefaultOnError` option is provided it will either fallback or fatally log & exit.
func MustFromEnvOrDefault[T Parseable](envVar string, defaultVal T, opts ...EnvParseOption) (dest T) {
	parsed, err := FromEnvOrDefault(envVar, defaultVal, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return parsed
}

// FromEnvOrDefault attempts to parse the environment variable provided. If it is empty or missing, the default value is used.
//
// If an error is encountered, depending on whether the `WithFallbackToDefaultOnError` option is provided it will either fallback or return the error back to the client.
func FromEnvOrDefault[T Parseable](envVar string, defaultVal T, opts ...EnvParseOption) (dest T, err error) {
	envStr := os.Getenv(envVar)
	if envStr == "" {
		return defaultVal, nil
	}

	parseOpts := &defaultParseOptions
	for _, opt := range opts {
		opt(parseOpts)
	}

	var (
		v any
	)
	switch any(dest).(type) {
	case string:
		v = envStr
	case bool:
		v, err = strconv.ParseBool(envVar)
	case int:
		v, err = strconv.Atoi(envStr)
	case uint:
		var i uint64
		i, err = strconv.ParseUint(envStr, 10, 64)
		v = uint(i)
	case int64:
		v, err = strconv.ParseInt(envStr, 10, 64)
	case uint64:
		v, err = strconv.ParseUint(envStr, 10, 64)
	case float64:
		v, err = strconv.ParseFloat(envStr, 64)
	case time.Duration:
		v, err = time.ParseDuration(envStr)
	case []string:
		v = strings.Split(envStr, parseOpts.separator)
	case []bool:
		vs := make([]bool, 0)
		for i, at := range strings.Split(envStr, parseOpts.separator) {
			parsed, innerErr := strconv.ParseBool(at)
			if innerErr != nil {
				err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
				break
			}
			vs = append(vs, parsed)
		}
		v = vs
	case []int:
		vs := make([]int, 0)
		for i, at := range strings.Split(envStr, parseOpts.separator) {
			parsed, innerErr := strconv.Atoi(envStr)
			if innerErr != nil {
				err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
				break
			}
			vs = append(vs, parsed)
		}
		v = vs
	case []uint:
		vs := make([]uint, 0)
		for i, at := range strings.Split(envStr, parseOpts.separator) {
			parsed, innerErr := strconv.ParseUint(envStr, 10, 64)
			if innerErr != nil {
				err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
				break
			}
			vs = append(vs, uint(parsed))
		}
		v = vs
	case []int64:
		vs := make([]int64, 0)
		for i, at := range strings.Split(envStr, parseOpts.separator) {
			parsed, innerErr := strconv.ParseInt(envStr, 10, 64)
			if innerErr != nil {
				err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
				break
			}
			vs = append(vs, parsed)
		}
		v = vs
	case []uint64:
		vs := make([]uint64, 0)
		for i, at := range strings.Split(envStr, parseOpts.separator) {
			parsed, innerErr := strconv.ParseUint(envStr, 10, 64)
			if innerErr != nil {
				err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
				break
			}
			vs = append(vs, parsed)
		}
		v = vs
	case []float64:
		vs := make([]float64, 0)
		for i, at := range strings.Split(envStr, parseOpts.separator) {
			parsed, innerErr := strconv.ParseFloat(envStr, 64)
			if innerErr != nil {
				err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
				break
			}
			vs = append(vs, parsed)
		}
		v = vs
	}
	if err != nil {
		if parseOpts.defaultOnError {
			return defaultVal, nil
		}

		return dest, fmt.Errorf("failed to parse env %s to %T: %v", envVar, dest, err)
	}

	dest, ok := v.(T)
	if !ok {
		return dest, fmt.Errorf("failed to cast env %s (value: %v) to %T", envVar, v, dest)
	}
	return dest, nil
}
