package env

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	// Parseable represents the types the parser is capable of handling.
	// This includes built-in types and any type that can be handled by custom marshallers.
	Parseable interface {
		string | bool | int | uint | int64 | uint64 | float64 | time.Duration | time.Time | url.URL | []string | []bool | []int | []uint | []int64 | []uint64 | []float64 | []time.Duration | []time.Time | []url.URL
	}
)

// MustFromEnvOrDefault attempts to parse the environment variable provided. If it is empty or missing, the default value is used.
//
// If an error is encountered, depending on whether the `WithFallbackToDefaultOnError` option is provided it will either fallback or fatally log & exit.
func MustFromEnvOrDefault[T any](ctx context.Context, envVar string, defaultVal T, opts ...EnvParseOption) (dest T) {
	parsed, err := FromEnvOrDefault(ctx, envVar, defaultVal, opts...)
	if err != nil {
		slog.Default().ErrorContext(ctx, "failed to parse env var", slog.String("env_var", envVar), slog.String("error", err.Error()))
		os.Exit(1)
	}

	return parsed
}

// FromEnvOrDefault attempts to parse the environment variable provided. If it is empty or missing, the default value is used.
//
// If an error is encountered, depending on whether the `WithFallbackToDefaultOnError` option is provided it will either fallback or return the error back to the client.
func FromEnvOrDefault[T any](ctx context.Context, envVar string, defaultVal T, opts ...EnvParseOption) (dest T, err error) {
	parseOpts := defaultParseOptions
	// Create a copy of the custom marshallers map to avoid modifying the global state
	if parseOpts.customMarshallers != nil {
		parseOpts.customMarshallers = make(map[reflect.Type]CustomMarshaller)
		for k, v := range defaultParseOptions.customMarshallers {
			parseOpts.customMarshallers[k] = v
		}
	}

	for _, opt := range opts {
		if err := opt(&parseOpts); err != nil {
			return dest, fmt.Errorf("option error: %w", err)
		}
	}

	envStr := parseOpts.envLoader(envVar)
	if envStr == "" {
		return defaultVal, nil
	}

	var (
		v any
	)

	// Check for custom marshaller first
	destType := reflect.TypeOf(dest)
	if marshaller, exists := parseOpts.customMarshallers[destType]; exists {
		v, err = marshaller.UnmarshalEnv(envStr)
		if err != nil {
			if parseOpts.defaultOnError {
				return defaultVal, nil
			}
			return dest, fmt.Errorf("custom marshaller failed for env %s to %T: %w", envVar, dest, err)
		}
	} else {
		// Fall back to built-in type handling
		switch any(dest).(type) {
		case string:
			v = envStr
		case bool:
			v, err = strconv.ParseBool(envStr)
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
		case time.Time:
			v, err = time.Parse(parseOpts.timeLayout, envStr)
		case url.URL:
			v, err = url.Parse(envStr)
		case []string:
			v = strings.Split(envStr, parseOpts.separator)
		case []bool:
			vs := make([]bool, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
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
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := strconv.Atoi(at)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, parsed)
			}
			v = vs
		case []uint:
			vs := make([]uint, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := strconv.ParseUint(at, 10, 64)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, uint(parsed))
			}
			v = vs
		case []int64:
			vs := make([]int64, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := strconv.ParseInt(at, 10, 64)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, parsed)
			}
			v = vs
		case []uint64:
			vs := make([]uint64, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := strconv.ParseUint(at, 10, 64)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, parsed)
			}
			v = vs
		case []float64:
			vs := make([]float64, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := strconv.ParseFloat(at, 64)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, parsed)
			}
			v = vs
		case []time.Duration:
			vs := make([]time.Duration, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := time.ParseDuration(at)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, parsed)
			}
			v = vs
		case []time.Time:
			vs := make([]time.Time, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := time.Parse(parseOpts.timeLayout, at)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, parsed)
			}
			v = vs
		case []url.URL:
			vs := make([]url.URL, 0)
			for i, at := range splitAndTrim(envStr, parseOpts.separator) {
				parsed, innerErr := url.Parse(at)
				if innerErr != nil {
					err = fmt.Errorf("item %s (pos: %d) failed to parse: %w", at, i, innerErr)
					break
				}
				vs = append(vs, *parsed)
			}
			v = vs
		default:
			// If no built-in type matches and no custom marshaller, return error
			return dest, fmt.Errorf("unsupported type %T for env var %s", dest, envVar)
		}
	}

	if err != nil {
		if parseOpts.defaultOnError {
			return defaultVal, nil
		}

		return dest, fmt.Errorf("failed to parse env %s to %T: %v", envVar, dest, err)
	}

	dest, ok := v.(T)
	if !ok {
		return dest, fmt.Errorf("failed to cast env %s to %T", envVar, dest)
	}
	return dest, nil
}

func splitAndTrim(in string, sep string) []string {
	strs := strings.Split(in, sep)
	for i, str := range strs {
		strs[i] = strings.TrimSpace(str)
	}
	return strs
}
