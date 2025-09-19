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

### Complex command combinations

```console
$ concrun \
  -c "docker build -t myapp ." \
  -c "npm run test:unit" \
  -c "npm run test:integration"
# Runs Docker build and tests concurrently
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

