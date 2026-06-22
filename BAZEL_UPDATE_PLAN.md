# Bazel Modernization Plan

This repo currently builds with the pinned Bazel version in `.bazelversion`
(`7.0.2`), but it is not ready for modern Bazel defaults. The main issue is
that `MODULE.bazel` and `WORKSPACE` are both active: JavaScript/TypeScript and
some C++ dependencies are configured through Bzlmod, while Go, Python, protobuf,
and Ruff still come from `WORKSPACE`.

Bazel 8 disables `WORKSPACE` by default, and Bazel 9 removes it. A quick probe
with Bazel 8.7.0 fails immediately because `@rules_python` is introduced only
from `WORKSPACE`.

## Goal

Modernize the repo so that:

- `bazel build //...` works on Bazel 8 without `--enable_workspace`.
- The repo no longer depends on `WORKSPACE` for external repositories.
- The pinned Bazel version can be bumped safely, ideally to a current Bazel 9
  release after the Bazel 8 migration is clean.
- Deprecated and noisy BUILD patterns are cleaned up while keeping target
  behavior stable.

## Proposed Sequence

### 1. Finish the Bzlmod migration

Move the remaining external dependency setup from `WORKSPACE` into
`MODULE.bazel`.

Specific work:

- Replace `rules_go` and Gazelle `http_archive` setup with Bzlmod
  `bazel_dep(...)` entries and module extensions.
- Replace generated `go_repositories.bzl` usage with Gazelle's Bzlmod
  `go_deps.from_file(go_mod = "//:go.mod")` flow.
- Register the Go SDK through the `rules_go` Bzlmod extension instead of
  `go_register_toolchains(...)` from `WORKSPACE`.
- Replace `rules_python` `http_archive`, `python_register_toolchains`, and
  `pip_parse` setup with current `rules_python` Bzlmod extensions.
- Stop pinning `com_google_protobuf` 3.14.0 from `WORKSPACE`; use the BCR
  `protobuf` module if a direct protobuf dependency is still needed.
- Replace `aspect_rules_lint` `fetch_ruff()` from `WORKSPACE` with the current
  Bzlmod setup for Ruff/tooling.
- Reduce `WORKSPACE` to a compatibility stub or remove it once Bazel 8 builds
  without `--enable_workspace`.

Validation:

```sh
USE_BAZEL_VERSION=8.7.0 bazel build //...
```

The target state is that this command passes without `--enable_workspace`.

### 2. Update rule versions

Several rule versions are old or prerelease. Update them after the initial
Bzlmod migration, so failures are easier to attribute.

Current notable pins:

- `rules_go` 0.44.2 in `WORKSPACE`
- `bazel_gazelle` 0.35.0 in `WORKSPACE`
- `rules_python` 0.29.0 in `WORKSPACE`
- `aspect_rules_js` 2.0.0-rc1
- `aspect_rules_ts` 3.0.0-rc0
- `aspect_rules_jest` 0.22.0-rc0
- `aspect_rules_lint` 1.0.0-rc1
- `rules_nodejs` 6.1.0
- `google_benchmark` 1.8.3
- `googletest` direct pin at 1.12.1, though module resolution currently picks
  1.14.0.bcr.1

Use the Bazel Central Registry as the source of truth for current compatible
versions. Update one ecosystem at a time:

1. Go and Gazelle.
2. Python and pip dependencies.
3. JS/TS/Jest/Lint.
4. C++ support libraries.

Validation after each ecosystem:

```sh
bazel mod tidy
bazel mod graph
bazel build //...
bazel test //...
```

### 3. Bump Bazel in stages

Do not jump straight from Bazel 7 to Bazel 9 before the Bzlmod migration is
complete.

Recommended sequence:

1. Keep `.bazelversion` at `7.0.2` while moving dependencies to `MODULE.bazel`.
2. Verify with Bazel 8 using `USE_BAZEL_VERSION=8.7.0`.
3. Update `.bazelversion` to a stable Bazel 8 release and commit that state.
4. Probe Bazel 9.
5. Fix any remaining Starlarkification or removed-feature issues.
6. Update `.bazelversion` to Bazel 9 when `bazel build //...` and
   `bazel test //...` pass.

Useful checks:

```sh
USE_BAZEL_VERSION=8.7.0 bazel build //...
USE_BAZEL_VERSION=8.7.0 bazel test //...
USE_BAZEL_VERSION=9.x.x bazel build //...
USE_BAZEL_VERSION=9.x.x bazel test //...
```

## Cleanup Items

These are not the main blockers, but they should be cleaned up during the
migration.

- Add `npmrc = "//:.npmrc"` to `npm.npm_translate_lock(...)` in `MODULE.bazel`
  to suppress the implicit `.npmrc` warning.
- Update the direct `googletest` pin in `MODULE.bazel`; `bazel mod graph`
  currently warns that the root asks for 1.12.1 while resolution selects
  1.14.0.bcr.1.
- Remove unused root BUILD loads:
  - `py_binary`
  - `py_test`
  - `py_console_script_binary`
- Replace constant globs in `cpp/lib/BUILD.bazel`:
  - `glob(["word_generator.cc"])`
  - `glob(["word_generator.h"])`
- Revisit labels that reference `@go_sdk`, especially:
  - root `:gofmt`
  - `bazel/tools/format:format`
  - `bazel/tools/go:go`
  The Bzlmod Go setup may expose the SDK repository differently.
- Consider replacing the hand-written `bazel/tools/go` wrapper if current
  `rules_go` offers a simpler supported way to run the pinned Go toolchain.
- Run Buildifier after each migration step.

## Acceptance Criteria

The migration is complete when:

```sh
bazel mod tidy
bazel build //...
bazel test //...
USE_BAZEL_VERSION=8.7.0 bazel build //...
USE_BAZEL_VERSION=8.7.0 bazel test //...
```

all pass, and `WORKSPACE` is no longer required.

After that, Bazel 9 migration can be treated as a smaller follow-up focused on
remaining removed APIs or rule-version compatibility.
