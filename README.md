# cadlib for Go

Go bindings for the SimScale `cadlib` native library: load a CAD model input
to a workflow method and run geometry queries on it (volume, bounding box,
tessellation, …).

> [!NOTE]
> `cadlib` is meant to be used in the Docker image that implements the workflow
> method. It is not intended for validations.

The package is a thin cgo layer over the prebuilt `cadlib` native library.
You need Go 1.21 or newer and a C compiler with `cgo` enabled (`CGO_ENABLED=1`):
`gcc` or `clang`, and `mingw-w64 gcc` on Windows.

## Install

### 1. Install the native library

Each release ships two install scripts: `install-cadlib.sh` for Linux and macOS,
and `install-cadlib.ps1` for Windows. A script downloads the native library of
its own release, checks it against the SHA-256 digest embedded in the script,
and installs it. Use the scripts of the release you depend on.

#### Linux and macOS

```bash
curl -fsSL https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/v<version>/install-cadlib.sh | sh
```

It installs into `/usr/local` by default, which usually needs `sudo`
(`… | sudo sh`), and prints the follow-up step for your platform:
`sudo ldconfig` on Linux, or the environment variables below for any other
location. `--dry-run` prints the download URL, the digest and the install paths
without changing anything; `--help` lists the options.

To install somewhere else, for example without `sudo`:

```bash
curl -fsSL https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/v<version>/install-cadlib.sh | sh -s -- --prefix "$HOME/.local"
export CGO_LDFLAGS="-L$HOME/.local/lib"                                         # where the build links it
export LD_LIBRARY_PATH="$HOME/.local/lib${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"  # where programs load it
```

On macOS use `CGO_LDFLAGS="-L$HOME/.local/lib -Wl,-rpath,$HOME/.local/lib"` and
`DYLD_LIBRARY_PATH` instead.

#### Windows

`cgo` on Windows needs the `mingw-w64 gcc`; Visual Studio's compiler does not
work with `cgo`. In PowerShell:

```powershell
choco install -y mingw                     # or MSYS2, with its mingw64\bin on PATH
irm https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/v<version>/install-cadlib.ps1 | iex
```

It installs into `%LOCALAPPDATA%\cadlib` by default and prints the `CGO_LDFLAGS`
and `PATH` settings to use. `-Prefix` installs somewhere else and `-DryRun`
changes nothing, like `--prefix` and `--dry-run` above:

```powershell
iex "& { $(irm https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/v<version>/install-cadlib.ps1) } -Prefix C:\cadlib"
```

