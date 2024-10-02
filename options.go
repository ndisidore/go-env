package env

import (
	"errors"
	"os"
	"time"
)

type (
	envParseOpts struct {
		envLoader      EnvLoader
		separator      string
		defaultOnError bool
		timeLayout     string
		sensitive      bool
	}

	// EnvLoader is an alias for a function that loads values from the env. It mirrors the signature of os.Getenv.
	EnvLoader func(key string) string

	// EnvParseOption is a means to customize parse options via variadic parameters.
	EnvParseOption func(o *envParseOpts) error
)

var (
	defaultParseOptions = envParseOpts{
		envLoader:      os.Getenv,
		separator:      ",",
		defaultOnError: false,
		timeLayout:     time.RFC3339,
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
