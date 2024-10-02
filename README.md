# go-env

A minimal way to load typed data into a golang application from a system's environment utilizing generics.

## Why?

There are some excellent, mature systems for loading config in the golang ecosystem such as [cobra](https://cobra.dev/) and [urfave/cli](https://cli.urfave.org/).
These come with all kinds of fancy bells and whistles and the author can attest to them having used them in many applications with great success.

However, these do come with some extra weight (in the form of dependencies) and syntactic lock in. Often times, just using a simple [flag](https://pkg.go.dev/flag)
is the desired option. The ergonomics of flag begin to struggle where loading from env with a default value is involved, and one ends up having to write a custom
function for each type and can become a bit verbose where non-strings are involved.

This aims to bridge that gap. It was created with the intention of allowing environment and type variables to be used in parallel to configure an app.

## How.

### With flags.

Lets look at the default signature for a flag, here for `uint64`

```
func Uint64(name string, value uint64, usage string) *uint64
```

With `go-env` we can provide a value for, ahem, `value` that is a type checked env in a single function call (thanks generics!).

```go
import (
    "os"
    "flags"

    "github.com/ndisidore/go-env"
)

os.Setenv("FIRST_FLAG_ENV_VAR", "42")
var ctx = context.Background()
// since `FIRST_FLAG_ENV_VAR` env is set, this will parse and use that value
var myFlag1 = flags.Uint64("first-flag-cmd-line", env.MustFromEnvOrDefault(ctx, "FIRST_FLAG_ENV_VAR", 7), "an example uint64 flag") *uint64
// since `SECOND_FLAG_ENV_VAR` env is not set, this will fallback to the default value
var myFlag2 = flags.Uint64("second-flag-cmd-line", env.MustFromEnvOrDefault(ctx, "SECOND_FLAG_ENV_VAR", 7), "an example uint64 flag") *uint64
fmt.Printf("%d; %d", *myFlag1, *myFlag2) // outputs -> 42; 7
```

Since `MustFromEnvOrDefault` (and its error returning counterpart `FromEnvOrDefault`) use generics and analyze the type dynamically, the call looks very similar for other flag data types.

```go
var ctx = context.Background()
var myStrFlag = flags.String("my-str-flag", env.MustFromEnvOrDefault(ctx, "MY_STR", "a string"), "an example string flag") *string
var myBoolFlag = flags.Bool("my-bool-flag", env.MustFromEnvOrDefault(ctx, "MY_BOOL", true), "an example bool flag") *bool
var myDurFlag = flags.Duration("my-duration-flag", env.MustFromEnvOrDefault(ctx, "MY_DURATION", time.Second * 5), "an example duration flag") *time.Duration
var myFloatFlag = flags.Float64("my-float64-flag", env.MustFromEnvOrDefault(ctx, "MY_FLOAT64", 7.11), "an example float64 flag") *float64
```

### Without flags.

Of course, usage with [flag](https://pkg.go.dev/flag) is not strictly necessary here. If one doesn't require cli flag support and just needs to load from env, that can be done directly.

```go
var (
    ctx       = context.Background()
    myStrVar  = env.MustFromEnvOrDefault(ctx, "MY_STR", "a string")
    myBoolVar = env.MustFromEnvOrDefault(ctx, "MY_BOOL", true)
)

// --- or ---

myStrVar2, err := env.FromEnvOrDefault(ctx, "MY_STR", "a string")
if err != nil { ... }
myBoolVar2, err := env.FromEnvOrDefault(ctx, "MY_BOOL", true)
if err != nil { ... }
```
