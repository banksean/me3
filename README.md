# me3: A monorepo for personal projects

See [ABOUT.md](./ABOUT.md) for background information, why this repo exists, why it's set up this way, etc.

## Use `bazel` to build, test and execute all targets

Build targets and dependencies in this repo are managed by [bazel](https://bazel.build/), so there are `BUILD.bazel` files located in most directories. 

## Use `make` helpers for a few selected tasks
I've included a top-level [`Makefile`](./Makefile) with some helper targets:
`make gazelle`, `make gofmt` and so on just to save some typing at the terminal.

## Setup
- git clone this repo
- install [bazelisk](https://github.com/bazelbuild/bazelisk?tab=readme-ov-file#installation)
- Now `bazel build //...` should work

### OSX Notes
- `brew install coreutils` (in order to install the `realpath` cli command, which is used by `bazel run //:go`).

You should now be able to `bazel build //...`, `bazel run //...` etc whatever targets are defined in this repo.

## Common development workflow tasks

### Build/test/run something

Commands for updating, building, testing and running things in this repo:

- build: `bazel build //<target>`
- test: `bazel test //<target>`
- run: `bazel run //<target>`

### Running a standard `go` command
Instead of running

```go <command> <args...>```

run this form:

```bazel run //:go -- <command> <args...>``` 

Doing so will make sure your `go` command is use the same toolchain and environment that `bazel` does when it deals with your go targets.

### Manage go package dependencies
This is a little more complicated than it might be in a purely go-based project repo, but it's still pretty straightforward:

Suppose you want to use an external go package. We'll use `github.com/urfave/cli/v2` as an example:

- run `bazel run //:go -- get github.com/urfave/cli/v2` to update `go.mod`
- add `"github.com/urfave/cli/v2"` to your `.go` file's `import`s
- run `make gazelle` to update `go_modules.bzl` based on the aforementioned update to `go.mod`

Note that `bazel build` does not actually look at the contents of `go.mod` - it uses `go_dependencies.bzl`, which `make gazelle` generates from `go.mod`.

### Manage NPM package dependencies
Use `bazel run @pnpm` instead of `npm`:

```
bazel run @pnpm -- install --dir $(pwd) <package name> [options...]
```

e.g.:
```
bazel run @pnpm --  install --dir $(pwd) -D prettier -w 
```

### Manage python package dependencies

Edit `requirements.in` (NOT `requirements.txt`) to reflect your updated *direct* dependencies - packages your code actually imports, directly, by name.

- Add the name you would normally use for `pip install <package-name>` to `requirements.in`.
- You probably need to figure out which version of the package you need and also specify it on that same line, separated by `==`. See `requirements.in` for examples.
- Run `bazel run //:requirements.update` in the root directory of this repo.  This will generate the full `requirements.txt` file including all the direct dependency's transitive dependencies as well.
- In your `py_binary` or `py_library` etc target, you'll need to add new lines to the `deps=[...]` list: `requirement("<package name>"),`. (You may need to add `load("@pip_deps//:requirements.bzl", "requirement")` to that BUILD file if it's not there already.)

### Other helpful bits

Maintaining `BUILD` targets by hand can be a pain. To automatically update `BUILD.bazel` files based on source import statements (and protobuf options, etc): run `bazel run //:gazelle`.