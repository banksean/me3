import * as vscode from "vscode";
import { getAcceptedSuggestionsForFile, AcceptLogLine } from "./acceptlog";

const MIN_PARTIAL_MATCH_SIZE = 8;

// this method is called when vs code is activated
export function activateDecorators(context: vscode.ExtensionContext) {
  console.log("decorator sample is activated");

  const toggleCallback = async function name(params: any) {
    console.log("toggle", params);
  };

  const toggle = vscode.commands.registerCommand(
    "coverage-gutters.toggleCoverage",
    toggleCallback,
  );

  context.subscriptions.push(toggle);

  let timeout: string | number | NodeJS.Timeout | undefined = undefined;

  const findMatchingRanges = function (
    searchString: string,
    activeEditor: vscode.TextEditor,
  ): vscode.Range[] {
    const ranges: vscode.Range[] = [];
    let acceptTextIndexInDoc = activeEditor.document
      .getText()
      .indexOf(searchString);
    while (acceptTextIndexInDoc != -1) {
      ranges.push(
        new vscode.Range(
          activeEditor.document.positionAt(acceptTextIndexInDoc),
          activeEditor.document.positionAt(
            acceptTextIndexInDoc + searchString.length,
          ),
        ),
      );
      acceptTextIndexInDoc = activeEditor?.document
        .getText()
        .indexOf(searchString, acceptTextIndexInDoc + 1);
    }
    return ranges;
  };

  const inlineDecorationType = vscode.window.createTextEditorDecorationType({
    rangeBehavior: vscode.DecorationRangeBehavior.ClosedClosed,
    // Make the AI-generated text glow like it's magical, or radioactive :)
    textDecoration:
      "opacity:1.0; text-shadow: 0 0 6px rgba(255, 255, 255, 0.75)",
    overviewRulerColor: "blue",
    overviewRulerLane: vscode.OverviewRulerLane.Right,
    light: {
      borderWidth: "1px",
      borderStyle: "solid",
      borderColor: "darkblue",
      backgroundColor: "lightgray",
    },
    dark: {},
  });

  let activeEditor = vscode.window.activeTextEditor;

  function updateDecorations() {
    if (!activeEditor) {
      return;
    }
    const text = activeEditor.document.getText();
    const inlineDecorations: vscode.DecorationOptions[] = [];
    const gutterBlameAnnotations: vscode.DecorationOptions[] = [];

    // Inline annotations to hilight ranges of AI-generated code.
    // This is currently very broken and does not work well at all.
    const accepts = getAcceptedSuggestionsForFile(
      activeEditor.document.fileName,
    );
    for (let i = 0; i < accepts.length; i++) {
      const acceptLine = accepts[i];
      const acceptedLines = acceptLine.text.split("\n");
      const endPos = new vscode.Position(
        acceptLine.position.line + acceptedLines.length,
        acceptedLines?.pop()?.length || 0,
      );
      const range = new vscode.Range(acceptLine.position, endPos);

      // At this point, we have the originally logged suggestion text
      // and the position where it was inserted in the original file.
      //
      // However, this is not enough information to let us higlight the
      // *current* location of the suggested text in the file, if it is still
      // there at all. For instance:
      //
      // - The user may accept a suggestion, and subsequently edit or remove
      //   the inserted text altogether.
      // - The user may accept a suggestion, and then make edits inserting new
      //   lines at a previous line in the file, such that the logged
      // . line number of the inserted text is no longer accurate.
      // - The user may accept a suggestion and then delete some prefix or
      //   suffix of the suggested text.
      //
      // This code does nothing fancy to find fuzzy matches etc.  It looks
      // for an exact match at the original location, and as a series of
      // fall-back searches if no matches are found so far:
      //
      //  1. looks for exact matches at any location (besides the original) in the document,
      //  2. looks for prefix matches (user deleted code from the end of the suggestion)
      //  3. looks for suffix matches (user deleted code from the beginning of the suggestion)
      //
      // This is admittedly inexact and subject to a lot of false negatives (in
      // particular, if the user made edits somewhere in the middle of a block)
      // and false positives (user accepted a rote and trivial suggestion like
      // `if err == nil { return err }` in Go.

      // Look for exact matches of the accepted text at the original location:
      const actualTextAfterEdits = activeEditor.document.getText(range);
      if (actualTextAfterEdits == acceptLine.text) {
        const decoration: vscode.DecorationOptions = {
          range: range,
          hoverMessage: new vscode.MarkdownString(
            "## This is AI-generated code:\n```" +
              JSON.stringify(acceptLine.inferenceConfig) +
              "```\n" +
              `originally at ${JSON.stringify(acceptLine.position)}\n` +
              "## Raw text from the model:\n```\n" +
              acceptLine.text +
              "\n```\n",
          ),
        };
        inlineDecorations.push(decoration);
      } else {
        // We didn't the exact text at the original location, so try this series of back-ups.
        // TODO: research other approaches, e.g. edit distance, suffix arrays etc.  This is
        // all naive brute force and probably very far from optimal.

        let prefixMatchRanges: vscode.Range[] = [];
        let suffixMatchRanges: vscode.Range[] = [];

        // Look for exact matches of the accepted text at any other location in the file:
        let ranges = findMatchingRanges(acceptLine.text, activeEditor);

        // Seach for prefixes of acceptLine.text, if there were no exact matches:
        for (
          let len = acceptLine.text.length - 2;
          len > MIN_PARTIAL_MATCH_SIZE && prefixMatchRanges.length == 0;
          len--
        ) {
          const prefix = acceptLine.text.substring(0, len);
          prefixMatchRanges = findMatchingRanges(prefix, activeEditor);
        }

        // Seach for suffixes of acceptLine.text, if there were no exact matches or prefix matches:
        for (
          let len = 1;
          len < acceptLine.text.length - MIN_PARTIAL_MATCH_SIZE &&
          suffixMatchRanges.length == 0;
          len++
        ) {
          const suffix = acceptLine.text.substring(len);
          suffixMatchRanges = findMatchingRanges(suffix, activeEditor);
        }
        if (ranges.length == 0) {
          ranges = suffixMatchRanges.concat(prefixMatchRanges);
        }
        for (let i = 0; i < ranges.length; i++) {
          const decoration: vscode.DecorationOptions = {
            range: ranges[i],
            hoverMessage: new vscode.MarkdownString(
              "## This is AI-generated code\nInference config:```" +
                `${JSON.stringify(acceptLine.inferenceConfig)}` +
                "```\n" +
                `Originally inserted at position: ${JSON.stringify(acceptLine.position)}\n` +
                "## Raw text from the model:\n```\n" +
                acceptLine.text +
                "\n```\n",
            ),
          };
          inlineDecorations.push(decoration);
        }
      }

      // Line-by-line gutter annotations
      const textLines = text.split("\n");
      for (let lineNumber = 0; lineNumber < textLines.length; lineNumber++) {
        const linePos = new vscode.Position(lineNumber, 0);
        const gutterAnnotation: vscode.DecorationOptions = {
          range: new vscode.Range(linePos, linePos),
          renderOptions: {
            before: {
              contentText: `blaim line ${lineNumber}`,
            },
          },
        };
        gutterBlameAnnotations.push(gutterAnnotation);
      }

      // Clear any existing decorations.
      activeEditor.setDecorations(inlineDecorationType, []);
      activeEditor.setDecorations(inlineDecorationType, inlineDecorations);
    }
  }

  function triggerUpdateDecorations(throttle = false) {
    if (timeout) {
      clearTimeout(timeout);
      timeout = undefined;
    }
    if (throttle) {
      timeout = setTimeout(updateDecorations, 500);
    } else {
      updateDecorations();
    }
  }

  if (activeEditor) {
    triggerUpdateDecorations();
  }

  vscode.window.onDidChangeActiveTextEditor(
    (editor) => {
      activeEditor = editor;
      if (editor) {
        triggerUpdateDecorations();
      }
    },
    null,
    context.subscriptions,
  );

  vscode.workspace.onDidChangeTextDocument(
    (event) => {
      if (activeEditor && event.document === activeEditor.document) {
        triggerUpdateDecorations(true);
      }
    },
    null,
    context.subscriptions,
  );
}
