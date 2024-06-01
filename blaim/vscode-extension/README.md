# VS Code Extension for blAIm proof-of-concept

This extension demonstrates how one may instrument their development environment to log the necessary UI events to later determine which parts of code changes came from AI-generated suggestions.

It's not published anywhere (it's just a PoC) so you'll have to run it from your dev environment.

To debug this extension:
- From VS Code: Open a new window with this directory as its root.
- From the menu bar, Run > Start Debugging (or F5).
- This launches in a separate Extension Development Host window which looks like a new VS Code window.
- On startup, the extension will tell you where it's logging the accepted code suggestions, in output like this (check your VS Code debug console, not the Extension Development Host window):

``
blaim-completion started
writing accepted suggestions logs: export ACCEPT_LOG=/<...>/accepted.suggestions.log
``

- Make some edits to `playground.js` (using the Extension Development Host, not your main VS Code window), and make sure to accept some generated code suggestions.
- Copy that `export ACCEPT_LOG=...` line from the debug console
- Paste it into a shell window and run `git diff | bazel run //blaim/cmd` to verify that it can identify the AI-generated parts of the current set of code changes.

