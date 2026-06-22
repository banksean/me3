# Bazel Modernization Plan

This repo originally built with the pinned Bazel version in `.bazelversion`
(`7.0.2`), but it was not ready for modern Bazel defaults. The main issue was
that `MODULE.bazel` and `WORKSPACE` were both active: JavaScript/TypeScript and
some C++ dependencies were configured through Bzlmod, while Go, Python,
protobuf, and Ruff still came from `WORKSPACE`.

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

## Progress

- 2026-06-22: Completed baseline hygiene cleanup. Made `.npmrc` explicit in
  `npm.npm_translate_lock(...)`, aligned the direct `googletest` pin with the
  resolved module version, removed unused root BUILD loads, and replaced
  constant C++ globs with explicit file lists.
- 2026-06-22: Migrated Go, Gazelle, Python, pip dependencies, and Ruff binaries
  from `WORKSPACE` to `MODULE.bazel`, then bumped `.bazelversion` to Bazel
  8.7.0. The generated `go_repositories.bzl` macro is no longer used; Go
  dependencies now come from `go.mod` through Gazelle's
  `go_deps.from_file(...)` extension. `WORKSPACE` is now only a compatibility
  stub. The first attempt used a downloaded Go 1.24.12 SDK because current
  Gazelle builds helper tools from a `go.work` requiring Go 1.24.12 or newer,
  but SDK archive downloads were too slow in this environment. The migration
  now uses the host Go 1.26.4 SDK through `go_sdk.host(...)`; this unblocks
  Bazel 8 builds, but is less reproducible than a downloaded SDK and should be
  revisited after the rest of the migration is green.
- 2026-06-22: Validated the Bzlmod migration with `bazel mod tidy`,
  `bazel query @go_sdk//:ROOT`, `bazel build //cpp/...`, `bazel build //journ`,
  `bazel build //codenames-go //gencommitmsg //blaim/...`, and a workspace
  Buildifier check using the prebuilt Buildifier binary on Bazel 8.7.0.
- 2026-06-22: Updated the JS/TS/Jest/Lint rule stack off prerelease pins while
  keeping compatibility with the existing pnpm lockfile v6. `aspect_rules_js`
  3.x currently requires pnpm lockfile v9, and this container's system npm/pnpm
  tooling is not usable for regenerating `pnpm-lock.yaml`, so the lockfile v9
  migration remains a follow-up. Validated with `bazel mod tidy`, `bazel mod
  graph`, `bazel build //pool/...`, `bazel test //pool:test`, `bazel build
  //bazel/tools/format:format`, `bazel build //:format`, and a workspace
  Buildifier check.
- 2026-06-22: Updated C++ support libraries: `googletest` to 1.17.0.bcr.2,
  `google_benchmark` to 1.9.5, and `abseil-cpp` to 20250512.1. Validated with
  `bazel mod tidy`, `bazel test //cpp/...`, and a workspace Buildifier check.

### 1. Finish the Bzlmod migration

Move the remaining external dependency setup from `WORKSPACE` into
`MODULE.bazel`.

Specific work:

- Completed: Replace `rules_go` and Gazelle `http_archive` setup with Bzlmod
  `bazel_dep(...)` entries and module extensions.
- Completed: Replace generated `go_repositories.bzl` usage with Gazelle's Bzlmod
  `go_deps.from_file(go_mod = "//:go.mod")` flow.
- Completed: Register the Go SDK through the `rules_go` Bzlmod extension instead
  of `go_register_toolchains(...)` from `WORKSPACE`.
- Completed: Replace `rules_python` `http_archive`, `python_register_toolchains`, and
  `pip_parse` setup with current `rules_python` Bzlmod extensions.
- Completed: Stop pinning `com_google_protobuf` 3.14.0 from `WORKSPACE`; use the BCR
  `protobuf` module if a direct protobuf dependency is still needed.
- Completed: Replace `aspect_rules_lint` `fetch_ruff()` from `WORKSPACE` with the current
  Bzlmod setup for Ruff/tooling.
- Completed: Reduce `WORKSPACE` to a compatibility stub or remove it once Bazel 8 builds
  without `--enable_workspace`.

Validation:

```sh
USE_BAZEL_VERSION=8.7.0 bazel build //...
```

The target state is that this command passes without `--enable_workspace`.
Current targeted validation passes for C++ and Go targets; full-repo build and
test validation are still pending because Python and JS/TS targets have not yet
been exercised after the Bzlmod migration.

### 2. Update rule versions

Several rule versions are old or prerelease. Update them after the initial
Bzlmod migration, so failures are easier to attribute.

Current notable pins:

- `rules_go` 0.61.1 in `MODULE.bazel`
- `gazelle` 0.51.3 in `MODULE.bazel`
- `rules_python` 2.0.2 in `MODULE.bazel`
- `buildifier_prebuilt` 6.4.0 in `MODULE.bazel`
- `aspect_bazel_lib` 2.22.0
- `aspect_rules_js` 2.4.2
- `aspect_rules_ts` 3.6.3
- `aspect_rules_jest` 0.25.2
- `aspect_rules_lint` 2.7.1
- `rules_nodejs` 6.7.3
- `google_benchmark` 1.9.5
- `googletest` 1.17.0.bcr.2
- `abseil-cpp` 20250512.1

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

1. Completed: Move dependencies to `MODULE.bazel`.
2. Completed: Verify with Bazel 8 using `USE_BAZEL_VERSION=8.7.0`.
3. Completed: Update `.bazelversion` to a stable Bazel 8 release and commit that state.
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

- Completed: Add `npmrc = "//:.npmrc"` to `npm.npm_translate_lock(...)` in `MODULE.bazel`
  to suppress the implicit `.npmrc` warning.
- Completed: Update the direct `googletest`, `google_benchmark`, and
  `abseil-cpp` pins in `MODULE.bazel`.
- Completed: Remove unused root BUILD loads:
  - `py_binary`
  - `py_test`
  - `py_console_script_binary`
- Completed: Replace constant globs in `cpp/lib/BUILD.bazel`:
  - `glob(["word_generator.cc"])`
  - `glob(["word_generator.h"])`
- Completed: Revisit labels that reference `@go_sdk`, especially:
  - root `:gofmt`
  - `bazel/tools/format:format`
  - `bazel/tools/go:go`
  The Bzlmod Go setup may expose the SDK repository differently.
- Consider replacing the hand-written `bazel/tools/go` wrapper if current
  `rules_go` offers a simpler supported way to run the pinned Go toolchain.
- Revisit `go_sdk.host(...)` and move back to `go_sdk.download(...)` once Go SDK
  archive downloads are reliable in the developer/CI environment.
- Completed: Move `//:buildifier` from the Go `buildtools` dependency to the
  prebuilt Buildifier binary. The old Go dependency exposed an internal
  `@rules_proto` repo-mapping failure under Bzlmod.
- Completed: Move the JS/TS/Jest/Lint rules off prerelease pins:
  - `aspect_rules_js` 2.4.2
  - `aspect_rules_ts` 3.6.3
  - `aspect_rules_jest` 0.25.2
  - `aspect_rules_lint` 2.7.1
- Regenerate `pnpm-lock.yaml` with pnpm 9 or newer, then move to
  `aspect_rules_js` 3.x and compatible `aspect_rules_ts`/`rules_nodejs` pins.
- Run Buildifier after each migration step. With the current prebuilt alias, run
  `bazel build //:buildifier` and then invoke the generated binary over the
  desired file list.

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