To run your program on another machine, ship `cadlib.dll` next to the `.exe`;
that machine also needs the
[VC++ 2015–2022 redistributable](https://learn.microsoft.com/cpp/windows/latest-supported-vc-redist).

#### In CI

Install the library of the module version your `go.mod` resolves to. On a Linux
runner:

```bash
V=$(go list -m -f '{{.Version}}' github.com/SimScaleGmbH/simscale-go-cadlib)  # the module version is the git tag
curl -fsSL "https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/$V/install-cadlib.sh" | sh
ldconfig
```

CI containers usually run as root; without root, use `--prefix` and the exports
above. On a Windows runner, in one step:

```powershell
$V = go list -m -f '{{.Version}}' github.com/SimScaleGmbH/simscale-go-cadlib
irm "https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/$V/install-cadlib.ps1" | iex
$env:CGO_LDFLAGS = "-L$($env:LOCALAPPDATA -replace '\\', '/')/cadlib/lib"
$env:PATH = "$env:LOCALAPPDATA\cadlib\lib;$env:PATH"
```

To download from a mirror of the release assets, set `CADLIB_DOWNLOAD_BASE`
(the script fetches `$CADLIB_DOWNLOAD_BASE/v<version>/cadlib-<platform>.tar.gz`);
the embedded digest still applies.

#### About the script

- **Needs** `curl` or `wget`, and `sha256sum` or `shasum`; on Windows, the
  `tar.exe` that ships with Windows 10 1803 and later.
- **Installs** `lib/libcadlib.so` (`libcadlib.dylib` on macOS, `cadlib.dll` and
  `cadlib.lib` on Windows), `include/cadlib.h`, and the licenses in
  `share/doc/cadlib/`, all under the prefix. To uninstall, delete them.
- **Updating** means running the new release's script. On a checksum mismatch it
  installs nothing. On Linux and macOS, running programs keep the old library
  until they restart; on Windows, stop them first.
- **To read it before running it**, download it in place of piping it to `sh`.
  A commit hash in place of the tag in the URL pins the script and its digests.
- **Errors:** "Intel Macs are not supported" also shows in a Rosetta shell on an
  Apple silicon Mac. "cannot write to /usr/local" needs `sudo` or `--prefix`.

### 2. Get the module

```bash
go get github.com/SimScaleGmbH/simscale-go-cadlib@v<version>
```

```go
import cadlib "github.com/SimScaleGmbH/simscale-go-cadlib"
```

The C header ships inside the module.

### Versions must match

Each release of the module accepts exactly one native library: the one released
with it. On the first `Load` the package compares the two and fails with an
error pointing at the install scripts of the right release.

## Usage

```go
// Load the CAD model from the CAD data map entry of the method's input parameters.
cad, err := cadlib.Load(request.DataMap.Cad)
if err != nil {
	return err
}
defer cad.Close()

result, err := cad.Volume(&cadlib.VolumeParameters{EntityIDs: []string{"B1_TE5"}})
```

- **Load once, reuse.** Each `Load` rebuilds the model, which is expensive, so
  load once and run every query on the same `*Cad`.
- **Close** releases the native resources; call it when you no longer need the
  model.
- **Not safe for concurrent use.** Queries on one `*Cad` block and must not run
  from several goroutines at once.
- **Errors** are either a `*cadlib.CadExpectedError`, a known domain error with
  an `ErrorCode`, or a `*cadlib.CadUnexpectedError`. Use `errors.As` to tell them
  apart. The codes are listed in [ERRORS.md](ERRORS.md). The exception is a
  failure in the bindings to marshal an argument or unmarshal a result, which
  is returned as a plain `error` wrapping the `encoding/json` error.

The API reference is on
[pkg.go.dev](https://pkg.go.dev/github.com/SimScaleGmbH/simscale-go-cadlib).

## Docker

Use a glibc base image (see [Supported platforms](#supported-platforms)). The
library has to be in both stages: the build links it, the runtime loads it.

```dockerfile
FROM golang:1.21-bookworm AS builder
ARG CADLIB_VERSION
RUN curl -fsSL "https://raw.githubusercontent.com/SimScaleGmbH/simscale-go-cadlib/v${CADLIB_VERSION}/install-cadlib.sh" | sh \
 && ldconfig
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o /app/myservice .

FROM debian:bookworm-slim
COPY --from=builder /usr/local/lib/libcadlib.so /usr/local/lib/
RUN ldconfig
COPY --from=builder /app/myservice /usr/local/bin/myservice
ENTRYPOINT ["myservice"]
```

Build it with `--build-arg CADLIB_VERSION=<version>`, the module version you
depend on without the leading `v`. [In CI](#in-ci) shows how to read it from your `go.mod`.

Without `ldconfig` the library is on disk but not found: `cannot find -lcadlib`
when building, `error while loading shared libraries: libcadlib.so` when running.

## Supported platforms

| OS | Architecture | Native library | Requires |
|----|--------------|----------------|----------|
| Linux | x86_64, aarch64 | `libcadlib.so` | glibc 2.28 or newer (Debian 10+, Ubuntu 20.04+, RHEL 8+) |
| macOS | aarch64 (Apple silicon) | `libcadlib.dylib` | macOS 11 or newer |
| Windows | x86_64 | `cadlib.dll` + `cadlib.lib` | Windows 10 or newer, VC++ 2015–2022 redistributable |

Alpine and other musl-based distributions are not supported: the library is
built for glibc and cannot load on musl. Intel Macs and Windows on Arm are not
supported either.

## License

MIT, see [LICENSE](LICENSE).
