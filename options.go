package env

import (
	"errors"
	"os"
	"reflect"
	"time"
)

type (
	envParseOpts struct {
		envLoader         EnvLoader
		separator         string
		defaultOnError    bool
		timeLayout        string
		sensitive         bool
		customMarshallers map[reflect.Type]CustomMarshaller
	}

	// EnvLoader is an alias for a function that loads values from the env. It mirrors the signature of os.Getenv.
	EnvLoader func(key string) string

	// CustomMarshaller defines the interface for custom type marshalling from environment variables.
	CustomMarshaller interface {
		UnmarshalEnv(value string) (any, error)
	}

	// CustomMarshallerFunc is a function type that implements CustomMarshaller.
	CustomMarshallerFunc func(value string) (any, error)

	// EnvParseOption is a means to customize parse options via variadic parameters.
	EnvParseOption func(o *envParseOpts) error
)

// UnmarshalEnv implements CustomMarshaller for CustomMarshallerFunc.
func (f CustomMarshallerFunc) UnmarshalEnv(value string) (any, error) {
	return f(value)
}

var (
	defaultParseOptions = envParseOpts{
		envLoader:         os.Getenv,
		separator:         ",",
		defaultOnError:    false,
		timeLayout:        time.RFC3339,
		customMarshallers: make(map[reflect.Type]CustomMarshaller),
	}
)

// WithEnvLoader allows overriding how env vars are loaded.
//
// Primarily used for testing, but feel free to get creative.
func WithEnvLoader(loader EnvLoader) EnvParseOption {
	return func(o *envParseOpts) error {
		if loader == nil {
			return errors.New("env loader function cannot be nil")
		}

		o.envLoader = loader
		return nil
	}
}

// WithEnvParseSeparator allows overriding the separated used to parse arrays/slices of a given type.
func WithEnvParseSeparator(sep string) EnvParseOption {
	return func(o *envParseOpts) error {
		if sep == "" {
			return errors.New("separator cannot be empty string")
		}

		o.separator = sep
		return nil
	}
}

// WithFallbackToDefaultOnError informs the parser that if an error is encountered during parsing, it should fallback to the default value.
func WithFallbackToDefaultOnError(fallback bool) EnvParseOption {
	return func(o *envParseOpts) error {
		o.defaultOnError = fallback
		return nil
	}
}

// WithTimeLayout allows overriding the time layout used to parse time.Time values. Default is RFC3339.
func WithTimeLayout(layout string) EnvParseOption {
	return func(o *envParseOpts) error {
		if layout == "" {
			return errors.New("time layout cannot be empty string")
		}

		o.timeLayout = layout
		return nil
	}
}

// WithSensitive informs the parser that the value being parsed is sensitive and should not be logged.
func WithSensitive(sensitive bool) EnvParseOption {
	return func(o *envParseOpts) error {
		o.sensitive = sensitive
		return nil
	}
}

// WithCustomMarshaller registers a custom marshaller for a specific type.
func WithCustomMarshaller[T any](marshaller CustomMarshaller) EnvParseOption {
	return func(o *envParseOpts) error {
		if marshaller == nil {
			return errors.New("custom marshaller cannot be nil")
		}

		var zero T
		typ := reflect.TypeOf(zero)
		if o.customMarshallers == nil {
			o.customMarshallers = make(map[reflect.Type]CustomMarshaller)
		}
		o.customMarshallers[typ] = marshaller
		return nil
	}
}

// WithCustomMarshallerFunc registers a custom marshaller function for a specific type.
func WithCustomMarshallerFunc[T any](marshallerFunc func(string) (T, error)) EnvParseOption {
	return WithCustomMarshaller[T](CustomMarshallerFunc(func(value string) (any, error) {
		result, err := marshallerFunc(value)
		return result, err
	}))
}
