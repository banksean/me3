# blaim

The name `blaim` is a play on `git blame`, with the spelling altered to emphasize the role of AI.

This package implements a CLI tool for processing the "accepted suggestions" logs from a custom VS Code extension implemented [here](./vscode-extension/)


The [`//blaim/cmd`](../cmd/) CLI tool can compare changes in the current git diff against the accepted suggestions log to identify text ranges for AI-generated code:

```
> export ACCEPT_LOG=[...]
> git diff | bazel run //blaim/cmd
[...]
blaim/vscode-extension/playground.js: 9 accept events, 1 diff hunks
found a matching accept log entry for "onacci(n:" starting at position 12 on line 3 of blaim/vscode-extension/playground.js:
function fibonacci(n: number) {
    if (n <= 1) return n;
    return fibonacci(n - 2) + fibonacci(n - 1);
}
[...]
```