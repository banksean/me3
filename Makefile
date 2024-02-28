
#gazelle:
#	bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=go_repositories.bzl%go_repositories -prune

# BAZEL defines which executable to run.
#
# If BAZEL isn't defined then try to define it as bazelisk (if the executable exists).
# Otherwise, define it as bazel.
ifeq ($(BAZEL),)
	ifneq ($(strip $(shell which bazelisk)),)
		BAZEL := bazelisk
	else
		BAZEL := bazel
	endif
endif

.PHONY: update-go-bazel-files
update-go-bazel-files:
	$(BAZEL) run //:gazelle -- update ./

.PHONY: update-go-bazel-deps
update-go-bazel-deps:
	$(BAZEL) run //:gazelle -- update-repos -from_file=go.mod -to_macro=go_repositories.bzl%go_repositories -prune

.PHONY: gazelle
gazelle: update-go-bazel-deps update-go-bazel-files

.PHONY: update-py-deps
update-py-deps:
	$(BAZEL) run //:requirements.update

.PHONY: buildifier
buildifier:
	$(BAZEL) run //:buildifier

.PHONY: gofmt
gofmt:
	$(BAZEL) run //:gofmt -- -s -w .

.PHONY: gencommitmsg
gencommitmsg:
	$(BAZEL) run //gencommitmsg -- $(PWD)
