# blaim

The name `blaim` is a play on `git blame`, with the spelling altered to emphasize the role of AI.

This package implements a CLI tool for processing the "accepted suggestions" logs from a custom VS Code extension. 

The extension is a [`InlineCompletionItemProvider`](https://code.visualstudio.com/api/references/vscode-api#InlineCompletionItemProvider) that queries a generative model (e.g. OLlama) for code snippets to present as inline suggestions, and logs some information every time the user accepts one of these suggestions.

[TODO: move the VS Code extension source to this repo]. 

