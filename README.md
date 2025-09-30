# concrun

Run commands concurrently.

## Features

- Run multiple commands concurrently
- Detailed execution reports with timing information

## Usage

### Basic usage

```console
$ concrun -c "echo hello" -c "echo world"
--------------------------------------------------
▶ echo world
world
---- [ exit code: 0, excution time: 10.425333ms ]
--------------------------------------------------
▶ echo hello
hello
---- [ exit code: 0, excution time: 10.968042ms ]
```

### Fail-fast execution

```console
$ concrun --fail-fast -c "make test" -c "make lint" -c "make build"
```

### Retry on failure

```console
$ concrun --max-retries-per-command 3 -c "flaky-test.sh" -c "another-test.sh"
```

### Complex command combinations

```console
$ concrun \
  -c "docker build -t myapp ." \
  -c "npm run test:unit" \
  -c "npm run test:integration"
# Runs Docker build and tests concurrently
```

## Options

### `--max-retries-per-command`

Specifies the maximum number of times to automatically retry each command if it fails. The default is 0 (no retries).

- Commands will be executed repeatedly until they succeed (exit code 0) or until the specified number of retries is reached
- All retry attempts are logged

```console
$ concrun --max-retries-per-command 3 -c "npm test" -c "npm run lint"
```

### `--fail-fast`

Cancels other running commands if any command fails. The default is false (all commands run to completion).

- The determination is made **after all retries for each command are completed**
- This means that for each command type, failure is determined only after all of its retries have been exhausted

```console
$ concrun --fail-fast -c "make test" -c "make lint"
```

### Option Interactions

When both options are used together:

1. When a command fails, it first retries up to the specified number of times
2. After all retries are complete, the final exit code is checked
3. If the exit code is non-zero, other commands are cancelled

```console
$ concrun --fail-fast --max-retries-per-command 2 -c "npm test" -c "npm run build"
```

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/concrun
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/concrun/releases)

**go install:**

```console
$ go install github.com/k1LoW/concrun@latest
```

**deb:**

``` console
$ export CONCRUN_VERSION=X.X.X
$ curl -o concrun.deb -L https://github.com/k1LoW/concrun/releases/download/v$CONCRUN_VERSION/concrun_$CONCRUN_VERSION-1_amd64.deb
$ dpkg -i concrun.deb
```

**RPM:**

``` console
$ export CONCRUN_VERSION=X.X.X
$ yum install https://github.com/k1LoW/concrun/releases/download/v$CONCRUN_VERSION/concrun_$CONCRUN_VERSION-1_amd64.rpm

